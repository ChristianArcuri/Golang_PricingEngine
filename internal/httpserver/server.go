package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"golang_pricingengine/internal/models"
	"golang_pricingengine/internal/pricing"
)

type Server struct {
	svc *pricing.Service
}

func New(svc *pricing.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Get("/health", s.health)

	r.Route("/PricingEngine", func(r chi.Router) {
		r.Get("/GetDeliveryMethodAndStyleOptions", s.getDeliveryMethodAndStyleOptions)
		r.Get("/GetFees", s.getFees)
		r.Get("/GetExchangeRates", s.getExchangeRates)
		r.Get("/GetExchangeRatesByPayer", s.getExchangeRatesByPayer)
		r.Get("/GetOriginAmount", s.getOriginAmount)
		r.Get("/GetFeesExchangeRates", s.getFeesExchangeRates)
		r.Get("/GetBestRatesByCountry", s.getBestRatesByCountry)
		r.Get("/GetCountriesAndCurrencies", s.getCountriesAndCurrencies)
	})
	return r
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) getDeliveryMethodAndStyleOptions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	in := &models.InputDeliveryMethodStylesModel{
		OriginCountryName:       qGet(q, "OriginCountryName", "originCountryName"),
		OriginCountryCode:       qGet(q, "OriginCountryCode", "originCountryCode"),
		OriginStateCode:         qGet(q, "OriginStateCode", "originStateCode"),
		OriginCurrencyCode:      qGet(q, "OriginCurrencyCode", "originCurrencyCode"),
		DestinationCountryName:  qGet(q, "DestinationCountryName", "destinationCountryName"),
		DestinationCountryCode:  qGet(q, "DestinationCountryCode", "destinationCountryCode"),
		DestinationCurrencyCode: qGet(q, "DestinationCurrencyCode", "destinationCurrencyCode"),
		PartnerID:               parseInt(q, "PartnerId", "partnerId"),
	}
	out, err := s.svc.GetDeliveryMethodAndStyleOptions(r.Context(), in)
	s.respondPricing(w, out, err)
}

func priceEngineFromQuery(q url.Values) *models.PriceEngineModel {
	return &models.PriceEngineModel{
		OriginCountryName:       qGet(q, "OriginCountryName", "originCountryName"),
		OriginCountryCode:       qGet(q, "OriginCountryCode", "originCountryCode"),
		OriginStateCode:         qGet(q, "OriginStateCode", "originStateCode"),
		OriginCurrencyCode:      qGet(q, "OriginCurrencyCode", "originCurrencyCode"),
		DestinationCountryName:  qGet(q, "DestinationCountryName", "destinationCountryName"),
		DestinationCurrencyCode: qGet(q, "DestinationCurrencyCode", "destinationCurrencyCode"),
		DestinationCountryCode:  qGet(q, "DestinationCountryCode", "destinationCountryCode"),
		PartnerID:               parseInt(q, "PartnerId", "partnerId"),
		StyleID:                 parseIntPtr(q, "StyleId", "styleId"),
		ReceivingOptionID:       parseIntPtr(q, "ReceivingOptionId", "receivingOptionId"),
		PayerID:                 parseIntPtr(q, "PayerId", "payerId"),
		OriginAmount:            parseFloatPtr(q, "OriginAmount", "originAmount"),
		DestinationAmount:       parseFloatPtr(q, "DestinationAmount", "destinationAmount"),
		FormsOfPaymentsTypeID:   parseIntPtr(q, "FormsOfPaymentsTypeId", "formsOfPaymentsTypeId"),
		IsEmployee:              parseBool(q, "IsEmployee", "isEmployee"),
		GetPayersFromBanks:      nil,
		DestinationState:        qGet(q, "DestinationState", "destinationState"),
		DestinationCity:         qGet(q, "DestinationCity", "destinationCity"),
	}
}

func (s *Server) getFees(w http.ResponseWriter, r *http.Request) {
	m := priceEngineFromQuery(r.URL.Query())
	out, err := s.svc.GetFees(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getExchangeRates(w http.ResponseWriter, r *http.Request) {
	m := priceEngineFromQuery(r.URL.Query())
	out, err := s.svc.GetExchangeRates(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getFeesExchangeRates(w http.ResponseWriter, r *http.Request) {
	m := priceEngineFromQuery(r.URL.Query())
	out, err := s.svc.GetFeesExchangeRates(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getOriginAmount(w http.ResponseWriter, r *http.Request) {
	m := priceEngineFromQuery(r.URL.Query())
	out, err := s.svc.GetOriginAmount(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getExchangeRatesByPayer(w http.ResponseWriter, r *http.Request) {
	m := priceEngineFromQuery(r.URL.Query())
	out, err := s.svc.GetExchangeRatesByPayer(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getBestRatesByCountry(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	m := &models.BestRateModel{
		OriginCountryName:       qGet(q, "OriginCountryName", "originCountryName"),
		OriginCountryCode:       qGet(q, "OriginCountryCode", "originCountryCode"),
		OriginCurrencyCode:      qGet(q, "OriginCurrencyCode", "originCurrencyCode"),
		DestinationCountryName:  qGet(q, "DestinationCountryName", "destinationCountryName"),
		DestinationCurrencyCode: qGet(q, "DestinationCurrencyCode", "destinationCurrencyCode"),
		DestinationCountryCode:  qGet(q, "DestinationCountryCode", "destinationCountryCode"),
		PartnerID:               parseInt(q, "PartnerId", "partnerId"),
		StyleID:                 parseInt(q, "StyleId", "styleId"),
		OriginAmount:            parseFloatPtr(q, "OriginAmount", "originAmount"),
	}
	out, err := s.svc.GetBestRatesByCountry(r.Context(), m)
	s.respondPricing(w, out, err)
}

func (s *Server) getCountriesAndCurrencies(w http.ResponseWriter, r *http.Request) {
	err := s.svc.GetCountriesAndCurrencies(r.Context())
	if errors.Is(err, pricing.ErrNotImplemented) {
		writeJSON(w, http.StatusNotImplemented, models.ErrorBody{Error: err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorBody{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (s *Server) respondPricing(w http.ResponseWriter, out interface{}, err error) {
	if errors.Is(err, pricing.ErrNotImplemented) {
		writeJSON(w, http.StatusNotImplemented, models.ErrorBody{Error: err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorBody{Error: err.Error()})
		return
	}
	if out == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}
