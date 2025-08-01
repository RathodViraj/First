package main

import (
	cachingservice "First/cachingservice"
	"First/db"
	"First/handler"
	"First/middleware"
	"First/notification"
	"First/repository"
	"First/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	dbConn, err := db.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	RDB := db.InitRedis()
	defer RDB.Close()
	cachingservice.SetRedies(RDB)

	db.InitNeo4j()
	defer db.Driver.Close(context.Background())

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err = session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(context.Background(), `RETURN "Connected Successfully!"`, nil)
		if err != nil {
			return nil, err
		}
		if res.Next(context.Background()) {
			return res.Record().Values[0], nil
		}
		return nil, res.Err()
	})

	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	postRepo := repository.NewPostRepository(dbConn)
	commentRepo := repository.NewCommentRepo(dbConn)
	graphRepo := repository.NewGraph(db.Driver)
	connectionRepo := repository.NewConnectionRepo(dbConn, graphRepo)
	likeRepo := repository.NewLikeRepo(dbConn, graphRepo, userRepo)
	searchRepo := repository.NewSearchRepo(dbConn)
	notificationRepo := repository.NewNotificationRepo(dbConn)

	hub := notification.NewHub()
	go hub.Run()

	userService := service.NewUserService(userRepo, graphRepo)
	postService := service.NewPostService(postRepo, graphRepo)
	commentsService := service.NewCommentsService(commentRepo, graphRepo)
	connectionService := service.NewConnectionService(connectionRepo, graphRepo)
	authService := service.NewAuthService(userService)
	likeService := service.NewLikeService(likeRepo)
	searchService := service.NewSearchService(searchRepo)
	notificationService := service.NewNotificationService(notificationRepo)

	userHandler := handler.NewUserHandler(userService, postService)
	postHandler := handler.NewPostHandler(postService, connectionService)
	commentsHandler := handler.NewCommentsHandler(commentsService, hub, graphRepo, notificationRepo)
	connectionHandler := handler.NewConnectionHandler(connectionService, hub, graphRepo, notificationRepo)
	authHandler := handler.NewAuthHandler(authService)
	likeHandler := handler.NewLikeHandler(likeService, hub, graphRepo, notificationRepo)
	searchHandler := handler.NewSearchHandler(searchService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Zap logger middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		zap.L().Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		)
	})
	router.Use(gin.Recovery())
	router.Use(middleware.GinMiddleware(authService))

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hey, this is home page!")
	})
	router.POST("/login", authHandler.Login)
	router.RouterGroup.POST("/register", authHandler.Register)

	// Home feed (recent posts from everyone)
	router.GET("/home", postHandler.RecentPosts)

	// Posts routes
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts/:id", postHandler.GetPost)
	router.DELETE("/posts/:id", postHandler.DeletePost)

	// like routes
	router.POST("/posts/:id/like", likeHandler.LikePost)
	router.DELETE("/posts/:id/like", likeHandler.UnlikePost)
	router.GET("/posts/:id/likes", likeHandler.GetLikes)

	router.GET("/posts/:id/comments", commentsHandler.GetAllComments)
	router.POST("/posts/:id/comments", commentsHandler.AddComment)

	// User-related routes
	users := router.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.GET("/:id/home", userHandler.GetFeed)
	}

	// Conection routes
	router.GET("/users/:id/followers", connectionHandler.GetFollowers)
	router.GET("/users/:id/followings", connectionHandler.GetFollowings)
	router.GET("/users/:id/mutual", connectionHandler.GetMutual)
	router.POST("/follow/:follower_id/:following_id", connectionHandler.FollowUser)
	router.DELETE("/unfollow/:follower_id/:following_id", connectionHandler.UnfollowUser)

	// Search routes
	router.GET("/search/users", searchHandler.SearchUser)
	router.GET("/search/posts", searchHandler.SearchPost)

	router.GET("/users/:id/notifications", notificationHandler.GetNotificationsByUser)

	router.GET("/ws", func(ctx *gin.Context) {
		handler.ServeWS(hub, ctx)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		zap.L().Info("Server starting", zap.String("addr", ":8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Server unexpectedly failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exited gracefully.")
}
