package grpt

import (
	ctx "context"
	"errors"
	"fmt"
	"image"
	"math"
	"path"
	"slices"

	"github.com/signintech/gopdf"
)

type rendererState struct {
	StrokeWidth float64
	StrokeColor Color
	Font        Font
	LineStyle   LineStyle
}

type RendererOptions struct {
	PageSize Size
	Padding  EdgeInsets
}

type DocumentRenderer struct {
	options      RendererOptions
	engine       gopdf.GoPdf
	currentState rendererState

	addingPageHooks []func(*DocumentRenderer)

	bodyHeight float64

	hasHeader        bool
	headerHeight     float64
	headerInAllPages bool

	hasFooter           bool
	footerHeight        float64
	footerInAllPages    bool
	footerCallback      func(*DocumentRenderer)
	addingFooterAttemps int

	context ctx.Context
}

func StartNewDocument(options RendererOptions) *DocumentRenderer {
	renderer := &DocumentRenderer{}
	renderer.options = options
	renderer.engine.Start(gopdf.Config{
		Unit:     gopdf.UnitPT,
		PageSize: *options.PageSize.ToRect(),
	})

	renderer.AddMultiFontFamilies(standardFontFamilies...)
	renderer.SetFont(standardFont)

	renderer.engine.SetMargins(
		options.Padding.Left,
		options.Padding.Top,
		options.Padding.Right,
		options.Padding.Bottom,
	)
	renderer.engine.AddPage()

	renderer.currentState.Font = standardFont

	renderer.context = SetAvailableSpace(
		ctx.Background(),
		options.PageSize.WithPadding(options.Padding),
	)

	return renderer
}

func (r *DocumentRenderer) Context() ctx.Context {
	return r.context
}

func (r *DocumentRenderer) WritePDF(path string) error {
	if r.footerCallback != nil {
		r.footerCallback(r)
	}
	return r.engine.WritePdf(path)
}

func (r *DocumentRenderer) GetCurrentOffset() Offset {
	return Offset{X: r.engine.GetX(), Y: r.engine.GetY()}
}

func (r *DocumentRenderer) GetPageSize() Size {
	return r.options.PageSize
}

func (r *DocumentRenderer) GetPageSizeWithPadding() Size {
	return r.options.PageSize.WithPadding(r.options.Padding)
}

func (r *DocumentRenderer) GetPageWidth() float64 {
	return r.options.PageSize.Width
}

func (r *DocumentRenderer) GetPageHeight() float64 {
	return r.options.PageSize.Height
}

func (r *DocumentRenderer) GetPageAxisSize(axis Axis) float64 {
	return r.options.PageSize.WithPadding(r.options.Padding).GetAxis(axis)
}

func (r *DocumentRenderer) GetX() float64 {
	return r.engine.GetX()
}

func (r *DocumentRenderer) GetY() float64 {
	return r.engine.GetY()
}

func (r *DocumentRenderer) SetX(value float64) {
	r.engine.SetX(value)
}

func (r *DocumentRenderer) SetY(value float64) {
	r.engine.SetY(value)
}

func (r *DocumentRenderer) SetXY(x, y float64) {
	r.engine.SetXY(x, y)
}

func (r *DocumentRenderer) SetOffset(offset Offset) {
	r.engine.SetXY(offset.X, offset.Y)
}

func (r *DocumentRenderer) AddX(value float64) {
	r.engine.SetX(r.engine.GetX() + value)
}

func (r *DocumentRenderer) AddY(value float64) {
	r.engine.SetY(r.engine.GetY() + value)
}

func (r *DocumentRenderer) AddXY(x, y float64) {
	r.engine.SetXY(r.engine.GetX()+x, r.engine.GetY()+y)
}

func (r *DocumentRenderer) AddOffset(offset Offset) {
	r.AddXY(offset.X, offset.Y)
}

func (r *DocumentRenderer) AddOffsetFromAxis(offset Offset, axis Axis) {
	switch axis {
	case HorizontalAxis:
		r.AddX(offset.X)
	case VerticalAxis:
		r.AddY(offset.Y)
	default:
		panic(fmt.Errorf("DocumentRenderer.AddOffsetFromAxis: %w", ErrInvalidAxis))
	}
}

func (r *DocumentRenderer) AddToAxis(value float64, axis Axis) {
	switch axis {
	case HorizontalAxis:
		r.AddX(value)
	case VerticalAxis:
		r.AddY(value)
	default:
		panic(fmt.Errorf("DocumentRenderer.AddToAxis: %w", ErrInvalidAxis))
	}
}

func (r *DocumentRenderer) AddFontFamily(family FontFamily) error {
	var err error

	if family.HasRegularStyle() {
		filePath := path.Join(family.Path, family.Regular)
		style := gopdf.TtfOption{Style: gopdf.Regular}
		err = r.engine.AddTTFFontWithOption(family.Name, filePath, style)
		if err != nil {
			return err
		}
	}

	if family.HasItalicStyle() {
		filePath := path.Join(family.Path, family.Italic)
		style := gopdf.TtfOption{Style: gopdf.Italic}
		err = r.engine.AddTTFFontWithOption(family.Name, filePath, style)
		if err != nil {
			return err
		}
	}

	if family.HasBoldStyle() {
		filePath := path.Join(family.Path, family.Bold)
		style := gopdf.TtfOption{Style: gopdf.Bold}
		err = r.engine.AddTTFFontWithOption(family.Name, filePath, style)
		if err != nil {
			return err
		}
	}

	if family.HasBoldItalicStyle() {
		filePath := path.Join(family.Path, family.BoldItalic)
		style := gopdf.TtfOption{Style: gopdf.Bold | gopdf.Italic}
		err = r.engine.AddTTFFontWithOption(family.Name, filePath, style)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *DocumentRenderer) AddMultiFontFamilies(families ...FontFamily) error {
	for _, family := range families {
		err := r.AddFontFamily(family)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *DocumentRenderer) SetFont(font Font) error {
	return r.setFont(font, false)
}

func (r *DocumentRenderer) OnAddingPage(hooks ...func(*DocumentRenderer)) {
	r.addingPageHooks = append(r.addingPageHooks, hooks...)
}

func (r *DocumentRenderer) setFont(font Font, keepCurrentState bool) error {
	var family string = r.currentState.Font.Family
	var style int = r.currentState.Font.Style.Combine()
	var size float64 = r.currentState.Font.Size

	if len(font.Family) > 0 {
		family = font.Family
	}

	if font.Size > 0 {
		size = font.Size
	}

	if font.Style != nil {
		style = font.Style.Combine()
	}

	err := r.engine.SetFontWithStyle(family, style, size)

	if err != nil {
		return nil
	}

	if !keepCurrentState {
		r.currentState.Font = font
	}

	return nil
}

func (r *DocumentRenderer) SetFontFamily(family string) error {
	return r.setFontFamily(family, false)
}

func (r *DocumentRenderer) setFontFamily(family string, keepCurrentState bool) error {
	err := r.SetFont(Font{Family: family})
	if err != nil {
		return err
	}

	if !keepCurrentState {
		r.currentState.Font.Family = family
	}

	return nil
}

func (r *DocumentRenderer) SetFontStyle(style FontStyle) error {
	return r.setFontStyle(style, false)
}

func (r *DocumentRenderer) setFontStyle(style FontStyle, keepCurrentState bool) error {
	err := r.engine.SetFontWithStyle(
		r.currentState.Font.Family,
		style.Combine(),
		r.currentState.Font.Size,
	)

	if err != nil {
		return err
	}

	if !keepCurrentState {
		r.currentState.Font.Style = &style
	}

	return nil
}

func (r *DocumentRenderer) SetFontSize(size float64) error {
	return r.setFontSize(size, false)
}

func (r *DocumentRenderer) setFontSize(size float64, keepCurrentState bool) error {
	err := r.engine.SetFontSize(size)
	if err != nil {
		return err
	}

	if !keepCurrentState {
		r.currentState.Font.Size = size
	}

	return nil
}

func (r *DocumentRenderer) SetStrokeWidth(width float64) {
	r.setStrokeWidth(width, false)
}

func (r *DocumentRenderer) setStrokeWidth(width float64, keepCurrentState bool) {
	r.engine.SetLineWidth(width)

	if !keepCurrentState {
		r.currentState.StrokeWidth = width
	}
}

func (r *DocumentRenderer) SetStrokeColor(color Color) {
	r.setStrokeColor(color, false)

}

func (r *DocumentRenderer) setStrokeColor(color Color, keepCurrentState bool) {
	r.engine.SetStrokeColor(color.R, color.G, color.B)

	if !keepCurrentState {
		r.currentState.StrokeColor = color
	}
}

func (r *DocumentRenderer) SetLineStyle(style LineStyle) {
	r.setLineStyle(style, false)
}

func (r *DocumentRenderer) setLineStyle(style LineStyle, keepCurrentState bool) {
	r.engine.SetLineType(string(style))

	if !keepCurrentState {
		r.currentState.LineStyle = style
	}
}

func (r *DocumentRenderer) GetTextHeight(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	height, _ := r.engine.MeasureCellHeightByText(text)
	return height
}

func (r *DocumentRenderer) GetCurrentPage() int {
	return r.engine.GetNumberOfPages()
}

func (r *DocumentRenderer) FitsCurrentPage(elementHeight float64) bool {
	r.engine.UnitsToPointsVar(&elementHeight)

	offset := r.GetCurrentOffset()
	pageHeight := r.GetPageHeight()
	paddingBottom := r.options.Padding.Bottom

	return !(offset.Y+elementHeight > pageHeight-paddingBottom)
}

func (r *DocumentRenderer) FitsIn(
	element Size,
	parentPosition Offset,
	parentSize Size,
) bool {
	r.engine.UnitsToPointsVar(&element.Height)
	offset := r.GetCurrentOffset()

	limit := parentPosition.Y + parentSize.Height
	return offset.Y+element.Height <= limit
}

func (r *DocumentRenderer) BodyHeight() float64 {
	return r.bodyHeight
}

func (r *DocumentRenderer) SetCurrentBodyHeight(height float64) {
	r.bodyHeight = height
}

func (r *DocumentRenderer) HasHeader() bool {
	return r.hasHeader
}

func (r *DocumentRenderer) HeaderHeight() float64 {
	return r.headerHeight
}

func (r *DocumentRenderer) HeaderInAllPages() bool {
	return r.headerInAllPages
}

func (r *DocumentRenderer) SetHeader(
	height float64,
	header Element,
	repeat bool,
) {
	if header == nil || height == 0 {
		r.hasHeader = false
		r.headerInAllPages = false
		return
	}

	r.hasHeader = true
	r.headerInAllPages = repeat
	r.headerHeight = height

	pageSize := r.GetPageSizeWithPadding()
	header.Measure(NewSize(pageSize.Width, height), r)
	header.Render(r)
	r.AddY(height)

	if repeat {
		r.OnAddingPage(func(renderer *DocumentRenderer) {
			pageSize := r.GetPageSizeWithPadding()
			header.Measure(NewSize(pageSize.Width, height), renderer)
			header.Render(renderer)
			renderer.AddY(height)
		})
	}
}

func (r *DocumentRenderer) HasFooter() bool {
	return r.hasFooter
}

func (r *DocumentRenderer) FooterHeight() float64 {
	return r.footerHeight
}

func (r *DocumentRenderer) FooterInAllPages() bool {
	return r.footerInAllPages
}

func (r *DocumentRenderer) SetFooter(
	height float64,
	footer Element,
	repeat bool,
) {
	if footer == nil || height == 0 {
		r.hasFooter = false
		r.footerInAllPages = false
		return
	}

	r.hasFooter = true
	r.footerInAllPages = repeat
	r.footerHeight = height

	callback := func(renderer *DocumentRenderer) {
		defer r.SetOffset(r.GetCurrentOffset())

		pageSize := r.GetPageSizeWithPadding()

		footer.Measure(NewSize(pageSize.Width, height), renderer)
		availableSpace := pageSize
		if renderer.hasHeader {
			if renderer.headerInAllPages || renderer.GetCurrentPage() == 1 {
				availableSpace.Height -= renderer.headerHeight
			}
		}

		availableSpace.Height -= renderer.bodyHeight
		renderer.addingFooterAttemps += 1
		currentAttempt := renderer.addingFooterAttemps
		if !footer.GetSize().FitsContainerAxis(availableSpace, VerticalAxis) {
			renderer.AddPage()
		}

		if currentAttempt == 1 {
			r.SetY(pageSize.Height - height + renderer.options.Padding.Top)
			footer.Render(renderer)
			renderer.addingFooterAttemps = 0
		}
	}

	if repeat {
		callback(r)
		r.OnAddingPage(callback)
	} else {
		r.footerCallback = callback
	}
}

func (r *DocumentRenderer) AddPage() {
	r.engine.AddPage()
	r.SetY(r.options.Padding.Top)

	for _, hook := range r.addingPageHooks {
		if hook != nil {
			hook(r)
		}
	}
}

func (r *DocumentRenderer) DrawText(
	text string,
	size Size,
	style *TextStyle,
) error {
	offset := r.GetCurrentOffset()
	defer r.SetOffset(offset)

	if len(text) == 0 {
		return nil
	}

	if style == nil {
		return r.engine.Cell(nil, text)
	}

	if style.Font != nil {
		r.setFont(*style.Font, true)
		defer r.SetFont(r.currentState.Font)
	}

	if len(style.Borders) > 0 {
		r.DrawBoxWithBorders(size, style.Borders...)
	}

	var texts []string
	var err error
	if len(text) > 0 {
		texts, err = r.SplitText(text, size, style)
		if err != nil {
			return err
		}
	}

	r.engine.SetY(offset.Y + style.Padding.Top)
	paddedSize := size.WithPadding(style.Padding)

	if !style.Multiline {
		text := texts[0]

		if len(texts) > 1 && len(style.Overflow) > 0 {
			text = text[:len(text)-len(style.Overflow)] + style.Overflow
		}

		r.engine.SetX(offset.X + style.Padding.Left)
		return r.engine.CellWithOption(
			paddedSize.ToRect(),
			text,
			gopdf.CellOption{
				Align: int(style.Alignment),
			},
		)
	}

	textHeight, _ := r.engine.MeasureCellHeightByText(text)
	for index, text := range texts {
		if index > 0 {
			r.engine.Br(textHeight)
		}

		r.engine.SetX(offset.X + style.Padding.Left)
		r.engine.CellWithOption(
			paddedSize.ToRect(),
			text,
			gopdf.CellOption{
				Align: int(style.Alignment),
			},
		)
	}

	return nil
}

func (r *DocumentRenderer) MeasureText(
	text string,
	style *TextStyle,
	boundries Size,
) (Size, error) {
	if style.Font != nil {
		r.setFont(*style.Font, true)
		defer r.setFont(r.currentState.Font, false)
	}

	chunks, _ := r.SplitText(text, boundries, style)
	fontHeight, _ := r.engine.MeasureCellHeightByText(text)

	size := boundries
	if size.Width == 0 {
		greaterChunk := slices.Max(chunks)
		fontWidth, err := r.engine.MeasureTextWidth(greaterChunk)
		if err != nil {
			return Size{}, err
		}

		size.Width = fontWidth * fontWidth
	}

	if size.Height == 0 {
		size.Height = fontHeight * float64(len(chunks))
	}

	return size, nil
}

func (r *DocumentRenderer) SplitText(
	text string,
	size Size,
	style *TextStyle,
) ([]string, error) {
	var texts []string
	var err error

	boundries := size.WithPadding(style.Padding)

	maxRuneWidth := 0.0
	for _, char := range text {
		width, _ := r.engine.MeasureTextWidth(string(char))
		if width > maxRuneWidth {
			maxRuneWidth = width
		}
	}

	if boundries.Width < maxRuneWidth {
		boundries.Width = maxRuneWidth
	}

	if style.WordWrap {
		texts, err = r.engine.SplitTextWithWordWrap(text, boundries.Width)
	} else {
		texts, err = r.engine.SplitText(text, boundries.Width)
	}

	if len(texts) == 0 {
		return texts, err
	}

	textHeight, _ := r.engine.MeasureCellHeightByText(texts[0])
	limit := int(math.Trunc(boundries.Height / textHeight))

	if limit <= 0 {
		return nil, fmt.Errorf(
			"DocumentRenderer.SplitText: text '%s' does not fit in size %v because the text height is %v: %w",
			text,
			size,
			textHeight,
			ErrElementOverflow,
		)
	}

	if limit < len(texts) {
		texts = texts[:limit]
		if len(style.Overflow) > 0 {
			textsLength := len(texts) - 1

			lastText := texts[textsLength]
			lastTextLength := len(lastText)
			lastText = lastText[:lastTextLength-len(style.Overflow)] + style.Overflow

			texts[textsLength] = lastText
		}
	}

	return texts, err
}

func (r *DocumentRenderer) DrawLine(
	size Size,
	offset Offset,
	options *LineOptions,
) error {
	if size.IsZero() {
		return ErrInvalidSize
	}

	if options != nil {
		if options.StrokeWidth != r.currentState.StrokeWidth {
			r.setStrokeWidth(options.StrokeWidth, true)
			defer r.SetStrokeWidth(r.currentState.StrokeWidth)
		}

		if options.Color != nil &&
			!r.currentState.StrokeColor.IsEqual(*options.Color) {
			r.setStrokeColor(*options.Color, true)
			defer r.SetStrokeColor(r.currentState.StrokeColor)
		}
	}

	if size.Width == 0 {
		size.Width = r.engine.GetX()
	}

	if size.Height == 0 {
		size.Height = r.engine.GetY()
	}

	r.setLineStyle(options.Style, true)
	defer r.SetLineStyle(r.currentState.LineStyle)

	r.engine.Line(offset.X, offset.Y, size.Width, size.Height)
	return nil
}

func (r *DocumentRenderer) DrawHorizontalLine(
	width float64,
	options *LineOptions,
) error {
	return r.DrawLine(NewSize(width, 0), r.GetCurrentOffset(), options)
}

func (r *DocumentRenderer) DrawHorizontalLineWithOffset(
	width float64,
	offset Offset,
	options *LineOptions,
) error {
	return r.DrawLine(NewSize(width, offset.Y), offset, options)
}

func (r *DocumentRenderer) DrawVerticalLine(
	height float64,
	options *LineOptions,
) error {
	return r.DrawLine(NewSize(0, height), r.GetCurrentOffset(), options)
}

func (r *DocumentRenderer) DrawVerticalLineWithOffset(
	height float64,
	offset Offset,
	options *LineOptions,
) error {
	return r.DrawLine(NewSize(offset.X, height), offset, options)
}

func (r *DocumentRenderer) DrawBox(size Size, lineOptions *LineOptions) error {
	if size.IsZero() {
		err := errors.New("DrawBox can't have a zero Size")
		return ErrInvalidSize.Wrap(err)
	}

	currentOffset := r.GetCurrentOffset()
	if lineOptions == nil {
		lineOptions = &LineOptions{
			Style:       LineStyleSolid,
			StrokeWidth: 0,
			Color:       &Color{0, 0, 0},
		}
	}

	r.DrawVerticalLine(size.Height+currentOffset.Y, lineOptions)

	offset := NewOffset(size.Width+currentOffset.X, currentOffset.Y)
	r.DrawVerticalLineWithOffset(
		size.Height+currentOffset.Y,
		offset,
		lineOptions,
	)

	r.DrawHorizontalLine(size.Width+currentOffset.X, lineOptions)

	offset = NewOffset(currentOffset.X, size.Height+currentOffset.Y)
	r.DrawHorizontalLineWithOffset(
		size.Width+currentOffset.X,
		offset,
		lineOptions,
	)

	return nil
}

func (r *DocumentRenderer) DrawBoxWithBorders(
	size Size,
	borders ...Border,
) error {
	if size.IsZero() {
		err := errors.New("BoxWithBorders can't have a zero Size")
		return ErrInvalidSize.Wrap(err)
	}

	currentOffset := r.GetCurrentOffset()
	for _, border := range borders {
		if border.Side == 0 {
			continue
		}

		if border.Side&BorderLeft != 0 {
			r.DrawVerticalLine(size.Height+currentOffset.Y, &border.Options)
		}

		if border.Side&BorderRigh != 0 {
			offset := NewOffset(size.Width+currentOffset.X, currentOffset.Y)
			r.DrawVerticalLineWithOffset(
				size.Height+currentOffset.Y,
				offset,
				&border.Options,
			)
		}

		if border.Side&BorderTop != 0 {
			r.DrawHorizontalLine(size.Width+currentOffset.X, &border.Options)
		}

		if border.Side&BorderBottom != 0 {
			offset := NewOffset(currentOffset.X, size.Height+currentOffset.Y)
			r.DrawHorizontalLineWithOffset(
				size.Width+currentOffset.X,
				offset,
				&border.Options,
			)
		}
	}

	return nil
}

func (r *DocumentRenderer) DrawImage(path image.Image, size Size) error {
	offset := r.GetCurrentOffset()

	return r.engine.ImageFrom(path, offset.X, offset.Y, size.ToRect())
}
