# ChannelServer Docker (production)
#
# 此 Dockerfile 从 `Dockerfile.build` 重定向的管道输出作为stdin创建一个新的容器.
#
# 首先创建 builder image:
#
#   ```
#   docker build -t channel-server-builder -f Dockerfile.build .
#   ```
#
# 接下来运行 builder 容器, 把它的输出管道定向到 runner container的创建上:
#
#   ```
#   docker run --rm channel-server-builder | docker build -t channel-server -f Dockerfile.run -
#   ```
#
# image. 之后按如下运行容器:
#
#   ```
#   docker run --rm --name channel-server -p 8080:8080 -p 8443:8443 -v `pwd`:/srv/extra -i -t channel-server
#   ```
#
# 现在你可以使用一个前端代理如 Nginx 来提供 TLS 给 ChannelServer 并从 Docker 容器来运行 production 版本. 
# 为了方便开发测试, 容器还提供了一个内嵌的自签名的TLS listener运行在8443端口
#
# 为了方便使用自定义配置, 使用 `server.conf.in` 文件作为模版 并且从 [http] 和 [https] 小节移除 listeners . 之后可用 `-c` 参数指定配置文件来运行docker容器:
#
#   ```
#   docker run --rm --name channel-server -p 8080:8080 -p 8443:8443 \
#       -v `pwd`:/srv/extra -i -t channel-server \
#       -c /srv/extra/server.conf
#   ```
#
# 最后, 容器在启动时会检查环境变量 NEWCERT 和 NEWSECRETS. 在启动时把它们设置为 `1` 来重新生成新数据.
# 当前的 certificate 和 secrets 会在启动前打印出来，方便你用于其它服务. 
# 当然，如果你想持久化 cert 和 secrets, 容器需要被持久化, 不需要 `--rm` 参数.
#

FROM frolvlad/alpine-glibc:alpine-3.3_glibc-2.23
MAINTAINER edison <52388483@qq.com>

ENV LANG=C.UTF-8

# Add dependencies.
RUN apk add --no-cache \
	openssl

# Add channel server as provided by Dockerfile.run.
COPY srv/ /srv

# Move around stuff from tarball to their expected locations.
RUN mv /srv/channel-server/dist/loader/* /srv/channel-server && \
	mv /srv/channel-server/dist/www/html /srv/channel-server && \
	mv /srv/channel-server/dist/www/static /srv/channel-server

# Add entrypoint.
COPY docker_entrypoint.sh /srv/entrypoint.sh

# Create default config.
RUN cp -v /srv/channel-server/server.conf.in /srv/channel-server/default.conf && \
	sed -i 's|listen = 127.0.0.1:8080|listen = 0.0.0.0:8080|' /srv/channel-server/default.conf && \
	sed -i 's|;root = .*|root = /srv/channel-server|' /srv/channel-server/default.conf && \
	sed -i 's|;listen = 127.0.0.1:8443|listen = 0.0.0.0:8443|' /srv/channel-server/default.conf && \
	sed -i 's|;certificate = .*|certificate = /srv/cert.pem|' /srv/channel-server/default.conf && \
	sed -i 's|;key = .*|key = /srv/privkey.pem|' /srv/channel-server/default.conf && \
	touch /etc/channel-server.conf

# Cleanup.
RUN rm -rf /tmp/* /var/cache/apk/*

# Add mount point for extra things.
RUN mkdir /srv/extra
VOLUME /srv/extra

# Tell about our service.
EXPOSE 8080
EXPOSE 8443

# Define entry point with default command.
ENTRYPOINT ["/bin/sh", "/srv/entrypoint.sh", "-dc", "/srv/channel-server/default.conf"]
CMD ["-c", "/etc/channel-server.conf"]