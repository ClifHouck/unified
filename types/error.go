package types

type Error struct {
	Error string `json:"error"`
	Name  string `json:"name"`
	Cause struct {
		Error string `json:"error"`
		Name  string `json:"name"`
	} `json:"cause"`
}
