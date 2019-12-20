package controllers

import (
	"context"
	"fmt"

	core "github.com/red-gold/telar-core"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/data/mongodb"
)

// Start run startup operations
func Start(ctx context.Context) (interface{}, error) {
	core.InitConfig()

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:
		mongoClient, err := mongodb.NewMongoClient(ctx)
		if err != nil {
			return nil, err
		}
		return mongoClient, nil
	}

	return nil, fmt.Errorf("Please set valid database type in confing file!")
}
