package logic

import (
	"SAI-IM/apps/user/rpc/internal/svc"
	"SAI-IM/apps/user/rpc/user"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {
	type args struct {
		in *user.RegisterReq
	}

	// 定义一个匿名结构体 “数组/切片” 并同时进行初始化
	/**
		[]struct { ... }{
	    	{ ... }, // 第一个元素
	    	{ ... }, // 第二个元素
		}
	*/
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{ // 每个元素本身是一个结构体实例，每个元素都符合定义的结构
		{
			"1", args{in: &user.RegisterReq{
				Phone:    "13700001112",
				Nickname: "mingri",
				Password: "123456",
				Avatar:   "png.jpg",
				Sex:      1,
			}}, true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}

// 原生成测试模板，仅供参考
func TestRegisterLogic_Register1(t *testing.T) {
	type fields struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
		Logger logx.Logger
	}
	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.RegisterResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &RegisterLogic{
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
				Logger: tt.fields.Logger,
			}
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
