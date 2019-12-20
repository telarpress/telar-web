package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/users/setting/config"
	"github.com/red-gold/telar-web/src/controllers/users/setting/handlers"
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
		server.POST("/", handlers.CreateSettingGroupHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/", handlers.UpdateUserSettingHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/", handlers.DeleteUserAllSettingHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.GetAllUserSetting(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
