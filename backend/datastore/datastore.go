package datastore

import (
	"context"
	"fmt"
	"strings"
	"time"

	models "github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"

	cache "github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = utils.GetLogger("datastore")

// Config is the configuration for the datastore
type Config struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Hosts    []string `json:"hosts"`
}

// Datastore manage all the data in project
// It store data in database and cache.
type Datastore struct {
	cache *cache.Cache
	db    *mongo.Database
}

var ErrNotFound = mongo.ErrNoDocuments

// New create a datastore
func New(config Config) (*Datastore, error) {
	loginURL := ""
	if config.Username != "" {
		if config.Password != "" {
			loginURL = fmt.Sprintf("%s:%s@", config.Username, config.Password)
		} else {
			loginURL = fmt.Sprintf("%s@", config.Username)
		}
	}
	url := fmt.Sprintf("mongodb://%s%s/%s", loginURL, strings.Join(config.Hosts, ","), config.Name)
	log.Debugf("Connecting to database %s", url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo")
	}

	// session, err := mgo.Dial(url)
	// if err != nil {
	// 	return nil, errors.New("failed to create mongo session: \n" + err.Error())
	// }

	datastore := &Datastore{
		// db: session.DB(config.Name),
		// Create a cache with a default expiration time of 30 minutes, and which
		// purges expired items every 10 minutes
		cache: cache.New(30*time.Minute, 10*time.Minute),
		db:    client.Database(config.Name),
	}

	return datastore, nil
}

// SetInCache set given value in cache accessible by key
func (ds *Datastore) SetInCache(key string, value interface{}) error {
	return ds.cache.Add(key, value, cache.DefaultExpiration)
}

// GetUniqueCacheKey return a cache key not in use
func (ds *Datastore) GetUniqueCacheKey() string {

	key := string(models.NewUUID())

	for ds.ExistInCache(key) {
		key = string(models.NewUUID())
	}

	return key
}

// ExistInCache is used to test existance of a data in cache
func (ds *Datastore) ExistInCache(key string) bool {
	_, found := ds.cache.Get(key)
	return found
}

func (ds *Datastore) IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}
