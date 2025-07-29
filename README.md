- **Запуск: `docker compose up --d`**
- **Swagger: http://localhost:8080/swagger/index.html**
- **ТЗ: [Открыть задание (PDF)](docs/TestTask.pdf)**
```
.
├── Dockerfile
├── README.md
├── TestTask.pdf
├── cmd
│   └── main.go
├── config
│   └── config.yml
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── domain
│   │   ├── document.go
│   │   ├── response.go
│   │   └── user.go
│   ├── handler
│   │   ├── auth.go
│   │   ├── documents.go
│   │   ├── handler.go
│   │   └── middleware.go
│   ├── service
│   │   ├── auth.go
│   │   ├── documents.go
│   │   ├── redis.go
│   │   ├── service.go
│   │   └── users.go
│   └── storage
│       ├── documents_postgres.go
│       ├── postgres.go
│       ├── repository.go
│       └── users_postgres.go
├── migrations
│   ├── 1_init.down.sql
│   └── 1_init.up.sql
└── server.go
```
