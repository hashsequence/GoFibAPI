FROM golang:onbuild

EXPOSE 8080

RUN mkdir ./goFibApi

WORKDIR ./goFIbApi

COPY . .

RUN go build -o fibApi

CMD ["./fibApi"]