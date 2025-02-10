package main

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/database"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/kitex_gen/user/userservice"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: bluebell config.yaml")
		return
	}
	config.LoadConfig(os.Args[1])
	logging.InitLogger(config.AppConfig.App) // logging service
	// jaeger
	tp, err := tracing.SetTraceProvider(config.AppConfig.App.ServiceName)
	if err != nil {
		panic(err)
	}
	defer func(tp *sdkTrace.TracerProvider, ctx context.Context) {
		err = tp.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}(tp, context.Background())
	opts := kitexInit()

	userImpl := new(UserServiceImpl)
	// database
	db, err := database.InitMySQL(config.AppConfig.MySQL)
	if err != nil {
		panic(err)
	}
	userImpl.repo = repo.NewRepository(db)
	svr := userservice.NewServer(userImpl, opts...)

	err = svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// etcd
	r, err := etcd.NewEtcdRegistry(
		config.AppConfig.Etcd.Endpoints,
		etcd.WithAuthOpt(config.AppConfig.Etcd.Username, config.AppConfig.Etcd.Password),
	)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// address
	addr, err := net.ResolveTCPAddr("tcp", config.AppConfig.App.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: config.AppConfig.App.ServiceName,
	}))

	return
}
