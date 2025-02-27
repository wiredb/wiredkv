
<div align="center">
    <img src="cmd/wiredb-org.png" style="width: 86px; height: auto; display: inline-block;">
</div>

<p align="center">WireDB is a NoSQL that supports multiple data types based on Log-structured file system.</p>


---


[![Go Report Card](https://goreportcard.com/badge/github.com/auula/wiredb)](https://goreportcard.com/report/github.com/auula/wiredb)
[![Go Reference](https://pkg.go.dev/badge/github.com/auula/wiredb.svg)](https://pkg.go.dev/github.com/auula/wiredb)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/55bc449808ca4d0c80c0122f170d7313)](https://app.codacy.com/gh/auula/wiredb/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codecov](https://codecov.io/gh/wiredb/wiredb/graph/badge.svg?token=ekQ3KzyXtm)](https://codecov.io/gh/wiredb/wiredb)
[![DeepSource](https://app.deepsource.com/gh/wiredb/wiredb.svg/?label=active+issues&show_trend=true&token=sJBjq88ZxurlEgiOu_ukQ3O_)](https://app.deepsource.com/gh/wiredb/wiredb/?ref=repository-badge)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/wiredb/wiredb.svg)](https://github.com/wiredb/wiredb/releases)



---

[简体中文](#) | [English](#)

---

### 🌟 特 性

- 支持多种内置的数据结构
- 支持数据安全 IP 白名单访问功能
- 高吞吐量、低延迟、高效批量数据写入
- 支持磁盘数据存储和磁盘垃圾数据回收
- 支持磁盘数据静态加密和静态数据压缩
- 支持通过基于 RESTful API 协议操作数据

---

### 🚀 快速开始

使用容器管理工具可以快速部署 `wiredb:v1.0.0` 镜像来测试 WireDB 提供的服务。只需运行以下命令，即可拉取 WireDB 镜像并启动一个容器运行 WireDB 服务：

```shell
docker pull auula/wiredb:v1.0.0
```

运行 WireDB 镜像启动容器服务并且映射端口到外部主机网络，如何下面命令：

```shell
docker run -p 2668:2668 wiredb:v1.0.0
```

注意通过 RESTful API HTTP 协议操作数据时请在 HTTP 请求头中添加 `Auth-Token: xxxx` 访问密码。此默认访问密码为 WireDB 进程自动生成的，需要通过查看容器运行输出的日志信息获得，命令如下：

```shell
root@2c2m:~# docker logs 46ae91bc73a6
                         _            ____
                 _    __(_)______ ___/ / /
                | |/|/ / / __/ -_) _  / _ \
                |__,__/_/_/  \__/\_,_/_.__/  v1.0.0

  WireDB is a NoSQL database based on Log-structured file system.
  Software License: Apache 2.0  Website: https://wiredb.github.io

[WIREDB:C] 2025/02/27 10:07:01 [WARN]	The default password is: T9EHAvi5dcIpPK9G#ADlVj4NB 👈
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Logging output initialized successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Loading and parsing region data files...
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	Region compression activated successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	File system setup completed successfully
[WIREDB:C] 2025/02/27 10:07:01 [INFO]	HTTP server started at http://172.0.0.1:2668 🚀
```

推荐使用 Linux 发型版本来运行 WireDB 服务，WireDB 服务进程依赖配置文件中的参数，在运行 WireDB 服务之前将下面的配置内容写到 `config.yaml` 中：

```yaml
port: 2668                              # 服务 HTTP 协议端口
mode: "std"                             # 默认为 std 标准库
path: "/tmp/wiredb"                     # 数据库文件存储目录
auth: "Are we wide open to the world?"  # 访问 HTTP 协议的秘密
logpath: "/tmp/wiredb/out.log"          # WireDB 在运行时程序产生的日志存储文件
debug: false        # 是否开启 debug 模式
region:             # 数据区
    enable: true    # 是否开启数据压缩功能
    second: 1800    # 默认垃圾回收器执行周期单位为秒
    threshold: 3    # 默认个数据文件大小，单位 GB
encryptor:          # 是否开启静态数据加密功能
    enable: false
    secret: "your-static-data-secret!"
compressor:         # 是否开启静态数据压缩功能
    enable: false
allowip:            # 白名单 IP 列表，可以去掉这个字段，去掉之后白名单就不会开启
    - 192.168.31.221
    - 192.168.101.225
    - 127.0.0.1
```




