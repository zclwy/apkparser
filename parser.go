package apkparser

import (
	"errors"
	"image"
	"math"
	"os"

	ap "github.com/avast/apkparser"
	"github.com/avast/apkverifier"
)

const (
	androidExt = ".apk"
)

type AppInfo struct {
	Name             string      `json:"name,omitempty"`     // 应用名称
	BundleId         string      `json:"bundleId,omitempty"` // 包名
	Version          string      `json:"version,omitempty"`  // 版本名称
	Build            int64       `json:"build,omitempty"`    // 版本号
	Icon             image.Image `json:"icon,omitempty"`     // app icon
	Size             int64       `json:"size,omitempty"`     // app size in bytes
	CertInfo         *CertInfo   `json:"certInfo,omitempty"` // app 证书信息
	Md5              string      `json:"md5,omitempty"`      // app md5
	Sha1             string      `json:"sha1,omitempty"`
	Sha256           string      `json:"sha256,omitempty"`
	SupportOS64      bool        `json:"supportOS64,omitempty"` // 是否支持64位
	SupportOS32      bool        `json:"supportOS32,omitempty"` // 是否支持32位
	Permissions      []string    `json:"permissions,omitempty"` // 权限列表
	MinSdkVersion    int         `json:"minSdkVersion"`         // 最小兼容rom版本
	MaxSdkVersion    int         `json:"maxSdkVersion"`         // 最大兼容rom版本
	TargetSdkVersion int         `json:"targetSdkVersion"`      // 推荐rom版本
}
type CertInfo struct {
	Md5    string `json:"md5,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
	// ValidFrom          time.Time `json:"validFrom,omitempty"`
	// ValidTo            time.Time `json:"validTo,omitempty"`
	// Issuer             string    `json:"issuer,omitempty"`
	// Subject            string    `json:"subject,omitempty"`
	// SignatureAlgorithm string    `json:"signatureAlgorithm,omitempty"`
	// SerialNumber       *big.Int  `json:"serialNumber,omitempty"`
}

type Option struct {
	WithSignature        bool // 是否需要获取签名信息
	IgnoreSignatureError bool // 是否忽略签名错误，默认不忽略
	WithIcon             bool // 是否需要获取icon信息
	WithMd5              bool
	WithSha1             bool
	WithSha256           bool
}

func New(name string, option Option) (*AppInfo, error) {
	infoApk, err := openFile(name)
	if err != nil {
		return nil, err
	}
	// 释放资源
	defer infoApk.close()
	return getApkInfo(infoApk, option)
}

func NewFromReader(f *os.File, option Option) (*AppInfo, error) {
	infoApk, err := openReader(f)
	if err != nil {
		return nil, err
	}
	return getApkInfo(infoApk, option)
}

func getApkInfo(infoApk *apk, option Option) (*AppInfo, error) {
	getApkHash(infoApk, option)
	info := &AppInfo{
		Name:             infoApk.parseApkLabel(),
		BundleId:         infoApk.apkManifest.Package,
		Version:          infoApk.apkManifest.VersionName,
		Build:            infoApk.apkManifest.VersionCode,
		Size:             infoApk.size,
		Md5:              infoApk.md5,
		Sha1:             infoApk.sha1,
		Sha256:           infoApk.sha256,
		SupportOS64:      infoApk.supportOs64,
		SupportOS32:      infoApk.supportOs32,
		Permissions:      formatPermissions(infoApk.apkManifest.Permissions),
		MinSdkVersion:    infoApk.apkManifest.SDK.Min,
		MaxSdkVersion:    infoApk.apkManifest.SDK.Max,
		TargetSdkVersion: infoApk.apkManifest.SDK.Target,
	}

	// 获取证书信息
	if option.WithSignature {
		certInfo, errCert := getSignature(infoApk)
		if errCert != nil {
			if !option.IgnoreSignatureError {
				return nil, errCert
			}
		} else {
			info.CertInfo = certInfo
		}
	}
	if option.WithIcon {
		// 获取icon信息
		info.Icon = infoApk.parseApkIcon()
	}
	return info, nil
}

// 获取apk签名
func getSignature(apk *apk) (*CertInfo, error) {
	// res, err := apkverifier.Verify(apkPath, nil)
	optionalZip, err := ap.OpenZipReader(apk.f)
	if err != nil {
		return nil, err
	}
	defer optionalZip.Close()
	maxSdkVersion := apk.apkManifest.SDK.Max
	if maxSdkVersion == 0 {
		maxSdkVersion = math.MaxInt32
	}
	res, err := apkverifier.VerifyWithSdkVersionReader(
		apk.f,
		optionalZip,
		int32(apk.apkManifest.SDK.Min),
		int32(maxSdkVersion),
	)
	if err != nil {
		return nil, err
	}

	cert, _ := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		return nil, errors.New("no certificate found")
	}

	return &CertInfo{
		Md5:    cert.Md5,
		Sha1:   cert.Sha1,
		Sha256: cert.Sha256,
		// ValidFrom:          cert.ValidFrom,
		// ValidTo:            cert.ValidTo,
		// Issuer:             cert.Issuer,
		// Subject:            cert.Subject,
		// SignatureAlgorithm: cert.SignatureAlgorithm,
		// SerialNumber:       cert.SerialNumber,
	}, nil
}

func formatPermissions(permissions []permission) []string {
	var result []string
	for _, v := range permissions {
		result = append(result, v.Name)
	}

	return result
}

func getApkHash(apk *apk, option Option) {
	var hashOptions []string
	if option.WithMd5 {
		hashOptions = append(hashOptions, "md5")
	}
	if option.WithSha1 {
		hashOptions = append(hashOptions, "sha1")
	}
	if option.WithSha256 {
		hashOptions = append(hashOptions, "sha256")
	}
	apk.getApkHash(apk.f, hashOptions...)
}
