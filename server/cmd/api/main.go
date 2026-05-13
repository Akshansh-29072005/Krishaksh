package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	internalAI "github.com/aarcsx/krisho-backend/internal/ai"
	aiProviders "github.com/aarcsx/krisho-backend/internal/ai/providers"
	"github.com/aarcsx/krisho-backend/internal/config"
	"github.com/aarcsx/krisho-backend/internal/core/middleware"
	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/observability"

	adminHandler "github.com/aarcsx/krisho-backend/internal/modules/admin/handler"
	adminRepo "github.com/aarcsx/krisho-backend/internal/modules/admin/repository"
	adminRoutes "github.com/aarcsx/krisho-backend/internal/modules/admin/routes"
	adminService "github.com/aarcsx/krisho-backend/internal/modules/admin/service"
	adHandler "github.com/aarcsx/krisho-backend/internal/modules/advertisements/handler"
	adRepo "github.com/aarcsx/krisho-backend/internal/modules/advertisements/repository"
	adRoutes "github.com/aarcsx/krisho-backend/internal/modules/advertisements/routes"
	analyticsHandler "github.com/aarcsx/krisho-backend/internal/modules/analytics/handler"
	analyticsRepo "github.com/aarcsx/krisho-backend/internal/modules/analytics/repository"
	analyticsRoutes "github.com/aarcsx/krisho-backend/internal/modules/analytics/routes"
	analyticsService "github.com/aarcsx/krisho-backend/internal/modules/analytics/service"
	appHandler "github.com/aarcsx/krisho-backend/internal/modules/app/handler"
	appRepo "github.com/aarcsx/krisho-backend/internal/modules/app/repository"
	appRoutes "github.com/aarcsx/krisho-backend/internal/modules/app/routes"
	appService "github.com/aarcsx/krisho-backend/internal/modules/app/service"
	authHandler "github.com/aarcsx/krisho-backend/internal/modules/auth/handler"
	authRepo "github.com/aarcsx/krisho-backend/internal/modules/auth/repository"
	authRoutes "github.com/aarcsx/krisho-backend/internal/modules/auth/routes"
	authService "github.com/aarcsx/krisho-backend/internal/modules/auth/service"
	cartHandler "github.com/aarcsx/krisho-backend/internal/modules/cart/handler"
	cartRepo "github.com/aarcsx/krisho-backend/internal/modules/cart/repository"
	cartRoutes "github.com/aarcsx/krisho-backend/internal/modules/cart/routes"
	cartService "github.com/aarcsx/krisho-backend/internal/modules/cart/service"
	companyHandler "github.com/aarcsx/krisho-backend/internal/modules/company/handler"
	companyRoutes "github.com/aarcsx/krisho-backend/internal/modules/company/routes"
	companyService "github.com/aarcsx/krisho-backend/internal/modules/company/service"
	diseaseHandler "github.com/aarcsx/krisho-backend/internal/modules/diseases/handler"
	diseaseRepo "github.com/aarcsx/krisho-backend/internal/modules/diseases/repository"
	diseaseRoutes "github.com/aarcsx/krisho-backend/internal/modules/diseases/routes"
	orderHandler "github.com/aarcsx/krisho-backend/internal/modules/orders/handler"
	orderRepo "github.com/aarcsx/krisho-backend/internal/modules/orders/repository"
	orderRoutes "github.com/aarcsx/krisho-backend/internal/modules/orders/routes"
	orderService "github.com/aarcsx/krisho-backend/internal/modules/orders/service"
	paymentHandler "github.com/aarcsx/krisho-backend/internal/modules/payments/handler"
	razorpayClient "github.com/aarcsx/krisho-backend/internal/modules/payments/razorpay"
	paymentRepo "github.com/aarcsx/krisho-backend/internal/modules/payments/repository"
	paymentRoutes "github.com/aarcsx/krisho-backend/internal/modules/payments/routes"
	paymentService "github.com/aarcsx/krisho-backend/internal/modules/payments/service"
	productHandler "github.com/aarcsx/krisho-backend/internal/modules/products/handler"
	productRepo "github.com/aarcsx/krisho-backend/internal/modules/products/repository"
	productRoutes "github.com/aarcsx/krisho-backend/internal/modules/products/routes"
	recHandler "github.com/aarcsx/krisho-backend/internal/modules/recommendations/handler"
	recRoutes "github.com/aarcsx/krisho-backend/internal/modules/recommendations/routes"
	recService "github.com/aarcsx/krisho-backend/internal/modules/recommendations/service"
	scanHandler "github.com/aarcsx/krisho-backend/internal/modules/scans/handler"
	scanRepo "github.com/aarcsx/krisho-backend/internal/modules/scans/repository"
	scanRoutes "github.com/aarcsx/krisho-backend/internal/modules/scans/routes"
	scanService "github.com/aarcsx/krisho-backend/internal/modules/scans/service"
	userHandler "github.com/aarcsx/krisho-backend/internal/modules/users/handler"
	userRoutes "github.com/aarcsx/krisho-backend/internal/modules/users/routes"
	workers "github.com/aarcsx/krisho-backend/internal/workers"
	queue "github.com/aarcsx/krisho-backend/pkg/queue"
	"github.com/aarcsx/krisho-backend/pkg/s3"
)

func main() {
	observability.InitLogger()
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	if os.Getenv("RUN_DB_MIGRATIONS") == "true" {
		if err := database.RunMigrations(context.Background(), db, "db/migrations"); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	qClient := queue.NewQueueClient(cfg.RedisAddr)
	defer qClient.Close()

	r := gin.New()
	r.Use(middleware.RecoveryJSON(), middleware.CORS(), middleware.RequestResponseLogger(), middleware.RequestContext(), middleware.RequestTiming(), middleware.ErrorTracker(), middleware.ValidationGuard(), middleware.UploadSizeLimit(10<<20), middleware.RateLimit(240), middleware.AbusePrevention())

	r.GET("/api/v1/weather", func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")
		if lat == "" || lon == "" {
			c.JSON(400, gin.H{"status": "error", "message": "lat and lon are required"})
			return
		}

		// 1. Get real weather from Open-Meteo (free, no API key)
		weatherURL := fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&timezone=auto",
			lat, lon,
		)
		weatherResp, err := http.Get(weatherURL)
		if err != nil {
			c.JSON(200, gin.H{"status": "success", "data": gin.H{
				"temperature": "--°C", "condition": "Unknown", "humidity": "--%",
				"wind_speed": "-- km/h", "location_name": "Unknown Location",
			}})
			return
		}
		defer weatherResp.Body.Close()

		var weatherData struct {
			Current struct {
				Temperature float64 `json:"temperature_2m"`
				Humidity    float64 `json:"relative_humidity_2m"`
				WeatherCode int     `json:"weather_code"`
				WindSpeed   float64 `json:"wind_speed_10m"`
			} `json:"current"`
		}
		if err := json.NewDecoder(weatherResp.Body).Decode(&weatherData); err != nil {
			c.JSON(200, gin.H{"status": "success", "data": gin.H{
				"temperature": "--°C", "condition": "Unknown", "humidity": "--%",
				"wind_speed": "-- km/h", "location_name": "Unknown Location",
			}})
			return
		}

		// 2. Get location name from Open-Meteo geocoding (reverse)
		geoURL := fmt.Sprintf(
			"https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json&accept-language=en",
			lat, lon,
		)
		req, _ := http.NewRequest("GET", geoURL, nil)
		req.Header.Set("User-Agent", "Krisho-App/1.0")
		geoResp, err := http.DefaultClient.Do(req)
		locationName := "Your Location"
		if err == nil {
			defer geoResp.Body.Close()
			var geoData struct {
				Address struct {
					City    string `json:"city"`
					Town    string `json:"town"`
					Village string `json:"village"`
					County  string `json:"county"`
					State   string `json:"state"`
					Country string `json:"country"`
				} `json:"address"`
			}
			if err := json.NewDecoder(geoResp.Body).Decode(&geoData); err == nil {
				place := geoData.Address.City
				if place == "" {
					place = geoData.Address.Town
				}
				if place == "" {
					place = geoData.Address.Village
				}
				if place == "" {
					place = geoData.Address.County
				}
				region := geoData.Address.State
				if place != "" && region != "" {
					locationName = place + ", " + region
				} else if place != "" {
					locationName = place
				} else if region != "" {
					locationName = region
				}
			}
		}

		// 3. Map WMO weather code to condition string
		condition := "Clear"
		wc := weatherData.Current.WeatherCode
		switch {
		case wc == 0:
			condition = "Clear"
		case wc <= 3:
			condition = "Cloudy"
		case wc >= 45 && wc <= 48:
			condition = "Foggy"
		case wc >= 51 && wc <= 57:
			condition = "Drizzle"
		case wc >= 61 && wc <= 67:
			condition = "Rain"
		case wc >= 71 && wc <= 77:
			condition = "Snow"
		case wc >= 80 && wc <= 82:
			condition = "Showers"
		case wc >= 95:
			condition = "Thunderstorm"
		}

		c.JSON(200, gin.H{
			"status": "success",
			"data": gin.H{
				"temperature":   fmt.Sprintf("%.0f°C", weatherData.Current.Temperature),
				"condition":     condition,
				"humidity":      fmt.Sprintf("%.0f%%", weatherData.Current.Humidity),
				"wind_speed":    fmt.Sprintf("%.0f km/h", weatherData.Current.WindSpeed),
				"location_name": locationName,
			},
		})
	})

	authRepository := authRepo.NewAuthRepository(db)
	authSvc := authService.NewAuthService(authRepository)
	authHdlr := authHandler.NewAuthHandler(authSvc)
	userHdlr := userHandler.NewUserHandler(authRepository)

	s3Client, _ := s3.NewS3Client("ap-south-1", os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "krisho-scans"
	}

	scanRepository := scanRepo.NewScanRepository(db)
	scanSvc := scanService.NewScanService(scanRepository, s3Client, qClient, s3Bucket)
	scanHdlr := scanHandler.NewScanHandler(scanSvc)
	diseaseRepository := diseaseRepo.NewDiseaseRepository(db)
	diseaseHdlr := diseaseHandler.NewDiseaseHandler(diseaseRepository)
	productRepository := productRepo.NewProductRepository(db)
	productHdlr := productHandler.NewProductHandler(productRepository)
	adRepository := adRepo.NewAdvertisementRepository(db)
	adHdlr := adHandler.NewAdvertisementHandler(adRepository)
	recSvc := recService.NewRecommendationService(scanRepository, diseaseRepository, productRepository)
	recHdlr := recHandler.NewRecommendationHandler(recSvc)
	appRepository := appRepo.NewAppRepository(db)
	appSvc := appService.NewAppService(appRepository)
	appHdlr := appHandler.NewAppHandler(appSvc)

	cartRepository := cartRepo.NewCartRepository(db)
	cartSvc := cartService.NewCartService(cartRepository, productRepository)
	cartHdlr := cartHandler.NewCartHandler(cartSvc)
	orderRepository := orderRepo.NewOrderRepository(db)
	orderSvc := orderService.NewOrderService(orderRepository, cartRepository, productRepository)
	orderHdlr := orderHandler.NewOrderHandler(orderSvc)
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	rz := razorpayClient.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))
	paymentSvc := paymentService.NewPaymentService(orderRepository, paymentRepository, rz, qClient, os.Getenv("RAZORPAY_WEBHOOK_SECRET"), os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))
	paymentHdlr := paymentHandler.NewPaymentHandler(paymentSvc)

	analyticsRepository := analyticsRepo.NewAnalyticsRepository(db)
	analyticsSvc := analyticsService.NewAnalyticsService(analyticsRepository)
	analyticsH := analyticsHandler.NewAnalyticsHandler(analyticsSvc)
	adminSvc := adminService.NewAdminService(adminRepo.NewAdminRepository(db))
	adminH := adminHandler.NewAdminHandler(adminSvc)
	companySvc := companyService.NewCompanyService(analyticsRepository)
	companyH := companyHandler.NewCompanyHandler(companySvc)

	promptMgr := internalAI.NewPromptManager("internal/ai/prompts")
	aiMgr := internalAI.NewManager(
		[]aiProviders.VisionProvider{
			aiProviders.NewGeminiProvider(strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))),
		},
		promptMgr,
		"vision_v1",
		0.65,
	)
	aiSvc := internalAI.NewService(aiMgr)
	if os.Getenv("ENABLE_EMBEDDED_WORKER") == "true" {
		go func() {
			srv := asynq.NewServer(asynq.RedisClientOpt{Addr: cfg.RedisAddr}, asynq.Config{Concurrency: 20})
			mux := asynq.NewServeMux()
			mux.HandleFunc(queue.TypeScanAnalyze, workers.NewScanWorker(scanRepository, aiSvc).HandleScanAnalyzeTask)
			mux.HandleFunc(queue.TypePaymentEvent, workers.NewPaymentWorker(paymentRepository, orderRepository).HandlePaymentEventTask)
			mux.HandleFunc(queue.TypeRefundStart, workers.NewPaymentWorker(paymentRepository, orderRepository).HandleRefundPlaceholderTask)
			mux.HandleFunc(queue.TypeAnalyticsEvt, workers.NewAnalyticsWorker(analyticsRepository).HandleAnalyticsEventTask)
			if err := srv.Run(mux); err != nil {
				log.Fatalf("Asynq server failed: %v", err)
			}
		}()
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
	})
	r.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "db": false})
			return
		}
		if err := rdb.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "redis": false})
			return
		}
		observability.M.Inc("queue_ready_checks_total")
		c.JSON(http.StatusOK, gin.H{"ready": true, "db": true, "redis": true})
	})
	r.GET("/metrics", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; version=0.0.4")
		c.String(http.StatusOK, observability.M.PrometheusText())
	})

	v1 := r.Group("/api/v1")
	{
		authRoutes.RegisterAuthRoutes(v1, authHdlr)
		userRoutes.RegisterUserRoutes(v1, userHdlr)
		appRoutes.RegisterAppRoutes(v1, appHdlr)
		scanRoutes.RegisterScanRoutes(v1, scanHdlr)
		diseaseRoutes.RegisterDiseaseRoutes(v1, diseaseHdlr)
		productRoutes.RegisterProductRoutes(v1, productHdlr)
		recRoutes.RegisterRecommendationRoutes(v1, recHdlr)
		adRoutes.RegisterAdvertisementRoutes(v1, adHdlr)
		cartRoutes.RegisterCartRoutes(v1, cartHdlr)
		orderRoutes.RegisterOrderRoutes(v1, orderHdlr)
		paymentRoutes.RegisterPaymentRoutes(v1, paymentHdlr)
		analyticsRoutes.RegisterAnalyticsRoutes(v1, analyticsH)
		adminRoutes.RegisterAdminRoutes(v1, adminH)
		companyRoutes.RegisterCompanyRoutes(v1, companyH)
	}

	log.Printf("Starting Krisho API Server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
