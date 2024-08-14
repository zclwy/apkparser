package apkparser

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/avast/apkverifier"
)

const (
	androidExt = ".apk"
)

type AppInfo struct {
	Name        string      `json:"name,omitempty"`        // 应用名称
	BundleId    string      `json:"bundleId,omitempty"`    // 包名
	Version     string      `json:"version,omitempty"`     // 版本名称
	Build       int         `json:"build,omitempty"`       // 版本号
	Icon        image.Image `json:"icon,omitempty"`        // app icon
	Size        int64       `json:"size,omitempty"`        // app size in bytes
	CertInfo    CertInfo    `json:"certInfo,omitempty"`    // app 证书信息
	Md5         string      `json:"md5,omitempty"`         // app md5
	SupportOS64 bool        `json:"supportOS64,omitempty"` // 是否支持64位
	SupportOS32 bool        `json:"supportOS32,omitempty"` // 是否支持32位
	Permissions []string    `json:"permissions,omitempty"` // 权限列表
}
type CertInfo struct {
	Md5                string    `json:"md5,omitempty"`
	Sha1               string    `json:"sha1,omitempty"`
	Sha256             string    `json:"sha256,omitempty"`
	ValidFrom          time.Time `json:"validFrom,omitempty"`
	ValidTo            time.Time `json:"validTo,omitempty"`
	Issuer             string    `json:"issuer,omitempty"`
	Subject            string    `json:"subject,omitempty"`
	SignatureAlgorithm string    `json:"signatureAlgorithm,omitempty"`
	SerialNumber       *big.Int  `json:"serialNumber,omitempty"`
}

type androidManifest struct {
	Package     string       `xml:"package,attr"`
	VersionName string       `xml:"versionName,attr"`
	VersionCode string       `xml:"versionCode,attr"`
	Permissions []Permission `xml:"uses-permission"`
}
type Permission struct {
	Name string `xml:"name,attr"`
}

func New(name string) (*AppInfo, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	} else if filepath.Ext(stat.Name()) != androidExt {
		return nil, errors.New("unknown platform")
	}

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	var (
		xmlFile     *zip.File
		supportOS64 bool
		supportOS32 bool
		hasSoFile   bool
	)
	for _, f := range reader.File {
		switch f.Name {
		case "AndroidManifest.xml":
			xmlFile = f
		}
		if strings.HasSuffix(f.Name, ".so") {
			hasSoFile = true
		}
		if strings.HasPrefix(f.Name, "lib/arm64-v8a") {
			supportOS64 = true
		}
		if strings.HasPrefix(f.Name, "lib/armeabi") {
			supportOS32 = true
		}
	}
	info, errParse := parseApkFile(xmlFile)
	if errParse != nil {
		return nil, errParse
	}
	// 当前apk支持的系统位数
	if !hasSoFile && !supportOS64 && !supportOS32 {
		info.SupportOS64 = true
		info.SupportOS32 = true
	} else {
		info.SupportOS64 = supportOS64
		info.SupportOS32 = supportOS32
	}
	info.Md5, _ = getApkMd5(file)
	info.CertInfo, _ = getSignature(name)

	icon, label, errExtra := parseApkIconAndLabel(name)
	if errExtra != nil {
		return nil, errExtra
	}
	info.Name = label
	info.Icon = icon
	info.Size = stat.Size()

	return info, err
}

// 解析apk文件
func parseApkFile(xmlFile *zip.File) (*AppInfo, error) {
	if xmlFile == nil {
		return nil, errors.New("AndroidManifest.xml not found")
	}

	manifest, err := parseAndroidManifest(xmlFile)
	if err != nil {
		return nil, err
	}

	info := new(AppInfo)
	versionCode, _ := strconv.Atoi(manifest.VersionCode)

	info.BundleId = manifest.Package
	info.Version = manifest.VersionName
	info.Build = versionCode

	for _, v := range manifest.Permissions {
		info.Permissions = append(info.Permissions, v.Name)
	}

	return info, nil
}

// 解析AndroidManifest.xml文件
func parseAndroidManifest(xmlFile *zip.File) (*androidManifest, error) {
	rc, err := xmlFile.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rc.Close()
	}()

	buf, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	xmlContent, err := NewXMLFile(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	manifest := new(androidManifest)
	decoder := xml.NewDecoder(xmlContent.Reader())
	if err := decoder.Decode(manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

// 解析apk图标和名称
func parseApkIconAndLabel(name string) (image.Image, string, error) {
	pkg, err := openFile(name)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = pkg.close()
	}()

	icon, _ := pkg.icon(&ResTableConfig{
		Density: 720,
	})

	// label, _ := pkg.label(&ResTableConfig{
	// 	// Language: [2]uint8{'z', 'h'},
	// 	// Country: [2]uint8{'C', 'N'},
	// })
	label, _ := pkg.label(&ResTableConfig{})

	return icon, label, nil
}

// 获取apk md5
func getApkMd5(file *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%032x", hash.Sum(nil)), nil
}

// 获取apk签名
func getSignature(apkPath string) (CertInfo, error) {
	res, err := apkverifier.Verify(apkPath, nil)
	if err != nil {
		return CertInfo{}, err
	}

	cert, _ := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		return CertInfo{}, errors.New("No certificate found")
	}

	return CertInfo{
		Md5:                cert.Md5,
		Sha1:               cert.Sha1,
		Sha256:             cert.Sha256,
		ValidFrom:          cert.ValidFrom,
		ValidTo:            cert.ValidTo,
		Issuer:             cert.Issuer,
		Subject:            cert.Subject,
		SignatureAlgorithm: cert.SignatureAlgorithm,
		SerialNumber:       cert.SerialNumber,
	}, nil
}

func parseSignature(f *zip.File) {
	rc, err := f.Open()
	if err != nil {
		log.Printf("failed to open file %s: %v", f.Name, err)
		return
	}
	defer rc.Close()

	// 读取签名文件内容
	content, err := io.ReadAll(rc)
	if err != nil {
		log.Printf("failed to read file %s: %v", f.Name, err)
		return
	}

	// 解码 PEM 数据
	block, _ := pem.Decode(content)
	if block == nil {
		log.Printf("failed to decode PEM block from file %s", f.Name)
		return
	}

	// 解析证书
	certs, err := x509.ParseCertificates(content)
	if err != nil {
		log.Printf("failed to parse certificate from file %s: %v", f.Name, err)
		return
	}

	// 打印签名信息
	for _, cert := range certs {
		fmt.Printf("Certificate Subject: %s\n", cert.Subject)
		fmt.Printf("Issuer: %s\n", cert.Issuer)
		fmt.Printf("Serial Number: %s\n", cert.SerialNumber)
		fmt.Printf("Not Before: %s\n", cert.NotBefore)
		fmt.Printf("Not After: %s\n", cert.NotAfter)
		fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm)
		fmt.Printf("Public Key Algorithm: %s\n", cert.PublicKeyAlgorithm)
	}

}
