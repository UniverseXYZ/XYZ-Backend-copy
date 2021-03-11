package api

import (
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/state"
)

var log = logrus.WithField("module", "api")

type Config struct {
	Port           string
	DevCorsEnabled bool
	DevCorsHost    string
}

type API struct {
	config Config
	engine *gin.Engine

	db *sql.DB
}

func New(db *sql.DB, config Config) *API {
	err := state.Init(db)
	if err != nil {
		log.Fatal(err)
	}

	return &API{
		config: config,
		db:     db,
	}
}

func (a *API) Run() {
	a.engine = gin.Default()

	if a.config.DevCorsEnabled {
		a.engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{a.config.DevCorsHost},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	a.setRoutes()

	err := a.engine.Run(":" + a.config.Port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.NewTicker(1 * time.Minute)

		for {
			select {
			case <-t.C:
				err := state.Refresh()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()
}

func (a *API) Close() {
}
