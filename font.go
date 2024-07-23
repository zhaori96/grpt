package grpt

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

//go:embed assets/fonts
var fonts embed.FS

var (
	standardFontFamilies = []FontFamily{
		{
			Name:       "roboto",
			Source:     fonts,
			Path:       "assets/fonts/roboto",
			Regular:    "Roboto-Regular.ttf",
			Italic:     "Roboto-Italic.ttf",
			Bold:       "Roboto-Bold.ttf",
			BoldItalic: "Roboto-BoldItalic.ttf",
		},
		{
			Name:       "calibri",
			Source:     fonts,
			Path:       "assets/fonts/calibri",
			Regular:    "Calibri-Regular.ttf",
			Italic:     "Calibri-Italic.ttf",
			Bold:       "Calibri-Bold.ttf",
			BoldItalic: "Calibri-BoldItalic.ttf",
		},
	}

	standardFont = Font{
		Family: "calibri",
		Size:   8,
	}
)

func StandardFontFamilies() []FontFamily {
	return standardFontFamilies
}

type FontFamily struct {
	Name       string
	Source     fs.FS
	Path       string
	Regular    string
	Italic     string
	Bold       string
	BoldItalic string
}

func FontFamilyFromPath(name string, path string) (*FontFamily, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fontFamily := &FontFamily{
		Name: name,
		Path: path,
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".ttf") {
			continue
		}

		fileName := strings.ToLower(entry.Name())
		if strings.Contains(fileName, "regular") {
			fontFamily.Regular = fileName
			continue
		}

		if strings.Contains(fileName, "bolditalic") {
			fontFamily.BoldItalic = fileName
			continue
		}

		if strings.Contains(fileName, "italic") {
			fontFamily.Italic = fileName
			continue
		}

		if strings.Contains(fileName, "bold") {
			fontFamily.Bold = fileName
			continue
		}

	}

	return fontFamily, nil
}

func (f *FontFamily) HasRegularStyle() bool {
	return len(f.Regular) > 0
}

func (f *FontFamily) HasItalicStyle() bool {
	return len(f.Italic) > 0
}

func (f *FontFamily) HasBoldStyle() bool {
	return len(f.Bold) > 0
}

func (f *FontFamily) HasBoldItalicStyle() bool {
	return len(f.BoldItalic) > 0
}

type Font struct {
	Family string
	Size   float64
	Style  *FontStyle
}

type FontStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
}

func NewFontStyle(bold, italic, underline bool) *FontStyle {
	return &FontStyle{
		Bold:      bold,
		Italic:    italic,
		Underline: underline,
	}
}

func (f *FontStyle) Combine() int {
	if f == nil {
		return 0
	}

	var combination int
	if f.Bold {
		combination = combination | gopdf.Bold
	}

	if f.Italic {
		combination = combination | gopdf.Italic
	}

	if f.Underline {
		combination = combination | gopdf.Underline
	}

	return combination
}
