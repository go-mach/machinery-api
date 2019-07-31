// Copyright 2019 Luca Stasio <joshuagame@gmail.com>
// Copyright 2019 IT Resources s.r.l.
//
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package apigear

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-mach/machinery/pkg/machinery"
	"github.com/mitchellh/mapstructure"
)

type (
	// APIConf is the API Gear configuration structure.
	APIConf struct {
		Endpoint struct {
			Port            int
			BaseRoutingPath string
		}
		Security struct {
			Enabled bool
			Jwt     struct {
				Secret     string
				Expiration struct {
					Enabled bool
					Minutes int32
				}
			}
		}
	}

	// APIInfo is the structure returned by the starter func.
	// These params will be used by the APIGear to start the http server.
	APIInfo struct {
		Router      *chi.Mux
		Middlewares []func(http.Handler) http.Handler
	}

	// APICompositionFunc is the api root composition func in which specify API configuration.
	// Here the developer can create its routes and manually dependency wire into
	// controllers and services.2
	APICompositionFunc func(*machinery.Machinery) APIInfo

	// APIGear is the gear structure derived from the Machinery BaseGear.
	APIGear struct {
		machinery.BaseGear
		router  *chi.Mux
		compose APICompositionFunc
		config  APIConf
	}
)

// NewAPIGear creates a new APIGear instance.
func NewAPIGear(uname string, compose APICompositionFunc) *APIGear {
	return &APIGear{
		BaseGear: machinery.BaseGear{
			Uname:  uname,
			Logger: nil,
		},
		router:  chi.NewRouter(),
		compose: compose,
		config:  APIConf{},
	}
}

// Start will start the APIGear runtime.
// It creates the root api router, setup middelwares and mount the app routing.
func (apigear *APIGear) Start(machine *machinery.Machinery) {
	apigear.router = chi.NewRouter()

	// compose api calling provided composition func to get router and middlewares
	api := apigear.compose(machine)

	apigear.router.Use(api.Middlewares...)
	apigear.router.Route(apigear.config.Endpoint.BaseRoutingPath, func(r chi.Router) {
		r.Mount("/", api.Router)
	})

	addr := fmt.Sprintf(":%d", apigear.config.Endpoint.Port)
	log.Fatal(http.ListenAndServe(addr, apigear.router))
}

// Configure get the configuration map and struct it into APIConf structure.
func (apigear *APIGear) Configure(config interface{}) {
	configMap := config.(map[string]interface{})
	mapstructure.Decode(configMap, &apigear.config)
}

// Shutdown .
func (apigear *APIGear) Shutdown() {
	log.Printf("%s SHUT DOWN", apigear.Uname)
}

// Use setup middlewares.
func (apigear *APIGear) Use(middlewares ...func(http.Handler) http.Handler) {

}
