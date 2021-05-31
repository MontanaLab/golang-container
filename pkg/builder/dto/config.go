package dto

type Config struct {
	Defaults Defaults  `json:"defaults"`
	Packages []Package `json:"packages"`
}
