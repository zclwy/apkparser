package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zclwy/apkparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: apkparser <apk-file>")
		os.Exit(1)
	}

	apkFile := os.Args[1]
	app, err := apkparser.New(apkFile, Option{
		WithIcon:             true,
		WithSignature:        true,
		IgnoreSignatureError: true,
	})
	if err != nil {
		log.Fatalf("Failed to parse APK: %v", err)
	}

	fmt.Printf("Name: %v\n", app.Name)
	fmt.Printf("BundleId: %v\n", app.BundleId)
	fmt.Printf("Version: %v\n", app.Version)
	fmt.Printf("Build: %v\n", app.Build)
	fmt.Printf("Md5: %v\n", app.Md5)
	fmt.Printf("Signature md5: %v\n", app.CertInfo.Md5)
}
