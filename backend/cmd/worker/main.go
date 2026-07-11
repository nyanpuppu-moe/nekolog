package main

import (
	"nekolog/internal/config"
	"nekolog/internal/database"
	"nekolog/internal/handler"
	"nekolog/internal/log"
	"nekolog/internal/middleware"
	"nekolog/internal/model"
	"nekolog/internal/repository"
	"nekolog/internal/service"
	"nekolog/internal/storage"
	"nekolog/internal/web"
)

func main() {
	// 사용자의 설정을 불러옵니다
	// configs/config.yaml 을 불러옵니다
	config, err := config.Load()
	if err != nil {
		log.Error("Faild to load config: %v", err)
		return
	}

	// DataBase 에 연결을 시도합니다
	// config.Database.Path 에서 db 파일을 찾습니다
	db, err := database.Connect(config.Database.Path)
	if err != nil {
		log.Error("Faild to can not connect sqlite: %v", err)
	}

	// DataBase 를 자동으로 마이그레이션 합니다
	// 서비스를 시작할떄마다 항상 마이그레이션 합니다
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Error("Faild to sqlite migration: %v", err)
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

	router := web.NewRouter()

	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/users/:username", userHandler.Get)
		apiRoutes.GET("/articles/:username/:title", articleHandler.Get)
		apiRoutes.GET("/assets/:id", assetHandler.Get)
	}

	authRoutes := apiRoutes.Group("/auth")
	{
		authRoutes.POST("/register", sessionHandler.Register)
		authRoutes.POST("/login", sessionHandler.Login)
		authRoutes.POST("/logout", sessionHandler.Logout)
	}

	protectedRoutes := apiRoutes.Group("/protected")
	protectedRoutes.Use(middleware.AuthRequired)
	{
		protectedRoutes.POST("/article", articleHandler.Post)
		protectedRoutes.PATCH("/article/:username/:title", articleHandler.Patch)
		protectedRoutes.DELETE("/article/:username/:title", articleHandler.Delete)

		protectedRoutes.POST("/asset", assetHandler.Post)
	}

	router.Serve(config.Server.Port)
}
