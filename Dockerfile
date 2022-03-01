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