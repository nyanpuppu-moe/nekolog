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
	"nekolog/internal/storage"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	dbPath := config.Database.Path
	db, err := database.Connect(dbPath)
	if err != nil {
		log.Fatalf("Faild to can not connect sqlite: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Faild to sqlite migration: %v", err)
	}

	assetsStorage := storage.NewAssetStorage(config.Storage.AssetsPath)
	contentStorage := storage.NewContentStorage(config.Storage.ContentsPath)

	articleRepository := repository.NewArticleRepository(db, contentStorage)
	assetRepository := repository.NewAssetRepository(db, assetsStorage)
	userRepository := repository.NewUserRepository(db)

	sessionService := service.NewSessionService(userRepository)
	sessionHandler := handler.NewSessionHandler(sessionService)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	articleService := service.NewArticleService(
		articleRepository,
		userRepository,
	)
	articleHandler := handler.NewArticleHandler(articleService, userService)

	assetService := service.NewAssetService(
		assetRepository,
		userRepository,
	)
	assetHandler := handler.NewAssetHandler(assetService)

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
