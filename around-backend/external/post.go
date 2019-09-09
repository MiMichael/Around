package external

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Post struct {
	User     string   `json:"user"`
	Message  string   `json:"message"`
	Location Location `json:"location"`
	Url      string   `json:"url"`
}
