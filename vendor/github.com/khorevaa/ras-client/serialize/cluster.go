package serialize

import (
	"context"
	"github.com/khorevaa/ras-client/protocol/codec"
	uuid "github.com/satori/go.uuid"
	"io"
)

type ClusterInfoGetter interface {
	GetClusterInfo(ctx context.Context, cluster uuid.UUID) (ClusterInfo, error)
}

type ClusterInfo struct {
	UUID                       uuid.UUID `rac:"cluster" json:"uuid" example:"6d6958e1-a96c-4999-a995-698a0298161e"`
	Host                       string    `json:"host" example:"host"`
	Port                       int16     `json:"port" example:"1542"`
	Name                       string    `json:"name" example:"Новый кластер"`
	ExpirationTimeout          int       `json:"expiration_timeout" example:"0"`
	LifetimeLimit              int       `json:"lifetime_limit" example:"0"`
	MaxMemorySize              int       `json:"max_memory_size" example:"0"`
	MaxMemoryTimeLimit         int       `json:"max_memory_time_limit" example:"0"`
	SecurityLevel              int       `json:"security_level" example:"0"`
	SessionFaultToleranceLevel int       `json:"session_fault_tolerance_level" example:"0"`
	LoadBalancingMode          int       `json:"load_balancing_mode" example:"0"` // performance
	ErrorsCountThreshold       int       `json:"errors_count_threshold" example:"0"`
	KillProblemProcesses       bool      `json:"kill_problem_processes" example:"true"`
	KillByMemoryWithDump       bool      `json:"kill_by_memory_with_dump" example:"true"`
	LifeTimeLimit              int       `json:"life_time_limit" example:"0"`
}

func (i *ClusterInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.IntPtr(&i.ExpirationTimeout, r) // expirationTimeout
	decoder.StringPtr(&i.Host, r)
	decoder.IntPtr(&i.LifeTimeLimit, r)
	decoder.ShortPtr(&i.Port, r)
	decoder.IntPtr(&i.MaxMemorySize, r)
	decoder.IntPtr(&i.MaxMemoryTimeLimit, r)
	decoder.StringPtr(&i.Name, r)
	decoder.IntPtr(&i.SecurityLevel, r)
	decoder.IntPtr(&i.SessionFaultToleranceLevel, r)
	decoder.IntPtr(&i.LoadBalancingMode, r)
	decoder.IntPtr(&i.ErrorsCountThreshold, r)
	decoder.BoolPtr(&i.KillProblemProcesses, r)

	if version > 8 {
		decoder.BoolPtr(&i.KillByMemoryWithDump, r)
	}

}

func (i *ClusterInfo) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(i.UUID, w)
	encoder.Int(i.ExpirationTimeout, w) // expirationTimeout
	encoder.String(i.Host, w)
	encoder.Int(i.LifeTimeLimit, w)
	encoder.Short(i.Port, w)
	encoder.Int(i.MaxMemorySize, w)
	encoder.Int(i.MaxMemoryTimeLimit, w)
	encoder.String(i.Name, w)
	encoder.Int(i.SecurityLevel, w)
	encoder.Int(i.SessionFaultToleranceLevel, w)
	encoder.Int(i.LoadBalancingMode, w)
	encoder.Int(i.ErrorsCountThreshold, w)
	encoder.Bool(i.KillProblemProcesses, w)

	if version > 8 {
		encoder.Bool(i.KillByMemoryWithDump, w)
	}
}
