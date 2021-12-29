package serialize

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	uuid "github.com/satori/go.uuid"
	"io"
)

type AssignmentsList []*AssignmentInfo

func (l AssignmentsList) Each(fn func(info *AssignmentInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l *AssignmentsList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls AssignmentsList

	for i := 0; i < count; i++ {

		info := &AssignmentInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type AssignmentInfo struct {
	Uuid           uuid.UUID `json:"id" `
	ObjectType     string    `json:"object_type"`
	InfobaseName   string    `json:"infobase_name"`
	Type           int       `json:"type"`
	ApplicationExt string    `json:"application_ext"`
	Priority       int       `json:"priority"`
}

func (i *AssignmentInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.Uuid, r)
	decoder.StringPtr(&i.ObjectType, r)
	decoder.StringPtr(&i.InfobaseName, r)
	decoder.IntPtr(&i.Type, r)
	decoder.StringPtr(&i.ApplicationExt, r)
	decoder.IntPtr(&i.Priority, r)

}

func (i *AssignmentInfo) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(i.Uuid, w)
	encoder.String(i.ObjectType, w)
	encoder.String(i.InfobaseName, w)
	encoder.Int(i.Type, w)
	encoder.String(i.ApplicationExt, w)
	encoder.Int(i.Priority, w)
}
