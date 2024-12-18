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
	"hr-system/internal/cache"
	"hr-system/internal/common"
	employee_cache "hr-system/internal/employees/cache"
	"hr-system/internal/employees/handler"
	"hr-system/internal/employees/repo"
	"hr-system/internal/employees/service"
	"hr-system/internal/middleware"
)

var cachePrefixEmployee = "employee"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config, cause: %v", err)
	}
	fmt.Printf("config: %+v\n", cfg)

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
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully!")
	commonCache := cache.NewCache(rdb)

	employeeRepo, err := repo.NewEmployeeRepo(db)
	if err != nil {
		logger.Fatalf("Failed to New EmployeeRepo, cause: %v", err)
	}

	employeeService := service.NewEmployeeService(logger, employeeRepo,
		employee_cache.NewEmployeeCache(commonCache, cachePrefixEmployee))
	employeeHandler := handler.NewEmployeeHandler(logger, employeeService)

	r := gin.Default()
	r.Use(middleware.ContextMiddleware())

	r.POST("api/v1/employees", employeeHandler.CreateEmployee)
	r.GET("api/v1/employees/:id", employeeHandler.GetEmployeeByID)
	r.GET("api/v1/employees", employeeHandler.GetEmployees)

	logger.Fatalf(r.Run(fmt.Sprintf(":%s", cfg.RestServerPort)).Error())
}
