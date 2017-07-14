FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    git-core \
    gcc \
    binutils \
    npm \
    nodejs

RUN curl -s https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz | tar -v -C /usr/local -xz
RUN mkdir -p /go
RUN ln -s /usr/bin/nodejs /usr/bin/node

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:/go/bin:/usr/local/bin:$PATH
WORKDIR $GOPATH

RUN go get -u github.com/EUDAT-GEF/GEF/gefserver
WORKDIR $GOPATH/src/github.com/EUDAT-GEF/GEF
RUN mkdir -p build

CMD ["go", "build", "-o", "./build/gef_linux", "./gefserver"]

