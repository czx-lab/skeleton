package bootstrap

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"skeleton/internal/mongo"
	"skeleton/internal/variable"
)

func InitMongo() error {
	if !variable.Config.GetBool("Database.Mongo.Enable") {
		return nil
	}
	var err error
	mongoConfig := variable.Config.Get("Database.Mongo").(map[string]any)
	opts := options.Client().
		ApplyURI(mongoConfig["uri"].(string)).
		SetMaxPoolSize(mongoConfig["maxpoolsize"].(uint64)).
		SetMinPoolSize(mongoConfig["minpoolsize"].(uint64))
	if variable.MongoDB, err = mongo.New(opts); err != nil {
		return err
	}
	return nil
}
