package serialize

import (
	"fmt"
	"github.com/khorevaa/ras-client/protocol/codec"
	uuid "github.com/satori/go.uuid"
	"io"
	"strconv"
	"strings"
)

type ServerInfo struct {
	UUID                                 uuid.UUID `rac:"server" json:"uuid" example:"82b8f05a-898e-48ec-9a5b-461bdf66b7d0"`
	AgentHost                            string    `json:"agent_host" example:"app"`
	AgentPort                            int       `json:"agent_port" example:"1540"`
	PortRange                            []string  `json:"port_range" example:"1560:1591"`
	Name                                 string    `json:"name" example:"Центральный сервер"`
	MainServer                           bool      `json:"main_server" example:"true"`
	DedicateManagers                     bool      `json:"dedicate_managers" example:"false"`
	InfobasesLimit                       int       `json:"infobases_limit" example:"8"`
	MemoryLimit                          int64     `json:"memory_limit" example:"0"`
	ConnectionsLimit                     int       `json:"connections_limit" example:"128"`
	SafeWorkingProcessesMemoryLimit      int64     `json:"safe_working_processes_memory_limit" example:"0"`
	SafeCallMemoryLimit                  int64     `json:"safe_call_memory_limit" example:"0"`
	ClusterPort                          int       `json:"cluster_port" example:"1541"`
	CriticalTotalMemory                  int64     `json:"critical_total_memory" example:"0"`
	TemporaryAllowedTotalMemory          int64     `json:"temporary_allowed_total_memory" example:"0"`
	TemporaryAllowedTotalMemoryTimeLimit int64     `json:"temporary_allowed_total_memory_time_limit" example:"300"`
	ClusterID                            uuid.UUID `json:"cluster_id" example:"0e588a25-8354-4344-b935-53442312aa30"`
}

func (info *ServerInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&info.UUID, r)
	decoder.StringPtr(&info.AgentHost, r)
	decoder.IntPtr(&info.AgentPort, r)
	decoder.StringPtr(&info.Name, r)
	decoder.BoolPtr(&info.MainServer, r)
	decoder.LongPtr(&info.SafeWorkingProcessesMemoryLimit, r)
	decoder.LongPtr(&info.SafeCallMemoryLimit, r)
	decoder.IntPtr(&info.InfobasesLimit, r)
	decoder.LongPtr(&info.MemoryLimit, r)
	decoder.IntPtr(&info.ConnectionsLimit, r)
	decoder.IntPtr(&info.ClusterPort, r)
	decoder.BoolPtr(&info.DedicateManagers, r)

	count := decoder.Size(r)

	for i := 0; i < count; i++ {
		highBound := decoder.Short(r)
		lowBound := decoder.Short(r)
		info.PortRange = append(info.PortRange, fmt.Sprintf("%d:%D", lowBound, highBound))
	}

	if version >= 8 {
		decoder.LongPtr(&info.CriticalTotalMemory, r)
		decoder.LongPtr(&info.TemporaryAllowedTotalMemory, r)
		decoder.LongPtr(&info.TemporaryAllowedTotalMemoryTimeLimit, r)
	}
}

func (info *ServerInfo) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(info.UUID, w)
	encoder.String(info.AgentHost, w)
	encoder.Int(info.AgentPort, w)
	encoder.String(info.Name, w)
	encoder.Bool(info.MainServer, w)
	encoder.Long(info.SafeWorkingProcessesMemoryLimit, w)
	encoder.Long(info.SafeCallMemoryLimit, w)
	encoder.Int(info.InfobasesLimit, w)
	encoder.Long(info.MemoryLimit, w)
	encoder.Int(info.ConnectionsLimit, w)
	encoder.Int(info.ClusterPort, w)
	encoder.Bool(info.DedicateManagers, w)

	encoder.Size(len(info.PortRange), w)
	for _, s := range info.PortRange {

		bounds := strings.Split(s, ":")
		lowBound, _ := strconv.ParseInt(bounds[0], 10, 64)
		highBound, _ := strconv.ParseInt(bounds[0], 10, 64)
		encoder.Short(int16(highBound), w)
		encoder.Short(int16(lowBound), w)
	}

	if version >= 8 {
		encoder.Long(info.CriticalTotalMemory, w)
		encoder.Long(info.TemporaryAllowedTotalMemory, w)
		encoder.Long(info.TemporaryAllowedTotalMemoryTimeLimit, w)
	}

}
