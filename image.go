package grpt

import (
	"bytes"
	"errors"
	"image"
	"io"
	"os"

	"github.com/signintech/gopdf"
)

type ImageFormat string

const (
	ImageFormatDefault ImageFormat = ""
	ImageFormatJPEG    ImageFormat = "jpeg"
	ImageFormatPNG     ImageFormat = "png"
)

type ImageBlendMode string

const (
	BlendModeHue             = ImageBlendMode(gopdf.Hue)
	BlendModeColor           = ImageBlendMode(gopdf.Color)
	BlendModeNormalBlendMode = ImageBlendMode(gopdf.NormalBlendMode)
	BlendModeDarken          = ImageBlendMode(gopdf.Darken)
	BlendModeScreen          = ImageBlendMode(gopdf.Screen)
	BlendModeOverlay         = ImageBlendMode(gopdf.Overlay)
	BlendModeLighten         = ImageBlendMode(gopdf.Lighten)
	BlendModeMultiply        = ImageBlendMode(gopdf.Multiply)
	BlendModeExclusion       = ImageBlendMode(gopdf.Exclusion)
	BlendModeColorBurn       = ImageBlendMode(gopdf.ColorBurn)
	BlendModeHardLight       = ImageBlendMode(gopdf.HardLight)
	BlendModeSoftLight       = ImageBlendMode(gopdf.SoftLight)
	BlendModeDifference      = ImageBlendMode(gopdf.Difference)
	BlendModeSaturation      = ImageBlendMode(gopdf.Saturation)
	BlendModeLuminosity      = ImageBlendMode(gopdf.Luminosity)
	BlendModeColorDodge      = ImageBlendMode(gopdf.ColorDodge)
)

type ImageCrop struct {
	Offset
	Size
}

type ImageTransparency struct {
	Alpha     float64
	BlendMode ImageBlendMode
}

type ImageMask struct {
	ImageOptions
	BoundBox *[4]float64
	Source   any
}

type ImageOptions struct {
	Format         ImageFormat
	DegreeAngle    float64
	VerticalFlip   bool
	HorizontalFlip bool
	Crop           *ImageCrop
	Mask           *ImageMask
	Transparency   *ImageTransparency
}

func (i ImageOptions) toGopdf() gopdf.ImageOptions {
	options := gopdf.ImageOptions{}
	options.DegreeAngle = i.DegreeAngle
	options.HorizontalFlip = i.HorizontalFlip
	options.VerticalFlip = i.VerticalFlip

	if i.Transparency != nil {
		options.Transparency = &gopdf.Transparency{
			Alpha:         i.Transparency.Alpha,
			BlendModeType: gopdf.BlendModeType(i.Transparency.BlendMode),
		}
	}

	if i.Crop != nil {
		options.Crop = &gopdf.CropOptions{
			X:      i.Crop.X,
			Y:      i.Crop.Y,
			Width:  i.Crop.Width,
			Height: i.Crop.Height,
		}
	}

	if i.Mask != nil {
		var img gopdf.ImageHolder
		if i.Mask.Source != nil {
			var err error
			img, err = sourceToImageHolder(i.Mask.Source, i.Mask.Format)
			if err != nil {
				panic(err)
			}
		}
		options.Mask = &gopdf.MaskOptions{
			ImageOptions: i.Mask.toGopdf(),
			BBox:         i.Mask.BoundBox,
			Holder:       img,
		}
	}

	return options
}

type Image struct {
	Size    Size
	Options ImageOptions
	Source  any
}

func NewImage(source any, size Size) *Image {
	return &Image{
		Size:   size,
		Source: source,
	}
}

func NewImageWithOptions(source any, size Size, options ImageOptions) *Image {
	return &Image{
		Size:    size,
		Options: options,
		Source:  source,
	}
}

func (i Image) GetSize() Size {
	return i.Size
}

func (i *Image) Measure(boundries Size, renderer *DocumentRenderer) {
	i.Size = i.Size.Merge(boundries)
	if i.Size.HasZeroValue() {
		imageSize, err := i.getImageSize()
		if err != nil {
			panic(err)
		}
		i.Size = i.Size.Merge(imageSize)
	}
}

func (i *Image) Render(renderer *DocumentRenderer) error {
	defer renderer.SetOffset(renderer.GetCurrentOffset())
	return renderer.DrawImage(i.Source, i.Size, i.Options)
}

func (i Image) getImageSize() (Size, error) {
	var img image.Image
	var err error

	switch v := i.Source.(type) {
	case string:
		var file *os.File
		file, err = os.Open(v)
		if err != nil {
			return Size{}, err
		}
		defer file.Close()
		img, _, err = image.Decode(file)
	case []byte:
		reader := bytes.NewReader(v)
		img, _, err = image.Decode(reader)
	case io.Reader:
		img, _, err = image.Decode(v)
	case image.Image:
		img = v
	default:
		err = errors.New("unsupported input type")
	}

	if err != nil {
		return Size{}, err
	}

	bounds := img.Bounds()
	return NewSize(float64(bounds.Dx()), float64(bounds.Dy())), nil
}
