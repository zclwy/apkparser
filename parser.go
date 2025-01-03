package apkparser

import (
	"errors"
	"image"

	"github.com/avast/apkverifier"
)

const (
	androidExt = ".apk"
)

type AppInfo struct {
	Name             string      `json:"name,omitempty"`        // 应用名称
	BundleId         string      `json:"bundleId,omitempty"`    // 包名
	Version          string      `json:"version,omitempty"`     // 版本名称
	Build            int64       `json:"build,omitempty"`       // 版本号
	Icon             image.Image `json:"icon,omitempty"`        // app icon
	Size             int64       `json:"size,omitempty"`        // app size in bytes
	CertInfo         CertInfo    `json:"certInfo,omitempty"`    // app 证书信息
	Md5              string      `json:"md5,omitempty"`         // app md5
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

func New(name string) (*AppInfo, error) {
	infoApk, err := openFile(name)
	if err != nil {
		return nil, err
	}

	certInfo, _ := getSignature(name)

	return &AppInfo{
		Name:             infoApk.parseApkLabel(),
		BundleId:         infoApk.apkManifest.Package,
		Version:          infoApk.apkManifest.VersionName,
		Build:            infoApk.apkManifest.VersionCode,
		Icon:             infoApk.parseApkIcon(),
		Size:             infoApk.size,
		CertInfo:         certInfo,
		Md5:              infoApk.md5,
		SupportOS64:      infoApk.supportOs64,
		SupportOS32:      infoApk.supportOs32,
		Permissions:      formatPermissions(infoApk.apkManifest.Permissions),
		MinSdkVersion:    infoApk.apkManifest.SDK.Min,
		MaxSdkVersion:    infoApk.apkManifest.SDK.Max,
		TargetSdkVersion: infoApk.apkManifest.SDK.Target,
	}, nil
}

// 获取apk签名
func getSignature(apkPath string) (CertInfo, error) {
	res, err := apkverifier.Verify(apkPath, nil)
	if err != nil {
		return CertInfo{}, err
	}

	cert, _ := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		return CertInfo{}, errors.New("no certificate found")
	}

	return CertInfo{
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
