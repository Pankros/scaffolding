package service

import (
	"context"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
)

type PaymentMethodTypeDAO interface {
	Get(ctx context.Context, id int64) (model.PaymentMethodType, error)
	List(ctx context.Context) ([]model.PaymentMethodType, error)
	Create(ctx context.Context, entity model.PaymentMethodType) (int64, error)
	Update(ctx context.Context, entity model.PaymentMethodType) error
	Delete(ctx context.Context, id int64) error
}

type PaymentMethodTypeService struct {
	dao PaymentMethodTypeDAO
}

func NewPaymentMethodTypeService(dao PaymentMethodTypeDAO) PaymentMethodTypeService {
	return PaymentMethodTypeService{dao: dao}
}

func (s PaymentMethodTypeService) List(ctx context.Context) ([]model.PaymentMethodTypeOutput, error) {
	entities, err := s.dao.List(ctx)

	if err != nil {
		return nil, err
	}

	return toListDTO(entities), nil
}

func (s PaymentMethodTypeService) Get(ctx context.Context, id int64) (model.PaymentMethodTypeOutput, error) {
	entity, err := s.dao.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return toDTO(entity), nil
}

func (s PaymentMethodTypeService) Create(ctx context.Context, dto model.PaymentMethodTypeCreate) (model.PaymentMethodTypeOutput, error) {
	entity := toEntity(0, dto)

	id, err := s.dao.Create(ctx, entity)

	if err != nil {
		return model.PaymentMethodTypeOutput{}, err
	}

	return s.Get(ctx, id)
}

func (s PaymentMethodTypeService) Update(ctx context.Context, id int64, dto model.PaymentMethodTypeCreate) (model.PaymentMethodTypeOutput, error) {
	if id == 0 {
		return model.PaymentMethodTypeOutput{}, "can't update PaymentMethodType type without ID"
	}

	entity := toEntity(id, dto)

	err := s.dao.Update(ctx, entity)

	if err != nil {
		return model.PaymentMethodTypeOutput{}, err
	}

	return s.Get(ctx, id)
}

func (s PaymentMethodTypeService) Delete(ctx context.Context, id int64) error {
	err := s.dao.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
