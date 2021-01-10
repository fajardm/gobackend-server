package domain

type Sort struct {
	Field     FieldName `json:"field"`
	Direction Direction `json:"directon"`
}
