package engine

type Configuration struct {
	RootPath     string   `json:"-"`
	MpqLoadOrder []string `json:"mpqLoadOrder"`
}
