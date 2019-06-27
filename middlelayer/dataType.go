// dataType
package middlelayer

// 剪裁图结构
type CropImageStruct struct {
	SrcImg   Image
	CropRect Rect
	CropImg  string
}

type Image struct {
	ImageBuf string
	Size     int
	Format   int
	Width    int
	Height   int
}

type Rect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}
