package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type ManagerInfo struct {
	UUID        uuid.UUID `rac:"manager" json:"uuid" example:"0e588a25-8354-4344-b935-53442312aa30"`
	PID         string    `json:"pid" example:"3366"`
	Using       string    `json:"using" example:"normal"`
	Host        string    `json:"host" example:"srv"`
	MainManager int       `json:"main_manager" example:"1"`
	Port        int16     `json:"port" example:"1541"`
	Description string    `rac:"descr" json:"descr" example:"Главный менеджер кластера"`
	ClusterID   uuid.UUID `json:"cluster_id" example:"0e588a25-8354-4344-b935-53442312aa30"`
}

func (i *ManagerInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Host, r)
	decoder.IntPtr(&i.MainManager, r)
	decoder.ShortPtr(&i.Port, r) // expirationTimeout
	decoder.StringPtr(&i.PID, r)

}
