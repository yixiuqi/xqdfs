package needle

func NeedleOffset(offset int64) uint32 {
	return uint32(offset / PaddingSize)
}

func BlockOffset(offset uint32) int64 {
	return int64(offset) * PaddingSize
}

func align(d int32) int32 {
	return (d + PaddingAlign) & ^PaddingAlign
}

func Size(n int) int {
	return int(align(HeaderSize + int32(n) + FooterSize))
}
