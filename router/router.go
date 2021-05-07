package router

import (
	"codelodon.com/akarimarisa/app-management-platform-back-end/controller"
	"codelodon.com/akarimarisa/app-management-platform-back-end/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Static("/", "./public")

	api := app.Group("/api/v1")

	webApi := api.Group("/web")

	// 用户登陆
	webApi.Post("/login", controller.Login)
	// 获取系统参数
	webApi.Get("/systemParam", controller.GetSystemParam)
	// 根据短URL获取应用信息
	webApi.Get("/appInfoByUrl", controller.GetAppInfoByUrl)
	// 下载应用
	webApi.Get("/appUpdate/file", controller.DownloadApp)

	webProtectApi := webApi.Group("", handler.Protected())
	// 修改用户密码
	webProtectApi.Put("/user/password", controller.ChangeUserPassword)
	// 更新系统参数
	webProtectApi.Put("/systemParam", controller.UpdateSystemParam)
	// 获取应用信息列表
	webProtectApi.Get("/appInfo/list", controller.GetAppInfoList)
	// 根据ID获取应用信息
	webProtectApi.Get("/appInfo", controller.GetAppInfoById)
	// 删除应用信息及其对应的更新版本
	webProtectApi.Delete("/appInfo", controller.AbandonApp)
	// 获取应用更新版本信息列表
	webProtectApi.Get("/appUpdate/list", controller.GetAppUpdateList)
	// 获取应用下载次数
	webProtectApi.Get("/appDownload/count", controller.GetAppDownloadCount)
	// 更新应用更新版本信息的日志
	webProtectApi.Put("/appUpdate/log", controller.UpdateAppUpdateLog)
	// 标记应用更新版本上线
	webProtectApi.Put("/appUpdate/online", controller.MarkAppUpdateOnline)
	// 检查应用信息是否存在
	webProtectApi.Get("/appInfo/exist", controller.CheckAppInfoExsits)
	// 生成短链接
	webProtectApi.Get("/shortUrl/generate", controller.GenerateShortUrl)
	// 新增应用(第一次上传应用)
	webProtectApi.Post("/appInfo", controller.CreateAppInfo)
	// 上传应用更新
	webProtectApi.Post("/appUpdate", controller.UpdateApp)

	clientApi := api.Group("/client")
	// 客户端获取应用信息
	clientApi.Get("/appInfo", controller.ClientRetrieveAppInfo)
	// uniapp 客户端获取应用信息
	clientApi.Get("/appInfo/uniapp", controller.ClientRetrieveAppInfoUniApp)
	// 客户端检查应用更新
	clientApi.Get("/updates", controller.ClientCheckUpdate)
}
