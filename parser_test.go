package apkparser

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestNewAppParser(t *testing.T) {
	apkFile := "testdata/helloworld.apk"
	app, err := NewAppParser(apkFile)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Printf("BundleId: %v\n", app.BundleId)
	fmt.Printf("Md5: %v\n", app.Md5)
	fmt.Printf("Signature md5: %v\n", app.CertInfo.Md5)

	if app.Icon != nil {
		// 生成png的icon
		pngFile, err := os.Create("testdata/helloworld.png")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = pngFile.Close()
		}()
		// 将 img 保存为 PNG 格式的图片文件
		err = png.Encode(pngFile, app.Icon)
		if err != nil {
			log.Fatal(err)
		}
	}
}
