package apk

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/andrianbdn/iospng"
	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
)

var (
	ErrNoIcon = errors.New("icon not found")
)

const (
	androidExt = ".apk"
)

type AppInfo struct {
	Name     string
	BundleId string
	Version  string
	Build    string
	Icon     image.Image
	Size     int64
}

type androidManifest struct {
	Package     string `xml:"package,attr"`
	VersionName string `xml:"versionName,attr"`
	VersionCode string `xml:"versionCode,attr"`
}

func NewAppParser(name string) (*AppInfo, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}
	
	var xmlFile *zip.File
	for _, f := range reader.File {
		switch {
		case f.Name == "AndroidManifest.xml":
			xmlFile = f
		}
	}
	
	ext := filepath.Ext(stat.Name())
	
	if ext == androidExt {
		info, err := parseApkFile(xmlFile)
		icon, label, err := parseApkIconAndLabel(name)
		info.Name = label
		info.Icon = icon
		info.Size = stat.Size()
		return info, err
	}
	
	return nil, errors.New("unknown platform")
}

func parseAndroidManifest(xmlFile *zip.File) (*androidManifest, error) {
	rc, err := xmlFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	
	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	
	xmlContent, err := androidbinary.NewXMLFile(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	
	manifest := new(androidManifest)
	decoder := xml.NewDecoder(xmlContent.Reader())
	if err := decoder.Decode(manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

func parseApkFile(xmlFile *zip.File) (*AppInfo, error) {
	if xmlFile == nil {
		return nil, errors.New("AndroidManifest.xml not found")
	}
	
	manifest, err := parseAndroidManifest(xmlFile)
	if err != nil {
		return nil, err
	}
	
	info := new(AppInfo)
	info.BundleId = manifest.Package
	info.Version = manifest.VersionName
	info.Build = manifest.VersionCode
	
	return info, nil
}

func parseApkIconAndLabel(name string) (image.Image, string, error) {
	pkg, err := apk.OpenFile(name)
	if err != nil {
		return nil, "", err
	}
	defer pkg.Close()
	
	icon, _ := pkg.Icon(&androidbinary.ResTableConfig{
		Density: 720,
	})
	if icon == nil {
		return nil, "", ErrNoIcon
	}
	
	label, _ := pkg.Label(nil)
	
	return icon, label, nil
}

func parseIpaIcon(iconFile *zip.File) (image.Image, error) {
	if iconFile == nil {
		return nil, ErrNoIcon
	}
	
	rc, err := iconFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	
	var w bytes.Buffer
	iospng.PngRevertOptimization(rc, &w)
	
	return png.Decode(bytes.NewReader(w.Bytes()))
}
