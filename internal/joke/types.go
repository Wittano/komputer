package joke

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Type     string
	Category string
)

const (
	Single  Type = "single"
	TwoPart Type = "twopart"
)

const (
	PROGRAMMING Category = "Programming"
	MISC        Category = "Misc"
	DARK        Category = "Dark"
	YOMAMA      Category = "YoMama"
	Any         Category = "Any"
)

type DbModel struct {
	ID       primitive.ObjectID `bson:"_id"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
	Type     Type               `bson:"type"`
	Category Category           `bson:"category"`
	GuildID  string             `bson:"guild_id"`
}
