package connection

type Message struct {
	RecipientId string `json:"r"`
	Secret      string `json:"s"`
	Text        string `json:"t"`
}
