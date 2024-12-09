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
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getunlocksinfo"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/searchforplayers"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/verifyplayer"
	"github.com/cetteup/gasp/cmd/gasp/internal/options"
	awardsql "github.com/cetteup/gasp/internal/domain/award/sql"
	playersql "github.com/cetteup/gasp/internal/domain/player/sql"
	unlocksql "github.com/cetteup/gasp/internal/domain/unlock/sql"
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
	awardRecordRepository := awardsql.NewRecordRepository(db)
	unlockRepository := unlocksql.NewRepository(db)
	unlockRecordRepository := unlocksql.NewRecordRepository(db)
	gaih := getawardsinfo.NewHandler(awardRecordRepository)
	gbih := getbackendinfo.NewHandler(unlockRepository)
	guih := getunlocksinfo.NewHandler(playerRepository, awardRecordRepository, unlockRecordRepository)
	sfph := searchforplayers.NewHandler(playerRepository)
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
			err = c.NoContent(he.Code)
		} else {
			err = c.String(code, asp.NewErrorResponseWithMessage(code, message).Serialize())
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
	g.GET("/getunlocksinfo.aspx", guih.HandleGET)
	g.GET("/searchforplayers.aspx", sfph.HandleGET)
	g.GET("/VerifyPlayer.aspx", vph.HandleGET)

	if err = e.Start(opts.ListenAddr); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().
			Err(err).
			Str("address", opts.ListenAddr).
			Msg("Failed to start server")
	}
}
