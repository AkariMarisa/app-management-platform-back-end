package handler

import (
	"database/sql"
	"strings"

	"codelodon.com/akarimarisa/app-management-platform-back-end/db"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		// 验证用户身份 (Token)

		authorization := ctx.Get("Authorization")

		if authorization == "" {
			return ctx.Status(401).SendString("")
		}

		splits := strings.Split(authorization, "Bearer ")
		if len(splits) < 2 {
			return ctx.Status(401).SendString("")
		}

		token := splits[1]

		conn := db.GetConnection()

		stmt, _ := conn.Prepare(db.GetTokenRecord)

		defer stmt.Close()

		row := stmt.QueryRow(
			sql.Named("token", token),
		)

		var id uint
		row.Scan(&id)

		if id == 0 {
			return ctx.Status(401).SendString("")
		}

		return ctx.Next()
	}
}
