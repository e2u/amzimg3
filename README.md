
# 图片缓存以及缩放服务

提供对图片的缓存以及动态缩放服务

[TOC]


## 部署方式

* 编译 `make build`
* 创建 Docker 镜像 `REPOSITORY=<REPOSITORY:TAG> make build-docker`


## 运行方法

```
$ builds/amzimg3 -h

Usage of main:
  -address string
    	listen address (default "0.0.0.0")
  -allow string
    	allow srouce list file (default "/opt/amzimg3/etc/allow_sources.txt")
  -data string
    	cache image storage directory (default "/var/data")
  -port uint
    	listen port (default 8085)

```




## 配置文件 allow_sources

每行一个域名或 IP+PORT,如果源目标的主机不是 标准 http[s] 端口，则配置文件中的源需要注明端口

范例:

所有 `#` 开头的行都会被忽略

```
# 允许访问的源地址,忽略协议
127.0.0.1:9005
www.iconfinder.com
cdn0.iconfinder.com

```



## Docker 运行方法

* 临时运行: `docker run --rm -p 8085:8085 -v $(PWD)/etc:/opt/amzimg3/etc <REPOSITORY:TAG>`
* 运行为服务: `docker run -d --name amzimg3 -p 8085:8085 -v $(PWD)/etc/:/opt/amzimg3/etc <REPOSITORY:TAG>`



## 请求方法


### 原图缓存

```GET /<源图片地址>```

例:

* `http://host/http://www.domain.com:port/path/path1/file.jpg`
* `http://host/www.domain.com/file.jpg`

### 默认宽度缩放

默认宽度: 240px

```GET /r/<源图片地址>```

例:

* `http://host/r/http://www.domain.com:port/path/path1/file.jpg`
* `http://host/r/www.domain.com/path/path1/file.jpg`

### 自定宽度缩放

```GET /<宽度>/<源图片地址>```

例:

* `http://host/190/http://www.domain.com:port/path/path1/file.jpg`
* `http://host/190/www.domain.com/path/path1/file.jpg`


### 强制更新源图

Http Header reload: true

例:

* `curl -H "reload: true" http://host/190/http://www.domain.com:port/path/path1/file.jpg`
* `curl -H "reload: true" http://host/190/www.domain.com/path/path1/file.jpg`
