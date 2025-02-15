package main

import (
	"log"
	"net/http"

	"github.com/amiulam/music-catalog/internal/configs"
	membershipHandler "github.com/amiulam/music-catalog/internal/handler/memberships"
	trackHandler "github.com/amiulam/music-catalog/internal/handler/tracks"
	"github.com/amiulam/music-catalog/internal/models/memberships"
	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	membershipRepo "github.com/amiulam/music-catalog/internal/repository/memberships"
	spotifyRepo "github.com/amiulam/music-catalog/internal/repository/spotify"
	trackActivitiesRepo "github.com/amiulam/music-catalog/internal/repository/track_activities"
	membershipSvc "github.com/amiulam/music-catalog/internal/services/memberships"
	trackSvc "github.com/amiulam/music-catalog/internal/services/tracks"
	"github.com/amiulam/music-catalog/pkg/httpclient"
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
	db.AutoMigrate(&trackactivities.TrackActivity{})

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	httpClient := httpclient.NewClient(&http.Client{})

	// Repository
	membershipRepo := membershipRepo.NewRepository(db)
	spotifyOutbound := spotifyRepo.NewSpotifyOutbound(cfg, httpClient)
	trackActivityRepo := trackActivitiesRepo.NewRepository(db)

	// Services
	membershipSvc := membershipSvc.NewService(cfg, membershipRepo)
	trackSvc := trackSvc.NewService(spotifyOutbound, trackActivityRepo)

	// Handlers
	membershipHandler := membershipHandler.NewHandler(r, membershipSvc)
	trackHandler := trackHandler.NewHandler(r, trackSvc)

	// Routes
	membershipHandler.RegisterRoute()
	trackHandler.RegisterRoute()

	if err != nil {
		log.Fatalf("fail to connect to database, err: %+v", err)
	}

	r.Run(cfg.Service.Port)
}
