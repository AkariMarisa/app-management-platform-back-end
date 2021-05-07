package controller

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"codelodon.com/akarimarisa/app-management-platform-back-end/config"
	"codelodon.com/akarimarisa/app-management-platform-back-end/db"
	"codelodon.com/akarimarisa/app-management-platform-back-end/model"
	"codelodon.com/akarimarisa/app-management-platform-back-end/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

// 登陆
// 请求参数:
// username: 用户名
// password: 密码
func Login(ctx *fiber.Ctx) error {

	username := ctx.Query("username")

	if strings.TrimSpace(username) == "" {
		return ctx.Status(400).SendString("username不能为空")
	}

	password := ctx.Query("password")

	if strings.TrimSpace(password) == "" {
		return ctx.Status(400).SendString("password不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	stmt, _ := tx.Prepare(db.GetUserByName)

	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("username", username),
	)

	user := model.User{}
	row.Scan(&user.Id, &user.Username, &user.Password)

	if user.Id == 0 {
		return ctx.Status(400).SendString("用户名不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ctx.Status(400).SendString("用户名或密码错误")
	}

	// 登陆成功, 生成 Token
	token := util.GenerateSecureToken(128)

	deleteStmt, _ := tx.Prepare(db.DeleteUserOldTokens)

	defer deleteStmt.Close()

	deleteStmt.Exec(
		sql.Named("userId", user.Id),
	)

	insertStmt, _ := tx.Prepare(db.InsertToken)

	defer insertStmt.Close()

	insertStmt.Exec(
		sql.Named("token", token),
		sql.Named("userId", user.Id),
	)

	tx.Commit()

	result := make(map[string]string)
	result["token"] = token
	result["userId"] = strconv.Itoa(int(user.Id))

	return ctx.JSON(result)
}

// 修改用户密码
// 请求参数:
// userId: 用户ID
// oldPassword: 原密码
// newPassword: 新密码
// reNewPassword: 重复输入的新密码
func ChangeUserPassword(ctx *fiber.Ctx) error {

	userId := ctx.Query("userId")

	if strings.TrimSpace(userId) == "" {
		return ctx.Status(400).SendString("userId不能为空")
	}

	oldPassword := ctx.Query("oldPassword")

	if strings.TrimSpace(oldPassword) == "" {
		return ctx.Status(400).SendString("oldPassword不能为空")
	}

	newPassword := ctx.Query("newPassword")

	if strings.TrimSpace(newPassword) == "" {
		return ctx.Status(400).SendString("newPassword不能为空")
	}

	reNewPassword := ctx.Query("reNewPassword")

	if strings.TrimSpace(reNewPassword) == "" {
		return ctx.Status(400).SendString("reNewPassword不能为空")
	}

	if newPassword != reNewPassword {
		return ctx.Status(400).SendString("两次密码不一致")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	getUserStmt, _ := tx.Prepare(db.GetUserById)

	defer getUserStmt.Close()

	row := getUserStmt.QueryRow(
		sql.Named("id", userId),
	)

	user := model.User{}
	row.Scan(&user.Id, &user.Username, &user.Password)

	if user.Id == 0 {
		return ctx.Status(400).SendString("用户名不存在")
	}

	// TODO 这里可以加个限制, 用户原始密码输入错误多少次之后就锁定账户

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ctx.Status(400).SendString("原始密码错误")
	}

	stmt, _ := tx.Prepare(db.UpdateUserPassword)

	defer stmt.Close()

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	stmt.Exec(
		sql.Named("id", userId),
		sql.Named("password", string(hashedPassword)),
	)

	deleteStmt, _ := tx.Prepare(db.DeleteUserOldTokens)

	defer deleteStmt.Close()

	deleteStmt.Exec(
		sql.Named("userId", userId),
	)

	tx.Commit()

	return ctx.SendString("修改成功")
}

// 获取系统参数
// 请求参数:
// key: Key
func GetSystemParam(ctx *fiber.Ctx) error {

	querySQL := "SELECT Id, Key, Value FROM SystemParam "
	conditionSQL := "WHERE 1=1 "
	orderSQL := "ORDER BY Id DESC "

	var params []interface{}

	if key := ctx.Query("key"); strings.TrimSpace(key) != "" {
		conditionSQL += "AND Key = :key "
		params = append(params, sql.Named("key", key))
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(querySQL + conditionSQL + orderSQL)

	defer stmt.Close()

	rows, _ := stmt.Query(params...)

	defer rows.Close()

	var result []model.SystemParam
	for rows.Next() {
		systemParam := model.SystemParam{}

		rows.Scan(&systemParam.Id, &systemParam.Key, &systemParam.Value)
		result = append(result, systemParam)
	}

	return ctx.JSON(result)
}

// 更新系统参数
// 请求参数:
// key: Key
// value: Value
func UpdateSystemParam(ctx *fiber.Ctx) error {

	key := ctx.Query("key")

	if strings.TrimSpace(key) == "" {
		return ctx.Status(400).SendString("key不能为空")
	}

	value := ctx.Query("value")

	if strings.TrimSpace(value) == "" {
		return ctx.Status(400).SendString("value不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	stmt, _ := tx.Prepare(db.UpdateSystemParamByKey)

	defer stmt.Close()

	result, _ := stmt.Exec(
		sql.Named("key", key),
		sql.Named("value", value),
	)

	tx.Commit()

	rowsAffected, _ := result.RowsAffected()

	return ctx.SendString(fmt.Sprintf("已更新%d条记录", rowsAffected))
}

// 获取应用信息列表
// 请求参数:
// name: 应用名称, 支持模糊查询
// appType: 应用类型
func GetAppInfoList(ctx *fiber.Ctx) error {

	querySQL := `
	SELECT
		ai.Id,
		ai.Name,
		ai.PackageName,
		ai.Type,
		ai.Icon,
		ai.ShortUrl,
		ai.VersionName,
		ai.VersionCode,
		ai.Env,
		ai.FileSize,
		ai.CreatedAt,
		au.Id,
		au.VersionName,
		au.VersionCode,
		au.Env,
		au.ProvisionedDevices,
		au.MinimumOSVersion,
		au.UpdateLog,
		au.FileSize,
		au.FileName
	FROM AppInfo ai LEFT JOIN AppUpdate au ON au.AppInfoId = ai.Id AND au.IsOnlineVersion = 1
	`
	conditionSQL := " WHERE 1=1 "
	orderSQL := "ORDER BY ai.Id DESC "

	var params []interface{}

	if name := ctx.Query("name"); strings.TrimSpace(name) != "" {
		conditionSQL += "AND ai.Name like :name "
		params = append(params, sql.Named("name", "%"+name+"%"))
	}

	if appType := ctx.Query("appType"); strings.TrimSpace(appType) != "" {
		types := strings.Split(appType, ",")

		conditionSQL += "AND ai.Type IN (?" + strings.Repeat(",?", len(types)-1) + ")"
		for _, t := range types {
			params = append(params, t)
		}
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(querySQL + conditionSQL + orderSQL)

	defer stmt.Close()

	rows, _ := stmt.Query(params...)

	defer rows.Close()

	var result []model.AppInfo
	for rows.Next() {
		appInfo := model.AppInfo{}
		appUpdate := model.AppUpdate{}

		rows.Scan(
			&appInfo.Id,
			&appInfo.Name,
			&appInfo.PackageName,
			&appInfo.Type,
			&appInfo.Icon,
			&appInfo.ShortUrl,
			&appInfo.VersionName,
			&appInfo.VersionCode,
			&appInfo.Env,
			&appInfo.FileSize,
			&appInfo.CreatedAt,
			&appUpdate.Id,
			&appUpdate.VersionName,
			&appUpdate.VersionCode,
			&appUpdate.Env,
			&appUpdate.ProvisionedDevices,
			&appUpdate.MinimumOSVersion,
			&appUpdate.UpdateLog,
			&appUpdate.FileSize,
			&appUpdate.FileName,
		)

		appInfo.CurrentUpdate = &appUpdate

		result = append(result, appInfo)
	}

	return ctx.JSON(result)
}

// 根据ID获取应用信息
// 请求参数:
// appInfoId: 应用信息ID
func GetAppInfoById(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")
	if strings.TrimSpace(appInfoId) == "" {
		return ctx.Status(400).SendString("appInfoId不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppInfoById)

	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("id", appInfoId))

	appInfo := model.AppInfo{}
	row.Scan(
		&appInfo.Id,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.Type,
		&appInfo.Icon,
		&appInfo.ShortUrl,
		&appInfo.VersionName,
		&appInfo.VersionCode,
		&appInfo.Env,
		&appInfo.FileSize,
		&appInfo.CreatedAt,
	)

	return ctx.JSON(appInfo)
}

// 删除应用
// 请求参数:
// appInfoId: 应用信息ID
func AbandonApp(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")

	if strings.TrimSpace(appInfoId) == "" {
		return ctx.Status(400).SendString("appInfoId不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	// 获取所有应用文件名称

	stmt, _ := tx.Prepare(db.GetAppFiles)

	defer stmt.Close()

	rows, _ := stmt.Query(
		sql.Named("appInfoId", appInfoId),
	)

	var fileNames []string
	for rows.Next() {
		var fileName string

		rows.Scan(&fileName)

		fileNames = append(fileNames, fileName)
	}

	// 删除应用更新版本记录

	deleteAppUpdateStmt, _ := tx.Prepare(db.DeleteAppUpdatesByAppInfoId)

	defer deleteAppUpdateStmt.Close()

	deleteAppUpdateStmt.Exec(sql.Named("appInfoId", appInfoId))

	// 删除应用记录

	deleteAppInfoStmt, _ := tx.Prepare(db.DeleteAppInfoById)

	defer deleteAppInfoStmt.Close()

	deleteAppInfoStmt.Exec(sql.Named("id", appInfoId))

	// 删除应用文件
	for _, fileName := range fileNames {
		os.Remove(fmt.Sprintf("%s/%s", config.GetConfiguration().AppFileStorePath, fileName))
	}

	tx.Commit()

	return ctx.SendString("删除成功")
}

// 根据短URL获取应用信息
// 请求参数:
// shortUrl: 应用短URL
func GetAppInfoByUrl(ctx *fiber.Ctx) error {

	shortUrl := ctx.Query("shortUrl")

	if strings.TrimSpace(shortUrl) == "" {
		return ctx.Status(400).SendString("shortUrl不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppInfoByUrl)

	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("shortUrl", shortUrl))

	appUpdate := model.AppUpdate{}
	appInfo := model.AppInfo{}
	row.Scan(
		&appUpdate.Id,
		&appUpdate.VersionName,
		&appUpdate.VersionCode,
		&appUpdate.Env,
		&appUpdate.MinimumOSVersion,
		&appUpdate.UpdateLog,
		&appUpdate.FileSize,
		&appUpdate.FileName,
		&appUpdate.CreatedAt,
		&appInfo.Id,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.Type,
		&appInfo.Icon,
		&appInfo.ShortUrl,
	)

	appUpdate.AppInfo = &appInfo

	// 如果当前应用没有上线的更新信息, 则报错
	if appUpdate.Id == 0 {
		return ctx.Status(500).SendString("当前应用不存在或没有上线版本")
	}

	return ctx.JSON(appUpdate)
}

// 获取应用更新版本信息列表
// 请求参数:
// appInfoId: 应用信息ID
func GetAppUpdateList(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")

	if strings.TrimSpace(appInfoId) == "" {
		return ctx.Status(400).SendString("appInfoId不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppUpdates)

	defer stmt.Close()

	rows, _ := stmt.Query(sql.Named("appInfoId", appInfoId))

	defer rows.Close()

	var result []model.AppUpdate
	for rows.Next() {
		appUpdate := model.AppUpdate{}

		rows.Scan(
			&appUpdate.Id,
			&appUpdate.VersionName,
			&appUpdate.VersionCode,
			&appUpdate.Env,
			&appUpdate.ProvisionedDevices,
			&appUpdate.MinimumOSVersion,
			&appUpdate.UpdateLog,
			&appUpdate.IsOnlineVersion,
			&appUpdate.FileName,
			&appUpdate.FileSize,
			&appUpdate.CreatedAt,
		)
		result = append(result, appUpdate)
	}

	return ctx.JSON(result)
}

// 获取应用下载次数
// 请求参数:
// appInfoId: 应用信息ID
func GetAppDownloadCount(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")

	querySQL := db.GetDownloadCounts

	var params []interface{}

	if strings.TrimSpace(appInfoId) != "" {
		querySQL = db.GetDownloadCountsByApp
		params = append(params, sql.Named("appInfoId", appInfoId))
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(querySQL)

	defer stmt.Close()

	row := stmt.QueryRow(params...)

	var count int
	row.Scan(&count)

	return ctx.SendString(strconv.Itoa(count))
}

// 更新应用更新版本信息的日志
// 请求参数:
// appUpdateId: 应用更新版本信息ID
// log: 应用更新日志
func UpdateAppUpdateLog(ctx *fiber.Ctx) error {

	appUpdateId := ctx.Query("appUpdateId")

	if strings.TrimSpace(appUpdateId) == "" {
		return ctx.Status(400).SendString("appUpdateId不能为空")
	}

	log := ctx.FormValue("log")

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	stmt, _ := tx.Prepare(db.UpdateLog)

	defer stmt.Close()

	result, _ := stmt.Exec(
		sql.Named("appUpdateId", appUpdateId),
		sql.Named("log", log),
	)

	tx.Commit()

	rowsAffected, _ := result.RowsAffected()

	return ctx.SendString(fmt.Sprintf("已更新%d条记录", rowsAffected))
}

// 标记应用更新版本上线
// 请求参数:
// appInfoId: 应用信息ID
// appUpdateId: 应用更新版本信息ID
func MarkAppUpdateOnline(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")

	if strings.TrimSpace(appInfoId) == "" {
		return ctx.Status(400).SendString("appInfoId不能为空")
	}

	appUpdateId := ctx.Query("appUpdateId")

	if strings.TrimSpace(appUpdateId) == "" {
		return ctx.Status(400).SendString("appUpdateId不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	// 首先将所有应用对应更新版本切换到下线
	offlineStmt, _ := tx.Prepare(db.OfflineSQL)

	defer offlineStmt.Close()

	offlineStmt.Exec(sql.Named("appInfoId", appInfoId))

	// 然后再更新对应版本的上线状态
	onlineStmt, _ := tx.Prepare(db.OnlineSQL)

	defer onlineStmt.Close()

	onlineResult, _ := onlineStmt.Exec(sql.Named("appUpdateId", appUpdateId))

	tx.Commit()

	rowsAffected, _ := onlineResult.RowsAffected()

	return ctx.SendString(fmt.Sprintf("已更新%d条记录", rowsAffected))
}

// 检查应用信息是否存在
// 请求参数:
// name: Name
// packageName: PackageName
// appType: AppType
func CheckAppInfoExsits(ctx *fiber.Ctx) error {

	name := ctx.Query("name")

	if strings.TrimSpace(name) == "" {
		return ctx.Status(400).SendString("name不能为空")
	}

	packageName := ctx.Query("packageName")

	if strings.TrimSpace(packageName) == "" {
		return ctx.Status(400).SendString("packageName不能为空")
	}

	appType := ctx.Query("appType")

	if strings.TrimSpace(appType) == "" {
		return ctx.Status(400).SendString("appType不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppInfoByPackageInfo)

	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("name", name),
		sql.Named("packageName", packageName),
		sql.Named("appType", appType),
	)

	var appInfo model.AppInfo
	row.Scan(
		&appInfo.Id,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.Type,
		&appInfo.Icon,
		&appInfo.ShortUrl,
		&appInfo.VersionName,
		&appInfo.VersionCode,
		&appInfo.Env,
		&appInfo.FileSize,
		&appInfo.CreatedAt,
	)

	return ctx.JSON(appInfo)
}

func generateShortUrl() (string, error) {

	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 4)

	// 首先生成一个短链接给前端展示用
	// 每次生成一个短链接时, 也要到数据库里查一下, 如果存在的话, 需要重新生成

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.CheckShortUrlExists)

	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("shortUrl", id))

	var count int
	row.Scan(&count)

	if count > 0 {
		return generateShortUrl()
	}

	return id, nil
}

// 生成短链接
func GenerateShortUrl(ctx *fiber.Ctx) error {

	shortUrl, _ := generateShortUrl()

	return ctx.SendString(shortUrl)

}

// 新增应用(第一次上传应用)
func CreateAppInfo(ctx *fiber.Ctx) error {

	form, _ := ctx.MultipartForm()

	appId := form.Value["appInfo.appId"][0]
	name := form.Value["appInfo.name"][0]
	packageName := form.Value["appInfo.packageName"][0]
	appType := form.Value["appInfo.type"][0]
	icon := form.Value["appInfo.icon"][0]
	shortUrl := form.Value["appInfo.shortUrl"][0]

	versionName := form.Value["appUpdate.versionName"][0]
	versionCode := form.Value["appUpdate.versionCode"][0]
	env := form.Value["appUpdate.env"][0]
	provisionedDevices := form.Value["appUpdate.provisionedDevices"][0]
	minimunOSVersion := form.Value["appUpdate.minimunOSVersion"][0]
	updateLog := form.Value["appUpdate.updateLog"][0]

	if name == "" && packageName == "" && appType == "" {
		return ctx.Status(400).SendString("应用信息不能为空")
	}

	appFile := form.File["file"][0]

	fileName := appFile.Filename
	fileNameSplits := strings.Split(fileName, ".")
	fileNameSplitsLength := len(fileNameSplits)
	if fileNameSplitsLength <= 1 {
		return ctx.Status(400).SendString("文件名称非法")
	}

	fileExtensionName := strings.ToLower(fileNameSplits[fileNameSplitsLength-1])
	if fileExtensionName != "apk" && fileExtensionName != "ipa" {
		return ctx.Status(400).SendString("仅支持apk或ipa文件")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	// 生成短链接的思路
	// 使用nano-id库生成短链接(github.com/matoous/go-nanoid/v2)
	// 首先生成一个短链接给前端展示用
	// 每次生成一个短链接时, 也要到数据库里查一下, 如果存在的话, 需要重新生成
	// 如果前端没有修改的话, 就会用这个短链接
	// 此时后端需要在重新查询一下数据库里有没有这条短链接
	// 如果有, 则提醒已重复
	// 如果前端修改了短链接, 则一样发给后端
	// 后端一样先检测有没有, 如果有则一样提醒

	// 检查前端传过来的短链接是否已存在
	checkShortUrlStmt, _ := tx.Prepare(db.CheckShortUrlExists)

	defer checkShortUrlStmt.Close()

	row := checkShortUrlStmt.QueryRow(sql.Named("shortUrl", shortUrl))

	var count int
	row.Scan(&count)

	if count > 0 {
		return ctx.Status(400).SendString("短链接已被使用")
	}

	// 保存应用信息
	insertAppInfoStmt, _ := tx.Prepare(db.InsertAppInfoSQL)

	defer insertAppInfoStmt.Close()

	nowStr := time.Now().Format("2006-01-02 15:04:05")
	insertAppInfoResult, _ := insertAppInfoStmt.Exec(
		sql.Named("appId", appId),
		sql.Named("name", name),
		sql.Named("packageName", packageName),
		sql.Named("type", appType),
		sql.Named("icon", icon),
		sql.Named("shortUrl", shortUrl),
		sql.Named("versionName", versionName),
		sql.Named("versionCode", versionCode),
		sql.Named("env", env),
		sql.Named("fileSize", appFile.Size),
		sql.Named("createdAt", nowStr),
	)

	appInfoId, _ := insertAppInfoResult.LastInsertId()

	// 保存应用更新版本信息
	insertAppUpdateStmt, _ := tx.Prepare(db.InsertAppUpdateSQL)

	defer insertAppUpdateStmt.Close()

	localFileName := fmt.Sprintf("%s.%s", uuid.NewString(), fileExtensionName)
	insertAppUpdateStmt.Exec(
		sql.Named("versionName", versionName),
		sql.Named("versionCode", versionCode),
		sql.Named("env", env),
		sql.Named("provisionedDevices", provisionedDevices),
		sql.Named("minimumOSVersion", minimunOSVersion),
		sql.Named("updateLog", updateLog),
		sql.Named("fileSize", appFile.Size),
		sql.Named("createdAt", nowStr),
		sql.Named("appInfoId", appInfoId),
		sql.Named("fileName", localFileName),
	)

	// 最后保存应用文件到本地
	ctx.SaveFile(appFile, fmt.Sprintf("%s/%s", config.GetConfiguration().AppFileStorePath, localFileName))

	tx.Commit()

	return ctx.SendString(strconv.Itoa(int(appInfoId)))
}

// 上传应用更新
// 请求Body:
// appUpdate: 应用更新版本信息
func UpdateApp(ctx *fiber.Ctx) error {

	form, _ := ctx.MultipartForm()

	appInfoId := form.Value["appInfo.id"][0]

	versionName := form.Value["appUpdate.versionName"][0]
	versionCode := form.Value["appUpdate.versionCode"][0]
	env := form.Value["appUpdate.env"][0]
	provisionedDevices := form.Value["appUpdate.provisionedDevices"][0]
	minimunOSVersion := form.Value["appUpdate.minimunOSVersion"][0]
	updateLog := form.Value["appUpdate.updateLog"][0]

	if appInfoId == "" {
		return ctx.Status(400).SendString("应用信息ID不能为空")
	}

	appFile := form.File["file"][0]

	fileName := appFile.Filename
	fileNameSplits := strings.Split(fileName, ".")
	fileNameSplitsLength := len(fileNameSplits)
	if fileNameSplitsLength <= 1 {
		return ctx.Status(400).SendString("文件名称非法")
	}

	fileExtensionName := strings.ToLower(fileNameSplits[fileNameSplitsLength-1])
	if fileExtensionName != "apk" && fileExtensionName != "ipa" {
		return ctx.Status(400).SendString("仅支持apk或ipa文件")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	// 首先将所有一你敢用对应更新版本切换到下线
	offlineStmt, _ := tx.Prepare(db.OfflineSQL)

	defer offlineStmt.Close()

	offlineStmt.Exec(sql.Named("appInfoId", appInfoId))

	// 然后插入新的应用更新版本信息
	insertAppUpdateStmt, _ := tx.Prepare(db.InsertAppUpdateSQL)

	defer insertAppUpdateStmt.Close()

	nowStr := time.Now().Format("2006-01-02 15:04:05")
	localFileName := fmt.Sprintf("%s.%s", uuid.NewString(), fileExtensionName)
	insertAppUpdateStmt.Exec(
		sql.Named("versionName", versionName),
		sql.Named("versionCode", versionCode),
		sql.Named("env", env),
		sql.Named("provisionedDevices", provisionedDevices),
		sql.Named("minimumOSVersion", minimunOSVersion),
		sql.Named("updateLog", updateLog),
		sql.Named("fileSize", appFile.Size),
		sql.Named("createdAt", nowStr),
		sql.Named("appInfoId", appInfoId),
		sql.Named("fileName", localFileName),
	)

	// 然后更新应用信息
	updateAppInfoStmt, _ := tx.Prepare(db.SyncAppInfoVersion)

	defer updateAppInfoStmt.Close()

	updateAppInfoStmt.Exec(
		sql.Named("versionName", versionName),
		sql.Named("versionCode", versionCode),
		sql.Named("env", env),
		sql.Named("fileSize", appFile.Size),
		sql.Named("id", appInfoId),
	)

	// 最后保存应用文件到本地
	ctx.SaveFile(appFile, fmt.Sprintf("%s/%s", config.GetConfiguration().AppFileStorePath, localFileName))

	tx.Commit()

	return ctx.SendString("上传成功")
}

// 下载应用
// 请求参数:
// appUpdateId: 应用更新版本信息ID
func DownloadApp(ctx *fiber.Ctx) error {

	appUpdateId := ctx.Query("appUpdateId")

	if strings.TrimSpace(appUpdateId) == "" {
		return ctx.Status(400).SendString("appUpdateId不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	stmt, _ := tx.Prepare(db.GetAppFileName)

	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("id", appUpdateId))

	var fileName string
	var appInfoId uint32
	row.Scan(&fileName, &appInfoId)

	ctx.Download(fmt.Sprintf("%s/%s", config.GetConfiguration().AppFileStorePath, fileName))

	// 更新下载次数
	updateStmt, _ := tx.Prepare(db.InsertDownloadRecord)

	defer updateStmt.Close()

	updateStmt.Exec(
		sql.Named("appInfoId", appInfoId),
		sql.Named("createdAt", time.Now().Format("2006-01-02 15:04:05")),
	)

	tx.Commit()

	return nil
}

// 客户端获取应用信息
// 请求参数:
// name: 应用名称
// packageName: 应用包名
// appType: 应用类型
func ClientRetrieveAppInfo(ctx *fiber.Ctx) error {

	name := ctx.Query("name")
	if strings.TrimSpace(name) == "" {
		return ctx.Status(400).SendString("name不能为空")
	}

	packageName := ctx.Query("packageName")
	if strings.TrimSpace(packageName) == "" {
		return ctx.Status(400).SendString("packageName不能为空")
	}

	appType := ctx.Query("appType")
	if strings.TrimSpace(appType) == "" {
		return ctx.Status(400).SendString("appType不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppInfoByPackageInfo)

	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("name", name),
		sql.Named("packageName", packageName),
		sql.Named("appType", appType),
	)

	var appInfo model.AppInfo
	row.Scan(
		&appInfo.Id,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.Type,
		&appInfo.Icon,
		&appInfo.ShortUrl,
		&appInfo.VersionName,
		&appInfo.VersionCode,
		&appInfo.Env,
		&appInfo.FileSize,
		&appInfo.CreatedAt,
	)

	return ctx.JSON(appInfo)
}

// 客户端获取应用信息
// 请求参数:
// name: 应用名称
// appId: AppId
// appType: 应用类型
func ClientRetrieveAppInfoUniApp(ctx *fiber.Ctx) error {

	name := ctx.Query("name")
	if strings.TrimSpace(name) == "" {
		return ctx.Status(400).SendString("name不能为空")
	}

	appId := ctx.Query("appId")
	if strings.TrimSpace(appId) == "" {
		return ctx.Status(400).SendString("appId不能为空")
	}

	appType := ctx.Query("appType")
	if strings.TrimSpace(appType) == "" {
		return ctx.Status(400).SendString("appType不能为空")
	}

	conn := db.GetConnection()

	stmt, _ := conn.Prepare(db.GetAppInfoByUniApp)

	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("name", name),
		sql.Named("appId", appId),
		sql.Named("appType", appType),
	)

	var appInfo model.AppInfo
	row.Scan(
		&appInfo.Id,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.Type,
		&appInfo.Icon,
		&appInfo.ShortUrl,
		&appInfo.VersionName,
		&appInfo.VersionCode,
		&appInfo.Env,
		&appInfo.FileSize,
		&appInfo.CreatedAt,
	)

	return ctx.JSON(appInfo)
}

// 客户端检查应用更新
// 请求参数:
// appInfoId: 应用信息ID
// versionName: 版本名称
// versionCode: 版本代号
func ClientCheckUpdate(ctx *fiber.Ctx) error {

	appInfoId := ctx.Query("appInfoId")

	if strings.TrimSpace(appInfoId) == "" {
		return ctx.Status(400).SendString("appInfoId不能为空")
	}

	versionName := ctx.Query("versionName")

	if strings.TrimSpace(versionName) == "" {
		return ctx.Status(400).SendString("versionName不能为空")
	}

	versionCode := ctx.Query("versionCode")

	if strings.TrimSpace(versionCode) == "" {
		return ctx.Status(400).SendString("versionCode不能为空")
	}

	conn := db.GetConnection()

	tx, _ := conn.Begin()

	defer tx.Rollback()

	stmt, _ := tx.Prepare(db.GetNewerVersion)

	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("appInfoId", appInfoId),
		sql.Named("versionName", versionName),
		sql.Named("versionCode", versionCode),
	)

	var appUpdate model.AppUpdate
	row.Scan(
		&appUpdate.Id,
		&appUpdate.VersionName,
		&appUpdate.VersionCode,
		&appUpdate.Env,
		&appUpdate.ProvisionedDevices,
		&appUpdate.MinimumOSVersion,
		&appUpdate.UpdateLog,
		&appUpdate.IsOnlineVersion,
		&appUpdate.FileName,
		&appUpdate.FileSize,
		&appUpdate.CreatedAt,
	)

	// 如果应用没有上线的更新版本, 则返回空
	if appUpdate.Id == 0 {
		return ctx.SendString("当前应用暂无新版本")
	}

	ctx.Download(fmt.Sprintf("%s/%s", config.GetConfiguration().AppFileStorePath, appUpdate.FileName))

	// 更新下载次数
	updateStmt, _ := tx.Prepare(db.InsertDownloadRecord)

	defer updateStmt.Close()

	updateStmt.Exec(
		sql.Named("appInfoId", appInfoId),
		sql.Named("createdAt", time.Now().Format("2006-01-02 15:04:05")),
	)

	tx.Commit()

	return nil
}
