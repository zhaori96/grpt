package grpt

import (
	"image"
	"os"
)

type Image struct {
	Size   Size
	Source image.Image
}

func NewImage(source image.Image, size Size) *Image {
	return &Image{
		Size:   size,
		Source: source,
	}
}

func NewImageFromFile(file *os.File, size Size) *Image {
	source, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return &Image{
		Size:   size,
		Source: source,
	}
}

func NewImageFrom(path string, size Size) *Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return NewImageFromFile(file, size)
}

func (i Image) GetSize() Size {
	return i.Size
}

func (i *Image) Measure(boundries Size, renderer *DocumentRenderer) {
	i.Size = i.Size.Merge(boundries)
	if i.Size.HasZeroValue() {
		imageSize := i.Source.Bounds().Size()
		i.Size = i.Size.Merge(NewSize(float64(imageSize.X), float64(imageSize.Y)))
	}
}

func (i *Image) Render(renderer *DocumentRenderer) error {
	return renderer.DrawImage(i.Source, i.Size)
}
