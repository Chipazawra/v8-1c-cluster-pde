package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ connectionsApi = (*Client)(nil)

func (c *Client) GetClusterConnections(ctx context.Context, id uuid.UUID) (serialize.ConnectionShortInfoList, error) {

	req := &messages.GetConnectionsShortRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetConnectionsShortResponse)

	response.Connections.Each(func(info *serialize.ConnectionShortInfo) {
		info.ClusterID = id
	})
	return response.Connections, err
}

func (c *Client) DisconnectConnection(ctx context.Context, cluster uuid.UUID, process uuid.UUID, connection uuid.UUID, infobase uuid.UUID) error {

	req := &messages.DisconnectConnectionRequest{
		ClusterID:    cluster,
		ProcessID:    process,
		ConnectionID: connection,
		InfobaseID:   infobase,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	return err
}

func (c *Client) GetInfobaseConnections(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error) {

	req := &messages.GetInfobaseConnectionsShortRequest{ClusterID: cluster, InfobaseID: infobase}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobaseConnectionsShortResponse)

	response.Connections.Each(func(info *serialize.ConnectionShortInfo) {
		info.ClusterID = cluster
		info.InfobaseID = infobase
	})

	return response.Connections, nil
}
