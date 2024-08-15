package grpt

type Currency string

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
	EUR Currency = "EUR"
	JPY Currency = "JPY"
	GBP Currency = "GBP"
	AUD Currency = "AUD"
	CAD Currency = "CAD"
	CHF Currency = "CHF"
	CNY Currency = "CNY"
	SEK Currency = "SEK"
	NZD Currency = "NZD"
	SGD Currency = "SGD"
	MXN Currency = "MXN"
	INR Currency = "INR"
	ZAR Currency = "ZAR"
	RUB Currency = "RUB"
	KRW Currency = "KRW"
	TRY Currency = "TRY"
	TWD Currency = "TWD"
	HKD Currency = "HKD"
	ILS Currency = "ILS"
	ARS Currency = "ARS"
	COP Currency = "COP"
	MYR Currency = "MYR"
	PHP Currency = "PHP"
	IDR Currency = "IDR"
	PLN Currency = "PLN"
	HUF Currency = "HUF"
	SAR Currency = "SAR"
	AED Currency = "AED"
	EGP Currency = "EGP"
	CLP Currency = "CLP"
	PEN Currency = "PEN"
)

func CurrencyFormatter(currency Currency) *NumericFormatter {
	if formatter, ok := currencyFormatters[currency]; ok {
		return &formatter
	}
	return &NumericFormatter{Precision: 2}
}

var currencyFormatters = map[Currency]NumericFormatter{
	BRL: {
		Prefix:                  "R$ ",
		Suffix:                  "",
		DecimalSeparator:        ",",
		ThousandSeparator:       ".",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	USD: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	EUR: {
		Prefix:                  "",
		Suffix:                  " €",
		DecimalSeparator:        ",",
		ThousandSeparator:       ".",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          NegativeParenthesesFormat,
		HideSign:                false,
	},
	JPY: {
		Prefix:                  "¥",
		Suffix:                  "",
		DecimalSeparator:        "",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: true,
		Precision:               0,
		NegativeFormat:          "",
		HideSign:                false,
	},
	GBP: {
		Prefix:                  "£ ",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	AUD: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	CAD: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	CHF: {
		Prefix:                  "CHF ",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	CNY: {
		Prefix:                  "¥",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	SEK: {
		Prefix:                  "kr",
		Suffix:                  "",
		DecimalSeparator:        ",",
		ThousandSeparator:       " ",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          NegativeParenthesesFormat,
		HideSign:                false,
	},
	NZD: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	SGD: {
		Prefix:                  "S$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	MXN: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	INR: {
		Prefix:                  "₹",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	ZAR: {
		Prefix:                  "R",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	RUB: {
		Prefix:                  "₽",
		Suffix:                  "",
		DecimalSeparator:        ",",
		ThousandSeparator:       " ",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	KRW: {
		Prefix:                  "₩",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               0,
		NegativeFormat:          "",
		HideSign:                false,
	},
	TRY: {
		Prefix:                  "₺",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	TWD: {
		Prefix:                  "NT$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	HKD: {
		Prefix:                  "HK$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	ILS: {
		Prefix:                  "₪",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	ARS: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ",",
		ThousandSeparator:       ".",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	COP: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	MYR: {
		Prefix:                  "RM",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	PHP: {
		Prefix:                  "₱",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	IDR: {
		Prefix:                  "Rp",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	PLN: {
		Prefix:                  "zł",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	HUF: {
		Prefix:                  "Ft",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	SAR: {
		Prefix:                  "ر.س",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	AED: {
		Prefix:                  "د.إ",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	EGP: {
		Prefix:                  "ج.م",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	CLP: {
		Prefix:                  "$",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
	PEN: {
		Prefix:                  "S/",
		Suffix:                  "",
		DecimalSeparator:        ".",
		ThousandSeparator:       ",",
		WithoutDecimalSeparator: false,
		Precision:               2,
		NegativeFormat:          "",
		HideSign:                false,
	},
}
