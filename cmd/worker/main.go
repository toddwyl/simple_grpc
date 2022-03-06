package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"simple_grpc/cmd/global"
	"simple_grpc/cmd/worker/tcp_service"
	"simple_grpc/proto"
	"time"
)

var (
	RPCPort = ":8003"
)

func initLog() {
	//
	logrus.SetReportCaller(true)
	//设置输出样式，自带的只有两种样式logger.JSONFormatter{}和logger.TextFormatter{}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, path.Base(frame.File)
		},
	})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.DebugLevel)
	//logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	initLog()
	fmt.Println("Finished init logger...")
	global.DBEngine = tcp_service.InitDB()
	fmt.Println("Finished init mysql...")
	global.CacheClient = tcp_service.InitRedis()
	fmt.Println("Finished init redis...")

	//opts := []grpc.ServerOption{
	//	grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
	//		middleware.Recovery,
	//	)),
	//}
	s := grpc.NewServer()
	c := context.Background()
	rpcServer := tcp_service.NewUserService(c)
	proto.RegisterUserServiceServer(s, &rpcServer)
	lis, err := net.Listen("tcp", RPCPort)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("net.Serve err: %v", err)
	}

}
