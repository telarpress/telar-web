package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/admin/config"
	"github.com/red-gold/telar-web/src/controllers/admin/handlers"
)

func init() {

	cf.InitConfig()
}

// Cache state
var server *coreServer.ServerRouter
var db interface{}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Start
	if db == nil {
		var startErr error
		db, startErr = controllers.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {
		server = coreServer.NewServerRouter()
		server.POSTWR("/setup", handlers.SetupHandler(), coreServer.RouteProtectionAdmin)
		server.GET("/setup", handlers.SetupPageHandler, coreServer.RouteProtectionAdmin)
		server.GET("/login", handlers.LoginPageHandler, coreServer.RouteProtectionPublic)
		server.POSTWR("/login", handlers.LoginAdminHandler(db), coreServer.RouteProtectionPublic)
	}
	server.ServeHTTP(w, r)
}
