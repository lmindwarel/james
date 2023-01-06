package models

const (
	ParamCurrentSpotifyCredential = UUID("current_spotify_credential")
)

type Parameter struct {
	BaseModel `bson:",inline"`
	Value     interface{} `json:"value" bson:"value"`
}

type ParameterPatch struct {
	Value interface{} `json:"value"`
}
