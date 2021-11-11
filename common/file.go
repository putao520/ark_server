package common

import (
	"mime/multipart"
	"os"
)

const (
	FileMaxBuffer = int64(0x40000)
)

func ReadMultipartToOther(source multipart.File, info *multipart.FileHeader, target *os.File) bool {
	sourceLength := info.Size
	buffer := make([]byte, FileMaxBuffer)

	offset := int64(0)
	length := FileMaxBuffer
	for offset < sourceLength {
		furtherLength := offset + FileMaxBuffer
		if furtherLength > sourceLength {
			length = FileMaxBuffer - (furtherLength - sourceLength)
		}
		_, err := source.ReadAt(buffer[:length], offset)
		if err != nil {
			// panic(err)
			return false
		}
		_, err = target.WriteAt(buffer, offset)
		if err != nil {
			// panic(err)
			return false
		}
		offset += length
	}
	return true
}
