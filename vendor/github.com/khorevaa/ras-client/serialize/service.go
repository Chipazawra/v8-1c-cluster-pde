package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type ServiceInfo struct {
	Name        string      `json:"name" example:"EventLogService"`
	MainOnly    int         `json:"main_only" example:"0"`
	Manager     []uuid.UUID `json:"manager" example:"[ad2754ad-9415-4689-9559-74dc36b11592]"`
	Description string      `rac:"descr" json:"descr" example:"Сервис журналов регистрации"`
	ClusterID   uuid.UUID   `json:"cluster_id" example:"0e588a25-8354-4344-b935-53442312aa30"`
}

func (i *ServiceInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&i.Name, r)
	decoder.StringPtr(&i.Description, r)
	decoder.IntPtr(&i.MainOnly, r)

	idCount := decoder.Size(r)

	for ii := 0; ii < idCount; ii++ {
		i.Manager = append(i.Manager, decoder.Uuid(r))
	}
}
