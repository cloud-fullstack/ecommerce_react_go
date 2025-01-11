package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"shopa/db"
	"shopa/handler"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv" // Import godotenv for .env file support
)

var staticDir string // Define staticDir at the package level

func main() {
	// Set up logging
	log.SetHandler(cli.Default)

	// Load environment variables from .env file (if it exists)
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development" // Default to development
	}

	// Load the appropriate .env file
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Warnf("Error loading .env.%s file: %v", env, err)
	}

// Set default values and initialize router
port := getEnv("API_PORT", "8080") // Default port is 8080 if API_PORT is not set
staticDir := getEnv("STATIC_PATH", "/app/dist") // Default static files directory
frontendURL := getEnv("REACT_APP_DOMAIN_NAME", "https://rezav.gitlab.io/rezaverse") // Use REACT_APP_DOMAIN_NAME for CORS

// Log the environment variables for debugging
log.Infof("API_PORT: %s", port)
log.Infof("STATIC_PATH: %s", staticDir)
log.Infof("REACT_APP_DOMAIN_NAME: %s", frontendURL)

	// Initialize Gin and CORS configuration
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Recovery())

	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{frontendURL} // Use frontendURL for CORS
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Add routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Serve Static Files (Frontend)
	fmt.Println("Serving static content from", staticDir)
	r.Use(static.Serve("/", static.LocalFile(staticDir, true)))

	// Serve from index.html if route not found
	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", staticDir))
	})

	// Database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	dbPool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer dbPool.Close()

	// Set up logging to PostgreSQL
	log.SetHandler(postGresLogHandler(dbPool))

	// S3 configuration
	spacesKey := os.Getenv("SPACES_ACCESS_KEY")
	spacesSecret := os.Getenv("SPACES_SECRET_KEY")
	if spacesKey == "" || spacesSecret == "" {
		log.Fatal("SPACES_ACCESS_KEY or SPACES_SECRET_KEY is not set")
	}
	log.Infof("SPACES_ACCESS_KEY: %s", spacesKey)
	log.Infof("SPACES_SECRET_KEY: %s", spacesSecret)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:         aws.String("https://nyc3.digitaloceanspaces.com"),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatalf("Failed to create S3 session: %v", err)
	}
	s3Client := s3.New(newSession)

	// API Routes
	apiRouter := r.Group("/api")
	apiRouter.Use(db.MiddlewareDB(dbPool))
	{
		apiRouter.GET("/cors", handler.CorsAnywhere)
		apiRouter.GET("/avatar-blogs/:avatarKey", handler.BlogsByAvatarKey)
		apiRouter.GET("/most-loved-recent-blogs/", handler.MostLovedRecentBlogs)
		apiRouter.GET("/frontpage-product-previews/", handler.FrontpageProductPreviews)
		apiRouter.GET("/discounted-products-frontpage/", handler.DiscountedProductsFrontpages)
		apiRouter.GET("/store-details/:storeID", handler.StoreDetails)
		apiRouter.POST("/get-avatar-product/", handler.GetAvatarProduct)
		apiRouter.POST("/get-product-inventory-items/", handler.GetProductInventoryItems)
		apiRouter.POST("/gen-token/", handler.GenToken)
		apiRouter.POST("/profile-picture/", handler.GrabProfilePicture)

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

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Errorf("api-server ended unexpectedly: %s", err)
	}
}

// Helper function to retrieve environment variables with a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// serveSPA serves the Single Page Application for frontend routes
func serveSPA(c *gin.Context, urlPrefix string, fs static.ServeFileSystem) {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	if fs.Exists(urlPrefix, c.Request.URL.Path) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

// postGresLogHandler logs entries to the PostgreSQL database
func postGresLogHandler(conn *pgxpool.Pool) log.HandlerFunc {
	return func(entry *log.Entry) error {
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

		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		err := encoder.Encode(fields)
		if err != nil {
			fmt.Println("dbLogHandler: error encoding json fields", err)
			return nil
		}
		jsonMsg := buf.String()

		_, err = conn.Exec(context.Background(), "INSERT INTO log VALUES(DEFAULT, to_timestamp($1), 'api_server', $2, $3)",
			time.Now().Unix(), level, jsonMsg)
		if err != nil {
			fmt.Println("DB Logging error:", err)
			return nil
		}

		return nil
	}
}

// MiddlewareS3 attaches the S3 client to the request context
func MiddlewareS3(s3Client *s3.S3) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("s3Client", s3Client)
		c.Next() // call the next handler
	}
}