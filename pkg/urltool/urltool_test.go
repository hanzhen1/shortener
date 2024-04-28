package urltool

import "testing"

func TestGetBasePath(t *testing.T) {
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "基本示例",
			args:    args{targetUrl: "https://www.liwenzhou.com/posts/Go/golang-menu/"},
			want:    "golang-menu",
			wantErr: false,
		},
		{name: "相对路径url示例", args: args{targetUrl: "/Go/golang-menu/"}, want: "", wantErr: true},
		{name: "空字符串", args: args{targetUrl: ""}, want: "", wantErr: true},
		{name: "带query的url示例", args: args{targetUrl: "https://www.liwenzhou.com/posts/Go/golang-menu/?a=1&b=2"}, want: "golang-menu", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBasePath(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBasePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestGetBasePath(t *testing.T) {
//	convey.Convey("基础示例", t, func() {
//		url := "https://www.liwenzhou.com/posts/Go/unit-test-5/"
//		want := "unit-test-5"
//		got, _ := GetBasePath(url)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want) //断言got和want是否一致
//	})
//	convey.Convey("相对路径url示例", t, func() {
//		url := "/Go/unit-test-5/"
//		want := ""
//		got, _ := GetBasePath(url)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want) //断言got和want是否一致
//	})
//	convey.Convey("空字符串", t, func() {
//		url := ""
//		want := ""
//		got, _ := GetBasePath(url)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want) //断言got和want是否一致
//	})
//	convey.Convey("带query的url示例", t, func() {
//		url := "https://www.liwenzhou.com/posts/Go/unit-test-5/?name=hz"
//		want := "unit-test-5"
//		got, _ := GetBasePath(url)
//		// 断言
//		convey.So(got, convey.ShouldResemble, want) //断言got和want是否一致
//	})
//}
