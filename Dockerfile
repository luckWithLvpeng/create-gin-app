FROM scratch

WORKDIR $GOPATH/src/eme
COPY . $GOPATH/src/eme

EXPOSE 80
ENTRYPOINT ["./eme"]
