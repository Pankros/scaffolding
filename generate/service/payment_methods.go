package service

import (
	"context"
	"errors"
	dao "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/dao"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
	utils "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/platform/utils"
)

type PaymentMethodsDAO interface {
	Get(ctx context.Context, tx dao.DBConnection, id int64) (model.PaymentMethods, error)
	List(ctx context.Context, tx dao.DBConnection) ([]model.PaymentMethods, error)
	Create(ctx context.Context, tx dao.DBConnection, entity model.PaymentMethods) (int64, error)
	Update(ctx context.Context, tx dao.DBConnection, entity model.PaymentMethods) error
	Delete(ctx context.Context, tx dao.DBConnection, id int64) error
}

type PaymentMethodsMapper interface {
	ToListDTO(entities []model.PaymentMethods) []model.PaymentMethodsOutput
	ToDTO(entity model.PaymentMethods) model.PaymentMethodsOutput
	ToEntity(id int64, dto model.PaymentMethodsCreate) model.PaymentMethods
}

type PaymentMethodsService struct {
	dao       PaymentMethodsDAO
	txService TransactionManager
	mapper    PaymentMethodsMapper
}

func NewPaymentMethodsService(dao PaymentMethodsDAO, txService TransactionManager, mapper PaymentMethodsMapper) PaymentMethodsService {
	return PaymentMethodsService{
		dao:       dao,
		mapper:    mapper,
		txService: txService,
	}
}

func (s PaymentMethodsService) List(ctx context.Context) ([]model.PaymentMethodsOutput, error) {
	tx, err := s.txService.Open(ctx)
	if err != nil {
		return nil, err
	}

	entities, err := s.dao.List(ctx, tx)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToListDTO(entities), nil
}

func (s PaymentMethodsService) Get(ctx context.Context, id int64) (model.PaymentMethodsOutput, error) {
	tx, err := s.txService.Open(ctx)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	entity, err := s.dao.Get(ctx, tx, id)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	return s.mapper.ToDTO(entity), nil
}

func (s PaymentMethodsService) Create(ctx context.Context, dto model.PaymentMethodsCreate) (resp model.PaymentMethodsOutput, err error) {
	entity := s.mapper.ToEntity(0, dto)

	tx, err := s.txService.OpenTx(ctx)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	defer func() {
		deferErr := s.txService.CloseTx(ctx, tx, err)
		if deferErr != nil {
			err = utils.WrapOrCreateError(err, deferErr)
			resp = model.PaymentMethodsOutput{}
		}
	}()

	id, err := s.dao.Create(ctx, tx, entity)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	created, err := s.dao.Get(ctx, tx, id)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s PaymentMethodsService) Update(ctx context.Context, id int64, dto model.PaymentMethodsCreate) (resp model.PaymentMethodsOutput, err error) {
	if id == 0 {
		return model.PaymentMethodsOutput{}, errors.New("can't update PaymentMethods type without ID")
	}

	tx, err := s.txService.OpenTx(ctx)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	defer func() {
		deferErr := s.txService.CloseTx(ctx, tx, err)
		if deferErr != nil {
			err = utils.WrapOrCreateError(err, deferErr)
			resp = model.PaymentMethodsOutput{}
		}
	}()

	entity := s.mapper.ToEntity(id, dto)

	err = s.dao.Update(ctx, tx, entity)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	updated, err := s.dao.Get(ctx, tx, id)
	if err != nil {
		return model.PaymentMethodsOutput{}, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s PaymentMethodsService) Delete(ctx context.Context, id int64) error {
	tx, err := s.txService.OpenTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		deferErr := s.txService.CloseTx(ctx, tx, err)
		if deferErr != nil {
			err = utils.WrapOrCreateError(err, deferErr)
		}
	}()
	err = s.dao.Delete(ctx, tx, id)

	if err != nil {
		return err
	}

	return nil
}
