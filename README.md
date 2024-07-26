# apkparser

apk parser written in golang, aims to extract app information

[![GoDoc](https://pkg.go.dev/github.com/zclwy/apkparser?status.svg)](https://pkg.go.dev/github.com/zclwy/apkparser)
[![Build Status](https://travis-ci.org/phinexdaz/ipapk.svg?branch=master)](https://travis-ci.org/phinexdaz/ipapk)

## 简介

`apkparser` 项目是一个安卓`apk`文件解析器， 从`apk`文件中，获取以下信息

1. name - app 名称
1. BundleId - app 包名
1. Version - app 版本名称
1. Build - app 版本号
1. Icon - app 图标
1. Size - app 大小
1. CertInfo - app 签名
    1. Md5
    2. Sha1
    3. Sha256
    4. ...
1. Md5 - app md5
1. SupportOS64 - 是否支持 64 位操作系统
1. SupportOS32 - 是否支持 32 位操作系统
1. Permissions

apk 证书相关信息的获取， 来自 [avast/apkparser](https://github.com/avast/apkparser)，本项目整合了 `avast/apkparser`的能力。

### 依赖

-   操作系统： 不限
-   编程语言： golang、shell
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
