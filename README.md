
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

### 特 性

- 支持多种内置的数据结构
- 支持数据安全 IP 白名单访问功能
- 高吞吐量、低延迟、高效批量数据写入
- 支持磁盘数据存储和磁盘垃圾数据回收
- 支持磁盘数据静态加密和静态数据压缩
- 支持通过网络 RESTful API 数据操作协议

---

### 快速开始

使用容器管理工具可以快速部署 `wiredb:latest` 容器来测试 WireDB 提供的服务，使用下面命令即可快速部署运行 WireDB 镜像：

```shell
docker run -d --restart unless-stopped --name wiredb \
    -p 2668:2668 \
    -v $(pwd)/data:/tmp/wiredb \
    wiredb:latest
```

然而，对于长期运行的环境，建议直接在裸机上运行，以减少容器化带来的额外开销，并确保更好的性能和稳定性。
