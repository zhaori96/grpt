package grpt

type Template struct {
	Child Element
}

func (t Template) GetSize() Size {
	return t.Child.GetSize()
}

func (t *Template) Measure(boundries Size, renderer *DocumentRenderer) {
	t.Child.Measure(boundries, renderer)
}

func (t *Template) Render(renderer *DocumentRenderer) error {
	return t.Child.Render(renderer)
}
