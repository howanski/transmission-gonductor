FROM golang:1.16

WORKDIR /go/src/app
RUN go get -v github.com/gin-gonic/gin@v1.7.2 gorm.io/gorm@v1.20.7 gorm.io/driver/sqlite@v1.1.4
# RUN go install -v ./...
CMD ["/go/src/app/server_run.sh"]
