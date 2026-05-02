package pricing

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/shopspring/decimal"

	"golang_pricingengine/internal/config"
	"golang_pricingengine/internal/models"
	"golang_pricingengine/internal/repo"
	"golang_pricingengine/internal/validate"
)

var ErrNotImplemented = errors.New("this endpoint requires EF Core / WireTransac data paths not ported to Go")

type Service struct {
	cfg *config.Config
	r   *repo.Repo
}

func New(cfg *config.Config, r *repo.Repo) *Service {
	return &Service{cfg: cfg, r: r}
}

func (s *Service) GetFees(ctx context.Context, m *models.PriceEngineModel) (*models.OutputFeeModel, error) {
	if err := validate.Join(validate.OriginAmountRules(m)); err != nil {
		return nil, err
	}
	setFixedAmounts(m)
	raw, err := s.fetchFeesRaw(ctx, m, feesOpts{formsV2: formsFromModel})
	if err != nil {
		return nil, err
	}
	filtered := repo.FilterFeeRows(raw, m.OriginStateCode)
	if len(filtered) == 0 || allRatesMissing(filtered) {
		return nil, nil
	}
	return buildOutputFee(filtered, m)
}

func (s *Service) GetExchangeRates(ctx context.Context, m *models.PriceEngineModel) (*models.OutputExchangeRateModel, error) {
	if err := validate.Join(validate.OriginAmountRules(m)); err != nil {
		return nil, err
	}
	setFixedAmounts(m)
	raw, err := s.fetchFeesRaw(ctx, m, feesOpts{formsV2: formsFromModel})
	if err != nil {
		return nil, err
	}
	filtered := repo.FilterFeeRows(raw, m.OriginStateCode)
	if len(filtered) == 0 || allRatesMissing(filtered) {
		return nil, nil
	}
	row := pickBestRateRow(filtered, m.FormsOfPaymentsTypeID)
	if row == nil {
		return nil, nil
	}
	return buildOutputExchange(row, m), nil
}

func (s *Service) GetFeesExchangeRates(ctx context.Context, m *models.PriceEngineModel) (*models.OutputFeesRatesModel, error) {
	if err := validate.Join(validate.OriginAmountRules(m)); err != nil {
		return nil, err
	}
	setFixedAmounts(m)
	raw, err := s.fetchFeesRaw(ctx, m, feesOpts{formsV2: formsNil})
	if err != nil {
		return nil, err
	}
	filtered := repo.FilterFeeRows(raw, m.OriginStateCode)
	if len(filtered) == 0 || allRatesMissing(filtered) {
		return nil, nil
	}
	if !s.cfg.UseV2ForPartner(m.PartnerID) && m.FormsOfPaymentsTypeID != nil {
		ok := false
		for _, r := range filtered {
			if r.IsFormOfPaymentAvailable() == 1 && r.FormsOfPaymentsTypeID() == *m.FormsOfPaymentsTypeID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, nil
		}
	}
	fee, err := buildOutputFee(filtered, m)
	if err != nil {
		return nil, err
	}
	row := pickBestRateRow(filtered, m.FormsOfPaymentsTypeID)
	if row == nil {
		return nil, nil
	}
	ex := buildOutputExchange(row, m)
	return &models.OutputFeesRatesModel{OutputFeeModel: fee, OutputExchangeRate: ex}, nil
}

func (s *Service) GetDeliveryMethodAndStyleOptions(ctx context.Context, in *models.InputDeliveryMethodStylesModel) (*models.OutputDeliveryMethodStylesModel, error) {
	if err := validate.Join(validate.DeliveryMethodsRules(in)); err != nil {
		return nil, err
	}
	if !s.cfg.DeliveryMethodsUseSP {
		return nil, ErrNotImplemented
	}
	originAmt := 1.0
	raw, err := s.r.DeliveryMethodsAndStyleOptions(ctx, in.PartnerID, &originAmt,
		in.OriginCountryName, in.OriginCountryCode, in.OriginStateCode,
		in.DestinationCountryName, in.DestinationCurrencyCode, in.DestinationCountryCode)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	styleSeen := map[int]struct{}{}
	var styles []models.StyleOptionModel
	recvSeen := map[int]struct{}{}
	var recvs []receivingOpt
	for _, row := range raw {
		sid := repo.MapInt(row, "styleid", "StyleId")
		if _, ok := styleSeen[sid]; !ok && sid != 0 {
			styleSeen[sid] = struct{}{}
			styles = append(styles, models.StyleOptionModel{
				StyleID:     sid,
				StyleName:   styleName(sid),
				IsAvailable: true,
			})
		}
		rid := repo.MapInt(row, "receivingoptionid", "ReceivingOptionId")
		if _, ok := recvSeen[rid]; !ok && rid != 0 {
			recvSeen[rid] = struct{}{}
			name, code := receivingOptionEnum(rid)
			if name != "" {
				recvs = append(recvs, receivingOpt{name: strings.ToUpper(name), code: code})
			}
		}
	}
	sort.Slice(styles, func(i, j int) bool { return styles[i].StyleID < styles[j].StyleID })
	delivery := buildDeliveryMethodsList(recvs)
	return &models.OutputDeliveryMethodStylesModel{
		StyleOptionList:     styles,
		DeliveryMethodsList: delivery,
	}, nil
}

func (s *Service) GetOriginAmount(ctx context.Context, m *models.PriceEngineModel) (*models.OutputOriginAmountModel, error) {
	if err := validate.Join(validate.CalculateRules(m)); err != nil {
		return nil, err
	}
	setFixedAmounts(m)
	if !s.cfg.OriginAmountUseSP {
		return nil, ErrNotImplemented
	}
	v2 := s.cfg.UseV2ForPartner(m.PartnerID)
	var raw []map[string]interface{}
	var err error
	if v2 {
		raw, err = s.r.OriAmountCalculatorV2(ctx, m.PartnerID, m.DestinationAmount,
			m.OriginCountryName, m.OriginCountryCode, m.OriginStateCode,
			m.DestinationCountryName, m.DestinationCurrencyCode, m.DestinationCountryCode,
			m.StyleID, m.ReceivingOptionID, m.PayerID, m.FormsOfPaymentsTypeID)
	} else {
		br := s.cfg.BaseRateOption(m.DestinationCountryName, m.DestinationCurrencyCode)
		raw, err = s.r.OriAmountCalculator(ctx, m.PartnerID, m.DestinationAmount,
			m.OriginCountryName, m.OriginCountryCode, m.OriginStateCode,
			m.DestinationCountryName, m.DestinationCurrencyCode, m.DestinationCountryCode,
			m.StyleID, m.ReceivingOptionID, m.PayerID, m.FormsOfPaymentsTypeID,
			br, m.IsEmployee)
	}
	if err != nil {
		return nil, err
	}
	type scored struct {
		row repo.FeeSPRow
		rate float64
	}
	var list []scored
	for _, row := range raw {
		fr := repo.FeeSPRow{Map: row}
		oa := fr.OriginAmount()
		rt := fr.Rate()
		if oa == nil || rt == nil {
			continue
		}
		list = append(list, scored{row: fr, rate: *rt})
	}
	if len(list) == 0 {
		return nil, nil
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rate > list[j].rate })
	best := list[0].row
	oa := best.OriginAmount()
	up := best.UpToAmount()
	if oa == nil {
		return nil, nil
	}
	if up != nil && *oa > *up {
		return nil, fmt.Errorf("destinationAmount is too large for configured bands")
	}
	return &models.OutputOriginAmountModel{
		OriginAmount: roundDecimalAway(*oa, 2),
	}, nil
}

func (s *Service) GetExchangeRatesByPayer(ctx context.Context, m *models.PriceEngineModel) ([]models.OutputExchangeRateByPayerModel, error) {
	if err := validate.Join(validate.ExchangeRateByPayerRules(m)); err != nil {
		return nil, err
	}
	if !s.cfg.ExchangeRateByPayerUseSP {
		return nil, ErrNotImplemented
	}
	setFixedAmounts(m)
	mm := *m
	if s.cfg.UseV2ForPartner(m.PartnerID) {
		if mm.OriginAmount == nil || *mm.OriginAmount == 0 {
			if mm.DestinationAmount == nil {
				return nil, nil
			}
			oa, err := s.GetOriginAmount(ctx, &mm)
			if err != nil || oa == nil {
				return nil, err
			}
			v := oa.OriginAmount
			mm.OriginAmount = &v
		}
		raw, err := s.r.ExchangesRatesByPayersV2(ctx, mm.PartnerID, mm.OriginAmount,
			mm.OriginCountryName, mm.OriginCountryCode, mm.OriginStateCode,
			mm.DestinationCountryName, mm.DestinationCurrencyCode, mm.DestinationCountryCode,
			mm.StyleID, mm.ReceivingOptionID, mm.PayerID, mm.FormsOfPaymentsTypeID,
			mm.DestinationState, mm.DestinationCity, mm.DestinationAmount)
		if err != nil {
			return nil, err
		}
		return mapExchangeByPayerV2(raw, mm.OriginStateCode, *mm.OriginAmount)
	}
	if mm.OriginAmount == nil || *mm.OriginAmount == 0 {
		if mm.DestinationAmount == nil {
			return nil, nil
		}
		return nil, ErrNotImplemented
	}
	br := s.cfg.BaseRateOption(mm.DestinationCountryName, mm.DestinationCurrencyCode)
	raw, err := s.r.ExchangesRatesByPayers(ctx, mm.PartnerID, mm.OriginAmount,
		mm.OriginCountryName, mm.OriginCountryCode, mm.OriginStateCode,
		mm.DestinationCountryName, mm.DestinationCurrencyCode, mm.DestinationCountryCode,
		mm.StyleID, mm.ReceivingOptionID, mm.PayerID, mm.FormsOfPaymentsTypeID,
		br, mm.IsEmployee)
	if err != nil {
		return nil, err
	}
	return mapExchangeByPayerV1(raw, mm.OriginStateCode)
}

func (s *Service) GetBestRatesByCountry(_ context.Context, m *models.BestRateModel) (*models.OutputExchangeRateByPayerModel, error) {
	if err := validate.Join(validate.BestRateRules(m)); err != nil {
		return nil, err
	}
	return nil, ErrNotImplemented
}

func (s *Service) GetCountriesAndCurrencies(_ context.Context) error {
	return ErrNotImplemented
}

// --- helpers ---

type feesOpts struct {
	formsV2 formsMode
}

type formsMode int

const (
	formsNil formsMode = iota
	formsFromModel
)

func (s *Service) fetchFeesRaw(ctx context.Context, m *models.PriceEngineModel, opt feesOpts) ([]map[string]interface{}, error) {
	v2 := s.cfg.UseV2ForPartner(m.PartnerID)
	var forms *int
	if v2 && opt.formsV2 == formsFromModel {
		forms = m.FormsOfPaymentsTypeID
	}
	if v2 {
		return s.r.FeesAndExchangesByFormOfPaymentV2(ctx, m.PartnerID, m.OriginAmount,
			m.OriginCountryName, m.OriginCountryCode, m.OriginStateCode,
			m.DestinationCountryName, m.DestinationCurrencyCode, m.DestinationCountryCode,
			m.StyleID, m.ReceivingOptionID, m.PayerID, forms, m.IsEmployee)
	}
	br := s.cfg.BaseRateOption(m.DestinationCountryName, m.DestinationCurrencyCode)
	return s.r.FeesAndExchangesByFormOfPayment(ctx, m.PartnerID, m.OriginAmount,
		m.OriginCountryName, m.OriginCountryCode, m.OriginStateCode,
		m.DestinationCountryName, m.DestinationCurrencyCode, m.DestinationCountryCode,
		m.StyleID, m.ReceivingOptionID, forms, m.PayerID,
		br, m.IsEmployee)
}

func setFixedAmounts(m *models.PriceEngineModel) {
	if m.OriginAmount != nil {
		v := roundTrunc(*m.OriginAmount, 2)
		m.OriginAmount = &v
	}
	if m.DestinationAmount != nil {
		v := roundTrunc(*m.DestinationAmount, 2)
		m.DestinationAmount = &v
	}
}

func roundTrunc(f float64, places int32) float64 {
	d := decimal.NewFromFloat(f)
	return d.Truncate(places).InexactFloat64()
}

func roundDecimalAway(f float64, places int32) float64 {
	d := decimal.NewFromFloat(f)
	return d.Round(places).InexactFloat64()
}

func roundTowardZero(f float64, places int32) float64 {
	d := decimal.NewFromFloat(f)
	return d.RoundDown(places).InexactFloat64()
}

func allRatesMissing(rows []repo.FeeSPRow) bool {
	for _, r := range rows {
		if r.Rate() != nil {
			return false
		}
	}
	return true
}

func pickBestRateRow(rows []repo.FeeSPRow, formsID *int) *repo.FeeSPRow {
	var candidates []repo.FeeSPRow
	for _, r := range rows {
		if r.IsFormOfPaymentAvailable() != 1 || r.Rate() == nil {
			continue
		}
		if formsID != nil && r.FormsOfPaymentsTypeID() != *formsID {
			continue
		}
		candidates = append(candidates, r)
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return *candidates[i].Rate() > *candidates[j].Rate()
	})
	return &candidates[0]
}

func buildOutputExchange(row *repo.FeeSPRow, m *models.PriceEngineModel) *models.OutputExchangeRateModel {
	origin := roundTowardZero(*m.OriginAmount, 2)
	rate := roundTowardZero(*row.Rate(), 4)
	var rbPtr *float64
	if row.RateBase() != nil {
		rb := roundTowardZero(*row.RateBase(), 4)
		rbPtr = &rb
	}
	dest := roundTowardZero(rate*(*m.OriginAmount), 2)
	out := models.OutputExchangeRateModel{
		OutputExchangeRateByPayerModel: models.OutputExchangeRateByPayerModel{
			OriginAmount: origin,
			Rate:         fp(rate),
			RateBase:     rbPtr,
			PayerCode:    row.PayerCodeExchangeRate(),
			PayerID:      row.PayerIDExchangeRate(),
			DestinationAmount: dest,
		},
		StyleID: row.StyleID(),
	}
	return &out
}

func fp(f float64) *float64 { return &f }

func buildOutputFee(rows []repo.FeeSPRow, m *models.PriceEngineModel) (*models.OutputFeeModel, error) {
	out := &models.OutputFeeModel{
		OriginAmount: roundTowardZero(*m.OriginAmount, 2),
		PaymentMethods: []models.OutputPaymentMethodModel{},
	}
	if m.FormsOfPaymentsTypeID != nil {
		for _, r := range rows {
			if r.FormsOfPaymentsTypeID() != *m.FormsOfPaymentsTypeID {
				continue
			}
			if r.Fee() != nil {
				out.FeeAmount = roundTowardZero(*r.Fee(), 2)
				if r.TotalAmount() != nil {
					out.TotalAmount = roundTowardZero(*r.TotalAmount(), 2)
				}
			} else if r.Rate() != nil {
				out.FeeAmount = 0
				out.TotalAmount = 0
			}
			break
		}
	}
	for _, pm := range []struct {
		id   int
		name string
	}{
		{1, "ACH"},
		{3, "DebitCard"},
		{2, "CreditCard"},
	} {
		var row *repo.FeeSPRow
		for i := range rows {
			if rows[i].FormsOfPaymentsTypeID() == pm.id {
				row = &rows[i]
				break
			}
		}
		pay, err := paymentMethodFromRow(row, pm.id, pm.name)
		if err != nil {
			return nil, err
		}
		out.PaymentMethods = append(out.PaymentMethods, pay)
	}
	return out, nil
}

func paymentMethodFromRow(row *repo.FeeSPRow, id int, name string) (models.OutputPaymentMethodModel, error) {
	if row == nil {
		return models.OutputPaymentMethodModel{
			FeeAmount:               0,
			IsAvailable:             false,
			SenderPaymentMethodID:   id,
			SenderPaymentMethodName: name,
			StyleID:                 0,
		}, nil
	}
	pm := models.OutputPaymentMethodModel{
		IsAvailable:             row.IsFormOfPaymentAvailable() == 1,
		SenderPaymentMethodID:   id,
		SenderPaymentMethodName: name,
		StyleID:                 row.StyleID(),
	}
	if row.Fee() != nil {
		pm.FeeAmount = roundTowardZero(*row.Fee(), 2)
		return pm, nil
	}
	if row.Rate() != nil {
		return pm, nil
	}
	return pm, fmt.Errorf("price definition not found for payment method %s", name)
}

func styleName(styleID int) string {
	switch styleID {
	case 3:
		return "BestRate"
	case 2:
		return "FastDelivery"
	case 1:
		return "LowerFee"
	default:
		return ""
	}
}

func receivingOptionEnum(id int) (string, int) {
	switch id {
	case 1:
		return "CashPickup", 1
	case 2:
		return "HomeDelivery", 2
	case 3:
		return "BankDeposit", 3
	case 4:
		return "DigitalWallet", 4
	case 5:
		return "DepositDebitCard", 5
	default:
		return "", 0
	}
}

type receivingOpt struct {
	name string
	code int
}

func buildDeliveryMethodsList(recvs []receivingOpt) []models.DeliveryMethodModel {
	var out []models.DeliveryMethodModel
	for _, item := range recvs {
		n := strings.ReplaceAll(item.name, " ", "")
		switch strings.ToUpper(n) {
		case "CASHPICKUP":
			out = append(out, models.DeliveryMethodModel{
				TransactionTypeID:   1,
				DeliveryMethod:      "Wire",
				TransactionTypeName: item.name,
			})
		case "BANKDEPOSIT":
			out = append(out, models.DeliveryMethodModel{
				TransactionTypeID:   3,
				DeliveryMethod:      "Wire",
				TransactionTypeName: item.name,
			})
		case "HOMEDELIVERY":
			out = append(out, models.DeliveryMethodModel{
				TransactionTypeID:   1,
				DeliveryMethod:      "Delivery",
				TransactionTypeName: item.name,
			})
		case "DEPOSITDEBITCARD":
			out = append(out, models.DeliveryMethodModel{
				TransactionTypeID:   3,
				DeliveryMethod:      "DepositDebitCard",
				TransactionTypeName: item.name,
			})
		}
	}
	return out
}

func mapExchangeByPayerV1(raw []map[string]interface{}, originState string) ([]models.OutputExchangeRateByPayerModel, error) {
	filtered := repo.FilterExchangePayerRows(raw, originState)
	var out []models.OutputExchangeRateByPayerModel
	for _, m := range filtered {
		r := repo.FeeSPRow{Map: m}
		code := r.PayerCodeExchangeRate()
		if strings.TrimSpace(code) == "" {
			continue
		}
		var rate *float64
		if r.Rate() != nil {
			rate = fp(roundTowardZero(*r.Rate(), 4))
		}
		var rb *float64
		if r.RateBase() != nil {
			rb = fp(roundTowardZero(*r.RateBase(), 4))
		}
		oa, _ := toFloat(repo.Pick(r.Map, "originamount", "OriginAmount"))
		out = append(out, models.OutputExchangeRateByPayerModel{
			OriginAmount:      oa,
			Rate:              rate,
			RateBase:          rb,
			PayerCode:         code,
			PayerID:           r.PayerIDExchangeRate(),
			DestinationAmount: 0,
		})
	}
	if len(out) == 0 {
		return nil, nil
	}
	seen := map[string]struct{}{}
	var dedup []models.OutputExchangeRateByPayerModel
	for _, x := range out {
		if _, ok := seen[x.PayerCode]; ok {
			continue
		}
		seen[x.PayerCode] = struct{}{}
		dedup = append(dedup, x)
	}
	return dedup, nil
}

func mapExchangeByPayerV2(raw []map[string]interface{}, originState string, originAmt float64) ([]models.OutputExchangeRateByPayerModel, error) {
	type row struct {
		payerCode string
		combined  float64
		m         map[string]interface{}
	}
	var list []row
	for _, m := range raw {
		state := repo.MapString(m, "originstatecode", "OriginStateCode")
		if strings.TrimSpace(originState) != "" {
			if !strings.EqualFold(strings.TrimSpace(state), strings.TrimSpace(originState)) {
				continue
			}
		} else {
			if !strings.EqualFold(strings.TrimSpace(state), "ALL") {
				continue
			}
		}
		rt := repo.MapFloatPtr(m, "rate", "Rate")
		if rt == nil {
			continue
		}
		bps := 0.0
		if v := repo.MapFloatPtr(m, "bpsfrombase", "BpsFromBase"); v != nil {
			bps = roundTowardZero(*v, 4)
		}
		afb := 0.0
		if v := repo.MapFloatPtr(m, "amountfromthebase", "AmountFromTheBase"); v != nil {
			afb = roundTowardZero(*v, 4)
		}
		pc := repo.MapString(m, "payercode", "PayerCode")
		list = append(list, row{payerCode: pc, combined: bps + afb, m: m})
	}
	if len(list) == 0 {
		return nil, nil
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].combined > list[j].combined })
	seen := map[string]struct{}{}
	var out []models.OutputExchangeRateByPayerModel
	for _, item := range list {
		if item.payerCode == "" {
			continue
		}
		if _, ok := seen[item.payerCode]; ok {
			continue
		}
		seen[item.payerCode] = struct{}{}
		rt := repo.MapFloatPtr(item.m, "rate", "Rate")
		fxb := repo.MapFloatPtr(item.m, "fxbase", "FxBase")
		pid := repo.MapIntPtr(item.m, "payerid", "PayerId")
		dest := 0.0
		if v := repo.MapFloatPtr(item.m, "destinationamount", "DestinationAmount"); v != nil {
			dest = *v
		}
		var rate *float64
		if rt != nil {
			rate = fp(roundTowardZero(*rt, 4))
		}
		var rb *float64
		if fxb != nil {
			rb = fp(roundTowardZero(*fxb, 4))
		}
		out = append(out, models.OutputExchangeRateByPayerModel{
			OriginAmount:      roundTowardZero(originAmt, 2),
			Rate:              rate,
			RateBase:          rb,
			PayerCode:         item.payerCode,
			PayerID:           pid,
			DestinationAmount: dest,
		})
	}
	return out, nil
}

func toFloat(v interface{}) (float64, bool) {
	return repo.ToFloat64(v)
}

