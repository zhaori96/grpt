package grpt

import "io"

type DocumentHeader struct {
	ShouldRepeat bool
	Height       float64
	Elements     Elements
}

type DocumentBody struct {
	Elements Elements
}

type DocumentFooter struct {
	ShouldRepeat bool
	Height       float64
	Elements     Elements
}

type Document struct {
	PageSize Size
	Padding  EdgeInsets
	Header   DocumentHeader
	Body     DocumentBody
	Footer   DocumentFooter
}

func (d *Document) Write() ([]byte, error) {
	renderer, err := d.build()
	if err != nil {
		return nil, err
	}

	return renderer.Write()
}

func (d *Document) WriteTo(writer io.Writer) (int64, error) {
	renderer, err := d.build()
	if err != nil {
		return 0, err
	}

	return renderer.WriteTo(writer)
}

func (d *Document) WritePDF(path string) (error) {
	renderer, err := d.build()
	if err != nil {
		return err
	}

	return renderer.WritePDF(path)
}

func (d *Document) build() (*DocumentRenderer, error) {
	renderer := StartNewDocument(RendererOptions{
		PageSize: d.PageSize,
		Padding:  d.Padding,
	})

	if len(d.Header.Elements) > 0 {
		header := &Column{
			Size:     NewSize(MaxSize, d.Header.Height),
			Children: d.Header.Elements,
		}
		renderer.SetHeader(d.Header.Height, header, d.Header.ShouldRepeat)
	}

	if len(d.Footer.Elements) > 0 {
		footer := &Column{
			Size:     NewSize(MaxSize, d.Footer.Height),
			Children: d.Footer.Elements,
		}
		renderer.SetFooter(d.Footer.Height, footer, d.Footer.ShouldRepeat)
	}

	paddedPageSize := renderer.GetPageSizeWithPadding()
	initialBodySize := paddedPageSize
	if renderer.HasHeader() {
		initialBodySize.Height -= renderer.HeaderHeight()
	}

	if renderer.FooterInAllPages() {
		initialBodySize.Height -= renderer.FooterHeight()
	}

	body := &Column{
		OverflowMode: OverflowModeContinueOnNextPage,
		Size:         NewMaxWidth(),
		Children:     d.Body.Elements,
	}

	renderer.OnAddingPage(func(renderer *DocumentRenderer) {
		newBodyBoundries := renderer.GetPageSizeWithPadding()
		if renderer.HeaderInAllPages() && renderer.GetCurrentPage() > 1 {
			newBodyBoundries.Height -= renderer.HeaderHeight()
		}

		if renderer.FooterInAllPages() {
			newBodyBoundries.Height -= renderer.FooterHeight()
		}

		body.Measure(newBodyBoundries, renderer)
		renderer.SetCurrentBodyHeight(body.GetSize().Height)
	})

	body.Measure(initialBodySize, renderer)
	if err := body.Render(renderer); err != nil {
		return nil, err
	}

	return renderer, nil
}
