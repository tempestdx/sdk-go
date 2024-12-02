package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1/appv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	appv1connect.UnimplementedAppServiceHandler

	resourceDefinitions []ResourceDefinition
	done                chan struct{}
}

func New(opts ...AppOption) *App {
	var options appOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &App{
		resourceDefinitions: options.resourceDefinitions,
	}
}

type AppOption func(*appOptions)

type appOptions struct {
	resourceDefinitions []ResourceDefinition
}

const ResourceTypePattern = `^[A-Za-z_][A-Za-z0-9_]*$`

var resourceTypeRegex = regexp.MustCompile(ResourceTypePattern)

func WithResourceDefinition(rd ResourceDefinition) AppOption {
	return func(o *appOptions) {
		if !resourceTypeRegex.MatchString(rd.Type) {
			panic(fmt.Sprintf("resource type '%s' does not match pattern %s", rd.Type, ResourceTypePattern))
		}

		for _, existing := range o.resourceDefinitions {
			if existing.Type == rd.Type {
				panic(fmt.Sprintf("ResourceDefinition with the same type '%s' already exists", rd.Type))
			}
		}

		o.resourceDefinitions = append(o.resourceDefinitions, rd)
	}
}

func WithResourceDefinitions(rds ...ResourceDefinition) AppOption {
	return func(o *appOptions) {
		for _, rd := range rds {
			WithResourceDefinition(rd)(o)
		}
	}
}

// Serve will start the ConnectRPC server and block until the done channel is closed.
// The "port" flag will be set by the Tempest App Server to the port on which the server should listen.
func (a *App) Serve() error {
	flagPort := flag.Int("port", 8080, "port on which to listen")
	flag.Parse()

	mux := http.NewServeMux()
	path, handler := appv1connect.NewAppServiceHandler(a)
	mux.Handle(path, handler)

	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(*flagPort)))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	server := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	<-a.done
	err = server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
