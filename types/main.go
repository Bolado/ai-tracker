package types

type Article struct {
	Title   string `json:"name"`
	Summary string `json:"summary"`
	Link    string `json:"link"`
	Date    string `json:"date"`
	Source  string `json:"source"`
	Image   string `json:"image"`
}
