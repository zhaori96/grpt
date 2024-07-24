package grpt

type Selector struct {
	Index    int
	Selector func(length int) int
	Elements Elements
}

func (s Selector) GetSize() Size {
	return s.GetSelected().GetSize()
}

func (s *Selector) Measure(boundries Size, renderer *DocumentRenderer) {
	element := s.GetSelected()
	element.Measure(boundries, renderer)
}

func (s *Selector) Render(renderer *DocumentRenderer) error {
	return s.GetSelected().Render(renderer)
}

func (s *Selector) GetSelected() Element {
	if s.Selector != nil {
		index := s.Selector(len(s.Elements))
		return s.Elements[index]
	}

	return s.Elements[s.Index]
}
