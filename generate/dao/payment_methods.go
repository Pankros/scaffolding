package dao

import (
	"context"
	"errors"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
)

type PaymentMethodsDAO struct {
	audit AuditService
}

func NewPaymentMethodsDAO(audit AuditService) PaymentMethodsDAO {
	return PaymentMethodsDAO{audit: audit}
}

func (p PaymentMethodsDAO) List(ctx context.Context, tx DBConnection) ([]model.PaymentMethods, error) {
	var query = "SELECT id, name, code, created_at, created_by, updated_at, updated_by FROM payment_methods"
	var rows []model.PaymentMethods

	err := tx.SelectContext(ctx, &rows, query)
	if err != nil {
		return []model.PaymentMethods{}, err
	}

	return rows, nil
}

func (p PaymentMethodsDAO) Get(ctx context.Context, tx DBConnection, id int64) (model.PaymentMethods, error) {
	var query = "SELECT id, name, code, created_at, created_by, updated_at, updated_by FROM payment_methods WHERE id = ?"
	var row model.PaymentMethods

	stmt, err := tx.PrepareContext(ctx, query)

	if err != nil {
		return model.PaymentMethods{}, err
	}

	err = stmt.Get(&row, id)
	if err != nil {
		return model.PaymentMethods{}, err
	}

	return row, nil
}

func (p PaymentMethodsDAO) Create(ctx context.Context, tx DBConnection, entity model.PaymentMethods) (int64, error) {
	entity.Audit = p.audit.GetAuditForCreate(ctx)
	var query = "INSERT INTO payment_methods (updated_by, name, code, created_at, created_by, updated_at) VALUES (:updated_by, :name, :code, :created_at, :created_by, :updated_at)"

	result, err := tx.NamedExecContext(ctx, query, &entity)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}

func (p PaymentMethodsDAO) Update(ctx context.Context, tx DBConnection, entity model.PaymentMethods) error {
	entity.Audit = p.audit.GetAuditForUpdate(ctx)
	var query = "UPDATE payment_methods SET (updated_by = :updated_by, name = :name, code = :code, updated_at = :updated_at) WHERE id = :id"

	if entity.ID == 0 {
		return errors.New("can't update an entity without ID")
	}

	_, err := tx.NamedExecContext(ctx, query, &entity)

	if err != nil {
		return err
	}

	return nil
}

func (p PaymentMethodsDAO) Delete(ctx context.Context, tx DBConnection, id int64) error {
	var query = "DELETE FROM payment_methods WHERE id = :id"

	if id == 0 {
		return errors.New("can't delete an entity without ID")
	}

	_, err := tx.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}

	return nil
}
