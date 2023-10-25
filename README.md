# mdns-server
简单的mdns服务提供，适用于局域网内服务器发布服务域名

备注: 不要在win下面跑，没成功过

# 打包

```shell
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./mdns-server
```