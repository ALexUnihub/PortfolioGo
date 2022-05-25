package author

type Author struct {
	UserID   string `json:"id" bson:"id"`
	UserName string `json:"username" bson:"username"`
}
