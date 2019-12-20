package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	appConfig "github.com/red-gold/telar-web/src/controllers/storage/config"
	"github.com/red-gold/telar-web/src/controllers/storage/handlers"
)

func init() {

	appConfig.InitConfig()
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
		server.POSTFILE("/:uid/:dir", handlers.UploadeHandle(db), coreServer.RouteProtectionCookie)
		// server.PUTWR("/", handlers.GetFileHandle(db), coreServer.RouteProtectionCookie)
		// server.DELETEWR("/file/:name", handlers.GetFileHandle(db), coreServer.RouteProtectionCookie)
		// server.DELETEWR("/dir/:dir", handlers.GetFileHandle(db), coreServer.RouteProtectionCookie)
		server.GETWR("/:uid/:dir/:name", handlers.GetFileHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
