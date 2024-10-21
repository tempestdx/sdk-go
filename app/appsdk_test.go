package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc        string
		app         *App
		options     []AppOption
		shouldPanic bool
	}{
		{
			desc: "OK - No Options",
			app:  &App{},
		},
		{
			desc: "OK - With Resource Definition",
			app: &App{
				resourceDefinitions: []ResourceDefinition{
					{
						Type: "example",
					},
				},
			},
			options: []AppOption{
				WithResourceDefinition(ResourceDefinition{
					Type: "example",
				}),
			},
		},
		{
			desc: "OK - With Resource Definitions",
			app: &App{
				resourceDefinitions: []ResourceDefinition{
					{
						Type: "example",
					},
					{
						Type: "example2",
					},
				},
			},
			options: []AppOption{
				WithResourceDefinitions(
					ResourceDefinition{
						Type: "example",
					},
					ResourceDefinition{
						Type: "example2",
					},
				),
			},
		},
		{
			desc:        "PANIC - With bad Resource Type",
			shouldPanic: true,
			options: []AppOption{
				WithResourceDefinition(ResourceDefinition{
					Type: "&--invalid",
				}),
			},
		},
		{
			desc:        "PANIC - With duplicate Resource Definitions",
			shouldPanic: true,
			options: []AppOption{
				WithResourceDefinition(ResourceDefinition{
					Type: "example",
				}),
				WithResourceDefinition(ResourceDefinition{
					Type: "example",
				}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.shouldPanic {
				require.Panics(t, func() {
					New(tc.options...)
				})
				return
			}

			assert.Equal(t, tc.app, New(tc.options...))
		})
	}
}
