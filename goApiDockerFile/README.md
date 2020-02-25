# Building DockerImage and pushing into dockerhub

```console
$ sudo docker build -t gofibapi .
Sending build context to Docker daemon  19.97kB
Step 1/7 : FROM golang:onbuild
# Executing 3 build triggers
 ---> Using cache
 ---> Using cache
 ---> Using cache
 ---> 086c846e4abd
Step 2/7 : EXPOSE 3000
 ---> Using cache
 ---> d3917f097d26
Step 3/7 : RUN mkdir ./goFibApi
 ---> Using cache
 ---> 0aa747756ad4
Step 4/7 : WORKDIR ./goFIbApi
 ---> Using cache
 ---> 458ca3a5d081
Step 5/7 : COPY . .
 ---> Using cache
 ---> d26d22136213
Step 6/7 : RUN go build -o fibApi
 ---> Using cache
 ---> 9ee0f22ba529
Step 7/7 : CMD ["./fibApi"]
 ---> Using cache
 ---> 986197f7fed8
Successfully built 986197f7fed8
Successfully tagged gofibapi:latest
$ sudo docker login --username=hashsequence
$ sudo docker push hashsequence/gofibapi


```

running dockerimage 

```console
sudo docker container run --rm -p 80:3000 hashsequence/gofibapi
```

go on browser and run:


http://localhost/current

http://localhost/prev

http://localhost/next