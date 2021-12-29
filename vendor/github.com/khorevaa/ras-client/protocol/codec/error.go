package codec

import (
	"fmt"
)

type DecoderError struct {
	fn        string
	err       error
	needBytes []byte
	readBytes int
}

func (e *DecoderError) Error() string {

	return fmt.Sprintf("decoder: fn<%s> err<%s>", e.fn, e.err.Error())

}
