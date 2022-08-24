FROM golang:alpine AS build
ENV CGO_ENABLED=0
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o /out/gateway main.go

FROM scratch AS bin
COPY --from=build /out/gateway /gateway
CMD [ "/gateway", "server" ]