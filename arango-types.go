package main

type Rec struct {
	Key          string   `json:"_key"`
	Role         string   `json:"role"` // "P", "D"
	Tels         []string `json:"tels"`
	Route        []string `json:"route"`
	ViberEventID int64    `json:"EventID"`
	ChatName     string   `json:"ChatName"`
	ClientName   string   `json:"ClientName"`
	TimeStamp    int64    `json:"TimeStamp"`
	Body         string   `json:"Body"`
	// "MType": 1,
	// "EType": 0,
	// "ContactID": 791,
	ChatToken   string `json:"ChatToken"`
	Cargo       bool   `json:"cargo"`
	CleanedBody string `json:"cleanedBody"`
	Src         string `json:"src"` // viber, vk
}

type User struct {
	Key          string `json:"_key"`             // =ID
	ID           int64  `json:"id"`               // tg
	IsBot        bool   `json:"is_bot,omitempty"` // tg
	AppRole      string `json:"approle"`          // "admin"
	ChatID       int64  `json:"chat_id"`          // tg, =ID
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	PptRole      string `json:"role"`       // "P", "D"
	StartDate    int64  `json:"start_date"` // 1635819850,
	Tel1         string `json:"tel1"`       // "89144787432",
	Tel2         string `json:"tel2"`       // "89144787432",
	Username     string `json:"username"`
}
