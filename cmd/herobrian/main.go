package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bdreece/herobrian/internal/renderer"
	"github.com/bdreece/herobrian/internal/validator"
	"github.com/bdreece/herobrian/pkg/auth"
	"github.com/bdreece/herobrian/pkg/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogEcho "github.com/samber/slog-echo"
)

func main() {
	defer quit()
	provider, err := loadConfig(args.ConfigPath)
	if err != nil {
		panic(err)
	}

	render, err := renderer.New()
	if err != nil {
		panic(err)
	}

	logger, err := createLogger(provider)
	if err != nil {
		panic(err)
	}

	linodeClient, err := createLinodeClient(provider)
	if err != nil {
		panic(err)
	}

	router := echo.New()
	router.Renderer = render
	router.Validator = validator.Default

    webOpts := new(struct {
		AppDir    string
		StaticDir string
	})

    if err = provider.Get("web").Populate(webOpts); err != nil {
        panic(err)
    }

	router.Use(
		middleware.BodyLimit("4M"),
		middleware.Decompress(),
		middleware.Gzip(),
		middleware.CSRF(),
		middleware.Secure(),
        middleware.Static(webOpts.StaticDir),
        middleware.Static(webOpts.AppDir),
		slogEcho.New(logger))

	authOpts := new(auth.Options)
	if err = provider.Get("auth").Populate(authOpts); err != nil {
		panic(err)
	}

	authmw := auth.NewMiddleware(authOpts)

	router.GET("/", routes.Home(linodeClient, logger), authmw)
    router.GET("/sse", routes.SSE(linodeClient), authmw)
    router.POST("/boot", routes.Boot(linodeClient), authmw)
    router.POST("/shutdown", routes.Shutdown(linodeClient), authmw)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	select {
	case err := <-launch(router):
		panic(err)
	case <-ctx.Done():
		break
	}

	if err := shutdown(router); err != nil {
		panic(err)
	}
}

func launch(router *echo.Echo) <-chan error {
	errch := make(chan error, 1)
	go func() {
		if err := router.Start(args.Addr()); err != nil {
			errch <- err
		}
	}()

	return errch
}

func shutdown(router *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func quit() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "unexpected panic occurred: %v", r)
		os.Exit(1)
	}
}
