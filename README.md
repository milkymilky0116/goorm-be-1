# 구름톤 연합스터디 - 백엔드 과제 1

이 프로젝트는 구름톤 백엔드 연합스터디 과제 1을 구현한 결과물입니다.

## 요구사항

- Docker (필수)
- go >= 1.22.3 (선택)
- goose >= 3.20.0 (선택)

## 실행전

해당 폴더에 `configuration` 폴더를 만든 뒤 아래와 같은 형식의 `config.prod.yaml` 파일을 생성합니다.

```yaml
applicationPort: 8080
database:
  username: admin
  password: 123
  host: db
  port: 5432
  dbName: goorm-class
jaeger:
  host: jaeger
  port: 4318
```

docker compose 를 올바르기 실행하기 위해, 같은 루트 폴더에 `.env` 파일을 생성하고 아래와 같은 형식의 환경변수를 설정합니다.

```bash
DB_USERNAME=admin
DB_PASSWORD=123
DB_NAME=goorm-class
DB_HOST=127.0.0.1
DB_PORT=5432
```

## 실행

시스템에 Docker가 설치되어 있을 경우, 아래 커맨드로 애플리케이션을 실행합니다.

```bash
docker compose up -d
```

로컬에서 웹서버를 실행시키려는 경우, 아래의 커맨드로 웹서버 애플리케이션을 실행합니다.

```bash
go run ./cmd/api
```
