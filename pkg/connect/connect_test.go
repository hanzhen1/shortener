package connect

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//	func TestGet(t *testing.T) {
//		type args struct {
//			url string
//		}
//		tests := []struct {
//			name string
//			args args
//			want bool
//		}{
//			{name: "基本示例", args: args{url: "https://study.163.com"}, want: true},
//			{name: "无效的url示例", args: args{url: "https://study.163.co"}, want: false},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				if got := Get(tt.args.url); got != tt.want {
//					t.Errorf("Get() = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
// GoConvey是一个非常非常好用的Go测试框架
func TestGet(t *testing.T) {
	convey.Convey("基础用例", t, func() {
		url := "https://www.liwenzhou.com"
		got := Get(url)
		// 断言
		convey.So(got, convey.ShouldEqual, true) //断言got是否为true 相等判断

	})
	convey.Convey("url请求不通示例", t, func() {
		url := "https://www.wenzhou.com/posts/Go/unit-test-5/"
		got := Get(url)
		// 断言
		convey.So(got, convey.ShouldEqual, false)
	})
}
