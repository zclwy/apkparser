# apk-parse
apk parser written in golang, aims to extract app information

[![Build Status](https://travis-ci.org/phinexdaz/ipapk.svg?branch=master)](https://travis-ci.org/phinexdaz/ipapk)

## 项目名称

apk-parse

## 简介

apk-parse项目是一个安卓apk文件解析器， 从apk文件中，获取以下信息

1. name - app 名称
1. BundleId - app 包名
1. Version - app 版本名称
1. Build - app 版本号
1. Icon - app 图标
1. Size - app 大小
1. Signature - app 签名
2. Md5 - app md5
3. SupportOS64 - 是否支持64位操作系统
4. SupportOS32 - 是否支持32位操作系统


### 依赖

- 操作系统： 不限
- 编程语言： golang、shell
- 库和框架： 不限

### 安装

以下是安装该项目的步骤：

    $ go get github.com/zclwy/apk-parser

## 使用指南

使用方式见 [parser_test.go](parser_test.go)

## 贡献和许可

apk-parse 项目欢迎任何形式的贡献。如果您想要为该项目贡献代码或者报告某个问题，请提交一个 issue 或者一个 pull request。

### 许可

该项目使用 （指定许可证），因此，请确保您已经了解该许可证的要求和限制。
© zclwy, 2023 ~ now

## 参考资料

以下是参考本项目的资料：

- （列出参考资料）

