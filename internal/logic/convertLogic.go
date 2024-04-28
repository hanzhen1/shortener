package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链业务逻辑：输入一个长链接-->转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1.校验输入的数据
	//1.1数据不能为空
	//if len(req.LongURL)==0{}
	//使用validator包来做参数校验
	//1.2 输入的长链接必须是能请求通的网址
	//用http.Get()方法请求看返回的状态码是不是200，不是200都不可达
	if ok := connect.Get(req.LongURL); !ok {
		return nil, errors.New("无效的链接")
	}
	//1.3判断之前是否已经转链过(数据库中是否已经存在该长链接)
	//1.3.1 给长链接生成md5
	md5Value := md5.Sum([]byte(req.LongURL)) //注意！这里使用的是项目中自己封装的pkg/md5 包
	//1.3.2 拿md5值去数据库中查是否存在
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil { //数据库查到了表示该长链接已被转过
			return nil, fmt.Errorf("该链接已被转为%s", u.Surl.String)
		}
		//查数据库失败了
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.Field("err", err))
		return nil, err
	}

	//1.4 输入的不能是一个短链接(避免循环转链)
	//输入的是一个完整的url q1mi.cn/1d12a?name=q1mi
	basePath, err := urltool.GetBasePath(req.LongURL)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "lurl", Value: req.LongURL}, logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil { //数据库查到了表示该长链接已被转过
			return nil, errors.New("该链接已经是短链了")
		}
		logx.Errorw("convert_ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, errors.New("内部错误")
	}
	var short string
	for {
		//2.取号 基于MySQL实现的发号器
		//每来一个转链请求，我们就使用replace into语句往sequence表插入一条数据并且取出主键id作为号码
		res, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		fmt.Println(res)
		//3. 号码转短链
		//3.1安全性
		short = base62.IntToString(res)
		fmt.Printf("short:%v\n", short)
		//3.2 短域名黑名单避免某些特殊词比如api,health,fuck等等
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break //生成不在黑名单里的短链接就跳出for循环
		}
	}

	//4.存储长短链映射关系
	if _, err = l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongURL, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		},
	); err != nil {
		logx.Errorw("ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//将生成的短链接加到布隆过滤器中 b.基于redis版本,go-zero自带的
	if err = l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw("BloomFilter.Add failed", logx.LogField{Key: "err", Value: err.Error()})
	}

	////将生成的短链接加到布隆过滤器中 a.基于内存版本服务重启之后就没了，所以每次重启就要重新加载一下已有的短链接(从数据库查询)
	//l.svcCtx.Filter.Add([]byte(short))

	//5.返回响应
	//5.1返回的是 短域名+短链接 q1mi/1En
	shortUrl := l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{ShortURL: shortUrl}, nil
}
