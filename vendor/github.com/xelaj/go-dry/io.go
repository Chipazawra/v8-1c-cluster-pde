package dry

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
)

type CountingReader struct {
	Reader    io.Reader
	BytesRead int
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.BytesRead += n
	return n, err
}

type CountingWriter struct {
	Writer       io.Writer
	BytesWritten int
}

func (r *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = r.Writer.Write(p)
	r.BytesWritten += n
	return n, err
}

type CountingReadWriter struct {
	ReadWriter   io.ReadWriter
	BytesRead    int
	BytesWritten int
}

func (rw *CountingReadWriter) Read(p []byte) (n int, err error) {
	n, err = rw.ReadWriter.Read(p)
	rw.BytesRead += n
	return n, err
}

func (rw *CountingReadWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ReadWriter.Write(p)
	rw.BytesWritten += n
	return n, err
}

// ReadBinary wraps binary.Read with a CountingReader and returns
// the acutal bytes read by it.
func ReadBinary(r io.Reader, order binary.ByteOrder, data interface{}) (n int, err error) {
	countingReader := CountingReader{Reader: r}
	err = binary.Read(&countingReader, order, data)
	return countingReader.BytesRead, err
}

// WriteFull calls writer.Write until all of data is written,
// or an is error returned.
func WriteFull(data []byte, writer io.Writer) (n int, err error) {
	dataSize := len(data)
	for n = 0; n < dataSize; {
		m, err := writer.Write(data[n:])
		n += m
		if err != nil {
			return n, err
		}
	}
	return dataSize, nil
}

// ReadLine reads unbuffered until a newline '\n' byte and removes
// an optional carriege return '\r' at the end of the line.
// In case of an error, the string up to the error is returned.
func ReadLine(reader io.Reader) (line string, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	p := make([]byte, 1)
	for {
		var n int
		n, err = reader.Read(p)
		if err != nil || p[0] == '\n' {
			break
		}
		if n > 0 {
			buffer.WriteByte(p[0])
		}
	}
	data := buffer.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\r' {
		data = data[:len(data)-1]
	}
	return string(data), err
}

// WaitForStdin blocks until input is available from os.Stdin.
// The first byte from os.Stdin is returned as result.
// If there are println arguments, then fmt.Println will be
// called with those before reading from os.Stdin.
func WaitForStdin(v ...interface{}) byte {
	if len(v) > 0 {
		fmt.Println(v...)
	}
	buffer := make([]byte, 1)
	_, _ = os.Stdin.Read(buffer)
	return buffer[0]
}

// ReaderFunc implements io.Reader as function type with a Read method.
type ReaderFunc func(p []byte) (int, error)

func (f ReaderFunc) Read(p []byte) (int, error) {
	return f(p)
}

// WriterFunc implements io.Writer as function type with a Write method.
type WriterFunc func(p []byte) (int, error)

func (f WriterFunc) Write(p []byte) (int, error) {
	return f(p)
}

// CancelableReader позволяет читать данные с контекстом
type CancelableReader struct {
	ctx  context.Context
	data chan []byte

	// размер сообщения, которое мы хотим получить. пока в sizeWant не пошлется длина, ридер не будет читать
	sizeWant chan int

	err error
	r   io.Reader
}

func (c *CancelableReader) begin() {
	for {
		buf := make([]byte, <-c.sizeWant)
		_, err := c.r.Read(buf)
		if err != nil {
			c.err = err
			close(c.data)
			return
		}
		c.data <- buf
	}
}

func (c *CancelableReader) Read(p []byte) (int, error) {
	c.sizeWant <- len(p)
	select {
	case <-c.ctx.Done():
		return 0, c.ctx.Err()
	case d, ok := <-c.data:
		if !ok {
			return 0, c.err
		}
		copy(p, d)
		return len(d), nil
	}
}

func (c *CancelableReader) ReadByte() (byte, error) {
	b := make([]byte, 1)

	n, err := c.Read(b)
	if err != nil {
		return 0x0, err
	}
	PanicIf(n != 1, "read more than 1 byte, got "+strconv.Itoa(n))

	return b[0], nil
}

func NewCancelableReader(ctx context.Context, r io.Reader) *CancelableReader {
	c := &CancelableReader{
		r:        r,
		ctx:      ctx,
		data:     make(chan []byte),
		sizeWant: make(chan int),
	}
	go c.begin()
	return c
}
