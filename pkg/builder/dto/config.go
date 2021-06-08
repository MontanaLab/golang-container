package dto

type Config struct {
	Imports  []string  `json:"imports"`
	Defaults Defaults  `json:"defaults"`
	Packages []Package `json:"packages"`
}
