package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ serversApi = (*Client)(nil)

func (c *Client) GetWorkingServers(ctx context.Context, clusterID uuid.UUID) ([]*serialize.ServerInfo, error) {
	req := &messages.GetWorkingServersRequest{
		ClusterID: clusterID,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetWorkingServersResponse)
	for _, server := range response.Servers {
		server.ClusterID = clusterID
	}

	return response.Servers, err
}

func (c *Client) GetWorkingServerInfo(ctx context.Context, clusterID, serverID uuid.UUID) (*serialize.ServerInfo, error) {
	req := &messages.GetWorkingServerInfoRequest{
		ClusterID: clusterID,
		ServerID:  serverID,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetWorkingServerInfoResponse)
	return response.Info, err
}

func (c *Client) RegWorkingServer(ctx context.Context, clusterID uuid.UUID, info *serialize.ServerInfo) (*serialize.ServerInfo, error) {

	req := &messages.RegWorkingServerRequest{
		ClusterID: clusterID,
		Info:      info,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.RegWorkingServerResponse)

	return c.GetWorkingServerInfo(ctx, clusterID, response.ServerID)
}

func (c *Client) UnRegWorkingServer(ctx context.Context, clusterID, serverID uuid.UUID) error {
	req := &messages.UnregWorkingServerRequest{
		ClusterID: clusterID,
		ServerID:  serverID,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}

	return nil
}
