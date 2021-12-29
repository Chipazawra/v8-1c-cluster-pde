package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

//
//GET_WORKING_PROCESSES_REQUEST
//GET_WORKING_PROCESSES_RESPONSE
//GET_WORKING_PROCESS_INFO_REQUEST
//GET_WORKING_PROCESS_INFO_RESPONSE
//GET_SERVER_WORKING_PROCESSES_REQUEST
//GET_SERVER_WORKING_PROCESSES_RESPONSE

var _ EndpointRequestMessage = (*GetWorkingProcessesRequest)(nil)

// GetWorkingProcessesRequest получение списка процессов кластера
//
//  type GET_WORKING_PROCESSES_REQUEST = 19
//  respond GetWorkingProcessesResponse
type GetWorkingProcessesRequest struct {
	ClusterID uuid.UUID
}

func (r *GetWorkingProcessesRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetWorkingProcessesRequest) Type() EndpointMessageType {
	return GET_WORKING_PROCESSES_REQUEST
}

func (r *GetWorkingProcessesRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)

}

// GetWorkingProcessesResponse содержит список процессов кластера
//  type GET_WORKING_PROCESS_INFO_RESPONSE = 20
//  Managers serialize.ProcessInfoList
type GetWorkingProcessesResponse struct {
	Processes serialize.ProcessInfoList
}

func (res *GetWorkingProcessesResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ProcessInfoList{}
	list.Parse(decoder, version, r)

	res.Processes = list

}

func (_ *GetWorkingProcessesResponse) Type() EndpointMessageType {
	return GET_WORKING_PROCESS_INFO_RESPONSE
}

var _ EndpointRequestMessage = (*GetWorkingProcessInfoRequest)(nil)

// GetWorkingProcessesRequest получение списка процессов кластера
//
//  type GET_WORKING_PROCESS_INFO_REQUEST = 19
//  respond GetWorkingProcessesResponse
type GetWorkingProcessInfoRequest struct {
	ClusterID uuid.UUID
	ProcessID uuid.UUID
}

func (r *GetWorkingProcessInfoRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetWorkingProcessInfoRequest) Type() EndpointMessageType {
	return GET_WORKING_PROCESS_INFO_REQUEST
}

func (r *GetWorkingProcessInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ProcessID, w)

}

// GetWorkingProcessesResponse содержит список процессов кластера
//  type GET_WORKING_PROCESS_INFO_RESPONSE = 20
//  Managers serialize.ProcessInfo
type GetWorkingProcessInfoResponse struct {
	Info *serialize.ProcessInfo
}

func (res *GetWorkingProcessInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := &serialize.ProcessInfo{}
	info.Parse(decoder, version, r)

	res.Info = info

}

func (_ *GetWorkingProcessInfoResponse) Type() EndpointMessageType {
	return GET_WORKING_PROCESS_INFO_RESPONSE
}
