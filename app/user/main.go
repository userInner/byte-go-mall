package main

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/database"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/kitex_gen/cart/cartservice"
	"byte-go-mall/kitex_gen/user/userservice"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	// 1. 配置文件检查
	if len(os.Args) < 2 {
		klog.Fatal("缺少配置文件。使用方式: ./program config.yaml")
		return
	}

	fmt.Println("正在加载配置文件...")
	config.LoadConfig(os.Args[1])
	fmt.Println("配置文件加载成功")

	// 2. 初始化日志服务
	fmt.Println("正在初始化日志服务...")
	logging.InitLogger(config.AppConfig.App)
	fmt.Println("日志服务初始化成功")

	// 3. 初始化链路追踪
	fmt.Println("正在初始化Jaeger追踪...")
	tp, err := tracing.SetTraceProvider(config.AppConfig.App.ServiceName)
	if err != nil {
		klog.Fatalf("Jaeger初始化失败: %v", err)
		return
	}
	fmt.Println("Jaeger追踪初始化成功")

	// 设置trace provider的清理
	defer func(tp *sdkTrace.TracerProvider, ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			klog.Errorf("Trace Provider关闭失败: %v", err)
		}
	}(tp, context.Background())

	// 4. 初始化Kitex配置
	fmt.Println("正在初始化Kitex配置...")
	userOpts, err := kitexInit(config.AppConfig.App.UserAddress) // 使用 User 服务端口
	if err != nil {
		fmt.Printf("User服务Kitex初始化失败，错误详情: %v\n", err)
		klog.Fatalf("User服务Kitex初始化失败: %v", err)
		return
	}
	fmt.Println("User服务Kitex配置初始化成功")

	cartOpts, err := kitexInit(config.AppConfig.App.CartAddress) // 使用 Cart 服务端口
	if err != nil {
		fmt.Printf("Cart服务Kitex初始化失败，错误详情: %v\n", err)
		klog.Fatalf("Cart服务Kitex初始化失败: %v", err)
		return
	}
	fmt.Println("Cart服务Kitex配置初始化成功")

	// 5. 初始化数据库
	fmt.Println("正在初始化数据库连接...")
	db, err := database.InitMySQL(config.AppConfig.MySQL)
	if err != nil {
		klog.Fatalf("数据库初始化失败: %v", err)
		return
	}
	fmt.Println("数据库连接初始化成功")

	// 6. 初始化服务实现
	userImpl := new(UserServiceImpl)
	userImpl.repo = repo.NewRepository(db)
	cartImpl := new(CartServiceImpl)
	cartImpl.repo = repo.NewRepository(db)

	// 7. 创建并启动服务器
	userSvr := userservice.NewServer(userImpl, userOpts...)
	cartSvr := cartservice.NewServer(cartImpl, cartOpts...)

	go func() {
		fmt.Printf("正在启动User服务器，地址: %s...\n", config.AppConfig.App.UserAddress)
		if err := userSvr.Run(); err != nil {
			klog.Fatalf("User服务运行失败: %v", err)
		}
	}()

	fmt.Printf("正在启动Cart服务器，地址: %s...\n", config.AppConfig.App.CartAddress)
	if err := cartSvr.Run(); err != nil {
		klog.Fatalf("Cart服务运行失败: %v", err)
	}
}

// kitexInit 初始化Kitex配置
func kitexInit(address string) ([]server.Option, error) {
	var opts []server.Option

	// 检查 ETCD 配置
	fmt.Printf("ETCD配置信息检查:\n")
	fmt.Printf("Endpoints: %v\n", config.AppConfig.Etcd.Endpoints)
	fmt.Printf("是否配置认证: %v\n", config.AppConfig.Etcd.Username != "")

	// 尝试连接 ETCD
	fmt.Println("正在连接ETCD...")
	etcdRegistry, err := etcd.NewEtcdRegistry(
		config.AppConfig.Etcd.Endpoints,
		etcd.WithAuthOpt(config.AppConfig.Etcd.Username, config.AppConfig.Etcd.Password),
	)
	if err != nil {
		fmt.Printf("ETCD连接错误详情: %v\n", err)
		return nil, fmt.Errorf("ETCD连接失败: %v", err)
	}

	// 使用 etcdRegistry
	opts = append(opts, server.WithRegistry(etcdRegistry))
	fmt.Println("ETCD连接成功")

	// 解析服务地址
	fmt.Printf("正在解析服务地址: %s\n", address)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("地址解析失败: %v", err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// 设置服务信息
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: config.AppConfig.App.ServiceName,
	}))

	return opts, nil
}
