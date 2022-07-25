package datastore

import (
	"context"

	models "github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAccount return the user for given id
func (ds *Datastore) GetAccount(id models.UUID) (account models.Account, err error) {
	err = ds.db.Collection(models.CollAccounts).FindOne(context.Background(), bson.M{"_id": id}).Decode(&account)
	return
}

// GetAccounts all accounts
func (ds *Datastore) GetAccounts() (accounts []models.Account, err error) {
	cursor, err := ds.db.Collection(models.CollAccounts).Find(context.Background(), bson.D{})
	if err != nil {
		return accounts, errors.Wrap(err, "failed to query accounts")
	}

	err = cursor.All(context.Background(), &accounts)

	return
}

// UpsertAccount add or update the given user in mongo database
func (ds *Datastore) UpsertAccount(account models.Account) (models.Account, error) {
	var err error
	if models.EmptyUUID(account.ID) {
		account.BaseModel = models.NewBaseModel()
		_, err = ds.db.Collection(models.CollAccounts).InsertOne(context.Background(), account)
	} else {
		account.BaseModel.Update()
		_, err = ds.db.Collection(models.CollAccounts).ReplaceOne(context.Background(), bson.M{"_id": account.ID}, account)

	}

	return account, err
}
