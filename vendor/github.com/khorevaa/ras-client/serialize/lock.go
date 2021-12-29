package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type LocksList []*LockInfo

func (l LocksList) Each(fn func(info *LockInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l *LocksList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls LocksList

	for i := 0; i < count; i++ {

		info := &LockInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type LockInfo struct {
	ClusterID  uuid.UUID `json:"cluster_id" example:"00000000-0000-0000-0000-000000000000"`
	InfobaseID uuid.UUID `json:"infobase_id" example:"8b8a0817-4cb1-4e13-9a8f-472dde1a3b47"`

	ConnectionID uuid.UUID `json:"connection_id" example:"00000000-0000-0000-0000-000000000000"`
	SessionID    uuid.UUID `json:"session_id" example:"8b8a0817-4cb1-4e13-9a8f-472dde1a3b47"`
	ObjectID     uuid.UUID `json:"object_id" example:"00000000-0000-0000-0000-000000000000"`
	LockedAt     time.Time `json:"locked_at" example:"2020-10-01T08:30:00"`
	Description  string    `rac:"descr" json:"descr" example:"БД(сеанс ,УППБоеваяБаза,разделяемая)"`
}

func (i *LockInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.ConnectionID, r)
	decoder.StringPtr(&i.Description, r)
	decoder.TimePtr(&i.LockedAt, r)
	decoder.UuidPtr(&i.ObjectID, r)
	decoder.UuidPtr(&i.SessionID, r)

}
