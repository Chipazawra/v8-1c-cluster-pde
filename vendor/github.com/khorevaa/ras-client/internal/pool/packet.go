package pool

import (
	"bytes"
	"io"
)

const MaxShift = 7

// Send buffer size determines how many bytes we send in a single TCP write call.
// This can be anything from 1 to 65495.
// A good default value for this can be readPacket from: /proc/sys/net/ipv4/tcp_wmem
const sendBufferSize = 16384

// Packet represents a single network message.
// It has a byte code indicating the type of the message
// and a data payload in the form of a byte slice.
type Packet struct {
	Type   byte
	Length int
	Data   []byte

	reader *bytes.Reader
}

// New creates a new packet.
// It expects a byteCode for the type of message and
// a data parameter in the form of a byte slice.
func NewPacket(byteCode byte, data []byte) *Packet {
	return &Packet{
		Type:   byteCode,
		Length: len(data),
		Data:   data,
		reader: bytes.NewReader(data),
	}
}

// Read read the packet data
func (packet *Packet) Read(p []byte) (n int, err error) {

	return packet.reader.Read(p)

}

// Write writes the packet to the IO device.
func (packet *Packet) Write(writer io.Writer) error {

	// Для типа 0 NEGOTIATE пишем только тело
	if packet.Type != 0 {

		buf := bytes.NewBuffer([]byte{packet.Type})
		encodeSize(packet.Length, buf)

		_, err := buf.WriteTo(writer)

		if err != nil {
			return err
		}
	}

	bytesWritten := 0
	writeUntil := 0

	for bytesWritten < len(packet.Data) {
		writeUntil = bytesWritten + sendBufferSize

		if writeUntil > len(packet.Data) {
			writeUntil = len(packet.Data)
		}

		n, err := writer.Write(packet.Data[bytesWritten:writeUntil])

		if err != nil {
			return err
		}

		bytesWritten += n
	}

	return nil
}

// Bytes returns the raw byte slice serialization of the packet.
func (packet *Packet) Bytes() []byte {
	result := []byte{packet.Type}
	size := bytes.NewBuffer([]byte{})
	encodeSize(packet.Length, size)
	result = append(result, size.Bytes()...)
	result = append(result, packet.Data...)
	return result
}

func encodeSize(val int, buf *bytes.Buffer) {
	var b1 int

	msb := val >> MaxShift
	if msb != 0 {
		b1 = -128
	} else {
		b1 = 0
	}

	buf.Write([]byte{byte(b1 | (val & 0x7F))})

	for val = msb; val > 0; val = msb {

		msb >>= MaxShift
		if msb != 0 {
			b1 = -128
		} else {
			b1 = 0
		}

		buf.Write([]byte{byte(b1 | (val & 0x7F))})

	}

}
