package model

// 应用信息, (名称 + 包名 + 类型)
type AppInfo struct {
	Id            uint32          `json:"id"`
	AppId         string          `json:"appId"`                           // UniApp 下的 appId
	Name          string          `json:"name" validate:"required"`        // 应用名称
	PackageName   string          `json:"packageName" validate:"required"` // 应用包名
	Type          AppType         `json:"type" validate:"required"`        // 应用类型, (Android: 安卓, Apple: 苹果)
	Icon          string          `json:"icon"`                            // 应用图标url
	ShortUrl      string          `json:"shortUrl"`                        // 应用短链接
	VersionName   string          `json:"versionName" validate:"required"` // 版本名称, e.g. 12.0.1
	VersionCode   uint64          `json:"versionCode" validate:"required"` // 版本编号, e.g. 1201
	Env           AppleEnviroment `json:"env"`                             // 应用环境, 仅苹果应用存在此属性, (Development: 开发版, Production: 正式版)
	FileSize      float32         `json:"fileSize" validate:"required"`    // 应用文件大小
	CreatedAt     string          `json:"createdAt"`                       // 创建时间
	CurrentUpdate *AppUpdate      `json:"currentUpdate" validate:"dive"`   // 当前应用的线上版本信息
}

// 应用类型
type AppType string

const (
	Android AppType = "android"
	Apple   AppType = "apple"
)

// 苹果应用环境
type AppleEnviroment string

const (
	Development AppleEnviroment = "development"
	Production  AppleEnviroment = "production"
)

// 应用更新版本信息
type AppUpdate struct {
	Id                 uint32          `json:"id"`
	VersionName        string          `json:"versionName" validate:"required"`               // 版本名称, e.g. 12.0.1
	VersionCode        uint64          `json:"versionCode" validate:"required"`               // 版本编号, e.g. 1201
	Env                AppleEnviroment `json:"env"`                                           // 应用环境, 仅苹果应用存在此属性, (Development: 开发版, Production: 正式版)
	ProvisionedDevices string          `json:"provisionedDevices"`                            // 应用开发设备UUID, 仅苹果应用且是开发版的时候存在此属性, 多个时用半角逗号隔开
	MinimumOSVersion   string          `json:"minimumOSVersion"`                              // 应用要求最低系统版本, 仅苹果应用存在此属性
	UpdateLog          string          `json:"updateLog"`                                     // 应用更新日志
	IsOnlineVersion    bool            `json:"isOnlineVersion" validate:"required"`           // 应用是否为线上版本
	FileSize           float32         `json:"fileSize" validate:"required,eq=True|eq=False"` // 应用文件大小
	CreatedAt          string          `json:"createdAt"`
	AppInfo            *AppInfo        `json:"appInfo" validate:"dive"`
	FileName           string          `json:"fileName" validate:"required"`
}

// 应用下载记录
type DownloadRecord struct {
	Id        uint32 `json:"id"`
	AppInfoId uint32 `json:"appInfoId" validate:"required"` // 应用ID
	CreatedAt string `json:"createdAt"`
}

// 系统参数表
type SystemParam struct {
	Id    uint32 `json:"id"`
	Key   string `json:"key" validate:"required"`
	Value string `json:"value"`
}
