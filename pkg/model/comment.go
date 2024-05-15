package model

type Comment struct {
	ID          string `bson:"_id"`
	OwnerEmail  string `bson:"owner_email"`
	FromEmail   string `bson:"from_email"`
	Cotent      string `bson:"content"`
	IsDisplayed bool   `bson:"is_displayed"`
	Order       int    `bson:"order"`
}
