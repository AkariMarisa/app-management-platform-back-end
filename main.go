package main

import (
	"log"
	"os"

	"codelodon.com/akarimarisa/app-management-platform-back-end/config"
	"codelodon.com/akarimarisa/app-management-platform-back-end/db"
	"codelodon.com/akarimarisa/app-management-platform-back-end/handler"
	"codelodon.com/akarimarisa/app-management-platform-back-end/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// 读取配置
	configuration := config.GetConfiguration()

	initDatabase()
	defer db.CloseConnection()

	// Create AppFileStore directory
	if err := os.MkdirAll(configuration.AppFileStorePath, 0777); err != nil {
		log.Fatal(err.Error())
	}

	// Startup http server
	app := fiber.New(fiber.Config{
		BodyLimit:     100 * 1024 * 1024, // 100MB
		CaseSensitive: true,
		ErrorHandler:  handler.HandleError,
	})

	app.Use(recover.New())

	app.Use(cors.New())

	router.SetupRoutes(app)

	log.Fatal(app.Listen(configuration.ServerAddress))
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
