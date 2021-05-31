package dto

type Service struct {
	Name      string   `json:"name"`
	Factory   []string `json:"factory"`
	Arguments []string `json:"arguments"`
	Pointer   *bool    `json:"pointer"`
	Public    *bool    `json:"public"`
}
