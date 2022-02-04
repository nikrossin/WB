FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.html ./

RUN go build -o /docker-task

EXPOSE 8080

CMD [ "/docker-task" ]