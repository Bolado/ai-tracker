package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Article struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Summary string             `json:"summary" bson:"summary"`
	Link    string             `json:"link" bson:"link"`
	Date    string             `json:"date" bson:"date"`
	Source  string             `json:"source" bson:"source"`
	Image   string             `json:"image" bson:"image"`
	Content string             `json:"content" bson:"content"`
}

type Website struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	MainElement     string `json:"mainElement"`
	TitleElement    string `json:"titleElement"`
	SubtitleElement string `json:"subtitleElement"`
	ImageElement    string `json:"imageElement"`
	ContentElement  string `json:"contentElement"`
	AnchorElement   string `json:"anchorElement"`
	DateElement     string `json:"dateElement"`
}
