package poeclient

type Profile struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Realm  string `json:"realm"`
	Locale string `json:"locale"`
}
