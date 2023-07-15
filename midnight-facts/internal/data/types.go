package data

import "database/sql"

type CatFact struct {
	Id              int    `json:"id"`
	Fact            string `json:"fact"`
    Length          int    `json:"length"`
	//User            string `json:"user"`
	//UpdatedAt       string `json:"updatedAt"`
	//SendDate        string `json:"sendDate"`
	//Deleted         bool   `json:"deleted"`
	//Source          string `json:"source"`
	//Type            string `json:"string"`
	//StatusVerified  string `json:"status.verified"`
	//StatusFeedback  string `json:"status.feedback"`
	//StatusSentCount string `json:"status.sentCount"`
}

type Models struct {
    Facts PostgresDB
}

func NewModels(db *sql.DB) Models {
    return Models{
        Facts: PostgresDB{DB: db},
    }
}
