package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds runtime settings mirroring feature flags / app settings from the .NET API.
type Config struct {
	HTTPAddr string

	SQLServerConnectionString string

	// UseFeesV2 selects FeesAndExchangesByFormOfPayment_v2 / OriAmountCalculator_v2 / ExchangesRatesByPayers_v2 paths.
	UseFeesV2 bool
	// PartnerIDs forcing V2 when non-empty (subset behavior); if empty, UseFeesV2 applies to all partners.
	V2PartnerIDs map[int]struct{}

	BaseRateOptionOverride int // 0 = compute from treasury/simulator lists or default 3

	SimulatorPairs []pairKey // → base rate option 1
	TreasuryPairs  []pairKey // → base rate option 2

	DeliveryMethodsUseSP    bool
	ExchangeRateByPayerUseSP bool
	OriginAmountUseSP       bool
}

type pairKey struct {
	Country  string
	Currency string
}

func Load() (*Config, error) {
	c := &Config{
		HTTPAddr:                 getenv("HTTP_ADDR", ":8080"),
		SQLServerConnectionString: strings.TrimSpace(os.Getenv("SQLSERVER_CONNECTION_STRING")),
		UseFeesV2:                getenvBool("PRICING_USE_FEES_V2", false),
		BaseRateOptionOverride:   getenvInt("BASE_RATE_OPTION", 0),
		DeliveryMethodsUseSP:     getenvBool("DELIVERY_METHODS_USE_SP", true),
		ExchangeRateByPayerUseSP: getenvBool("EXCHANGE_RATE_BY_PAYER_USE_SP", true),
		OriginAmountUseSP:        getenvBool("ORIGIN_AMOUNT_USE_SP", true),
		SimulatorPairs:           parsePairs(os.Getenv("PRICING_SIMULATOR_PAIRS")),
		TreasuryPairs:            parsePairs(os.Getenv("PRICING_TREASURY_PAIRS")),
	}
	if s := strings.TrimSpace(os.Getenv("PRICING_V2_PARTNER_IDS")); s != "" {
		c.V2PartnerIDs = map[int]struct{}{}
		for _, p := range strings.Split(s, ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			id, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("PRICING_V2_PARTNER_IDS: invalid id %q", p)
			}
			c.V2PartnerIDs[id] = struct{}{}
		}
	}
	return c, nil
}

func (c *Config) UseV2ForPartner(partnerID int) bool {
	if len(c.V2PartnerIDs) > 0 {
		_, ok := c.V2PartnerIDs[partnerID]
		return ok
	}
	return c.UseFeesV2
}

func (c *Config) BaseRateOption(destCountry, destCurrency string) int {
	if c.BaseRateOptionOverride > 0 {
		return c.BaseRateOptionOverride
	}
	dc := strings.TrimSpace(strings.ToUpper(destCountry))
	dcur := strings.TrimSpace(strings.ToUpper(destCurrency))
	for _, p := range c.SimulatorPairs {
		if strings.ToUpper(p.Country) == dc && strings.ToUpper(p.Currency) == dcur {
			return 1
		}
	}
	for _, p := range c.TreasuryPairs {
		if strings.ToUpper(p.Country) == dc && strings.ToUpper(p.Currency) == dcur {
			return 2
		}
	}
	return 3
}

func parsePairs(s string) []pairKey {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []pairKey
	for _, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "|", 2)
		if len(kv) != 2 {
			continue
		}
		out = append(out, pairKey{
			Country:  strings.TrimSpace(kv[0]),
			Currency: strings.TrimSpace(kv[1]),
		})
	}
	return out
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvBool(k string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(k)))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes"
}

func getenvInt(k string, def int) int {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
