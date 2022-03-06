package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"simple_grpc/cmd/global"
	"simple_grpc/cmd/web/router"
	"time"
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
	logrus.SetLevel(logrus.InfoLevel)
}

var (
	httpPort = ":8999"
	rpcHost  = "localhost:8003"
)

func setupRPCClient() error {
	ctx := context.Background()
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.DialContext(ctx, rpcHost, opt)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	global.GRPCClient = clientConn
	return nil
}

func main() {
	initLog()

	//client = zerorpc.NewGrettyClient(rpcServerAddr, 10, 10) // init rpc client instance
	//configRPC()                                             // config for rpc client
	err := setupRPCClient()
	if err != nil {
		logrus.Errorf("config rpc err:%v", err)
	}
	logrus.Infof("rpc config sucess, listen at:%s", rpcHost)

	httpServer := router.InitGin() // init http server with gin
	go func() {
		err := httpServer.Run(httpPort)
		if err != nil {
			logrus.Warnf("Gin run err:%v", err)
		}
	}() // listen and serve on 0.0.0.0:httpPort

	logrus.Infof("web app started, listen at:%s", httpPort)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		logrus.Info("web app graceful shutdown...")
		return
	}
}
