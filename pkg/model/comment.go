package model

type Comment struct {
	ID          string `bson:"_id" json:"id"`
	OwnerEmail  string `bson:"owner_email" json:"ownerEmail"`
	FromEmail   string `bson:"from_email" json:"fromEmail"`
	Content     string `bson:"content" json:"content"`
	IsDisplayed bool   `bson:"is_displayed" json:"isDisplayed"`
	Order       int    `bson:"order" json:"order"`
}
