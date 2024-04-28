package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Err404 = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 自己写缓存  surl->lurl
// go-zero自带的缓存 surl ->整个数据行
func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	//查看短链接，输入去q1mi/h1hu->重定向到真实的链接
	//req.ShortURL=h1hu
	//1.根据短链接查询原始的长链接
	//1.0布隆过滤器
	//不存在的短链接直接返回404，不需要后续处理
	////a.基于内存版本,服务重启之后就没了，所以每次重启就要重新加载一下已有的短链接(从数据库查询)
	//if ok := l.svcCtx.Filter.Test([]byte(req.ShortURL)); !ok {
	//	return nil, Err404
	//}
	//b.基于redis版本,go-zero自带的
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortURL))
	if err != nil {
		logx.Errorw(" Bloom Filter.Exists failed", logx.LogField{Key: "err", Value: err.Error()})
	}
	//不存在的短链接直接返回
	if !exist {
		return nil, Err404
	}
	fmt.Println("开始查询缓存和DB....")
	//1.1 查询数据库之前可增加缓存
	//go-zero缓存支持singleflight 同时100w个请求q1mi/jts4，恰巧缓存中jts4失效了,这就是缓存击穿
	//使用singleflight可以解决缓存击穿，合并并发大量的请求，第一个请求先去查db,并发的后999999请求等他拿到结果返回，不需要自己去查
	long, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortURL, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows { //没有查到
			return nil, Err404
		} //查数据库出错
		logx.Errorw(" show_ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, errors.New("内部错误")
	}
	//2.返回查询到的长链接，在调用handler层返回重定向的响应
	return &types.ShowResponse{LongURL: long.Lurl.String}, nil
}
