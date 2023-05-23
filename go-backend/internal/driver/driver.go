package driver

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestDB(d *mongo.Client) error {
	err := d.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}
