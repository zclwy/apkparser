# ipapk
ipa or apk parser written in golang, aims to extract app information

[![Build Status](https://travis-ci.org/phinexdaz/ipapk.svg?branch=master)](https://travis-ci.org/phinexdaz/ipapk)

## INSTALL
	$ go get github.com/zclwy/apk-parser
  
## USAGE
```go
package main

import (
	"fmt"
	"github.com/zclwy/apk-parser"
)

func main() {
	apk, _ := apk.NewAppParser("test.apk")
	fmt.Println(apk)
}
```
