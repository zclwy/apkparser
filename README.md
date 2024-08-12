# apkparser

apk parser written in golang, aims to extract app information

[![GoDoc](https://pkg.go.dev/badge/github.com/zclwy/apkparser)](https://pkg.go.dev/github.com/zclwy/apkparser)
[![Go](https://github.com/zclwy/apkparser/actions/workflows/go.yml/badge.svg)](https://github.com/zclwy/apkparser/actions/workflows/go.yml)

## 简介

`apkparser` 项目是一个安卓`apk`文件解析器， 从`apk`文件中，获取 `AppInfo`

```
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
```

apk 证书相关信息的获取， 来自 [avast/apkparser](https://github.com/avast/apkparser)，本项目整合了 `avast/apkparser`的能力。

### 依赖

-   操作系统： 不限
-   编程语言： golang
-   库和框架： 不限

### 安装

以下是安装该项目的步骤：

    $ go get github.com/zclwy/apkparser

## 使用指南

使用方式见 [parser_test.go](parser_test.go)

## 贡献和许可

apk-parse 项目欢迎任何形式的贡献。如果您想要为该项目贡献代码或者报告某个问题，请提交一个 issue 或者一个 pull request。

### 许可

© zclwy, 2023 ~ now [license](LICENSE)

## 参考资料

-   [shogo82148/androidbinary](https://github.com/shogo82148/androidbinary)
-   [avast/apkparser](https://github.com/avast/apkparser)
