package apkparser

import (
	"image/png"
	"log"
	"os"
	"testing"
)

func TestAppParser(t *testing.T) {
	apkFile := "testdata/helloworld.apk"
	app, err := New(apkFile, Option{
		WithIcon:             true,
		WithSignature:        true,
		IgnoreSignatureError: true,
		WithMd5:              true,
		WithSha1:             true,
		WithSha256:           true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Name: %v\n", app.Name)
	t.Logf("BundleId: %v\n", app.BundleId)
	t.Logf("Version: %v\n", app.Version)
	t.Logf("Build: %v\n", app.Build)
	t.Logf("MinSdkVersion: %v\n", app.MinSdkVersion)
	t.Logf("TargetSdkVersion: %v\n", app.TargetSdkVersion)
	t.Logf("MaxSdkVersion: %v\n", app.MaxSdkVersion)
	t.Logf("Md5: %v\n", app.Md5)
	t.Logf("Sha1: %v\n", app.Sha1)
	t.Logf("Sha256: %v\n", app.Sha256)
	t.Logf("Signature md5: %v\n", app.CertInfo.Md5)

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

func TestAppParser_NewFromReader(t *testing.T) {
	apkFile := "testdata/helloworld.apk"
	f, err := os.Open(apkFile)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	app, err := NewFromReader(f, Option{
		WithIcon:             true,
		WithSignature:        true,
		IgnoreSignatureError: true,
		WithMd5:              true,
		WithSha1:             true,
		WithSha256:           true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Name: %v\n", app.Name)
	t.Logf("BundleId: %v\n", app.BundleId)
	t.Logf("Version: %v\n", app.Version)
	t.Logf("Build: %v\n", app.Build)
	t.Logf("MinSdkVersion: %v\n", app.MinSdkVersion)
	t.Logf("TargetSdkVersion: %v\n", app.TargetSdkVersion)
	t.Logf("MaxSdkVersion: %v\n", app.MaxSdkVersion)
	t.Logf("Md5: %v\n", app.Md5)
	t.Logf("Sha1: %v\n", app.Sha1)
	t.Logf("Sha256: %v\n", app.Sha256)
	t.Logf("Signature md5: %v\n", app.CertInfo.Md5)

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
