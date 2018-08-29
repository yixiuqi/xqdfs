package tool

import "xqdfs/storage/needle"

func align(d int32) int32 {
	return (d + needle.PaddingAlign) & ^needle.PaddingAlign
}

func FileSizeCalc(size int32) int32{
	totalSize:= int32(needle.HeaderSize + size + needle.FooterSize)
	paddingSize:= align(totalSize) - totalSize
	totalSize += paddingSize
	return totalSize
}
