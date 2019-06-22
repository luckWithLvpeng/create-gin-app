# create gin app

gin 作为 server 初始化golang项目

## Features 特点
* `conf/app.con`配置app相关参数，可以设置日志保存目录，和当发生重大错误(panic)时，日志保存文件 `error.log`。
* 引入[fvbock/endless](https://github.com/fvbock/endless),支持热更新，优雅重启，无缝更新代码
* 日志 [go-ini/ini](https://github.com/go-ini/ini)，支持配置文件的读写，[文档地址](https://ini.unknwon.io/docs/intro/getting_started)
* [gin](https://github.com/gin-gonic/gin/blob/master/README.md)是golang目前最快的路由框架，[文档](https://gin-gonic.com/docs/)，[示例代码](https://github.com/gin-gonic/examples)
* 数据库操作 [jinzhu/gorm](https://github.com/jinzhu/gorm) , [english doc](https://gorm.io/docs/), [中文文档](https://jasperxu.github.io/gorm-zh/database.html#m)

## 编码需要注意的地方
* 数据库操作包，我们用了gorm, 由于gorm的链式操作，每次返回都是一个db结构体， 当在某一个函数中使用如 db.Where()时，会污染全局的db变量， 建议在函数内部，重新初始化一个db变量去操作数据库, 例如
```go
func GetUser() {
    var users []User
    DB := db
    DB = DB.Find(&users)
    if DB.Error != nil {
        ...
    }
    ...
}
```

### swagger 自动化api文档
1. 添加 [swag](https://github.com/swaggo/gin-swagger) 接口, [swag文档](https://swaggo.github.io/swaggo.io/declarative_comments_format/)
2. 安装 [Swag](https://github.com/swaggo/swag) 命令:
```sh
$ go install github.com/swaggo/swag/cmd/swag
```
3. 执行 [Swag](https://github.com/swaggo/swag) 命令，在有 `main.go` 文件的根目录, [Swag](https://github.com/swaggo/swag) 将会格式化注释并生成文档在`docs` 目录 和 `docs/doc.go`文件
```sh
$ swag init
```
4. 项目`go run main.go`启动后，可以在[http://localhost:8080/swagger/index.html]( http://localhost:8080/swagger/index.html) 看到文档， 路由配置参考[swag项目主页](https://github.com/swaggo/gin-swagger)
