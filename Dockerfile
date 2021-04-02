FROM golang:1.13 AS build

WORKDIR /go/src

COPY server ./server
COPY go.sum ./go.sum
COPY go.mod ./go.mod
COPY main.go .

ENV CGO_ENABLED=0
RUN go mod download

RUN go build -a -installsuffix cgo -o nxsiteman .

FROM scratch AS runtime
COPY --from=build /go/src/nxsiteman ./
EXPOSE 8080/tcp
ENTRYPOINT ["./nxsiteman"]
