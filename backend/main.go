package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"shopa/db"
	"shopa/handler"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Set up logging

	// Set default values and initialize router
	port := getEnv("API_PORT", "8080")                                           // Default port is 8080 if API_PORT is not set
	staticDir := getEnv("STATIC_PATH", "/app/dist")                              // Default static files directory
	frontendURL := getEnv("REACT_APP_DOMAIN_NAME", "https://rezav.onrender.com") // Use REACT_APP_DOMAIN_NAME for CORS

	// Initialize Gin and CORS configuration
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Recovery())

	// Configure CORS
	configureCORS(r, frontendURL)

	// Serve Static Files (Frontend)
	log.Infof("Serving static content from %s", staticDir)
	r.Use(static.Serve("/", static.LocalFile(staticDir, true)))

	// Serve from index.html if route not found
	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", staticDir))
	})

	// Initialize database connection
	dbPool, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Set up logging to PostgreSQL
	log.SetHandler(postgresLogHandler(dbPool))

	// Initialize S3 client
	s3Client, err := initS3()
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Set up API routes
	setupAPIRoutes(r, dbPool, s3Client)

	// Start the server
	log.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Errorf("Server ended unexpectedly: %v", err)
	}
}

// getEnv retrieves environment variables with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// configureCORS sets up CORS middleware
func configureCORS(r *gin.Engine, frontendURL string) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{frontendURL, "https://rezav.onrender.com"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowCredentials = false // Enable if your app sends credentials
	corsConfig.MaxAge = 12 * time.Hour  // Cache preflight requests for 12 hours
	r.Use(cors.New(corsConfig))
}

// initDB initializes the database connection
func initDB() (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""),
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", ""),
		getEnv("DB_SSLMODE", "require"),
	)

	return pgxpool.Connect(context.Background(), connString)
}

// initS3 initializes the S3 client
func initS3() (*s3.S3, error) {
	spacesKey := getEnv("SPACES_ACCESS_KEY", "")
	spacesSecret := getEnv("SPACES_SECRET_KEY", "")
	if spacesKey == "" || spacesSecret == "" {
		return nil, fmt.Errorf("SPACES_ACCESS_KEY or SPACES_SECRET_KEY is not set")
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:         aws.String("https://nyc3.digitaloceanspaces.com"),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %v", err)
	}
	return s3.New(newSession), nil
}

// setupAPIRoutes sets up all API routes
func setupAPIRoutes(r *gin.Engine, dbPool *pgxpool.Pool, s3Client *s3.S3) {
	apiRouter := r.Group("/api")
	apiRouter.Use(db.MiddlewareDB(dbPool))
	{
		apiRouter.GET("/frontpage-product-previews/", handler.FrontpageProductPreviews)
		apiRouter.GET("/store-details/:storeID", handler.StoreDetails)
		apiRouter.GET("/most-loved-recent-blogs/", handler.MostLovedRecentBlogs)
		apiRouter.GET("/discounted-products-frontpage/", handler.DiscountedProductsFrontpages)

		apiRouter.POST("/get-avatar-product/", handler.GetAvatarProduct)
		apiRouter.POST("/get-product-inventory-items/", handler.GetProductInventoryItems)
		apiRouter.POST("/gen-token/", handler.GenToken)
		apiRouter.POST("/profile-picture/", handler.GrabProfilePicture)
		apiRouter.GET("/cors", handler.CorsAnywhere)
		apiRouter.GET("/avatar-blogs/:avatarKey", handler.BlogsByAvatarKey)

		slOnlyRouter := apiRouter.Group("")
		{
			slOnlyRouter.POST("/cancel-order-sl/", handler.CancelOrder)
			slOnlyRouter.GET("/deliver-dropbox/", handler.DeliverDropbox)
			slOnlyRouter.POST("/complete-order/", handler.CompleteOrder)
			slOnlyRouter.POST("/insert-dropbox-repo/", handler.InsertDropboxRepo)
			slOnlyRouter.POST("/insert-dropbox/", handler.InsertDropbox)
			slOnlyRouter.POST("/update-dropbox-contents/", handler.UpdateDropboxContents)
		}
	}

	// Authenticated API Routes
	apiRouterAuth := apiRouter.Group("")
	apiRouterAuth.Use(handler.MiddlewareAuth())
	apiRouterAuth.Use(db.MiddlewareDB(dbPool))
	{
		apiRouterAuth.Use(MiddlewareS3(s3Client)).POST("/upload-picture/:pictureType", handler.UploadPicture)
		apiRouterAuth.DELETE("/answer/:answerID", handler.DeleteAnswer)
		apiRouterAuth.POST("/answer/", handler.InsertAnswer)
		apiRouterAuth.DELETE("/question/:questionID", handler.DeleteQuestion)
		apiRouterAuth.POST("/question/", handler.InsertQuestion)
		apiRouterAuth.DELETE("/blog-post-love/:postID", handler.DeleteBlogPostLove)
		apiRouterAuth.POST("/blog-post-love/:postID", handler.InsertBlogPostLove)
		apiRouterAuth.DELETE("/blog-post/:id", handler.DeleteBlogPost)
		apiRouterAuth.POST("/blog-post/", handler.InsertBlogPost)
		apiRouterAuth.DELETE("/product-picture-love/:pictureID", handler.DeleteProductPictureLove)
		apiRouterAuth.POST("/product-picture-love/:pictureID", handler.InsertProductPictureLove)
		apiRouterAuth.POST("/delete-product/:productID", handler.DeleteProduct)
		apiRouterAuth.POST("/delete-store/:storeID", handler.DeleteStore)
		apiRouterAuth.POST("/redeliver-order-product/", handler.RedeliverOrderProduct)
		apiRouterAuth.POST("/resend-order/", handler.ResendOrder)
		apiRouterAuth.POST("/cancel-order/", handler.CancelOrder)
		apiRouterAuth.POST("/get-notifications/", handler.GetNotification)
		apiRouterAuth.POST("/order-history/", handler.OrderHistory)
		apiRouterAuth.POST("/create-order/", handler.CreateOrder)
		apiRouterAuth.POST("/get-avatar-products/", handler.GetAvatarProducts)
		apiRouterAuth.POST("/get-avatar-stores/", handler.GetAvatarStores)
		apiRouterAuth.POST("/insert-store/", handler.InsertStore)
		apiRouterAuth.POST("/insert-product/", handler.InsertProduct)
		apiRouterAuth.POST("/get-avatar-dropbox-contents/", handler.GetAvatarDropboxContents)
		apiRouterAuth.POST("/get-avatar-dropboxes/", handler.GetAvatarDropboxes)
		apiRouterAuth.POST("/update-hud-heartbeat/", handler.UpdateHeartbeatURLHUD)
		apiRouterAuth.POST("/heartbeat-hud/", handler.PingHUD)
	}
}

// postgresLogHandler logs entries to the PostgreSQL database
func postgresLogHandler(conn *pgxpool.Pool) log.Handler {
	return log.HandlerFunc(func(entry *log.Entry) error {
		level := ""
		switch entry.Level {
		case log.InfoLevel:
			level = "inf"
		case log.ErrorLevel:
			level = "err"
		case log.WarnLevel:
			level = "war"
		}

		fields := make(map[string]interface{})
		for key, value := range entry.Fields {
			fields[key] = value
		}
		fields["message"] = entry.Message

		jsonMsg, err := json.Marshal(fields)
		if err != nil {
			log.Errorf("Error encoding JSON fields: %v", err)
			return nil
		}

		_, err = conn.Exec(context.Background(), "INSERT INTO log VALUES(DEFAULT, to_timestamp($1), 'api_server', $2, $3)",
			time.Now().Unix(), level, jsonMsg)
		if err != nil {
			log.Errorf("DB Logging error: %v", err)
			return nil
		}

		return nil
	})
}

// MiddlewareS3 attaches the S3 client to the request context
func MiddlewareS3(s3Client *s3.S3) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("s3Client", s3Client)
		c.Next()
	}
}
