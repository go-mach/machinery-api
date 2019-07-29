package apigear

import (
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
		Router *chi.Mux
		Addr   string
	}

	// APICompositionFunc is the api root composition func in which specify API configuration.
	// Here the developer can create its routes and manually dependency wire into
	// controllers and services.
	APICompositionFunc func(*machinery.Machinery) APIInfo

	// APIGear is the gear structure derived from the Machinery BaseGear.
	APIGear struct {
		machinery.BaseGear
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
		compose: compose,
		config:  APIConf{},
	}
}

// Start will start the APIGear runtime.
// It creates the root api router, setup middelwares and mount the app routing.
func (apigear *APIGear) Start(machine *machinery.Machinery) {
	router := chi.NewRouter()
	api := apigear.compose(machine)
	router.Route(apigear.config.Endpoint.BaseRoutingPath, func(r chi.Router) {
		r.Mount("/", api.Router)
	})

	log.Fatal(http.ListenAndServe(api.Addr, router))
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
