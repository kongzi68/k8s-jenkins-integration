# Compile stage
FROM golang:1.21.6 AS build-env
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /opt/betack
ADD . /opt/betack
RUN go build -o /opt/devops/bf-jenkins-api

# Final stage
FROM ubuntu:22.04
EXPOSE 8080
WORKDIR /opt/betack/
COPY config/*kube.config /opt/betack/config/
COPY --from=build-env /opt/devops/bf-jenkins-api /opt/betack/bf-jenkins-api
CMD ["/opt/betack/bf-jenkins-api"]