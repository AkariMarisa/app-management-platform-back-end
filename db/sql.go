package db

const (
	UpdateSystemParamByKey string = `
	UPDATE SystemParam SET
		Key = :key,
		Value = :value
	WHERE Key = :key
	`
	GetAppInfoById string = `
	SELECT
		Id,
		Name,
		PackageName,
		Type,
		Icon,
		ShortUrl,
		VersionName,
		VersionCode,
		Env,
		FileSize,
		CreatedAt
	FROM AppInfo
	WHERE Id = :id
	`
	GetAppInfoByPackageInfo string = `
	SELECT
		Id,
		Name,
		PackageName,
		Type,
		Icon,
		ShortUrl,
		VersionName,
		VersionCode,
		Env,
		FileSize,
		CreatedAt
	FROM AppInfo
	WHERE Name = :name
		AND PackageName = :packageName
		AND Type = :appType
	ORDER BY Id DESC
	`
	GetAppInfoByUniApp string = `
	SELECT
		Id,
		Name,
		PackageName,
		Type,
		Icon,
		ShortUrl,
		VersionName,
		VersionCode,
		Env,
		FileSize,
		CreatedAt
	FROM AppInfo
	WHERE Name = :name
		AND AppId = :appId
		AND Type = :appType
	ORDER BY Id DESC
	`
	GetAppInfoByUrl string = `
	SELECT
		au.Id,
		au.VersionName,
		au.VersionCode,
		au.Env,
		au.MinimumOSVersion,
		au.UpdateLog,
		au.FileSize,
		au.FileName,
		au.CreatedAt,
		ai.Id,
		ai.Name,
		ai.PackageName,
		ai.Type,
		ai.Icon,
		ai.ShortUrl
	FROM AppUpdate au
		INNER JOIN AppInfo ai
		on au.AppInfoId = ai.Id
	WHERE au.IsOnlineVersion = 1 AND ai.ShortUrl = :shortUrl
	ORDER BY au.Id DESC
	`
	GetAppUpdates string = `
	SELECT
		Id,
		VersionName,
		VersionCode,
		Env,
		ProvisionedDevices,
		MinimumOSVersion,
		UpdateLog,
		IsOnlineVersion,
		FileName,
		FileSize,
		CreatedAt
	FROM AppUpdate au
	WHERE AppInfoId = :appInfoId #{conditions}
	ORDER BY Id DESC
	LIMIT 50
	`
	GetDownloadCounts    string = "SELECT count(*) FROM DownloadRecord "
	InsertDownloadRecord string = `
	INSERT INTO DownloadRecord (
		AppInfoId,
		CreatedAt
	) VALUES (
		:appInfoId,
		:createdAt
	)
	`
	UpdateLog           string = "UPDATE AppUpdate SET UpdateLog = :log WHERE Id = :appUpdateId "
	CheckShortUrlExists string = `
	SELECT count(*) FROM AppInfo WHERE ShortUrl = :shortUrl
	`
	InsertAppInfoSQL string = `
	INSERT INTO AppInfo (
		AppId,
		Name,
		PackageName,
		Type,
		Icon,
		ShortUrl,
		VersionName,
		VersionCode,
		Env,
		FileSize,
		CreatedAt
	) VALUES (
		:appId,
		:name,
		:packageName,
		:type,
		:icon,
		:shortUrl,
		:versionName,
		:versionCode,
		:env,
		:fileSize,
		:createdAt
	)
	`
	InsertAppUpdateSQL string = `
	INSERT INTO AppUpdate (
		VersionName,
		VersionCode,
		Env,
		ProvisionedDevices,
		MinimumOSVersion,
		UpdateLog,
		IsOnlineVersion,
		FileSize,
		CreatedAt,
		AppInfoId,
		FileName
	) VALUES (
		:versionName,
		:versionCode,
		:env,
		:provisionedDevices,
		:minimumOSVersion,
		:updateLog,
		1,
		:fileSize,
		:createdAt,
		:appInfoId,
		:fileName
	)
	`
	SyncAppInfoVersion = `
	UPDATE AppInfo SET
		VersionName = :versionName,
		VersionCode = :versionCode,
		Env = :env,
		FileSize = :fileSize
	WHERE Id = :id
	`
	OfflineSQL      string = "UPDATE AppUpdate SET IsOnlineVersion = 0 WHERE AppInfoId = :appInfoId "
	OnlineSQL       string = "UPDATE AppUpdate SET IsOnlineVersion = 1 WHERE Id = :appUpdateId "
	GetAppFileName  string = "SELECT FileName, AppInfoId FROM AppUpdate WHERE Id = :id "
	GetNewerVersion string = `
	SELECT
		Id,
		VersionName,
		VersionCode,
		Env,
		ProvisionedDevices,
		MinimumOSVersion,
		UpdateLog,
		IsOnlineVersion,
		FileName,
		FileSize,
		CreatedAt
	FROM AppUpdate
	WHERE AppInfoId = :appInfoId
		AND IsOnlineVersion = 1
		AND VersionName != :versionName
		AND VersionCode != :versionCode
	ORDER BY Id DESC
	`
	GetUserByName string = `
	SELECT
		Id,
		Username,
		Password
	FROM User
	WHERE Username = :username
	ORDER BY Id DESC
	`
	GetUserById string = `
	SELECT
		Id,
		Username,
		Password
	FROM User
	WHERE Id = :id
	ORDER BY Id DESC
	`
	UpdateUserPassword string = `
	UPDATE User SET Password = :password WHERE Id = :id
	`
	DeleteUserOldTokens string = `
	DELETE FROM Token WHERE UserId = :userId
	`
	InsertToken string = `
	INSERT INTO Token (
		Token,
		UserId
	) VALUES (
		:token,
		:userId
	)
	`
	GetTokenRecord string = `
	SELECT Id FROM Token WHERE Token = :token
	`
	GetAppFiles string = `
	SELECT
		FileName
	FROM AppUpdate
	WHERE AppInfoId = :appInfoId
	`
	DeleteAppInfoById           string = "DELETE FROM AppInfo WHERE Id = :id"
	DeleteAppUpdatesByAppInfoId string = "DELETE FROM AppUpdate WHERE AppInfoId = :appInfoId"
)
