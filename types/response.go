package types

type ResponseGetUrl struct {
	Type      string         `json:"type"`
	Html      string         `json:"html"`
	Lastbatch bool           `json:"lastbatch"`
	Postflair map[string]int `json:"postflair"`
}
