FROM golang:1.14

WORKDIR /go/src/app

COPY . . 
COPY ./entrypoint.sh /entrypoint.sh

# wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

RUN go get -d -v ./...
RUN go install -v ./...
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/jinzhu/gorm/dialects/mysql

RUN go get github.com/githubnemo/CompileDaemon


ENTRYPOINT [ "sh", "/entrypoint.sh" ]


EXPOSE 8080
