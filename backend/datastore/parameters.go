package datastore

import (
	"context"
	"time"

	models "github.com/lmindwarel/james/backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

const CollParameters = "parameters"

func (ds *Datastore) GetParameter(id models.UUID) (param models.Parameter, err error) {
	err = ds.db.Collection(CollParameters).FindOne(context.Background(), bson.M{"_id": id}).Decode(&param)
	return
}

func (ds *Datastore) GetStringParameter(id models.UUID) (string, error) {
	param, err := ds.GetParameter(id)
	return param.Value.(string), err
}

func (ds *Datastore) GetUUIDParameter(id models.UUID) (models.UUID, error) {
	res, err := ds.GetStringParameter(id)
	return models.UUID(res), err
}

func (ds *Datastore) UpsertParameter(param models.Parameter) (models.Parameter, error) {
	var err error
	if param.DateCreated.IsZero() {
		now := time.Now()
		param.DateCreated = now
		param.DateUpdated = now
		_, err = ds.db.Collection(CollAccounts).InsertOne(context.Background(), param)
	} else {
		param.BaseModel.Update()
		_, err = ds.db.Collection(CollAccounts).ReplaceOne(context.Background(), bson.M{"_id": param.ID}, param)

	}

	return param, err
}
