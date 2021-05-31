package dto

type Package struct {
	Name     string    `json:"name"`
	Services []Service `json:"services"`
}
