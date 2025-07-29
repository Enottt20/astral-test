package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Enottt20/astral-test"
	"github.com/Enottt20/astral-test/internal/handler"
	"github.com/Enottt20/astral-test/internal/service"
	"github.com/Enottt20/astral-test/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Enottt20/astral-test/docs" // swag docs генерируются сюда
)

func main() {
	// Логгер в JSON формате
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	// Загрузка конфига
	if err := InitConfig(); err != nil {
		logrus.Fatalf("Failed to load config: %s", err.Error())
	}

	// Настройка Swagger BasePath (если есть префикс API)
	// Например, если все роуты начинаются с /api, то:
	// docs.SwaggerInfo.BasePath = "/api"
	// Если нет - оставь "/"
	// docs.SwaggerInfo.BasePath = "/"

	// Подключение к Postgres
	db, err := storage.NewPostgresDB(storage.PostgresConfig{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		DBName:   viper.GetString("postgres.dbname"),
		SSLMode:  viper.GetString("postgres.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Failed to connect Postgres DB: %s", err.Error())
	}
	defer db.Close()

	// Миграции базы
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Fatalf("Failed to create migrate driver: %s", err.Error())
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsPath := "file://" + filepath.Join(dir, "../migrations")

	migrations, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logrus.Fatalf("Failed to create migrate instance: %s", err.Error())
	}
	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Migrations error: %s", err.Error())
	}

	// Репозиторий и Redis
	repo := storage.NewRepository(db)

	intRedisDB, err := strconv.Atoi(viper.GetString("redis.db"))
	if err != nil {
		logrus.Fatal("Invalid value of redis.db")
	}

	cache, err := service.NewRedisClient(service.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       intRedisDB,
	})
	if err != nil {
		logrus.Fatalf("Failed to connect Redis: %s", err.Error())
	}
	defer cache.Close()

	// Инициализация сервисов, эндпоинтов и сервера
	services := service.NewService(repo, viper.GetString("adminToken"), cache)
	endp := handler.NewEndpoint(services)

	// Gin роутер
	router := endp.InitRoutes()

	// Добавляем Swagger UI роут
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &astraltest.Server{}

	// Запуск сервера в отдельной горутине
	go func() {
		port := viper.GetString("port")
		logrus.Infof("Swagger docs available at: http://localhost:%s/swagger/index.html", port)

		if err := server.Run(port, router); err != nil {
			logrus.Fatalf("Failed to run server: %s", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Server shutdown error: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
