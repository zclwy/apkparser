package apkparser

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // handle jpeg format
	_ "image/png"  // handle png format
	"io"
	"os"
	"strconv"
	"strings"
)

// apk is an application package file for android.
type apk struct {
	f           *os.File
	zipReader   *zip.Reader
	apkManifest apkManifest
	table       *TableFile
	supportOs32 bool
	supportOs64 bool
	md5         string
	sha1        string
	sha256      string
	size        int64
}

// openFile will open the file specified by filename and return apk
func openFile(filename string) (apk *apk, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	apk, err = openZipReader(f, fi.Size())
	if err != nil {
		f.Close()
		return nil, err
	}
	apk.f = f
	apk.size = fi.Size()
	return
}

func openReader(f *os.File) (apk *apk, err error) {
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	apk, err = openZipReader(f, fi.Size())
	if err != nil {
		f.Close()
		return nil, err
	}
	apk.f = f
	apk.size = fi.Size()
	return
}

// openZipReader has same arguments like zip.NewReader
func openZipReader(r *os.File, size int64) (*apk, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	apk := &apk{
		zipReader: zipReader,
	}

	if err = apk.parseResources(); err != nil {
		return nil, err
	}

	if err = apk.parseManifest(); err != nil {
		return nil, errors.New("parse-apkManifest:" + err.Error())
	}

	apk.parseOsSupport(zipReader)
	apk.getApkHash(r)

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

	// 新增临时解析逻辑
	var tmp tempManifest
	if err := xml.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 转换特定属性的类型
	hardwareAccelerated := k.parseBool(tmp.App.HardwareAccelerated)
	vMSafeMode := k.parseBool(tmp.App.VMSafeMode)
	largeHeap := k.parseBool(tmp.App.LargeHeap)

	k.apkManifest = apkManifest{
		Package:     tmp.Package,
		VersionCode: tmp.VersionCode,
		VersionName: tmp.VersionName,
		App: apkApplication{
			AllowTaskReParenting:  tmp.App.AllowTaskReParenting,
			AllowBackup:           tmp.App.AllowBackup,
			BackupAgent:           tmp.App.BackupAgent,
			Debuggable:            tmp.App.Debuggable,
			Description:           tmp.App.Description,
			Enabled:               tmp.App.Enabled,
			HasCode:               tmp.App.HasCode,
			HardwareAccelerated:   hardwareAccelerated,
			Icon:                  tmp.App.Icon,
			KillAfterRestore:      tmp.App.KillAfterRestore,
			Label:                 tmp.App.Label,
			Logo:                  tmp.App.Logo,
			ManageSpaceActivity:   tmp.App.ManageSpaceActivity,
			Name:                  tmp.App.Name,
			Permission:            tmp.App.Permission,
			Persistent:            tmp.App.Persistent,
			Process:               tmp.App.Process,
			RestoreAnyVersion:     tmp.App.RestoreAnyVersion,
			RequiredAccountType:   tmp.App.RequiredAccountType,
			RestrictedAccountType: tmp.App.RestrictedAccountType,
			SupportsRtl:           tmp.App.SupportsRtl,
			TaskAffinity:          tmp.App.TaskAffinity,
			TestOnly:              tmp.App.TestOnly,
			Theme:                 tmp.App.Theme,
			UIOptions:             tmp.App.UIOptions,
			Activities:            tmp.App.Activities,
			ActivityAliases:       tmp.App.ActivityAliases,
			VMSafeMode:            vMSafeMode,
			LargeHeap:             largeHeap,
		},
		Instrument:  tmp.Instrument,
		Permissions: tmp.Permissions,
		SDK:         tmp.SDK,
	}

	return nil
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

func (k *apk) getResourceToBool(id string, resConfig *ResTableConfig) bool {
	resID, err := ParseResID(id)
	if err != nil {
		return false
	}
	val, err := k.table.GetResource(resID, resConfig)
	if err != nil {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return false
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

func (k *apk) parseOsSupport(r *zip.Reader) {
	var (
		hasSoFile                bool
		isSupport32, isSupport64 bool
	)
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".so") {
			hasSoFile = true
		}
		// if strings.HasPrefix(f.Name, "lib/arm64-v8a") {
		// 	fmt.Println(f.Name)
		// 	isSupport64 = true
		// }
		// if strings.HasPrefix(f.Name, "lib/armeabi") {
		// 	fmt.Println(f.Name)
		// 	isSupport32 = true
		// }

		if strings.Contains(f.Name, "armeabi-v7a") || strings.Contains(f.Name, "x86") {
			isSupport32 = true
		} else if strings.Contains(f.Name, "arm64-v8a") || strings.Contains(f.Name, "x86_64") {
			isSupport64 = true
		}
	}
	// 当前apk支持的系统位数
	if !hasSoFile && !isSupport64 && !isSupport32 {
		k.supportOs32 = true
		k.supportOs64 = true
	} else {
		k.supportOs32 = isSupport32
		k.supportOs64 = isSupport64
	}
}

// 获取apk md5
func (k *apk) getApkMd5(file *os.File) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return
	}

	k.md5 = fmt.Sprintf("%032x", hash.Sum(nil))
}

func (k *apk) getApkHash(file *os.File) {
	md5Hasher := md5.New()
	sha1Hasher := sha1.New()
	sha256Hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(md5Hasher, sha1Hasher, sha256Hasher), file); err != nil {
		return
	}
	k.md5 = hex.EncodeToString(md5Hasher.Sum(nil))
	k.sha1 = hex.EncodeToString(sha1Hasher.Sum(nil))
	k.sha256 = hex.EncodeToString(sha256Hasher.Sum(nil))
}

// 解析apk名称
func (k *apk) parseApkLabel() string {
	label, _ := k.label(&ResTableConfig{
		Language: [2]uint8{'z', 'h'},
		Country:  [2]uint8{'C', 'N'},
	})
	// label, _ := pkg.label(&ResTableConfig{})

	return label
}

// 解析apk图标
func (k *apk) parseApkIcon() image.Image {
	icon, _ := k.icon(&ResTableConfig{
		Density: 720,
	})

	return icon
}

func (k *apk) parseBool(val string) bool {
	if IsResID(val) {
		// 资源ID解析
		return k.getResourceToBool(val, nil)
	}

	b, _ := strconv.ParseBool(val)
	return b
}
