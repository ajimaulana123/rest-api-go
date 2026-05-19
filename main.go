package main

import (
	"log"
	"os"

	"be/config"
	"be/database"
	"be/routes"
	"be/supabase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	sbClient, err := supabase.NewClient(cfg.Supabase)
	if err != nil {
		log.Fatalf("supabase client failed: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	routes.Setup(r, db, cfg, sbClient)

	addr := ":" + cfg.Port
	log.Printf("server running on http://localhost%s (supabase: %s, db: %s-%s pooler)",
		addr, cfg.Supabase.ProjectRef(), cfg.Supabase.Pooler, cfg.Supabase.Region)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
