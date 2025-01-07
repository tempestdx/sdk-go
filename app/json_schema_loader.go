package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	appSchemaURL      = "https://developer.tempestdx.com/schema/v1/tempest-app-schema.json"
	propertySchemaURL = "https://developer.tempestdx.com/schema/v1/tempest-property-schema.json"
)

type TempestSchemaLoader struct {
	client *http.Client
}

func (l *TempestSchemaLoader) Load(location string) (any, error) {
	switch location {
	case appSchemaURL:
		return l.load(appSchemaURL)
	case propertySchemaURL:
		return l.load(propertySchemaURL)
	default:
		return nil, fmt.Errorf("unknown schema location: %s", location)
	}
}

func (l *TempestSchemaLoader) load(url string) (any, error) {
	resp, err := l.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return jsonschema.UnmarshalJSON(resp.Body)
}

func NewTempestSchemaLoader() *TempestSchemaLoader {
	httpLoader := TempestSchemaLoader{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	return &httpLoader
}
