//go:generate go run github.com/Pankros/scaffolding github.com/Pankros/scaffolding/generate.PaymentMethodType
package generate

import "time"

type PaymentMethodType struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}
