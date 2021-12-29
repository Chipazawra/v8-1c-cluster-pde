package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type ProcessInfoList []*ProcessInfo

func (l ProcessInfoList) Each(fn func(info *ProcessInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l ProcessInfoList) Filter(fn func(info *ProcessInfo) bool) ProcessInfoList {

	return l.filter(fn, 0)

}

func (l ProcessInfoList) filter(fn func(info *ProcessInfo) bool, count int) (val ProcessInfoList) {

	n := 0

	for _, info := range l {

		if n == count && count > 0 {
			break
		}

		result := fn(info)

		if result {
			n += 1
			val = append(val, info)
		}

	}

	return

}

func (l *ProcessInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls ProcessInfoList

	for i := 0; i < count; i++ {

		info := &ProcessInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type ProcessInfo struct {
	UUID                uuid.UUID       `rac:"process" json:"uuid" example:"0e588a25-8354-4344-b935-53442312aa30"`
	Host                string          `json:"host" example:"srv"`
	Port                int16           `json:"port" example:"1564"`
	Pid                 string          `json:"pid" example:"3366"`
	Enable              bool            `rac:"is-enable" json:"enable" example:"true"`
	Running             bool            `json:"running" example:"true"`
	StartedAt           time.Time       `json:"started_at" example:"2018-03-29T11:16:02"`
	Use                 bool            `json:"use" example:"true"`
	AvailablePerfomance int             `json:"available_perfomance" example:"100"`
	Capacity            int             `json:"capacity" example:"1000"`
	Connections         int             `json:"connections" example:"7"`
	MemorySize          int             `json:"memory_size" example:"1518604"`
	MemoryExcessTime    int             `json:"memory_excess_time" example:"0"`
	SelectionSize       int             `json:"selection_size" example:"61341"`
	AvgBackCallTime     float64         `json:"avg_back_call_time" example:"0.000"`
	AvgCallTime         float64         `json:"avg_call_time" example:"0.483"`
	AvgDbCallTime       float64         `json:"avg_db_call_time" example:"0.124"`
	AvgLockCallTime     float64         `json:"avg_lock_call_time" example:"0.000"`
	AvgServerCallTime   float64         `json:"avg_server_call_time" example:"-0.265"`
	AvgThreads          float64         `json:"avg_threads" example:"0.281"`
	Reverse             bool            `json:"reverse" example:"true"`
	Licenses            LicenseInfoList `json:"licenses"`

	ClusterID uuid.UUID `json:"cluster_id" example:"0e588a25-8354-4344-b935-53442312aa30"`
}

func (i *ProcessInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)

	decoder.DoublePtr(&i.AvgBackCallTime, r)
	decoder.DoublePtr(&i.AvgCallTime, r)
	decoder.DoublePtr(&i.AvgDbCallTime, r)
	decoder.DoublePtr(&i.AvgLockCallTime, r)
	decoder.DoublePtr(&i.AvgServerCallTime, r)
	decoder.DoublePtr(&i.AvgThreads, r)
	decoder.IntPtr(&i.Capacity, r)
	decoder.IntPtr(&i.Connections, r)
	decoder.StringPtr(&i.Host, r)
	decoder.BoolPtr(&i.Enable, r)

	licenseList := LicenseInfoList{}
	licenseList.Parse(decoder, version, r)
	i.Licenses = licenseList

	decoder.ShortPtr(&i.Port, r)
	decoder.IntPtr(&i.MemoryExcessTime, r)
	decoder.IntPtr(&i.MemorySize, r)

	decoder.StringPtr(&i.Pid, r)

	running := decoder.Int(r)
	if running == 1 {
		i.Running = true
	}

	decoder.IntPtr(&i.SelectionSize, r)
	decoder.TimePtr(&i.StartedAt, r)

	use := decoder.Int(r)
	if use == 1 {
		i.Use = true
	}

	decoder.IntPtr(&i.AvailablePerfomance, r)

	if version >= 9 {
		decoder.BoolPtr(&i.Reverse, r)
	}

	i.Licenses.Each(func(info *LicenseInfo) {
		info.ProcessID = i.UUID
	})

}
