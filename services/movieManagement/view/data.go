package view

type MoviesCategory struct {
	CategoryID   int    `json:"categoryID"`
	CategoryName string `json:"categoryName"`
}

type Movies struct {
	MoviesID   int    `json:"MoviesID"`
	MovieName  string `json:"MovieName"`
	CategoryID int    `json:"CategoryID"`
	URL        string `json:"URL"`
	UserID     int    `json:"UserID"`
}
