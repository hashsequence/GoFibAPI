FROM golang:onbuild

EXPOSE 3000

RUN mkdir ./goFibApi

WORKDIR ./goFIbApi

COPY . .

RUN go build -o fibApi

CMD ["./fibApi"]