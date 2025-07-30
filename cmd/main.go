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

	_ "github.com/Enottt20/astral-test/docs"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	if err := InitConfig(); err != nil {
		logrus.Fatalf("Failed to load config: %s", err.Error())
	}

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

	services := service.NewService(repo, viper.GetString("adminToken"), cache)
	endp := handler.NewEndpoint(services)

	router := endp.InitRoutes()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &astraltest.Server{}

	go func() {
		port := viper.GetString("port")
		logrus.Infof("Swagger docs available at: http://localhost:%s/swagger/index.html", port)

		if err := server.Run(port, router); err != nil {
			logrus.Fatalf("Failed to run server: %s", err.Error())
		}
	}()

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
