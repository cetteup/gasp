package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/gasp/cmd/gasp/internal/config"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getawardsinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getbackendinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getleaderboard"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getplayerinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getrankinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getunlocksinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/ranknotification"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/searchforplayers"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/selectunlock"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/verifyplayer"
	"github.com/cetteup/gasp/cmd/gasp/internal/options"
	armysql "github.com/cetteup/gasp/internal/domain/army/sql"
	awardsql "github.com/cetteup/gasp/internal/domain/award/sql"
	fieldsql "github.com/cetteup/gasp/internal/domain/field/sql"
	killsql "github.com/cetteup/gasp/internal/domain/kill/sql"
	kitsql "github.com/cetteup/gasp/internal/domain/kit/sql"
	leaderboardsql "github.com/cetteup/gasp/internal/domain/leaderboard/sql"
	playersql "github.com/cetteup/gasp/internal/domain/player/sql"
	unlocksql "github.com/cetteup/gasp/internal/domain/unlock/sql"
	vehiclesql "github.com/cetteup/gasp/internal/domain/vehicle/sql"
	weaponsql "github.com/cetteup/gasp/internal/domain/weapon/sql"
	"github.com/cetteup/gasp/internal/sqlutil"
	"github.com/cetteup/gasp/pkg/asp"
)

var (
	buildVersion = "development"
	buildCommit  = "uncommitted"
	buildTime    = "unknown"
)

func main() {
	version := fmt.Sprintf("gasp %s (%s) built at %s", buildVersion, buildCommit, buildTime)
	opts := options.Init()

	// Print version and exit
	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    !opts.ColorizeLogs,
		TimeFormat: time.RFC3339,
	})
	if opts.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("config", opts.ConfigPath).
			Msg("Failed to read config file")
	}

	db := sqlutil.Connect(
		cfg.Database.Host,
		cfg.Database.DatabaseName,
		cfg.Database.Username,
		cfg.Database.Password,
	)
	defer func() {
		err2 := db.Close()
		if err2 != nil {
			log.Error().
				Err(err2).
				Msg("Failed to close database connection")
		}
	}()

	playerRepository := playersql.NewRepository(db)
	armyRecordRepository := armysql.NewRecordRepository(db)
	awardRecordRepository := awardsql.NewRecordRepository(db)
	fieldRecordRepository := fieldsql.NewRecordRepository(db)
	killHistoryRecordRepository := killsql.NewHistoryRecordRepository(db)
	kitRecordRepository := kitsql.NewRecordRepository(db)
	leaderboardRepository := leaderboardsql.NewRepository(db)
	vehicleRecordRepository := vehiclesql.NewRecordRepository(db)
	weaponRecordRepository := weaponsql.NewRecordRepository(db)
	unlockRepository := unlocksql.NewRepository(db)
	unlockRecordRepository := unlocksql.NewRecordRepository(db)
	gaih := getawardsinfo.NewHandler(awardRecordRepository)
	gbih := getbackendinfo.NewHandler(unlockRepository)
	glbh := getleaderboard.NewHandler(leaderboardRepository)
	gpih := getplayerinfo.NewHandler(
		playerRepository,
		armyRecordRepository,
		fieldRecordRepository,
		killHistoryRecordRepository,
		kitRecordRepository,
		vehicleRecordRepository,
		weaponRecordRepository,
	)
	grih := getrankinfo.NewHandler(playerRepository)
	guih := getunlocksinfo.NewHandler(playerRepository, awardRecordRepository, unlockRecordRepository)
	rnh := ranknotification.NewHandler(playerRepository)
	sfph := searchforplayers.NewHandler(playerRepository)
	suh := selectunlock.NewHandler(playerRepository, awardRecordRepository, unlockRecordRepository)
	vph := verifyplayer.NewHandler(playerRepository)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	// Error handler is strongly modeled after the default one
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		message := http.StatusText(code)
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
			message = http.StatusText(code)
		}

		// Send response
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			// Always return 200/OK to match original GameSpy behaviour.
			// Note: Logs will contain the "underlying" status code, not 200.
			err = c.String(http.StatusOK, asp.NewErrorResponseWithMessage(code, message).Serialize())
		}
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to send error response")
		}

	}
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Second * 10,
		ErrorMessage: asp.NewErrorResponseWithMessage(
			http.StatusServiceUnavailable,
			http.StatusText(http.StatusServiceUnavailable),
		).Serialize(),
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogError:     true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Err(v.Error).
				Str("remote", v.RemoteIP).
				Str("method", v.Method).
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.Truncate(time.Millisecond).String()).
				Str("agent", v.UserAgent).
				Msg("request")

			return nil
		},
	}))

	g := e.Group("/ASP")
	g.GET("/getawardsinfo.aspx", gaih.HandleGET)
	g.GET("/getbackendinfo.aspx", gbih.HandleGET)
	g.GET("/getleaderboard.aspx", glbh.HandleGET)
	g.GET("/getplayerinfo.aspx", gpih.HandleGET)
	g.GET("/getrankinfo.aspx", grih.HandleGET)
	g.GET("/getunlocksinfo.aspx", guih.HandleGET)
	g.GET("/ranknotification.aspx", rnh.HandleGET)
	g.GET("/searchforplayers.aspx", sfph.HandleGET)
	g.GET("/VerifyPlayer.aspx", vph.HandleGET)
	g.POST("/selectunlock.aspx", suh.HandlePOST)

	if err = e.Start(opts.ListenAddr); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().
			Err(err).
			Str("address", opts.ListenAddr).
			Msg("Failed to start server")
	}
}
