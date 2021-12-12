FROM alpine as build
# install build tools
RUN apk add go git
RUN go env -w GOPROXY=direct
# cache dependencies
COPY go.mod go.sum ./
RUN go mod download
# build
COPY *.go ./
RUN go build -o /main
# copy artifacts to a clean image
FROM alpine
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /
RUN chmod 755 /entry.sh
COPY --from=build /main /main
ENTRYPOINT [ "/entry.sh" ]
