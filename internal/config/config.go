package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	SequenceDB struct { //mysql数据库配置,除mysql外,可能还有mongo等其他数据库
		DSN string // mysql链接地址,满足 $user:$password@tcp($ip:$port)/$db?$queries
	}
	ShortUrlDB struct { //mysql数据库配置,除mysql外,可能还有mongo等其他数据库
		DSN string // mysql链接地址,满足 $user:$password@tcp($ip:$port)/$db?$queries
	}
	CacheRedis cache.CacheConf //redis缓存
	Redis      struct {        //redis数据库配置
		Host string
	}
	BaseString        string   //base62指定基础字符串
	ShortUrlBlackList []string //短域名黑名单列表
	ShortDomain       string   //短域名
}
