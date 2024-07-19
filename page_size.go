package grpt

var (
	DefaultPagePadding = NewEdgeInsets(14.1732, 14.1732, 14.1732, 0)
)

var (
	PageSizeLetter      Size = Size{Width: 612, Height: 792}
	PageSizeLetterSmall Size = Size{Width: 612, Height: 792}
	PageSizeTabloid     Size = Size{Width: 792, Height: 1224}
	PageSizeLedger      Size = Size{Width: 1224, Height: 792}
	PageSizeLegal       Size = Size{Width: 612, Height: 1008}
	PageSizeStatement   Size = Size{Width: 396, Height: 612}
	PageSizeExecutive   Size = Size{Width: 540, Height: 720}
	PageSizeA0          Size = Size{Width: 2384, Height: 3371}
	PageSizeA1          Size = Size{Width: 1685, Height: 2384}
	PageSizeA2          Size = Size{Width: 1190, Height: 1684}
	PageSizeA3          Size = Size{Width: 842, Height: 1190}
	PageSizeA4          Size = Size{Width: 595, Height: 842}
	PageSizeA4Landscape Size = Size{Width: 842, Height: 595}
	PageSizeA4Small     Size = Size{Width: 595, Height: 842}
	PageSizeA5          Size = Size{Width: 420, Height: 595}
	PageSizeB4          Size = Size{Width: 729, Height: 1032}
	PageSizeB5          Size = Size{Width: 516, Height: 729}
	PageSizeFolio       Size = Size{Width: 612, Height: 936}
	PageSizeQuarto      Size = Size{Width: 610, Height: 780}
	PageSize10x14       Size = Size{Width: 720, Height: 1008}
)
