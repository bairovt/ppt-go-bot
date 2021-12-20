package main

type RecDoc struct {
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

type UserDoc struct {
	Key         string `json:"_key"`    // =ID
	ID          int64  `json:"id"`      // tg
	IsBot       bool   `json:"is_bot"`  // tg
	AppRole     string `json:"approle"` // "admin"
	ChatID      int64  `json:"chat_id"` // tg, =ID
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"username"`
	PptRole     string `json:"role"`                 // "P", "D"
	StartDate   int64  `json:"start_date,omitempty"` // 1635819850,
	RestartDate int64  `json:"restart_date"`         // 1635819850,
	Tel1        string `json:"tel1"`                 // "89144787432",
	Tel2        string `json:"tel2"`                 // "89144787432",
}

type PointDoc struct {
  Name string `json:"name"`
  Names []string `json:"names"`
  Regex string `json:"regex"`
}