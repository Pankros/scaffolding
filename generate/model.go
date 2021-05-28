//go:generate go run github.com/Pankros/scaffolding/src github.com/Pankros/scaffolding/generate.PaymentMethodType payment_method_types
package generate

import "time"

//payment_method_types
type PaymentMethodType struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}
