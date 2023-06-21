package apkparser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // handle jpeg format
	_ "image/png"  // handle png format
	"io"
	"os"
	"strconv"
)

// apk is an application package file for android.
type apk struct {
	f           *os.File
	zipReader   *zip.Reader
	apkManifest apkManifest
	table       *TableFile
}

// openFile will open the file specified by filename and return apk
func openFile(filename string) (apk *apk, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	apk, err = openZipReader(f, fi.Size())
	if err != nil {
		return nil, err
	}
	apk.f = f
	return
}

// openZipReader has same arguments like zip.NewReader
func openZipReader(r io.ReaderAt, size int64) (*apk, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	apk := &apk{
		zipReader: zipReader,
	}
	if err = apk.parseManifest(); err != nil {
		return nil, errors.New("parse-apkManifest:" + err.Error())
	}
	if err = apk.parseResources(); err != nil {
		return nil, err
	}
	return apk, nil
}

// close is avaliable only if apk is created with openFile
func (k *apk) close() error {
	if k.f == nil {
		return nil
	}
	return k.f.Close()
}

// icon returns the icon image of the APK.
func (k *apk) icon(resConfig *ResTableConfig) (image.Image, error) {
	iconPath := k.getResource(k.apkManifest.App.Icon, resConfig)
	if IsResID(iconPath) {
		return nil, errors.New("unable to convert icon-id to icon path")
	}
	imgData, err := k.readZipFile(iconPath)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(bytes.NewReader(imgData))
	return m, err
}

// label returns the label of the APK.
func (k *apk) label(resConfig *ResTableConfig) (s string, err error) {
	s = k.getResource(k.apkManifest.App.Label, resConfig)
	if IsResID(s) {
		err = errors.New("unable to convert label-id to string")
	}
	return
}

// manifest returns the apkManifest of the APK.
func (k *apk) manifest() apkManifest {
	return k.apkManifest
}

// packageName returns the package name of the APK.
func (k *apk) packageName() string {
	return k.apkManifest.Package
}

// mainActivity returns the name of the main activity.
func (k *apk) mainActivity() (activity string, err error) {
	for _, act := range k.apkManifest.App.Activities {
		for _, intent := range act.IntentFilters {
			if intent.Action.Name == "android.intent.action.MAIN" &&
				intent.Category.Name == "android.intent.category.LAUNCHER" {
				return act.Name, nil
			}
		}
	}
	for _, act := range k.apkManifest.App.ActivityAliases {
		for _, intent := range act.IntentFilters {
			if intent.Action.Name == "android.intent.action.MAIN" &&
				intent.Category.Name == "android.intent.category.LAUNCHER" {
				return act.TargetActivity, nil
			}
		}
	}

	return "", errors.New("no main activity found")
}

func (k *apk) parseManifest() error {
	xmlData, err := k.readZipFile("AndroidManifest.xml")
	if err != nil {
		return errors.New("read-apkManifest.xml" + err.Error())
	}
	xmlFile, err := NewXMLFile(bytes.NewReader(xmlData))
	if err != nil {
		return errors.New("parse-xml:" + err.Error())
	}
	reader := xmlFile.Reader()
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, &k.apkManifest)
}

func (k *apk) parseResources() (err error) {
	resData, err := k.readZipFile("resources.arsc")
	if err != nil {
		return
	}
	k.table, err = NewTableFile(bytes.NewReader(resData))
	return
}

func (k *apk) getResource(id string, resConfig *ResTableConfig) string {
	resID, err := ParseResID(id)
	if err != nil {
		return id
	}
	val, err := k.table.GetResource(resID, resConfig)
	if err != nil {
		return id
	}
	return fmt.Sprintf("%s", val)
}

func (k *apk) readZipFile(name string) (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	for _, file := range k.zipReader.File {
		if file.Name != name {
			continue
		}
		rc, er := file.Open()
		if er != nil {
			fmt.Println("file.Open " + er.Error())
			err = er
			return
		}
		_, err = io.Copy(buf, rc)
		if err != nil {
			_ = rc.Close()
			return data, err
		}
		data = buf.Bytes()
		_ = rc.Close()
		break
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file %s not found", strconv.Quote(name))
	} else {
		return data, nil
	}
}
