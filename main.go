package main

import (
	"log"
	"os"

	"codelodon.com/akarimarisa/app-management-platform-back-end/config"
	"codelodon.com/akarimarisa/app-management-platform-back-end/controller"
	"codelodon.com/akarimarisa/app-management-platform-back-end/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	initDatabase()
	defer db.CloseConnection()

	// Create AppFileStore directory
	if err := os.MkdirAll(config.AppFileStorePath, 0777); err != nil {
		log.Fatal(err.Error())
	}

	// Startup http server
	app := fiber.New(fiber.Config{
		BodyLimit:     100 * 1024 * 1024, // 100MB
		CaseSensitive: true,
	})

	app.Use(cors.New())

	registerRoutes(app)

	log.Fatal(app.Listen(config.ServerAddress))
}

func initDatabase() {
	conn := db.GetConnection()

	config := &sqlite3.Config{}
	driver, _ := sqlite3.WithInstance(conn, config)

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"ql", driver)

	if err != nil {
		log.Fatal("database migration failed", err.Error())
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		m.Down()
		log.Fatal(err.Error())
	}
}

func registerRoutes(app *fiber.App) {

	// TODO 添加登陆接口
	// TODO 添加系统参数修改接口
	// TODO 需要对接口进行鉴权

	app.Static("/", "./public")

	api := app.Group("/api/v1")

	webApi := api.Group("/web")
	// 获取系统参数
	webApi.Get("/systemParam", controller.GetSystemParam)
	// 更新系统参数
	webApi.Put("/systemParam", controller.UpdateSystemParam)
	// 获取应用信息列表
	webApi.Get("/appInfo/list", controller.GetAppInfoList)
	// 根据ID获取应用信息
	webApi.Get("/appInfo", controller.GetAppInfoById)
	// 根据短URL获取应用信息
	webApi.Get("/appInfoByUrl", controller.GetAppInfoByUrl)
	// 获取应用更新版本信息列表
	webApi.Get("/appUpdate/list", controller.GetAppUpdateList)
	// 获取应用下载次数
	webApi.Get("/appDownload/count", controller.GetAppDownloadCount)
	// 更新应用更新版本信息的日志
	webApi.Put("/appUpdate/log", controller.UpdateAppUpdateLog)
	// 标记应用更新版本上线
	webApi.Put("/appUpdate/online", controller.MarkAppUpdateOnline)
	// 检查应用信息是否存在
	webApi.Get("/appInfo/exist", controller.CheckAppInfoExsits)
	// 生成短链接
	webApi.Get("/shortUrl/generate", controller.GenerateShortUrl)
	// 新增应用(第一次上传应用)
	webApi.Post("/appInfo", controller.CreateAppInfo)
	// 上传应用更新
	webApi.Post("/appUpdate", controller.UpdateApp)
	// 下载应用
	webApi.Get("/appUpdate/file", controller.DownloadApp)

	clientApi := api.Group("/client")
	// 客户端获取应用信息
	clientApi.Get("/appInfo", controller.ClientRetrieveAppInfo)
	// uniapp 客户端获取应用信息
	clientApi.Get("/appInfo/uniapp", controller.ClientRetrieveAppInfoUniApp)
	// 客户端检查应用更新
	clientApi.Get("/updates", controller.ClientCheckUpdate)
}
