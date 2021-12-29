package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetInfobasesShortRequest)(nil)

// GetInfobasesShortRequest получение списка инфорамационных баз кластера
//
//  type GET_INFOBASES_SHORT_REQUEST = 43
//  kind MESSAGE_KIND = 1
//  respond GetInfobasesShortResponse
type GetInfobasesShortRequest struct {
	ClusterID uuid.UUID
}

func (r GetInfobasesShortRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r GetInfobasesShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ GetInfobasesShortRequest) Type() EndpointMessageType {
	return GET_INFOBASES_SHORT_REQUEST
}

// GetInfobasesShortResponse
// type GET_INFOBASES_SHORT_RESPONSE = 44
type GetInfobasesShortResponse struct {
	Infobases serialize.InfobaseSummaryList
}

func (res *GetInfobasesShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.InfobaseSummaryList{}

	list.Parse(decoder, version, r)

	res.Infobases = list

}

func (_ *GetInfobasesShortResponse) Type() EndpointMessageType {
	return GET_INFOBASES_SHORT_RESPONSE
}

var _ EndpointRequestMessage = (*CreateInfobaseRequest)(nil)

// CreateInfobaseRequest запрос на создание новой базы
//
//  type CREATE_INFOBASE_REQUEST = 38
//  kind MESSAGE_KIND = 1
//  respond CreateInfobaseResponse
type CreateInfobaseRequest struct {
	ClusterID uuid.UUID
	Infobase  *serialize.InfobaseInfo
	Mode      int // Mode 1 - создавать базу на сервере, 0 - не создавать
}

func (r *CreateInfobaseRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *CreateInfobaseRequest) Type() EndpointMessageType {
	return CREATE_INFOBASE_REQUEST
}

func (r *CreateInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)

	r.Infobase.Format(encoder, version, w)

	encoder.Int(r.Mode, w)
}

// CreateInfobaseResponse ответ создания новой информационной базы
//  type CREATE_INFOBASE_RESPONSE = 39
//  return uuid.UUID созданной базы
type CreateInfobaseResponse struct {
	InfobaseID uuid.UUID
}

func (_ *CreateInfobaseResponse) Type() EndpointMessageType {
	return CREATE_INFOBASE_RESPONSE
}

func (res *CreateInfobaseResponse) Parse(decoder codec.Decoder, _ int, r io.Reader) {
	decoder.UuidPtr(&res.InfobaseID, r)
}

var _ EndpointRequestMessage = (*GetInfobaseInfoRequest)(nil)

// GetInfobaseInfoRequest запрос получение информации по информационной базе
//
//  type GET_INFOBASE_INFO_REQUEST = 49
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseInfoResponse
type GetInfobaseInfoRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
}

func (r *GetInfobaseInfoRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetInfobaseInfoRequest) Type() EndpointMessageType {
	return GET_INFOBASE_INFO_REQUEST
}

func (r *GetInfobaseInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

// GetInfobaseInfoResponse ответ с информацией о информационной базы
//  type GET_INFOBASE_INFO_RESPONSE = 50
//  return serialize.InfobaseInfo
type GetInfobaseInfoResponse struct {
	Infobase serialize.InfobaseInfo
}

func (_ *GetInfobaseInfoResponse) Type() EndpointMessageType {
	return GET_INFOBASE_INFO_RESPONSE
}

func (res *GetInfobaseInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := &serialize.InfobaseInfo{}
	info.Parse(decoder, version, r)

	res.Infobase = *info

}

var _ EndpointRequestMessage = (*DropInfobaseRequest)(nil)

// DropInfobaseRequest запрос удаление информационной базы
//
//  type DROP_INFOBASE_REQUEST = 42
//  kind MESSAGE_KIND = 1
//  respond nothing
type DropInfobaseRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	Mode       int
}

func (r *DropInfobaseRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *DropInfobaseRequest) Type() EndpointMessageType {
	return DROP_INFOBASE_REQUEST
}

func (r *DropInfobaseRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
	encoder.Int(r.Mode, w)
}

var _ EndpointRequestMessage = (*UpdateInfobaseRequest)(nil)

// UpdateInfobaseRequest запрос обновление данных по информационной базы
//
//  type UPDATE_INFOBASE_REQUEST = 40
//  kind MESSAGE_KIND = 1
//  respond nothing
type UpdateInfobaseRequest struct {
	ClusterID uuid.UUID
	Infobase  serialize.InfobaseInfo
}

func (r *UpdateInfobaseRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.Infobase.UUID)
}

func (_ *UpdateInfobaseRequest) Type() EndpointMessageType {
	return UPDATE_INFOBASE_REQUEST
}

func (r *UpdateInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	r.Infobase.Format(encoder, version, w)

}

var _ EndpointRequestMessage = (*UpdateInfobaseShortRequest)(nil)

// UpdateInfobaseShortRequest запрос обновление данных по информационной базы
//
//  type UPDATE_INFOBASE_REQUEST = 40
//  kind MESSAGE_KIND = 1
//  respond nothing
type UpdateInfobaseShortRequest struct {
	ClusterID uuid.UUID
	Infobase  serialize.InfobaseSummaryInfo
}

func (r *UpdateInfobaseShortRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.Infobase.UUID)
}

func (_ *UpdateInfobaseShortRequest) Type() EndpointMessageType {
	return UPDATE_INFOBASE_SHORT_REQUEST
}

func (r *UpdateInfobaseShortRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	r.Infobase.Format(encoder, version, w)

}
