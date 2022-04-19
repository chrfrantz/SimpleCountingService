FROM golang:1.16 as builder

LABEL maintainer="author@mail.com"

WORKDIR /

ADD ./go.mod /
ADD ./main.go /

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main


FROM scratch

WORKDIR /

COPY --from=builder /main .

EXPOSE 8080

CMD ["/main"]


