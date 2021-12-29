package codec

import (
	"encoding/binary"
	uuid "github.com/satori/go.uuid"
	"github.com/xelaj/go-dry"
	"io"
	"math"
	"time"
)

var _ Encoder = (*encoder)(nil)

type encoder struct {
	codec        Codec
	PanicOnError bool
}

func (e *encoder) EndpointId(val int, w io.Writer) {
	e.NullableSize(val, w)
}

func (e *encoder) Value(val interface{}, w io.Writer) {

	switch typed := val.(type) {

	case bool:
		e.Bool(typed, w)
	case float32:
		e.Float32(typed, w)
	case float64:
		e.Float64(typed, w)
	case int:
		e.Int(typed, w)
	case uint:
		e.Uint(typed, w)
	case int16:
		e.Int16(typed, w)
	case uint16:
		e.Uint16(typed, w)
	case int32:
		e.Int32(typed, w)
	case uint32:
		e.Uint32(typed, w)
	case int64:
		e.Int64(typed, w)
	case uint64:
		e.Uint64(typed, w)
	case string:
		e.String(typed, w)
	case byte:
		e.Byte(typed, w)
	case time.Time:
		e.Time(typed, w)
	default:
		dry.PanicIf(true, "errow encode typed value")
	}

}

func (e *encoder) Time(val time.Time, w io.Writer) {

	ticks := dateToTicks(val)
	e.Long(ticks, w)

}

func (e *encoder) Codec() Codec {
	return e.codec
}

func (e *encoder) Bool(val bool, w io.Writer) {
	if val {
		e.write(w, []byte{TRUE_BYTE})
	} else {
		e.write(w, []byte{FALSE_BYTE})
	}
}

func (e *encoder) Byte(val byte, w io.Writer) {

	e.write(w, []byte{val})

}

func (e *encoder) Char(val int, w io.Writer) {

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(val))
	e.write(w, buf)
}

func (e *encoder) Short(val int16, w io.Writer) {
	e.Int16(val, w)
}

func (e *encoder) Int(val int, w io.Writer) {
	e.Uint32(uint32(val), w)
}

func (e *encoder) Uint(val uint, w io.Writer) {

	e.Uint32(uint32(val), w)
}

func (e *encoder) Int16(val int16, w io.Writer) {
	e.Uint16(uint16(val), w)
}

func (e *encoder) Uint16(val uint16, w io.Writer) {

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	e.write(w, buf)

}

func (e *encoder) Int32(val int32, w io.Writer) {
	e.Uint32(uint32(val), w)
}

func (e *encoder) Uint32(val uint32, w io.Writer) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	e.write(w, buf)
}

func (e *encoder) Int64(val int64, w io.Writer) {
	e.Uint64(uint64(val), w)
}

func (e *encoder) Uint64(val uint64, w io.Writer) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	e.write(w, buf)
}

func (e *encoder) Long(val int64, w io.Writer) {
	e.Uint64(uint64(val), w)
}

func (e *encoder) Float32(val float32, w io.Writer) {
	e.Uint32(math.Float32bits(val), w)
}

func (e *encoder) Float64(val float64, w io.Writer) {
	e.Uint64(math.Float64bits(val), w)
}

func (e *encoder) Double(val float64, w io.Writer) {
	e.Float64(val, w)
}

func (e *encoder) Null(w io.Writer) {
	//e.Byte(NULL_BYTE, w)
	e.Byte(0x00, w)
}

func (e *encoder) String(val string, w io.Writer) {
	if len(val) == 0 {
		e.Null(w)
		return
	}

	b := []byte(val)
	e.NullableSize(len(b), w)
	e.write(w, b)
}

func (e *encoder) TypedValue(val interface{}, w io.Writer) {

	if val == nil {
		e.Null(w)
		return
	}

	valueType := detectType(val)
	e.Type(byte(valueType), w)

	switch valueType {

	case BOOLEAN:
		e.Bool(val.(bool), w)
	case INT:
		e.Int(val.(int), w)
	case LONG:
		e.Long(val.(int64), w)
	case BYTE:
		e.Byte(val.(byte), w)
	default:
		dry.PanicIf(true, "errow encode typed value")
	}
}

func (e *encoder) Uuid(val uuid.UUID, w io.Writer) {
	buf, _ := val.MarshalBinary()
	e.write(w, buf)
}

func (e *encoder) Size(val int, w io.Writer) {
	var b1 int

	msb := val >> MAX_SHIFT
	if msb != 0 {
		b1 = -128
	} else {
		b1 = 0
	}

	e.write(w, []byte{byte(b1 | (val & 0x7F))})

	for val = msb; val > 0; val = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = -128
		} else {
			b1 = 0
		}

		e.write(w, []byte{byte(b1 | (val & 0x7F))})

	}
}

func (e *encoder) NullableSize(val int, w io.Writer) {
	var b1 int

	msb := val >> NULL_SHIFT
	if msb != 0 {
		b1 = NULL_NEXT_MASK
	} else {
		b1 = 0
	}

	e.write(w, []byte{byte(b1 | (val & 0x7F))})

	for val = msb; val > 0; val = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = NEXT_MASK
		} else {
			b1 = 0
		}

		e.write(w, []byte{byte(b1 | (val & 0x7F))})
	}
}

func (e *encoder) Type(val byte, w io.Writer) {
	if val == NULL_BYTE {
		e.Null(w)
		return
	}
	e.Byte(val, w)
}

func (e *encoder) Bytes(val []byte, w io.Writer) {
	e.write(w, val)
}

func detectType(val interface{}) TypeInterface {

	switch val.(type) {
	case bool:
		return BOOLEAN
	case byte:
		return BYTE
	case int:
		return INT
	case int64:
		return LONG
	case uint64:
		return LONG
	default:
		return BYTE
	}

}

func (e *encoder) panicOnError(_ int, err error) {

	if err != nil && e.PanicOnError {
		panic(err)
	}
}

func (e *encoder) write(w io.Writer, p []byte) {
	e.panicOnError(w.Write(p))
}
