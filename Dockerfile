FROM golang:latest AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
ENV APP_ENV=prod
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o goorm-class ./cmd/api
COPY configuration configuration 

# 5. 최종 이미지 생성
FROM build

WORKDIR /app

# 필요한 런타임 의존성 복사
COPY --from=build /app/goorm-class .
COPY --from=build /app/configuration configuration

RUN chmod +x ./goorm-class

EXPOSE 8080 
# 어플리케이션 실행
CMD ["./goorm-class"]
