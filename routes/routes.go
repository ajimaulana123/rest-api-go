package routes

import (
	"be/config"
	"be/handlers"
	"be/middleware"

	"github.com/gin-gonic/gin"
	supago "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, cfg config.Config, sb *supago.Client) {
	authHandler := &handlers.AuthHandler{DB: db, Config: cfg}
	itemHandler := &handlers.ItemHandler{DB: db}
	_ = sb

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"supabase": gin.H{
				"project_id": cfg.Supabase.ProjectRef(),
				"url":        cfg.Supabase.URL,
			},
		})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			protected.GET("/profile", authHandler.Profile)

			items := protected.Group("/items")
			{
				items.GET("", itemHandler.List)
				items.GET("/:id", itemHandler.Get)
				items.POST("", itemHandler.Create)
				items.PUT("/:id", itemHandler.Update)
				items.DELETE("/:id", itemHandler.Delete)
			}
		}
	}
}
