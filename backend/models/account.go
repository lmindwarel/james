package models

// Account store infomations about a quizzbox account
type Account struct {
	BaseModel `bson:",inline"`
	// Name of the account
	Name string `json:"name" bson:"name"`
	// Icon of the account
	Icon string `json:"icon" bson:"icon"`
}

// AccountPatch
type AccountPatch struct {
	// Name of the account
	Name *string `json:"name" bson:"name"`
	// Icon of the account
	Icon *string `json:"icon" bson:"icon"`
}

const SpotifyPasswordHashKey = "example key 1234"
const SpotifyDeviceName = "James"

type SpotifyCredential struct {
	BaseModel      `bson:",inline"`
	User           string `json:"user"`
	HashedPassword string `json:"hashedPassword"`
}
