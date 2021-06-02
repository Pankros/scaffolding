package handler

import (
	"context"
	web "github.com/mercadolibre/fury_go-core/pkg/web"
	model "github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model"
	"net/http"
)

type PaymentMethodsService interface {
	Get(ctx context.Context, id int64) (model.PaymentMethodsOutput, error)
	List(ctx context.Context) ([]model.PaymentMethodsOutput, error)
	Create(ctx context.Context, entity model.PaymentMethodsCreate) (model.PaymentMethodsOutput, error)
	Update(ctx context.Context, id int64, entity model.PaymentMethodsCreate) (model.PaymentMethodsOutput, error)
	Delete(ctx context.Context, id int64) error
}

type PaymentMethodsHandler struct {
	service PaymentMethodsService
}

func NewPaymentMethodsHandler(service PaymentMethodsService) PaymentMethodsHandler {
	return PaymentMethodsHandler{service: service}
}

func (h PaymentMethodsHandler) List(w http.ResponseWriter, r *http.Request) error {
	dto, err := h.service.List(r.Context())
	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return web.RespondJSON(w, dto, http.StatusOK)
}

func (h PaymentMethodsHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id, err := web.Params(r).Int("id")
	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	dto, err := h.service.Get(r.Context(), int64(id))
	if err != nil {
		return web.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if dto == (model.PaymentMethodsOutput{}) {
		return web.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return web.RespondJSON(w, dto, http.StatusOK)
}

func (h PaymentMethodsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	dto := model.PaymentMethodsCreate{}
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

func (h PaymentMethodsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	dto := model.PaymentMethodsCreate{}
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

func (h PaymentMethodsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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
