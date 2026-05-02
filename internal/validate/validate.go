package validate

import (
	"fmt"
	"strings"

	"golang_pricingengine/internal/models"
)

func OriginAmountRules(m *models.PriceEngineModel) []string {
	var e []string
	add := func(cond bool, msg string) {
		if cond {
			e = append(e, msg)
		}
	}
	add(m.OriginAmount == nil, "OriginAmount is required")
	add(strings.TrimSpace(m.OriginCountryCode) == "", "OriginCountryCode is required")
	add(strings.TrimSpace(m.OriginCountryName) == "", "OriginCountryName is required")
	add(strings.TrimSpace(m.OriginCurrencyCode) == "", "OriginCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryName) == "", "DestinationCountryName is required")
	add(strings.TrimSpace(m.DestinationCountryCode) == "", "DestinationCountryCode is required")
	add(strings.TrimSpace(m.DestinationCurrencyCode) == "", "DestinationCurrencyCode is required")
	add(m.PartnerID <= 0, "PartnerId must be greater than 0")
	if m.ReceivingOptionID != nil {
		add(*m.ReceivingOptionID < 1 || *m.ReceivingOptionID > 5, "Invalid ReceivingOptionId")
	}
	if m.StyleID != nil {
		add(*m.StyleID < 1 || *m.StyleID > 3, "Invalid StyleId")
	}
	return e
}

func CalculateRules(m *models.PriceEngineModel) []string {
	var e []string
	add := func(cond bool, msg string) {
		if cond {
			e = append(e, msg)
		}
	}
	add(m.DestinationAmount == nil, "DestinationAmount is required")
	add(strings.TrimSpace(m.OriginCountryName) == "", "OriginCountryName is required")
	add(strings.TrimSpace(m.OriginCountryCode) == "", "OriginCountryCode is required")
	add(strings.TrimSpace(m.OriginCurrencyCode) == "", "OriginCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryName) == "", "DestinationCountryName is required")
	add(strings.TrimSpace(m.DestinationCountryCode) == "", "DestinationCountryCode is required")
	add(strings.TrimSpace(m.DestinationCurrencyCode) == "", "DestinationCurrencyCode is required")
	add(m.PartnerID <= 0, "PartnerId must be greater than 0")
	if m.ReceivingOptionID != nil {
		add(*m.ReceivingOptionID < 1 || *m.ReceivingOptionID > 4, "Invalid ReceivingOptionId")
	}
	if m.StyleID != nil {
		add(*m.StyleID < 1 || *m.StyleID > 3, "Invalid StyleId")
	}
	return e
}

func ExchangeRateByPayerRules(m *models.PriceEngineModel) []string {
	var e []string
	add := func(cond bool, msg string) {
		if cond {
			e = append(e, msg)
		}
	}
	add(strings.TrimSpace(m.OriginCountryName) == "", "OriginCountryName is required")
	add(strings.TrimSpace(m.OriginCountryCode) == "", "OriginCountryCode is required")
	add(strings.TrimSpace(m.OriginCurrencyCode) == "", "OriginCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryName) == "", "DestinationCountryName is required")
	add(strings.TrimSpace(m.DestinationCountryCode) == "", "DestinationCountryCode is required")
	add(strings.TrimSpace(m.DestinationCurrencyCode) == "", "DestinationCurrencyCode is required")
	add(m.PartnerID <= 0, "PartnerId must be greater than 0")
	oa := m.OriginAmount == nil || (m.OriginAmount != nil && *m.OriginAmount == 0)
	da := m.DestinationAmount == nil || (m.DestinationAmount != nil && *m.DestinationAmount == 0)
	add(oa && da, "Either OriginAmount or DestinationAmount is required")
	if m.ReceivingOptionID != nil {
		add(*m.ReceivingOptionID < 1 || *m.ReceivingOptionID > 5, "Invalid ReceivingOptionId")
	}
	if m.StyleID != nil {
		add(*m.StyleID < 1 || *m.StyleID > 3, "Invalid StyleId")
	}
	return e
}

func DeliveryMethodsRules(m *models.InputDeliveryMethodStylesModel) []string {
	var e []string
	add := func(cond bool, msg string) {
		if cond {
			e = append(e, msg)
		}
	}
	add(strings.TrimSpace(m.OriginCountryName) == "", "OriginCountryName is required")
	add(strings.TrimSpace(m.OriginStateCode) == "", "OriginStateCode is required")
	add(strings.TrimSpace(m.OriginCountryCode) == "", "OriginCountryCode is required")
	add(strings.TrimSpace(m.OriginCurrencyCode) == "", "OriginCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryName) == "", "DestinationCountryName is required")
	add(strings.TrimSpace(m.DestinationCurrencyCode) == "", "DestinationCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryCode) == "", "DestinationCountryCode is required")
	add(m.PartnerID <= 0, "PartnerId must be greater than 0")
	return e
}

func BestRateRules(m *models.BestRateModel) []string {
	var e []string
	add := func(cond bool, msg string) {
		if cond {
			e = append(e, msg)
		}
	}
	add(strings.TrimSpace(m.OriginCountryName) == "", "OriginCountryName is required")
	add(strings.TrimSpace(m.OriginCountryCode) == "", "OriginCountryCode is required")
	add(strings.TrimSpace(m.OriginCurrencyCode) == "", "OriginCurrencyCode is required")
	add(strings.TrimSpace(m.DestinationCountryName) == "", "DestinationCountryName is required")
	add(strings.TrimSpace(m.DestinationCountryCode) == "", "DestinationCountryCode is required")
	add(strings.TrimSpace(m.DestinationCurrencyCode) == "", "DestinationCurrencyCode is required")
	add(m.PartnerID <= 0, "PartnerId must be greater than 0")
	add(m.StyleID != 3, "StyleId must be 3 (BestRate)")
	add(m.OriginAmount == nil, "OriginAmount is required")
	return e
}

func Join(errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "; "))
}
