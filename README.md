# 应用分发平台后端

#### 请先访问主仓库[https://github.com/AkariMarisa/app-management-platform](https://github.com/AkariMarisa/app-management-platform)

## 项目开发环境
go: 1.16
make: 4.3
code-oss: 1.55.2
golang.go (vscode插件): 0.24.1


## 项目依赖安装
```
go get
```

### 项目编译
```
make build
```

### 项目压缩
需要安装 upx
```
make compress
```

### 项目打包
```
make package
```

### 项目目录结构
```
root
|-bin 项目编译后的二进制文件
|-config 项目配置(监听端口, 文件保存位置)
|-controller 控制器, API 入口
|-db 数据库操作与项目 SQL
|-handler fiber 中间件
|-migrations golang-migrate 目录, 数据库 VCS文件
|-model 模型
|-out 打包输出目录
|-public 前端静态文件目录, 前端打包生产文件复制到这里
|-router fiber 路由
|-util 工具
```

### TODOs

- [x] 测试手机端应用信息检测与更新接口
- [x] 添加登陆接口
- [x] 添加系统参数修改接口
- [x] 需要对接口进行鉴权
- [x] 数据分页