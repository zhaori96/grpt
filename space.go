package grpt

type Space struct {
	axis                   Axis
	size                   Size
	wasMeasuredAtLeastOnce bool
	originalSize           Size
}

func NewVerticalSpace(size float64) *Space {
	return &Space{axis: VerticalAxis, size: NewHeight(size)}
}

func NewHorizontalSpace(size float64) *Space {
	return &Space{axis: HorizontalAxis, size: NewWidth(size)}
}

func (s Space) GetSize() Size {
	return s.size
}

func (s *Space) Measure(boundries Size, _ *DocumentRenderer) {
	if s.wasMeasuredAtLeastOnce {
		s.size = s.originalSize
	} else {
		s.originalSize = s.size
	}
	s.wasMeasuredAtLeastOnce = true

	size := s.size.GetAxis(s.axis)
	if size == 0 || size == MaxSize {
		s.size = s.size.SetMainAxis(boundries.GetAxis(s.axis), s.axis)
	}
}

func (s *Space) Render(_ *DocumentRenderer) error {
	return nil
}
