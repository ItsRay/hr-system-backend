package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"

	"github.com/gin-gonic/gin"

	"hr-system/config"
	"hr-system/internal/cache"
	"hr-system/internal/common"
	employee_cache "hr-system/internal/employees/cache"
	employee_handler "hr-system/internal/employees/handler"
	employee_repo "hr-system/internal/employees/repo"
	employee_service "hr-system/internal/employees/service"
	leave_cache "hr-system/internal/leaves/cache"
	leave_handler "hr-system/internal/leaves/handler"
	leave_repo "hr-system/internal/leaves/repo"
	leave_service "hr-system/internal/leaves/service"
	"hr-system/internal/middleware"
)

var cachePrefixEmployee = "employee"
var cachePrefixLeave = "leave"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config, cause: %v", err)
	}

	logger := common.NewLogger()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDBName)
	maxRetries := 10
	db, err := connectMySqlWithRetry(logger, dsn, maxRetries, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL after %d attempts: %v", maxRetries, err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	commonCache := cache.NewCache(rdb)

	r := gin.Default()
	r.Use(middleware.ContextMiddleware())

	// API for employees
	employeeRepo, err := employee_repo.NewEmployeeRepo(db)
	if err != nil {
		logger.Fatalf("Failed to New EmployeeRepo, cause: %v", err)
	}
	if err = employeeRepo.SeedData(ctx); err != nil {
		logger.Fatalf("Failed to seed data, cause: %v", err)
	}

	employeeService := employee_service.NewEmployeeService(logger, employeeRepo,
		employee_cache.NewEmployeeCache(commonCache, cachePrefixEmployee))
	employeeHandler := employee_handler.NewEmployeeHandler(logger, employeeService)
	r.POST("api/v1/employees", employeeHandler.CreateEmployee)
	r.GET("api/v1/employees/:id", employeeHandler.GetEmployeeByID)
	r.GET("api/v1/employees", employeeHandler.GetEmployees)

	// API for leaves
	leaveRepo, err := leave_repo.NewLeaveRepo(db)
	if err != nil {
		logger.Fatalf("Failed to New leaveRepo, cause: %v", err)
	}
	if err = leaveRepo.SeedData(ctx, employeeRepo); err != nil {
		logger.Fatalf("Failed to seed data, cause: %v", err)
	}
	leaveService := leave_service.NewLeaveService(logger, leaveRepo, employeeRepo,
		leave_cache.NewLeaveCache(commonCache, cachePrefixLeave))
	leaveHandler := leave_handler.NewLeaveHandler(logger, leaveService)
	// TODO: add API for revoking leave
	r.POST("api/v1/leaves", leaveHandler.CreateLeave)
	r.POST("api/v1/leaves/:id/review", leaveHandler.ReviewLeave)
	r.GET("api/v1/leaves", leaveHandler.GetLeaves)
	r.GET("api/v1/leaves/:id", leaveHandler.GetLeaveByID)

	logger.Fatalf(r.Run(fmt.Sprintf(":%s", cfg.RestServerPort)).Error())
}

func connectMySqlWithRetry(logger *common.Logger, dsn string, maxRetries int, retryDelay time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			logger.Infof("Successfully connected to MySQL on attempt %d", i+1)
			return db, nil
		}

		logger.Warnf("Failed to connect to MySQL (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}
