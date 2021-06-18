//go:generate go run github.com/Pankros/scaffolding/src github.com/Pankros/scaffolding/generate.ConditionType condition_types
package generate

import "time"

type Audit struct {
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

type ConditionType struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
	Audit
}