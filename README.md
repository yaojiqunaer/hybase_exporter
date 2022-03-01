# HyBase Server Exporter

HyBase exporter for HyBase server metrics.

支持版本:
* HyBase >= 8.0

## 构建和运行

### 编译

编译GO环境 GO版本>=1.17

    go build -o hybase_exporter 

### 运行

使用命令行参数运行

    ./mysqld_exporter --hybase.address='192.168.100.100:5555'

使用环境变量方式运行

    export HYBASE_ADDRESS='192.168.100.100:5555'
    ./mysqld_exporter <flags>



### 采集指标

| Name             | Description   |
| ---------------- | ------------- |
| hybase.cpu.usage | 海贝Cpu使用率 |

### 参数与环境变量
| 命令参数       | 环境变量       | 说明     |
| -------------- | -------------- | -------- |
| hybase.address | HYBASE_ADDRESS | 海贝地址 |
| hybase.user    | HYBASE_USER    | 海贝账号 |
| hybase.pwd     | HYBASE_PWD     | 海贝密码 |

注：命令行的参数设置的优先级大于环境变量


## Docker

#### Dockerfile

```dockerfile
FROM harbor.trscd.com.cn/iso/ailpine:1.20
LABEL maintainer="Authors <mpteam@trs.com.cn>"
WORKDIR /TRS/APP/
# 默认值
COPY hybase_exporter /TRS/APP/
ENV HYBASE_ADDRESS="127.0.0.1:5555" HYBSAE_USER="admin" HYBASE_PWD="admin"
EXPOSE      9555
USER        hybase
# 指定的命令不能被覆盖，docker run指定的参数当做ENTRYPOINT指定命令的参数。
ENTRYPOINT  [ "/TRS/APP/hybase_exporter"]
```

#### 构建

```
docker build .
```

#### 运行

使用自构建镜像，或仓库 [trs/hybase_exporter](https://registry.hub.docker.com/r/prom/mysqld-exporter/) 镜像

启动：

```bash
docker pull trs/hybase_exporter
docker run -d \
  -p 9555:9555 \
  -e HYBASE_ADDRESS="192.168.100.100" \
  trs/hybase_exporter
```
