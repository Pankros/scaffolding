package dao

import (
	"context"
	sqlx "github.com/jmoiron/sqlx"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
)

type PaymentMethodTypeDAO struct {
	db    *sqlx.DB
	audit model.AuditService
}

func NewPaymentMethodTypeDAO(db *sqlx.DB, audit model.AuditService) PaymentMethodTypeDAO {
	return PaymentMethodTypeDAO{
		audit: audit,
		db:    db,
	}
}

func (p PaymentMethodTypeDAO) List(ctx context.Context) ([]model.PaymentMethodType, error) {
	var query = "SELECT id, name, code, created_at, created_by, updated_at, updated_by FROM payment_method_types"

	var rows []model.PaymentMethodType

	err := p.db.SelectContext(ctx, &rows, query)

	if err != nil {
		return []model.PaymentMethodType{}, err
	}

	return rows, nil
}

func (p PaymentMethodTypeDAO) Get(ctx context.Context, id int64) (model.PaymentMethodType, error) {
	var query = "SELECT id, name, code, created_at, created_by, updated_at, updated_by FROM payment_method_types WHERE id = ?"

	var row model.PaymentMethodType

	err := p.db.Get(&row, query, id)

	if err != nil {
		return model.PaymentMethodType{}, err
	}

	return row, nil
}

func (p PaymentMethodTypeDAO) Create(ctx context.Context, entity model.PaymentMethodType) (int64, error) {
	var query = "INSERT INTO payment_method_types (updated_by, name, code, created_at, created_by, updated_at) VALUES (:updated_by, :name, :code, :created_at, :created_by, :updated_at)"

	result, err := p.db.NamedExecContext(ctx, query, &entity)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p PaymentMethodTypeDAO) Update(ctx context.Context, entity model.PaymentMethodType) error {
	var query = "UPDATE payment_method_types SET (updated_by = :updated_by, name = :name, code = :code, updated_at = :updated_at) WHERE id = :id"

	_, err := p.db.NamedExecContext(ctx, query, &entity)

	if err != nil {
		return err
	}

	return nil
}

func (p PaymentMethodTypeDAO) Delete(ctx context.Context, id int64) error {
	var query = "DELETE FROM payment_method_types WHERE id = :id"

	_, err := p.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})

	if err != nil {
		return err
	}

	return nil
}
