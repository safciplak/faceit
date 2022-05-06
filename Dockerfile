FROM golang:1.17-alpine3.15

ENV TZ Europe/Istanbul
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go test --count=1 -coverprofile=coverage.out ./.. ; cat coverage.out | awk 'BEGIN {cov=0; stat=0;} \
                                             $3!="" { cov+=($3==1?$2:0); stat+=$2; } \
                                         END {printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...

EXPOSE 8080
CMD [ "restapi" ]
