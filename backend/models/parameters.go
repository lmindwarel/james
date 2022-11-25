package models

const (
	ParamSpotifyCredentials = UUID("spotify-credentials")
)

type Parameter struct {
	BaseModel `bson:",inline"`
	Value     interface{} `bson:"value"`
}
