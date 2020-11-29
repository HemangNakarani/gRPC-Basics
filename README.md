# gRPC-Basics
Kicking off with golang and gRPC 🤩🥳

# TechStack
- Golang, gRPC driver and Protocol Buffers
- MongoDb drivers for Golang
- Evans CLI for testing Go server without writing client code. (It will need to use Reflection API in server code)

# Instructions
- Download golang and set GOPATH in environment
- Get `protoc` compiler from its Repo using ```go get <LINK>```
- For Generating code from `.proto` visit `generate.sh` file in repo

- Clone Repo in Local Machine

In one Terminal
```
> cd gRPC-Basics
> go run hemangnakarani\calculator\clac_server\server.go
```
In another Terminal
```
> cd gRPC-Basics
> go run hemangnakarani\calculator\calc_client\client.go
```
