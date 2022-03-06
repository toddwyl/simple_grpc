go 安装

https://go.dev/doc/install

https://www.jetbrains.com/help/go/how-to-use-wsl-development-environment-in-product.html#wsl-general

建议还是不用vscode，用goland

https://www.jetbrains.com/go



简单的grpc应用

http://liuqh.icu/2022/01/20/go/rpc/03-grpc-ru-men/



安装protoc

```shell
sudo apt-get install autoconf automake libtool curl make g++ unzip unar
```

protoc

https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip

export protoc to zshrc



```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```



```shell
protoc -I=. --go_out=. --go-grpc_out=. ./user.proto 
```





```
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'todd123456';
```

