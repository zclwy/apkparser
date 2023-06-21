package apkparser

import (
	"image/png"
	"log"
	"os"
	"testing"
)

func TestNewAppParser(t *testing.T) {
	apkFile := "testdata/helloworld.apk"
	app, err := NewAppParser(apkFile, "keytool")
	t.Log(app)
	t.Log(err)

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
