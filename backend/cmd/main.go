package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/pgsql_repo"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/services"
	http_transport "github.com/Ilyusha206420/first-day-of-internship/backend/internal/transport/http"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	PGUser        string `env:"POSTGRES_USER,required"`
	PGPass        string `env:"POSTGRES_PASSWORD,required"`
	AdminUsername string `env:"ADMIN_USERNAME,required"`
	AdminPassword string `env:"ADMIN_PASSWORD,required"`
	JWTSecret     string `env:"JWT_SECRET,required"`
	DBG           bool   `env:"DEBUG"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg Config
	if err = env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	db, err := sqlx.Connect("pgx", "user="+cfg.PGUser+" password="+cfg.PGPass+" dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ur := pgsql_repo.NewUserRepository(db)

	admin, err := ur.FindByName(context.Background(), cfg.AdminUsername)
	if err != nil {
		log.Fatal(err)
	}

	if admin == nil {
		admin, err = domain.NewUser(cfg.AdminUsername, cfg.AdminPassword)
		if err != nil {
			log.Fatal(err)
		}
		err = ur.Add(context.Background(), admin)
		if err != nil {
			log.Fatal(err)
		}
	}

	as := services.NewAuthService(ur)

	tr := pgsql_repo.NewTokenRepo(db)

	ss := services.NewSessionService(ur, tr, cfg.JWTSecret, 15*time.Minute, 30*24*time.Hour)

	ah := http_transport.NewAuthHandler(as, ss, 15*time.Minute, 30*24*time.Hour, !cfg.DBG)

	r := gin.Default()

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", ah.Login)
		authGroup.POST("/refresh", ah.Refresh)
		authGroup.POST("/logout", ah.Logout)
	}

	protected := r.Group("/api")
	protected.Use(http_transport.AuthMiddleware(ss))
	{
		protected.GET("/me", func(c *gin.Context) {
			userID := c.MustGet("userPublicID")
			c.JSON(http.StatusOK, gin.H{
				"public_id": userID,
			})
		})
	}

	r.Run()
}
