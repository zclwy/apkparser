package apkparser

// apkInstrumentation is an application instrumentation code.
type apkInstrumentation struct {
	Name            string `xml:"name,attr"`
	Target          string `xml:"targetPackage,attr"`
	HandleProfiling bool   `xml:"handleProfiling,attr"`
	FunctionalTest  bool   `xml:"functionalTest,attr"`
}

// apkActivityAction is an action of an activity.
type apkActivityAction struct {
	Name string `xml:"name,attr"`
}

// apkActivityCategory is a category of an activity.
type apkActivityCategory struct {
	Name string `xml:"name,attr"`
}

// apkActivityIntentFilter is an intent filter of an activity.
type apkActivityIntentFilter struct {
	Action   apkActivityAction   `xml:"action"`
	Category apkActivityCategory `xml:"category"`
}

// apkAppActivity is an activity in an application.
type apkAppActivity struct {
	Theme         string                    `xml:"theme,attr"`
	Name          string                    `xml:"name,attr"`
	Label         string                    `xml:"label,attr"`
	IntentFilters []apkActivityIntentFilter `xml:"intent-filter"`
}

// apkAppActivityAlias https://developer.android.com/guide/topics/manifest/activity-alias-element
type apkAppActivityAlias struct {
	Name           string                    `xml:"name,attr"`
	Label          string                    `xml:"label,attr"`
	TargetActivity string                    `xml:"targetActivity,attr"`
	IntentFilters  []apkActivityIntentFilter `xml:"intent-filter"`
}

// apkApplication is an application in an APK.
type apkApplication struct {
	AllowTaskReParenting  bool                  `xml:"allowTaskReparenting,attr"`
	AllowBackup           bool                  `xml:"allowBackup,attr"`
	BackupAgent           string                `xml:"backupAgent,attr"`
	Debuggable            bool                  `xml:"debuggable,attr"`
	Description           string                `xml:"description,attr"`
	Enabled               bool                  `xml:"enabled,attr"`
	HasCode               bool                  `xml:"hasCode,attr"`
	HardwareAccelerated   bool                  `xml:"hardwareAccelerated,attr"`
	Icon                  string                `xml:"icon,attr"`
	KillAfterRestore      bool                  `xml:"killAfterRestore,attr"`
	Label                 string                `xml:"label,attr"`
	Logo                  string                `xml:"logo,attr"`
	ManageSpaceActivity   string                `xml:"manageSpaceActivity,attr"`
	Name                  string                `xml:"name,attr"`
	Permission            string                `xml:"permission,attr"`
	Persistent            bool                  `xml:"persistent,attr"`
	Process               string                `xml:"process,attr"`
	RestoreAnyVersion     bool                  `xml:"restoreAnyVersion,attr"`
	RequiredAccountType   string                `xml:"requiredAccountType,attr"`
	RestrictedAccountType string                `xml:"restrictedAccountType,attr"`
	SupportsRtl           bool                  `xml:"supportsRtl,attr"`
	TaskAffinity          string                `xml:"taskAffinity,attr"`
	TestOnly              bool                  `xml:"testOnly,attr"`
	Theme                 string                `xml:"theme,attr"`
	UIOptions             string                `xml:"uiOptions,attr"`
	Activities            []apkAppActivity      `xml:"activity"`
	ActivityAliases       []apkAppActivityAlias `xml:"activity-alias"`
	VMSafeMode            bool                  `xml:"vmSafeMode,attr"`
	LargeHeap             bool                  `xml:"largeHeap,attr"`
}

// tmpApkApplication 临时解析Application，用于类型转换
type tmpApkApplication struct {
	AllowTaskReParenting  bool                  `xml:"allowTaskReparenting,attr"`
	AllowBackup           bool                  `xml:"allowBackup,attr"`
	BackupAgent           string                `xml:"backupAgent,attr"`
	Debuggable            bool                  `xml:"debuggable,attr"`
	Description           string                `xml:"description,attr"`
	Enabled               bool                  `xml:"enabled,attr"`
	HasCode               bool                  `xml:"hasCode,attr"`
	HardwareAccelerated   string                `xml:"hardwareAccelerated,attr"` // 需要转换
	Icon                  string                `xml:"icon,attr"`
	KillAfterRestore      bool                  `xml:"killAfterRestore,attr"`
	Label                 string                `xml:"label,attr"`
	Logo                  string                `xml:"logo,attr"`
	ManageSpaceActivity   string                `xml:"manageSpaceActivity,attr"`
	Name                  string                `xml:"name,attr"`
	Permission            string                `xml:"permission,attr"`
	Persistent            bool                  `xml:"persistent,attr"`
	Process               string                `xml:"process,attr"`
	RestoreAnyVersion     bool                  `xml:"restoreAnyVersion,attr"`
	RequiredAccountType   string                `xml:"requiredAccountType,attr"`
	RestrictedAccountType string                `xml:"restrictedAccountType,attr"`
	SupportsRtl           bool                  `xml:"supportsRtl,attr"`
	TaskAffinity          string                `xml:"taskAffinity,attr"`
	TestOnly              bool                  `xml:"testOnly,attr"`
	Theme                 string                `xml:"theme,attr"`
	UIOptions             string                `xml:"uiOptions,attr"`
	Activities            []apkAppActivity      `xml:"activity"`
	ActivityAliases       []apkAppActivityAlias `xml:"activity-alias"`
	VMSafeMode            string                `xml:"vmSafeMode,attr"` // 需要转换
	LargeHeap             string                `xml:"largeHeap,attr"`  // 需要转换
}

// apkUsesSDK is target SDK version.
type apkUsesSDK struct {
	Min    int `xml:"minSdkVersion,attr"`
	Target int `xml:"targetSdkVersion,attr"`
	Max    int `xml:"maxSdkVersion,attr"`
}

// apkManifest is a apkManifest of an APK.
type apkManifest struct {
	Package     string             `xml:"package,attr"`
	VersionCode int64              `xml:"versionCode,attr"`
	VersionName string             `xml:"versionName,attr"`
	App         apkApplication     `xml:"application"`
	Instrument  apkInstrumentation `xml:"instrumentation"`
	Permissions []permission       `xml:"uses-permission"`
	SDK         apkUsesSDK         `xml:"uses-sdk"`
}

// tempManifest 临时解析Manifest，用于类型转换
type tempManifest struct {
	Package     string             `xml:"package,attr"`
	VersionCode int64              `xml:"versionCode,attr"`
	VersionName string             `xml:"versionName,attr"`
	App         tmpApkApplication  `xml:"application"`
	Instrument  apkInstrumentation `xml:"instrumentation"`
	Permissions []permission       `xml:"uses-permission"`
	SDK         apkUsesSDK         `xml:"uses-sdk"`
}

type permission struct {
	Name string `xml:"name,attr"`
}
