package main

import (
	"log"

	"nekolog/internal/config"
	"nekolog/internal/database"
	"nekolog/internal/handler"
	"nekolog/internal/middleware"
	"nekolog/internal/model"
	"nekolog/internal/repository"
	"nekolog/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := "runtime/local.db"
	db, err := database.Connect(dbPath)
	if err != nil {
		log.Fatalf("faild to can not connect sqlite: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("faild to sqlite migration: %v", err)
	}

	articleRepository := repository.NewArticleRepository(db)
	assetRepository := repository.NewAssetRepository(db)
	userRepository := repository.NewUserRepository(db)

	sessionService := service.NewSessionService(userRepository)
	sessionHandler := handler.NewSessionHandler(sessionService)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	articleService := service.NewArticleService(
		articleRepository,
		userRepository,
	)
	articleHandler := handler.NewArticleHandler(articleService)

	assetService := service.NewAssetService(
		assetRepository,
		userRepository,
	)
	assetHandler := handler.NewAssetHandler(assetService)

	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	var store sessions.Store
	if config.Environment.Mode == "development" {
		store = service.InitDevelopmentSessionStore(config.Server.SessionStore.PrivateKey)
	} else {
		store = service.InitSessionStore(config.Server.SessionStore.PrivateKey)
	}

	router.Use(sessions.Sessions("NEKO_LOG_SESSION", store))

	router.GET("/users/:username", userHandler.Get)
	router.GET("/articles/:username/:title", articleHandler.Get)
	router.GET("/assets/:id", assetHandler.Get)

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", sessionHandler.Register)
		authRoutes.POST("/login", sessionHandler.Login)
		authRoutes.POST("/logout", sessionHandler.Logout)
	}

	protectedRoutes := router.Group("/api/protected")
	protectedRoutes.Use(middleware.AuthRequired())
	{
		protectedRoutes.POST("/article", articleHandler.Post)
		protectedRoutes.PATCH("/article/:username/:title", articleHandler.Patch)
		protectedRoutes.DELETE("/article/:username/:title", articleHandler.Delete)

		protectedRoutes.POST("/asset", assetHandler.Post)
	}

	router.Run(":" + config.Server.Port)
}
