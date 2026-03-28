package api

type Data struct {
	Value struct {
		Id   string `json:"id"`
		Self bool   `json:"self"`
		Jid  string `json:"jid"`
	} `json:"value"`
	Type int `json:"type"`
}
