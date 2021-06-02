package model

type PaymentMethodsOutput struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type PaymentMethodsCreate struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
