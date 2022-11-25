package datastore

import (
	"context"

	models "github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CollAccount is the collection name for accounts
const CollAccounts = "accounts"

// GetAccount return the user for given id
func (ds *Datastore) GetAccount(id models.UUID) (account models.Account, err error) {
	err = ds.db.Collection(CollAccounts).FindOne(context.Background(), bson.M{"_id": id}).Decode(&account)
	return
}

// GetAccounts all accounts
func (ds *Datastore) GetAccounts() (accounts []models.Account, err error) {
	cursor, err := ds.db.Collection(CollAccounts).Find(context.Background(), nil)
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
		_, err = ds.db.Collection(CollAccounts).InsertOne(context.Background(), account)
	} else {
		account.BaseModel.Update()
		_, err = ds.db.Collection(CollAccounts).ReplaceOne(context.Background(), bson.M{"_id": account.ID}, account)

	}

	return account, err
}

const CollSpotifyCredentials = "spotify_credentials"

func (ds *Datastore) GetSpotifyCredentials() (credentials []models.SpotifyCredential, err error) {
	cursor, err := ds.db.Collection(CollSpotifyCredentials).Find(context.Background(), nil)
	if err != nil {
		return credentials, errors.Wrap(err, "failed to query credentials")
	}

	err = cursor.All(context.Background(), &credentials)

	return
}

func (ds *Datastore) GetSpotifyCredential(id models.UUID) (credentials models.SpotifyCredential, err error) {
	err = ds.db.Collection(CollSpotifyCredentials).FindOne(context.Background(), bson.M{"_id": id}).Decode(&credentials)
	return
}

func (ds *Datastore) UpsertSpotifyCredential(credentials models.SpotifyCredential) (models.SpotifyCredential, error) {
	var err error
	if models.EmptyUUID(credentials.ID) {
		credentials.BaseModel = models.NewBaseModel()
		_, err = ds.db.Collection(CollAccounts).InsertOne(context.Background(), credentials)
	} else {
		credentials.BaseModel.Update()
		_, err = ds.db.Collection(CollSpotifyCredentials).ReplaceOne(context.Background(), bson.M{"_id": credentials.ID}, credentials)

	}

	return credentials, err
}
