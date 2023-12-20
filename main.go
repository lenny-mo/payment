package main

import (
	"fmt"
	"strconv"

	"github.com/lenny-mo/emall-utils/tracer"
	"github.com/lenny-mo/payment/conf"
	"github.com/lenny-mo/payment/domain/dao"
	"github.com/lenny-mo/payment/domain/models"
	"github.com/lenny-mo/payment/domain/services"
	"github.com/lenny-mo/payment/handler"
	"github.com/lenny-mo/payment/middleware"
	"github.com/lenny-mo/payment/proto/payment"
	"github.com/lenny-mo/payment/utils"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1 配置中心 找到配置文件
	consulCof, err := conf.GetConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		fmt.Println(err)
		fmt.Println("获取配置失败")
		panic(err)
	}

	// 2. 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 3 链路追踪
	serviceName := "go.micro.service.payment"
	// 3 链路追踪
	err = tracer.InitTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tracer.Closer.Close()
	opentracing.SetGlobalTracer(tracer.Tracer)

	// 4. 获取mysql配置
	mysqlConf := conf.GetMysqlFromConsul(consulCof, "mysql")

	// 5. 初始化数据库连接
	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + strconv.FormatInt(mysqlConf.Port, 10) + ")/" + mysqlConf.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if !db.Migrator().HasTable(&models.Payment{}) {
		db.Migrator().CreateTable(&models.Payment{})
	}

	// 获取redis配置
	redisConf := conf.GetRedisFromConsul(consulCof, "redis")
	// 初始化redis连接
	middleware.StartRDB(redisConf.Host, redisConf.Port, redisConf.DB, redisConf.PoolSize, redisConf.Password)
	defer middleware.Close()

	// 6 设置prometheus
	utils.PrometheusBoot(9092)

	// 7 创建服务
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"), // 服务监听地址
		// 使用consul注册中心
		micro.Registry(consulRegistry),
		// uber 漏桶 添加限流 每秒处理1000·个请求
		micro.WrapHandler(ratelimit.NewHandlerWrapper(conf.QPS)),
		// 添加prometheus
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		// 客户端链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 客户端pro
		micro.WrapClient(prometheus.NewClientWrapper()),
	)

	service.Init()

	// 8. 创建service 和 handler 并且注册服务
	paymentDAO := dao.NewPaymentDAO(db)
	paymentService := services.NewPaymentService(*paymentDAO.(*dao.PaymentDAO))
	err = payment.RegisterPaymentServiceHandler(service.Server(), &handler.PaymentHandler{
		PaymentService: *paymentService.(*services.PaymentService),
	})

	if err != nil {
		panic(err)
	}

	// 9. 启动service
	if err = service.Run(); err != nil {
		fmt.Println(err)
		panic(err)
	}

}
