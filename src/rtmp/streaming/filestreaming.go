package streaming

import (
	"bufio"
	"io"
	"os"
)

/*
文件流, 将数据写入文件
*/
type FileStreaming struct {
	fileName string
	file *os.File
	buffer io.Writer
}

func NewFileStreaming(fileName string) (*FileStreaming, error) {
	stream := FileStreaming {
		fileName: fileName,
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if (err != nil) {
		return nil, err
	}
	stream.file = file
	stream.buffer = bufio.NewWriter(stream.file) // 创建缓冲

	return &stream, nil
}

func (stream *FileStreaming) Write(byteArray []byte) (int, error) {
	return stream.buffer.Write(byteArray)
}


func (stream *FileStreaming) Close() {
	if (stream.buffer != nil) {
		stream.buffer = nil
	}

	if (stream.file != nil) {
		stream.file.Close()
		stream.file = nil
	}
}