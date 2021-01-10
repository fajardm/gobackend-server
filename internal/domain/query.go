package domain

type QueryComparator struct {
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

type Query map[FieldName]QueryComparator

type QueryOption struct {
	Skip            int         `json:"skip"`
	Limit           int         `json:"limit"`
	ACL             []string    `json:"acl"`
	Sort            Sort        `json:"sort"`
	Count           bool        `json:"count"`
	Keys            []string    `json:"keys"`
	Op              string      `json:"op"`
	Distinct        bool        `json:"distinct"`
	Pipeline        interface{} `json:"pipeline"`
	ReadPreference  string      `json:"readPreference"`
	Hint            interface{} `json:"hint"`
	Explain         bool        `json:"explain"`
	CaseInsensitive bool        `json:"caseInsensitive"`
	Action          string      `json:"action"`
	AddsField       bool        `json:"addsField"`
}

type UpdateQueryOption struct {
	Many   bool `json:"many"`
	Upsert bool `json:"upsert"`
}
