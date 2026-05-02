# Golang_PricingEngine

Puerto del API **PricingEngine** (.NET) a Go: mismas rutas bajo `/PricingEngine/...`, llamadas a los stored procedures `[dbo]` documentados en el proyecto original y la lógica principal basada en SP (sin EF Core, LaunchDarkly ni bases WireTransac).

## Requisitos

- Go 1.22+
- SQL Server con la base y los SP del PricingEngine original

## Configuración

Variables de entorno principales:

| Variable | Descripción |
|----------|-------------|
| `SQLSERVER_CONNECTION_STRING` | Cadena ADO para `go-mssqldb` (obligatoria). Ejemplo: `sqlserver://user:pass@localhost:1433?database=PricingEngine&encrypt=disable` |
| `HTTP_ADDR` | Dirección del servidor (por defecto `:8080`) |
| `PRICING_USE_FEES_V2` | `true` para usar `FeesAndExchangesByFormOfPayment_v2`, `OriAmountCalculator_v2`, `ExchangesRatesByPayers_v2` |
| `PRICING_V2_PARTNER_IDS` | Lista `1,2,3`: solo esos partners usan V2; si está vacío aplica `PRICING_USE_FEES_V2` a todos |
| `BASE_RATE_OPTION` | `1` simulador, `2` tesorería, `3` diario; `0` = automático según pares siguientes |
| `PRICING_SIMULATOR_PAIRS` | Pares `País|MON`;ej. `Mexico|MXN;USA|USD` → opción base 1 |
| `PRICING_TREASURY_PAIRS` | Igual para opción base 2 |
| `DELIVERY_METHODS_USE_SP` | `true` (defecto): `DeliveryMethodsAndStyleOptions` |
| `EXCHANGE_RATE_BY_PAYER_USE_SP` | `true` (defecto): SP por pagador |
| `ORIGIN_AMOUNT_USE_SP` | `true` (defecto): `OriAmountCalculator` / `_v2` |

## Ejecución

```bash
set SQLSERVER_CONNECTION_STRING=sqlserver://...
go run ./cmd/server
```

Salud: `GET /health`

## Alcance no portado

- `GET /PricingEngine/GetBestRatesByCountry` → **501** (requiere EF / tablas de definición de precios y tasas base como en .NET).
- `GET /PricingEngine/GetCountriesAndCurrencies` → **501** (requiere `PriceDefinitionService` / EF).
- Rama EF de `GetExchangeRatesByPayer` cuando `PRICING_USE_FEES_V2=false` y falta `OriginAmount`: el .NET usa otro flujo; aquí se devuelve **501** si hace falta calcular origen sin la ruta V2.
