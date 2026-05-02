package models

// PriceEngineModel mirrors Intermex.PricingEngine.API.Services.Models.PriceEngineModel query binding.
type PriceEngineModel struct {
	OriginCountryName       string
	OriginCountryCode       string
	OriginStateCode         string
	OriginCurrencyCode      string
	DestinationCountryName  string
	DestinationCurrencyCode string
	DestinationCountryCode  string
	PartnerID               int
	StyleID                 *int
	ReceivingOptionID       *int
	PayerID                 *int
	OriginAmount            *float64
	DestinationAmount       *float64
	FormsOfPaymentsTypeID   *int
	IsEmployee              bool

	GetPayersFromBanks *bool
	DestinationState   string
	DestinationCity    string
}

// InputDeliveryMethodStylesModel mirrors InputDeliveryMethodStylesModel.
type InputDeliveryMethodStylesModel struct {
	OriginCountryName       string
	OriginCountryCode       string
	OriginStateCode         string
	OriginCurrencyCode      string
	DestinationCountryName  string
	DestinationCountryCode  string
	DestinationCurrencyCode string
	PartnerID               int
}

// BestRateModel mirrors BestRateModel.
type BestRateModel struct {
	OriginCountryName       string
	OriginCountryCode       string
	OriginCurrencyCode      string
	DestinationCountryName  string
	DestinationCurrencyCode string
	DestinationCountryCode  string
	PartnerID               int
	StyleID                 int
	OriginAmount            *float64
}

// OutputFeeModel mirrors OutputFeeModel.
type OutputFeeModel struct {
	OriginAmount    float64                  `json:"originAmount"`
	FeeAmount       float64                  `json:"feeAmount"`
	TotalAmount     float64                  `json:"totalAmount"`
	PaymentMethods  []OutputPaymentMethodModel `json:"paymentMethods"`
}

// OutputPaymentMethodModel mirrors OutputPaymentMethodModel.
type OutputPaymentMethodModel struct {
	StyleID               int     `json:"styleId"`
	SenderPaymentMethodID int     `json:"senderPaymentMethodId"`
	SenderPaymentMethodName string `json:"senderPaymentMethodName"`
	IsAvailable           bool    `json:"isAvailable"`
	FeeAmount             float64 `json:"feeAmount"`
}

// OutputExchangeRateByPayerModel mirrors OutputExchangeRateByPayerModel (Newtonsoft default: camelCase).
type OutputExchangeRateByPayerModel struct {
	OriginAmount                     float64  `json:"originAmount"`
	Rate                             *float64 `json:"rate,omitempty"`
	RateBase                         *float64 `json:"rateBase,omitempty"`
	PayerCode                        string   `json:"payerCode,omitempty"`
	PayerID                          *int     `json:"payerId,omitempty"`
	IsParticipatingPayerForExchangeRate bool `json:"isParticipatingPayerForExchangeRate"`
	DestinationAmount                float64  `json:"destinationAmount"`
}

// OutputExchangeRateModel adds StyleId (inherits fields in JSON).
type OutputExchangeRateModel struct {
	OutputExchangeRateByPayerModel
	StyleID int `json:"styleId"`
}

// OutputFeesRatesModel mirrors OutputFeesRatesModel.
type OutputFeesRatesModel struct {
	OutputExchangeRate *OutputExchangeRateModel `json:"outputExchangeRate,omitempty"`
	OutputFeeModel     *OutputFeeModel          `json:"outputFeeModel,omitempty"`
}

// StyleOptionModel mirrors StyleOptionModel.
type StyleOptionModel struct {
	StyleID     int    `json:"styleId"`
	StyleName   string `json:"styleName"`
	IsAvailable bool   `json:"isAvailable"`
}

// DeliveryMethodModel mirrors DeliveryMethodModel.
type DeliveryMethodModel struct {
	TransactionTypeID   int    `json:"transactionTypeId"`
	DeliveryMethod      string `json:"deliveryMethod"`
	TransactionTypeName string `json:"transactionTypeName"`
}

// OutputDeliveryMethodStylesModel mirrors OutputDeliveryMethodStylesModel.
type OutputDeliveryMethodStylesModel struct {
	StyleOptionList      []StyleOptionModel    `json:"styleOptionList"`
	DeliveryMethodsList  []DeliveryMethodModel `json:"deliveryMethodsList"`
}

// OutputOriginAmountModel mirrors output for GetOriginAmount.
type OutputOriginAmountModel struct {
	OriginAmount float64 `json:"originAmount"`
}

// ErrorBody for JSON errors.
type ErrorBody struct {
	Error string `json:"error"`
}
