package md5

import "testing"

func TestSum(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "基本示例", args: args{data: []byte("https://www.baidu.com")}, want: "f9751de431104b125f48dd79cc55822a"},
		{name: "示例", args: args{data: []byte("11111111")}, want: "1bbd886460827015e5d605ed44252251"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.data); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestSum(t *testing.T) {
//	convey.Convey("基础示例", t, func() {
//		data := []byte("https://www.liwenzhou.com/posts/Go/unit-test-5/")
//		want := "9c2970c150c78443d688c9602f8b1535"
//		got := Sum(data)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want) //断言got和want是否一致
//
//	})
//	convey.Convey("示例", t, func() {
//		data := []byte("111111111111112222222222")
//		want := "2fa1a6a7818bdb8aac6aecbfa2d65e87"
//		got := Sum(data)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want)
//	})
//}
