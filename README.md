# 短链接项目
# 搭建项目的骨架
1. 建库建表
新建发号器表
````sql
   CREATE TABLE `sequence` (
   `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
   `stub` VARCHAR(1) NOT NULL,
   `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `idx_uniq_stub` (`stub`)
   ) ENGINE=MYISAM DEFAULT CHARSET=utf8 COMMENT = '序号表';
````
新建长链接短链接映射表
````sql
CREATE TABLE `short_url_map` (
`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
`create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
`is_del` TINYINT UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',
`lurl` VARCHAR(2048) DEFAULT NULL COMMENT '长链接',
`md5` CHAR(32) DEFAULT NULL COMMENT '长链接MD5',
`surl` VARCHAR(11) DEFAULT NULL COMMENT '短链接',
PRIMARY KEY (`id`),
INDEX(`is_del`),
UNIQUE(`md5`),
UNIQUE(`surl`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '长短链映射表';
 ````
2. 搭建go-zero框架骨架
2.1 编写api文件，使用goctl命令生成代码
````api
syntax = "v1"

info(
    title: "短链接项目"
    desc: "短链接重定向跳转长链接"
    author: "hz"
    email: "@929983177@qq.com"
    version:"1.0"
)
type ConvertRequest{
    LongURL string `json:"longUrl" validate:"required"`
}
type ConvertResponse{
    ShortURL string  `json:"shortUrl"`
}
type ShowRequest{
    ShortURL string  `json:"shortUrl" validate:"required"`
}
type ShowResponse{
    LongURL string `json:"longUrl"`
}
@server (
    prefix :api
)

service shortener-api{
    @handler ConvertHandler
    post /convert (ConvertRequest)returns(ConvertResponse)
    @handler ShowHandler
    get /:shortUrl (ShowRequest)returns(ShowResponse)
}
````
2.2 根据api文件生成go代码
````bash
goctl api go -api shortener.api  -dir .  -style=goZero
````
3. 根据数据表生成model层代码
````bash
goctl model mysql datasource -url="root:abc123@tcp(127.0.0.1:3306)/mall" -table="sequence" -dir="./model"  

goctl model mysql datasource -url="root:abc123@tcp(127.0.0.1:3306)/mall" -table="short_url_map" -dir="./model" 
````
4. 下载项目依赖
````bash
go mod tidy
````
5. 运行项目
````bash
go run shortener.go
````
看到如下输出代表项目成功启动
````bash
Starting server at 0.0.0.0:8888...
````
6. 修改配置结构体和配置文件
注意：两边一定一定要对齐！

# 转链参数检验
1.go-zero使用validator
https://pkg.go.dev/github.com/go-playground/validator/v10
下载依赖：
````bash
go get github.com/go-playground/validator/v10
````
导入依赖：
import "github.com/go-playground/validator/v10"
在api中为结构体添加validat额tag 并添加校验规则

# 查看短链接
# 缓存版
有两种方式
1. 使用自己实现的缓存  surl->lurl 能够节省缓存空间，缓存的数据量小
2. 使用go-zero自带的缓存 surl ->数据行，不需要自己实现，开发量小
这使用第二种方案：
1. 添加缓存配置
 - 配置文件
 - 配置config结构体
2. 删除旧的model层代码
 -删除 shorturlmapmodel.go
3. 重新生成model层代码
````bash
goctl model mysql datasource -url="root:abc123@tcp(127.0.0.1:3306)/mall" -table="short_url_map" -dir="./model" -c
````
4.修改svc层 ServiceContext.go文件 代码文件