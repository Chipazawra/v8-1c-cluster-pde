package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

var _ EndpointRequestMessage = (*GetClustersRequest)(nil)

// GetClustersRequest получение списка кластеров
//
//  type GET_CLUSTERS_REQUEST = 11
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClustersRequest struct{}

func (_ *GetClustersRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (_ *GetClustersRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

func (_ GetClustersRequest) Type() EndpointMessageType {
	return GET_CLUSTERS_REQUEST
}

// GetClustersResponse ответ со списком кластеров
//
//  type GET_CLUSTERS_RESPONSE = 12
//  kind MESSAGE_KIND = 1
//  Clusters []*serialize.ClusterInfo
type GetClustersResponse struct {
	Clusters []*serialize.ClusterInfo
}

func (res *GetClustersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ClusterInfo{}
		info.Parse(decoder, version, r)
		res.Clusters = append(res.Clusters, info)
	}

}

func (_ *GetClustersResponse) Type() EndpointMessageType {
	return GET_CLUSTERS_RESPONSE
}

var _ EndpointRequestMessage = (*GetClusterInfoRequest)(nil)

// GetClusterInfoRequest получение информации о кластере
//
//  type GET_CLUSTER_INFO_REQUEST = 13
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClusterInfoRequest struct {
	ClusterID uuid.UUID
	response  *GetClusterInfoResponse
}

func (r *GetClusterInfoRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *GetClusterInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ *GetClusterInfoRequest) Type() EndpointMessageType {
	return GET_CLUSTER_INFO_REQUEST
}

// GetClustersResponse ответ со списком кластеров
//
//  type GET_CLUSTERS_RESPONSE = 14
//  Info serialize.ClusterInfo
type GetClusterInfoResponse struct {
	Info serialize.ClusterInfo
}

func (res *GetClusterInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := serialize.ClusterInfo{}
	info.Parse(decoder, version, r)
	res.Info = info
}

func (_ *GetClusterInfoResponse) Type() EndpointMessageType {
	return GET_CLUSTER_INFO_RESPONSE
}

// RegClusterRequest регистрация нового кластера
//
//  type REG_CLUSTER_REQUEST
//  respond GetClustersResponse
type RegClusterRequest struct {
	Info serialize.ClusterInfo
}

func (r *RegClusterRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *RegClusterRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	r.Info.Format(encoder, version, w)
}

func (_ *RegClusterRequest) Type() EndpointMessageType {
	return REG_CLUSTER_REQUEST
}

// GetClustersResponse ответ id созданного кластера кластеров
//
//  type REG_CLUSTER_RESPONSE
//  ClusterID uuid.UUID
type RegClusterResponse struct {
	ClusterID uuid.UUID
}

func (res *RegClusterResponse) Parse(decoder codec.Decoder, r io.Reader) {
	decoder.UuidPtr(&res.ClusterID, r)
}

func (_ *RegClusterResponse) Type() EndpointMessageType {
	return REG_CLUSTER_RESPONSE
}

// UnregClusterRequest регистрация нового кластера
//
//  type REG_CLUSTER_REQUEST
//  respond GetClustersResponse
type UnregClusterRequest struct {
	ClusterID uuid.UUID
}

func (r *UnregClusterRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *UnregClusterRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ *UnregClusterRequest) Type() EndpointMessageType {
	return UNREG_CLUSTER_REQUEST
}

// GetClusterAdminsRequest получение списка админов кластера
//
//  type GET_CLUSTER_ADMINS_REQUEST
//  respond GetClusterAdminsResponse
type GetClusterAdminsRequest struct {
	ClusterID uuid.UUID
}

func (r *GetClusterAdminsRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r *GetClusterAdminsRequest) Format(e codec.Encoder, _ int, w io.Writer) {
	e.Uuid(r.ClusterID, w)
}

func (_ *GetClusterAdminsRequest) Type() EndpointMessageType {
	return GET_CLUSTER_ADMINS_REQUEST
}

// GetAgentAdminsResponse ответ со списком админов агента кластера
//
//  type GET_CLUSTER_ADMINS_RESPONSE
//  Users serialize.UsersList
type GetClusterAdminsResponse struct {
	Users serialize.UsersList
}

func (res *GetClusterAdminsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.UsersList{}
	list.Parse(decoder, version, r)

	res.Users = list

}

func (_ *GetClusterAdminsResponse) Type() EndpointMessageType {
	return GET_CLUSTER_ADMINS_RESPONSE
}

// RegClusterAdminRequest регистрация админа кластера
//
//  type REG_CLUSTER_ADMIN_REQUEST
type RegClusterAdminRequest struct {
	ClusterID uuid.UUID
	User      serialize.UserInfo
}

func (r *RegClusterAdminRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *RegClusterAdminRequest) Format(e codec.Encoder, v int, w io.Writer) {
	e.Uuid(r.ClusterID, w)
	r.User.Format(e, v, w)
}

func (_ *RegClusterAdminRequest) Type() EndpointMessageType {
	return REG_CLUSTER_ADMIN_REQUEST
}

// UnregClusterAdminRequest удаление админа кластера
//
//  type REG_AGENT_ADMIN_REQUEST
type UnregClusterAdminRequest struct {
	ClusterID uuid.UUID
	User      string
}

func (r *UnregClusterAdminRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *UnregClusterAdminRequest) Format(e codec.Encoder, v int, w io.Writer) {
	e.Uuid(r.ClusterID, w)
	e.String(r.User, w)
}

func (_ *UnregClusterAdminRequest) Type() EndpointMessageType {
	return UNREG_CLUSTER_ADMIN_REQUEST
}
