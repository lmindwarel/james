package datastore

import (
	"context"
	"time"

	models "github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollParameters = "parameters"

func (ds *Datastore) GetParameters() (params []models.Parameter, err error) {
	cursor, err := ds.db.Collection(CollParameters).Find(context.Background(), bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return []models.Parameter{}, nil
		}

		return params, errors.Wrap(err, "failed to query accounts")
	}

	err = cursor.All(context.Background(), &params)

	return
}

func (ds *Datastore) GetParameter(id models.UUID) (param models.Parameter, err error) {
	err = ds.db.Collection(CollParameters).FindOne(context.Background(), bson.M{"_id": id}).Decode(&param)
	return
}

func (ds *Datastore) GetStringParameter(id models.UUID) (value string, err error) {
	param, err := ds.GetParameter(id)
	if err != nil {
		return value, errors.Wrap(err, "failed to get parameter")
	}

	value, isString := param.Value.(string)
	if !isString {
		return value, errors.New("parameter is not a string")
	}

	return
}

func (ds *Datastore) GetUUIDParameter(id models.UUID) (value models.UUID, err error) {
	paramString, err := ds.GetStringParameter(id)
	if err != nil {
		return value, errors.Wrap(err, "failed to get parameter")
	}

	return models.UUID(paramString), nil
}

func (ds *Datastore) UpsertParameter(param models.Parameter) (models.Parameter, error) {
	var err error
	if param.DateCreated.IsZero() {
		now := time.Now()
		param.DateCreated = now
		param.DateUpdated = now
		_, err = ds.db.Collection(CollParameters).InsertOne(context.Background(), param)
	} else {
		param.BaseModel.Update()
		_, err = ds.db.Collection(CollParameters).ReplaceOne(context.Background(), bson.M{"_id": param.ID}, param)

	}

	return param, err
}
