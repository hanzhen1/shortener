package svc

import (
	"errors"
	BloomFilter "github.com/bits-and-blooms/bloom/v3"
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"
)

type ServiceContext struct {
	Config            config.Config
	Sequence          sequence.Sequence      //sequence
	ShortUrlModel     model.ShortUrlMapModel //short_url_map
	ShortUrlBlackList map[string]struct{}
	//Filter            *BloomFilter.BloomFilter //bloom filter布隆过滤器 "github.com/bits-and-blooms/bloom/v3"
	Filter *bloom.Filter //布隆过滤器"github.com/zeromicro/go-zero/core/bloom"
}

func NewServiceContext(c config.Config) *ServiceContext {
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	//把配置文件中配置的黑名单加载到map，方便后续判断
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	sqlConn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	//初始化布隆过滤器
	//初始化redisBitSet
	store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	//声明一个bitSet,key="bloom_filter" bits是20*100万
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))
	//重新加载已有的短链接数据  b.基于redis版本,go-zero自带的
	if err := loadDataToFilter(sqlConn, filter); err != nil {
		logx.Errorw("loadDataToFilter failed", logx.LogField{Key: "err", Value: err.Error()})
	}
	//// 初始化bloomfilter
	//filter := BloomFilter.NewWithEstimates(1<<20, 0.01)
	////重新加载已有的短链接数据
	////a.基于内存版本 服务重启之后就没了，所以每次重启就要重新加载一下已有的短链接(从数据库查询)
	//if err := loadDataToBloomFilter(sqlConn, filter); err != nil {
	//	logx.Errorw("loadDataToBloomFilter failed", logx.LogField{Key: "err", Value: err.Error()})
	//}
	return &ServiceContext{
		Config:   c,
		Sequence: sequence.NewMySQL(c.SequenceDB.DSN),
		//Sequence:      sequence.NewRedis(c.Redis.Host),
		ShortUrlModel:     model.NewShortUrlMapModel(sqlConn, c.CacheRedis),
		ShortUrlBlackList: m, //短链接黑名单map
		Filter:            filter,
	}
}

// 注意导入的是这个bloom
//import "github.com/bits-and-blooms/bloom/v3"

// loadDataToBloomFilter 加载已有的短链接数据至布隆过滤器
func loadDataToBloomFilter(conn sqlx.SqlConn, filter *BloomFilter.BloomFilter) error {
	// 循环从数据库查询数据加载至filter
	if conn == nil || filter == nil {
		return errors.New("loadDataToBloomFilter invalid param")
	}

	// 查总数
	total := 0
	if err := conn.QueryRow(&total, "select count(*) from short_url_map where is_del=0"); err != nil {
		logx.Errorw("conn.QueryRowCount failed", logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	logx.Infow("total data", logx.LogField{Key: "total", Value: total})
	if total == 0 {
		logx.Info("no data need to load")
		return nil
	}
	pageTotal := 0
	pageSize := 20
	if total%pageSize == 0 {
		pageTotal = total / pageSize
	} else {
		pageTotal = total/pageSize + 1
	}
	logx.Infow("pageTotal", logx.LogField{Key: "pageTotal", Value: pageTotal})
	// 循环查询所有数据
	for page := 1; page <= pageTotal; page++ {
		offset := (page - 1) * pageSize
		surls := []string{}
		if err := conn.QueryRows(&surls, `select surl from short_url_map where is_del=0 limit ?,?`, offset, pageSize); err != nil {
			return err
		}

		for _, surl := range surls {
			filter.AddString(surl)
		}
	}
	logx.Info("load data to bloom success")
	return nil
}

// 注意导入的是这个bloom
//import "github.com/zeromicro/go-zero/core/bloom"

// loadDataToFilter 加载已有的短链接数据至布隆过滤器
func loadDataToFilter(conn sqlx.SqlConn, filter *bloom.Filter) error {
	// 循环从数据库查询数据加载至filter
	if conn == nil || filter == nil {
		return errors.New("loadDataToBloomFilter invalid param")
	}

	// 查总数
	total := 0
	if err := conn.QueryRow(&total, "select count(*) from short_url_map where is_del=0"); err != nil {
		logx.Errorw("conn.QueryRowCount failed", logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	logx.Infow("total data", logx.LogField{Key: "total", Value: total})
	if total == 0 {
		logx.Info("no data need to load")
		return nil
	}
	pageTotal := 0
	pageSize := 20
	if total%pageSize == 0 {
		pageTotal = total / pageSize
	} else {
		pageTotal = total/pageSize + 1
	}
	logx.Infow("pageTotal", logx.LogField{Key: "pageTotal", Value: pageTotal})
	// 循环查询所有数据
	for page := 1; page <= pageTotal; page++ {
		offset := (page - 1) * pageSize
		surls := []string{}
		if err := conn.QueryRows(&surls, `select surl from short_url_map where is_del=0 limit ?,?`, offset, pageSize); err != nil {
			return err
		}

		for _, surl := range surls {
			filter.Add([]byte(surl))
		}
	}
	logx.Info("load data to bloom success")
	return nil
}
