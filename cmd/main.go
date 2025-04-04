package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crs-backend/internal/database"
	"crs-backend/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func main() {
	database.ConnectDB()
	// 1. Initialize configuration
	initConfig()

	// 2. Initialize logger
	initLogger()

	// 3. Initialize database (اتصال برقرار شده و مهاجرت داخل تابع انجام می‌شود)
	initDatabase()

	// 4. Create Gin router with global middlewares
	router := initRouter()

	// 5. Register routes
	registerRoutes(router)

	// 6. Start server with graceful shutdown
	startServer(router)
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config: %v", err)
	}
}

func initLogger() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	if gin.Mode() == gin.DebugMode {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

func initDatabase() {
	// استفاده از تابع ConnectDB که در package database تعریف شده است
	database.ConnectDB()
	// در صورت نیاز می‌توانید به اتصال ایجاد شده دسترسی داشته باشید:
	// db := database.DB
}

func initRouter() *gin.Engine {
	if viper.GetString("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// استفاده از میانجی‌های عمومی شامل Logger, Recovery, CORS و SecureHeaders
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(middlewares.SecureHeaders())

	return router
}

func registerRoutes(r *gin.Engine) {
	// مسیرهای عمومی
	public := r.Group("/api")
	{
		public.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello from public endpoint!"})
		})
	}
	// مسیرهای محافظت‌شده با استفاده از JWTAuth
	protected := r.Group("/api")
	protected.Use(middlewares.JWTAuth())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello from protected endpoint!"})
		})
	}
}

func startServer(r *gin.Engine) {
	var tlsConfig *tls.Config
	certPath := viper.GetString("SSL_CERT_PATH")
	keyPath := viper.GetString("SSL_KEY_PATH")
	if certPath != "" && keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			logger.Fatalf("Failed to load TLS certificates: %v", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	serverAddr := viper.GetString("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		var err error
		if tlsConfig != nil {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()
	logger.Infof("Server started on %s", serverAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
