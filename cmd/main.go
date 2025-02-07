package main

import (
	"log"

	"github.com/amiulam/music-catalog/internal/configs"
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

	if err != nil {
		log.Fatalf("fail to connect to database, err: %+v", err)
	}

	r.Run(cfg.Service.Port)
}
