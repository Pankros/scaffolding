package mapper

import model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"

type PaymentMethodsMapper struct{}

func NewPaymentMethodsMapper() PaymentMethodsMapper {
	return PaymentMethodsMapper{}
}

func (m PaymentMethodsMapper) ToListDTO(entities []model.PaymentMethods) []model.PaymentMethodsOutput {
	var dto []model.PaymentMethodsOutput

	for _, entity := range entities {
		dto = append(dto, m.ToDTO(entity))
	}
	return dto
}

func (m PaymentMethodsMapper) ToDTO(entity model.PaymentMethods) model.PaymentMethodsOutput {
	return model.PaymentMethodsOutput{
		Code: entity.Code,
		ID:   entity.ID,
		Name: entity.Name,
	}
}

func (m PaymentMethodsMapper) ToEntity(dto model.PaymentMethodsOutput) model.PaymentMethods {
	return model.PaymentMethods{
		Code: dto.Code,
		ID:   dto.ID,
		Name: dto.Name,
	}
}
