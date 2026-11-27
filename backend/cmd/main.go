package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/pgsql_repo"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/sqlite_repo"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/seed"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/services"
	http_transport "github.com/Ilyusha206420/first-day-of-internship/backend/internal/transport/http"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	AdminUsername string `env:"ADMIN_USERNAME,required"`
	AdminPassword string `env:"ADMIN_PASSWORD,required"`

	DBType string `env:"DB_TYPE,required"`

	DBG bool `env:"DEBUG"`

	JWTSecret string `env:"JWT_SECRET,required"`

	PGUser string `env:"POSTGRES_USER"`
	PGPass string `env:"POSTGRES_PASSWORD"`

	SQLitePath string `env:"SQLITE_PATH"`
}

type user_repository interface {
	services.AuthUserRepository
	services.SessionUserRepository
	Add(ctx context.Context, user *domain.User) error
}

type token_repository interface {
	services.SessionTokenRepository
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("error with loading environment: %v", err))
	}

	var cfg Config
	if err = env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Errorf("error with pasing environment %v", err))
	}

	var db *sqlx.DB
	var ur user_repository
	var tr token_repository
	var gr seed.GraphSeederRepository

	switch cfg.DBType {
	case "postgres":
		db, err = sqlx.Connect("pgx", "user="+cfg.PGUser+" password="+cfg.PGPass+" dbname=mydb sslmode=disable")
		if err != nil {
			log.Fatal(fmt.Errorf("Error with connecting postgres database: %v", err))
		}

		ur = pgsql_repo.NewUserRepository(db)
		tr = pgsql_repo.NewTokenRepo(db, ur)
		gr = pgsql_repo.NewGraphRepo(db)
	case "sqlite":
		dsn := cfg.SQLitePath + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)"
		db, err = sqlx.Connect("sqlite", dsn)
		if err != nil {
			log.Fatal(fmt.Errorf("error with connecting databse: %v", err))
		}

		err = sqlite_repo.CreateTebles(db, "./db_scripts/sqlite/create.sql")
		if err != nil {
			log.Fatal(fmt.Errorf("error with creating tables: %v", err))
		}

		ur = sqlite_repo.NewUserRepository(db)
		tr = sqlite_repo.NewTokenRepo(db, ur)
		gr = sqlite_repo.NewGraphRepo(db)
	default:
		log.Fatal(fmt.Errorf("can't get DB_TYPE, variable getted: %v, excepted sqlite or postgres", cfg.DBType))
	}

	defer db.Close()

	if err := seed.CompareAndSeed(
		context.Background(),
		"./data/graph/dialogs.json",
		gr,
	); err != nil {
		log.Fatal(fmt.Errorf("error with synchronizing dialog graph: %w", err))
	}

	admin, err := ur.FindByName(context.Background(), cfg.AdminUsername)
	if err != nil {
		log.Fatal(fmt.Errorf("error with loading admin user: %v", err))
	}

	if admin == nil {
		admin, err = domain.NewUser(cfg.AdminUsername, cfg.AdminPassword)
		if err != nil {
			log.Fatal(fmt.Errorf("error with creating admin user: %v", err))
		}
		err = ur.Add(context.Background(), admin)
		if err != nil {
			log.Fatal(fmt.Errorf("error with saving admin user: %v", err))
		}
	}

	as := services.NewAuthService(ur)
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
