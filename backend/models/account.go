package models

// CollAccount is the collection name for accounts
const CollAccounts = "accounts"

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
