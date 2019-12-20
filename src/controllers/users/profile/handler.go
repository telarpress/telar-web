package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/users/profile/config"
	"github.com/red-gold/telar-web/src/controllers/users/profile/handlers"
	// handlers "github.com/red-gold/telar-web/src/controllers/users/auth/handlers"
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
		fmt.Println("Server is nil")
		server = coreServer.NewServerRouter()
		server.GET("/my", handlers.ReadMyProfileHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.QueryUserProfileHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/id/:userId", handlers.ReadProfileHandle(db), coreServer.RouteProtectionCookie)
		server.POST("/index", handlers.InitProfileIndexHandle(db), coreServer.RouteProtectionHMAC)
	}
	server.ServeHTTP(w, r)
}
