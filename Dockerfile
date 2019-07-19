FROM scratch

WORKDIR $GOPATH/src/create-gin-app
COPY . $GOPATH/src/create-gin-app

EXPOSE 80
ENTRYPOINT ["./create-gin-app"]
