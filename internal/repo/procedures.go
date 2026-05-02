package repo

import (
	"context"
	"database/sql"
	"strings"
)

// DeliveryMethodsAndStyleOptions calls [dbo].[DeliveryMethodsAndStyleOptions].
func (r *Repo) DeliveryMethodsAndStyleOptions(ctx context.Context, partnerID int, originAmount *float64,
	originCountryName, originCountryCode, originStateCode, destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[DeliveryMethodsAndStyleOptions] @PartnerId, @OriginAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode`
	var oa interface{}
	if originAmount != nil {
		oa = *originAmount
	}
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("OriginAmount", oa),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// FeesAndExchangesByFormOfPayment calls [dbo].[FeesAndExchangesByFormOfPayment].
func (r *Repo) FeesAndExchangesByFormOfPayment(ctx context.Context,
	partnerID int, originAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, formsOfPaymentsTypeID, payerID *int,
	baseRateOption int, isEmployee bool,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[FeesAndExchangesByFormOfPayment] @PartnerId, @OriginAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId, @BaseRateOption, @IsEmployee`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("OriginAmount", floatPtrIface(originAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
		sql.Named("BaseRateOption", baseRateOption),
		sql.Named("IsEmployee", isEmployee),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// FeesAndExchangesByFormOfPaymentV2 calls [dbo].[FeesAndExchangesByFormOfPayment_v2].
func (r *Repo) FeesAndExchangesByFormOfPaymentV2(ctx context.Context,
	partnerID int, originAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, payerID, formsOfPaymentsTypeID *int,
	isEmployee bool,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[FeesAndExchangesByFormOfPayment_v2] @PartnerId, @OriginAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId, @IsEmployee`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("OriginAmount", floatPtrIface(originAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
		sql.Named("IsEmployee", isEmployee),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// OriAmountCalculator calls [dbo].[OriAmountCalculator].
func (r *Repo) OriAmountCalculator(ctx context.Context,
	partnerID int, destinationAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, payerID, formsOfPaymentsTypeID *int,
	baseRateOption int, isEmployee bool,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[OriAmountCalculator] @PartnerId, @DestinationAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId, @BaseRateOption, @IsEmployee`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("DestinationAmount", floatPtrIface(destinationAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
		sql.Named("BaseRateOption", baseRateOption),
		sql.Named("IsEmployee", isEmployee),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// OriAmountCalculatorV2 calls [dbo].[OriAmountCalculator_v2].
func (r *Repo) OriAmountCalculatorV2(ctx context.Context,
	partnerID int, destinationAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, payerID, formsOfPaymentsTypeID *int,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[OriAmountCalculator_v2] @PartnerId, @DestinationAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("DestinationAmount", floatPtrIface(destinationAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// ExchangesRatesByPayers calls [dbo].[ExchangesRatesByPayers].
func (r *Repo) ExchangesRatesByPayers(ctx context.Context,
	partnerID int, originAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, payerID, formsOfPaymentsTypeID *int,
	baseRateOption int, isEmployee bool,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[ExchangesRatesByPayers] @PartnerId, @OriginAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId, @BaseRateOption, @IsEmployee`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("OriginAmount", floatPtrIface(originAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
		sql.Named("BaseRateOption", baseRateOption),
		sql.Named("IsEmployee", isEmployee),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

// ExchangesRatesByPayersV2 calls [dbo].[ExchangesRatesByPayers_v2].
func (r *Repo) ExchangesRatesByPayersV2(ctx context.Context,
	partnerID int, originAmount *float64,
	originCountryName, originCountryCode, originStateCode string,
	destinationCountryName, destinationCurrencyCode, destinationCountryCode string,
	styleID, receivingOptionID, payerID, formsOfPaymentsTypeID *int,
	destinationState, destinationCity string,
	destinationAmount *float64,
) ([]map[string]interface{}, error) {
	const q = `EXEC [dbo].[ExchangesRatesByPayers_v2] @PartnerId, @OriginAmount, @OriginCountryName, @OriginCountryCode, @OriginStateCode, @DestinationCountryName, @DestinationCurrencyCode, @DestinationCountryCode, @StyleId, @ReceivingOptionId, @PayerId, @FormsOfPaymentsTypeId, @DestinationState, @DestinationCity, @DestinationAmount`
	rows, err := r.db.QueryContext(ctx, q,
		sql.Named("PartnerId", partnerID),
		sql.Named("OriginAmount", floatPtrIface(originAmount)),
		sql.Named("OriginCountryName", nullStr(originCountryName)),
		sql.Named("OriginCountryCode", nullStr(originCountryCode)),
		sql.Named("OriginStateCode", nullStr(originStateCode)),
		sql.Named("DestinationCountryName", nullStr(destinationCountryName)),
		sql.Named("DestinationCurrencyCode", nullStr(destinationCurrencyCode)),
		sql.Named("DestinationCountryCode", nullStr(destinationCountryCode)),
		sql.Named("StyleId", intPtrIface(styleID)),
		sql.Named("ReceivingOptionId", intPtrIface(receivingOptionID)),
		sql.Named("PayerId", intPtrIface(payerID)),
		sql.Named("FormsOfPaymentsTypeId", intPtrIface(formsOfPaymentsTypeID)),
		sql.Named("DestinationState", nullStr(destinationState)),
		sql.Named("DestinationCity", nullStr(destinationCity)),
		sql.Named("DestinationAmount", floatPtrIface(destinationAmount)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func intPtrIface(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func floatPtrIface(p *float64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

// FeeSPRow normalizes a row from fees SP result sets.
type FeeSPRow struct {
	Map map[string]interface{}
}

func (row FeeSPRow) FormsOfPaymentsTypeID() int {
	return mapInt(row.Map, "formsofpaymentstypeid", "FormsOfPaymentsTypeId")
}

func (row FeeSPRow) IntermexDebitCard() bool {
	return row.FormsOfPaymentsTypeID() == 4
}

func (row FeeSPRow) OriginStateCode() string {
	return mapString(row.Map, "originstatecode", "OriginStateCode")
}

func (row FeeSPRow) Rate() *float64 {
	return mapFloatPtr(row.Map, "rate", "Rate")
}

func (row FeeSPRow) IsFormOfPaymentAvailable() int {
	return mapInt(row.Map, "isformofpaymentavailable", "IsFormOfPaymentAvailable")
}

func (row FeeSPRow) StyleID() int {
	return mapInt(row.Map, "styleid", "StyleId")
}

func (row FeeSPRow) PayerCodeExchangeRate() string {
	return mapString(row.Map, "payercodeexchangerate", "PayerCodeExchangeRate")
}

func (row FeeSPRow) PayerIDExchangeRate() *int {
	return mapIntPtr(row.Map, "payeridexchangerate", "PayerIdExchangeRate")
}

func (row FeeSPRow) RateBase() *float64 {
	return mapFloatPtr(row.Map, "ratebase", "RateBase")
}

func (row FeeSPRow) Fee() *float64 {
	return mapFloatPtr(row.Map, "fee", "Fee")
}

func (row FeeSPRow) TotalAmount() *float64 {
	return mapFloatPtr(row.Map, "totalamount", "TotalAmount")
}

func (row FeeSPRow) BPSFromBase() *float64 {
	return mapFloatPtr(row.Map, "bpsfrombase", "BPSFromBase")
}

func (row FeeSPRow) AmountFromTheBase() *float64 {
	return mapFloatPtr(row.Map, "amountfromthebase", "AmountFromTheBase")
}

func (row FeeSPRow) OriginAmount() *float64 {
	return mapFloatPtr(row.Map, "originamount", "OriginAmount")
}

func (row FeeSPRow) UpToAmount() *float64 {
	return mapFloatPtr(row.Map, "uptoamount", "UpToAmount")
}

// FilterFeeRows applies IntermexDebitCard exclusion and origin-state / ALL narrowing (fees SPs).
func FilterFeeRows(rows []map[string]interface{}, originState string) []FeeSPRow {
	const all = "ALL"
	var out []FeeSPRow
	for _, m := range rows {
		r := FeeSPRow{Map: m}
		if r.IntermexDebitCard() {
			continue
		}
		out = append(out, r)
	}
	if len(out) == 0 {
		return nil
	}
	var narrowed []FeeSPRow
	if originState != "" {
		for _, r := range out {
			if stringsEq(r.OriginStateCode(), originState) {
				narrowed = append(narrowed, r)
			}
		}
	}
	if len(narrowed) > 0 {
		out = narrowed
	} else {
		var allRows []FeeSPRow
		for _, r := range out {
			if stringsEq(r.OriginStateCode(), all) {
				allRows = append(allRows, r)
			}
		}
		out = allRows
	}
	return out
}

// FilterExchangePayerRows applies IsFormOfPaymentAvailable and origin-state / ALL narrowing (ExchangesRatesByPayers).
func FilterExchangePayerRows(rows []map[string]interface{}, originState string) []map[string]interface{} {
	var avail []map[string]interface{}
	for _, m := range rows {
		if mapInt(m, "isformofpaymentavailable", "IsFormOfPaymentAvailable") != 1 {
			continue
		}
		avail = append(avail, m)
	}
	if len(avail) == 0 {
		return nil
	}
	os := strings.TrimSpace(originState)
	if os != "" {
		var match []map[string]interface{}
		for _, m := range avail {
			if stringsEq(mapString(m, "originstatecode", "OriginStateCode"), os) {
				match = append(match, m)
			}
		}
		if len(match) > 0 {
			return match
		}
	}
	const all = "ALL"
	var allRows []map[string]interface{}
	for _, m := range avail {
		if stringsEq(mapString(m, "originstatecode", "OriginStateCode"), all) {
			allRows = append(allRows, m)
		}
	}
	return allRows
}

func stringsEq(a, b string) bool {
	return stringsTrimUpper(a) == stringsTrimUpper(b)
}

func stringsTrimUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
