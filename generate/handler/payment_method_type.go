package handler

import (
	"context"
	web "github.com/mercadolibre/fury_go-core/pkg/web"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
	"net/http"
)

type PaymentMethodTypeService interface {
	Get(ctx context.Context, id int64) (model.PaymentMethodTypeOutput, error)
	List(ctx context.Context) ([]model.PaymentMethodTypeOutput, error)
	Create(ctx context.Context, entity model.PaymentMethodTypeCreate) (model.PaymentMethodTypeOutput, error)
	Update(ctx context.Context, id int64, entity model.PaymentMethodTypeCreate) (model.PaymentMethodTypeOutput, error)
	Delete(ctx context.Context, id int64) error
}

type PaymentMethodTypeHandler struct {
	service PaymentMethodTypeService
}

func NewPaymentMethodTypeHandler(service PaymentMethodTypeService) PaymentMethodTypeHandler {
	return PaymentMethodTypeHandler{service: service}
}

func (h PaymentMethodTypeHandler) List(w http.ResponseWriter, r http.Request) error {
	dto, err := h.service.List(r.Context())

	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return web.RespondJSON(w, dto, http.StatusOK)
}

func (h PaymentMethodTypeHandler) Get(w http.ResponseWriter, r http.Request) error {
	id, err := web.Params(r).Int("id")

	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	dto, err := h.service.Get(r.Context(), id)

	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if dto == (model.PaymentMethodTypeOutput{}) {
		return web.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return web.RespondJSON(w, dto, http.StatusOK)
}

func (h PaymentMethodTypeHandler) Create(w http.ResponseWriter, r http.Request) error {
	dto := model.PaymentMethodTypeCreate{}
	err := web.Bind(r, &dto)

	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.service.Create(r.Context(), dto)

	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return web.RespondJSON(w, resp, http.StatusOK)
}

func (h PaymentMethodTypeHandler) Update(w http.ResponseWriter, r http.Request) error {
	dto := model.PaymentMethodTypeCreate{}
	err := web.Bind(r, &dto)

	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	id, err := web.Params(r).Int("id")

	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.service.Update(r.Context(), int64(id), dto)

	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return web.RespondJSON(w, resp, http.StatusOK)
}

func (h PaymentMethodTypeHandler) Delete(w http.ResponseWriter, r http.Request) error {
	id, err := web.Params(r).Int("id")

	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	err = h.service.Delete(r.Context(), int64(id))

	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return web.RespondJSON(w, nil, http.StatusNoContent)
}
