package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ clusterApi = (*Client)(nil)

func (c *Client) GetClusters(ctx context.Context) ([]*serialize.ClusterInfo, error) {

	req := &messages.GetClustersRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClustersResponse)

	return response.Clusters, err
}

func (c *Client) RegCluster(ctx context.Context, info serialize.ClusterInfo) (uuid.UUID, error) {

	req := &messages.RegClusterRequest{
		Info: info,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return uuid.Nil, err
	}

	response := resp.(*messages.RegClusterResponse)

	return response.ClusterID, err
}

func (c *Client) UnregCluster(ctx context.Context, clusterId uuid.UUID) error {

	req := &messages.UnregClusterRequest{
		ClusterID: clusterId,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetClusterInfo(ctx context.Context, cluster uuid.UUID) (serialize.ClusterInfo, error) {

	req := &messages.GetClusterInfoRequest{ClusterID: cluster}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return serialize.ClusterInfo{}, err
	}

	response := resp.(*messages.GetClusterInfoResponse)
	return response.Info, nil
}

func (c *Client) GetClusterManagers(ctx context.Context, id uuid.UUID) ([]*serialize.ManagerInfo, error) {

	req := &messages.GetClusterManagersRequest{ClusterID: id}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterManagersResponse)

	for _, manager := range response.Managers {
		manager.ClusterID = id
	}

	return response.Managers, err
}

func (c *Client) GetClusterServices(ctx context.Context, id uuid.UUID) ([]*serialize.ServiceInfo, error) {

	req := &messages.GetClusterServicesRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterServicesResponse)
	for _, service := range response.Services {
		service.ClusterID = id
	}
	return response.Services, err
}

func (c *Client) GetClusterInfobases(ctx context.Context, id uuid.UUID) (serialize.InfobaseSummaryList, error) {

	req := &messages.GetInfobasesShortRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobasesShortResponse)

	response.Infobases.Each(func(info *serialize.InfobaseSummaryInfo) {
		info.ClusterID = id
	})

	return response.Infobases, err
}

func (c *Client) GetClusterAdmins(ctx context.Context, clusterID uuid.UUID) (serialize.UsersList, error) {

	req := &messages.GetClusterAdminsRequest{
		ClusterID: clusterID,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterAdminsResponse)

	return response.Users, err
}

func (c *Client) RegClusterAdmin(ctx context.Context, clusterID uuid.UUID, user serialize.UserInfo) error {

	req := &messages.RegClusterAdminRequest{
		ClusterID: clusterID,
		User:      user,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UnregClusterAdmin(ctx context.Context, clusterID uuid.UUID, user string) error {

	req := &messages.UnregClusterAdminRequest{
		ClusterID: clusterID,
		User:      user,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}
	return nil
}
