package grpt

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NumericFormatterOption = func(formatter *NumericFormatter)

type Formatter = func(text any) string

type TextFormatter interface {
	Format(text any) string
}

type FormattedText interface {
	Formatted() string
}

const NegativeParenthesesFormat = "(#)"

var registeredFormatters = map[TextType]TextFormatter{
	TextTypeDate: &DateTimeFormatter{
		Pattern: time.DateOnly,
	},
	TextTypeTime: &DateTimeFormatter{
		Pattern: time.TimeOnly,
	},
	TextTypeDateTime: &DateTimeFormatter{
		Pattern: time.DateTime,
	},
	TextTypeReal: &NumericFormatter{
		Precision: 2,
	},
	TextTypePercentage: &NumericFormatter{
		Suffix:    "%",
		Precision: 2,
	},
}

func RegisterFormatter(formatter TextFormatter, types ...TextType) bool {
	if formatter == nil {
		return false
	}

	for _, textType := range types {
		registeredFormatters[textType] = formatter
	}

	return true
}

func RegisteredFormatter(textType TextType) TextFormatter {
	if formatter, ok := registeredFormatters[textType]; ok {
		return formatter
	}
	return nil
}

type NumericFormatter struct {
	Prefix                  string
	Suffix                  string
	DecimalSeparator        string
	ThousandSeparator       string
	WithoutDecimalSeparator bool
	Precision               int
	Transform               func(float64) float64
	NegativeFormat          string
	HideSign                bool
}

func NewNumericFormatter(
	options ...NumericFormatterOption,
) *NumericFormatter {
	formatter := &NumericFormatter{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(formatter)
	}
	return formatter
}

func (n *NumericFormatter) Format(text any) string {
	if text == nil {
		return ""
	}

	number, err := n.parseText(text)
	if err != nil {
		return "NaN"
	}

	if n.Transform != nil {
		number = n.Transform(number)
	}

	isNegative := number < 0

	formatted := fmt.Sprintf("%."+strconv.Itoa(n.Precision)+"f", number)
	formatted = strings.ReplaceAll(formatted, "-", "")
	if len(n.DecimalSeparator) > 0 || n.WithoutDecimalSeparator {
		formatted = strings.ReplaceAll(formatted, ".", n.DecimalSeparator)
	}

	if len(n.ThousandSeparator) > 0 {
		formatted = n.addThousendSeparators(formatted)
	}

	if isNegative {
		if !n.HideSign {
			formatted = "-" + formatted
		}
		if len(n.NegativeFormat) > 0 {
			formatted = strings.Replace(n.NegativeFormat, "#", formatted, -1)
		}
	}

	formatted = n.Prefix + formatted + n.Suffix
	return formatted
}

func (n NumericFormatter) With(
	options ...NumericFormatterOption,
) *NumericFormatter {
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&n)
	}
	return &n
}

func (n NumericFormatter) WithPrefix(prefix string) *NumericFormatter {
	n.Prefix = prefix
	return &n
}

func (n NumericFormatter) WithSuffix(suffix string) *NumericFormatter {
	n.Suffix = suffix
	return &n
}

func (n NumericFormatter) WithDecimalSeparator(
	separator string,
) *NumericFormatter {
	n.DecimalSeparator = separator
	return &n
}

func (n NumericFormatter) WithThousandSeparator(
	separator string,
) *NumericFormatter {
	n.ThousandSeparator = separator
	return &n
}

func (n NumericFormatter) WithPrecision(precision int) *NumericFormatter {
	n.Precision = precision
	return &n
}

func (n NumericFormatter) WithNegativeFormat(format string) *NumericFormatter {
	n.NegativeFormat = format
	return &n
}

func (n NumericFormatter) WithHideSign(hide bool) *NumericFormatter {
	n.HideSign = hide
	return &n
}

func (n NumericFormatter) WithTransform(
	transform func(float64) float64,
) *NumericFormatter {
	n.Transform = transform
	return &n
}

func (n *NumericFormatter) parseText(text any) (float64, error) {
	switch v := text.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	}

	var raw any
	if valuer, ok := text.(driver.Valuer); ok {
		raw, _ = valuer.Value()
	} else {
		raw = text
	}

	if value, err := strconv.ParseFloat(fmt.Sprint(raw), 64); err == nil {
		return value, nil
	} else {
		return 0, errors.New("invalid number")
	}
}

func (n *NumericFormatter) addThousendSeparators(text string) string {
	var parts []string
	if len(n.DecimalSeparator) > 0 {
		parts = strings.Split(text, n.DecimalSeparator)
	} else {
		parts = strings.Split(text, ".")
	}

	var result strings.Builder
	integerPart := parts[0]
	integerLength := len(integerPart)
	for i, digit := range integerPart {
		if i > 0 && (integerLength-i)%3 == 0 {
			result.WriteString(n.ThousandSeparator)
		}
		result.WriteRune(digit)
	}

	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	if len(decimalPart) > 0 {
		return result.String() + n.DecimalSeparator + decimalPart
	}
	return result.String()
}

type DateTimeFormatter struct {
	SourcePatterns []string
	Pattern        string
}

func NewDateTimeFormatter(
	pattern string,
	sourcePatterns ...string,
) *DateTimeFormatter {
	formatter := &DateTimeFormatter{
		Pattern:        pattern,
		SourcePatterns: sourcePatterns,
	}

	return formatter
}

func (d *DateTimeFormatter) Format(text any) string {
	if text == nil {
		return ""
	}

	dateTime, err := d.parseText(text)
	if err != nil {
		return "N/D"
	}

	if dateTime.IsZero() {
		return ""
	}

	return dateTime.Format(d.Pattern)
}

func (d *DateTimeFormatter) parseText(text any) (time.Time, error) {
	var raw any
	if valuer, ok := text.(driver.Valuer); ok {
		raw, _ = valuer.Value()
	} else {
		raw = text
	}

	if raw == nil {
		return time.Time{}, nil
	}

	if value, ok := raw.(time.Time); ok {
		return value, nil
	}

	formats := d.SourcePatterns
	if len(formats) == 0 {
		formats = []string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
			time.RFC3339Nano,
			time.Kitchen,
			time.Stamp,
			time.StampMilli,
			time.StampMicro,
			time.StampNano,
			time.DateTime,
			time.DateOnly,
			time.TimeOnly,
		}
	}

	rawString := fmt.Sprint(raw)
	for _, format := range formats {
		if value, err := time.Parse(format, rawString); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("invalid date/time")
}
