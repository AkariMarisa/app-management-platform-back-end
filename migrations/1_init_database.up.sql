CREATE TABLE AppInfo (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    PackageName TEXT NOT NULL,
    Type TEXT NOT NULL,
    Icon TEXT,
    ShortUrl TEXT NOT NULL,
    VersionName TEXT NOT NULL,
    VersionCode INTEGER NOT NULL,
    Env TEXT,
    FileSize REAL NOT NULL,
    CreatedAt TEXT
);

CREATE TABLE AppUpdate (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    VersionName TEXT NOT NULL,
    VersionCode INTEGER NOT NULL,
    Env TEXT,
    ProvisionedDevices TEXT,
    MinimumOSVersion TEXT,
    UpdateLog TEXT,
    IsOnlineVersion INTEGER NOT NULL,
    FileSize REAL NOT NULL,
    CreatedAt TEXT,
    AppInfoId INTEGER NOT NULL,
    FOREIGN KEY(AppInfoId) REFERENCES AppInfo(Id)
);

CREATE TABLE DownloadRecord (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    AppInfoId INTEGER NOT NULL,
    CreatedAt TEXT,
    FOREIGN KEY(AppInfoId) REFERENCES AppInfo(Id)
);

CREATE TABLE SystemParam (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT NOT NULL,
    Value TEXT
);