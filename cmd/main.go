package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"

	"github.com/gin-gonic/gin"

	"hr-system/config"
	"hr-system/internal/common"
	"hr-system/internal/employees/handler"
	"hr-system/internal/employees/repo"
	"hr-system/internal/employees/service"
)

func main() {
	fmt.Printf("main first line!!")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config, cause: %v", err)
	}
	fmt.Printf("config: %v\n", cfg)

	logger := common.NewLogger()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDBName)
	fmt.Printf("dsn: %s\n", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to connect to MySQL: %v", err)
	}
	fmt.Println("Connected to MySQL successfully!")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	// 驗證 Redis 連線
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully!")

	employeeRepo, err := repo.NewEmployeeRepo(db)
	if err != nil {
		logger.Fatalf("Failed to New EmployeeRepo, cause: %v", err)
	}

	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	r := gin.Default()
	r.POST("/employees", employeeHandler.CreateEmployee)

	logger.Fatal(r.Run(":8080"))
}
