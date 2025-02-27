package app

import (
	"fmt"
	"regexp"

	"github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1/appv1connect"
)

type App struct {
	appv1connect.UnimplementedAppServiceHandler

	resourceDefinitions []ResourceDefinition
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
