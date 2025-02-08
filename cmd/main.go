package main

import (
	"log"

	"github.com/amiulam/music-catalog/internal/configs"
	membershipHandler "github.com/amiulam/music-catalog/internal/handler/memberships"
	"github.com/amiulam/music-catalog/internal/models/memberships"
	membershipRepo "github.com/amiulam/music-catalog/internal/repository/memberships"
	membershipSvc "github.com/amiulam/music-catalog/internal/services/memberships"
	"github.com/amiulam/music-catalog/pkg/internalsql"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var (
		cfg *configs.Config
	)

	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	if err != nil {
		log.Fatal("Gagal inisialisasi config", err)
	}

	cfg = configs.Get()

	db, err := internalsql.Connect(cfg.Database.DatabaseSourceName)

	db.AutoMigrate(&memberships.User{})

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Repository
	membershipRepo := membershipRepo.NewRepository(db)

	// Services
	membershipSvc := membershipSvc.NewService(cfg, membershipRepo)

	// Handlers
	membershipHandler := membershipHandler.NewHandler(r, membershipSvc)

	// Routes
	membershipHandler.RegisterRoute()

	if err != nil {
		log.Fatalf("fail to connect to database, err: %+v", err)
	}

	r.Run(cfg.Service.Port)
}
