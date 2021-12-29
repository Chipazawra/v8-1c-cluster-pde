package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ locksApi = (*Client)(nil)

func (c *Client) GetClusterLocks(ctx context.Context, cluster uuid.UUID) (serialize.LocksList, error) {

	req := &messages.GetLocksRequest{
		ClusterID: cluster,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetLocksResponse)

	response.List.Each(func(info *serialize.LockInfo) {
		info.ClusterID = cluster
	})
	return response.List, err

}

func (c *Client) GetInfobaseLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error) {

	req := &messages.GetInfobaseLockRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobaseLockResponse)

	response.List.Each(func(info *serialize.LockInfo) {
		info.ClusterID = cluster
		info.InfobaseID = infobase
	})

	return response.List, err

}

func (c *Client) GetSessionLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error) {

	req := &messages.GetSessionLockRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
		SessionID:  session,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetSessionLockResponse)

	response.List.Each(func(info *serialize.LockInfo) {
		info.ClusterID = cluster
		info.InfobaseID = infobase
	})

	return response.List, err

}

func (c *Client) GetConnectionLocks(ctx context.Context, cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error) {

	req := &messages.GetConnectionLockRequest{
		ClusterID:    cluster,
		ConnectionID: connection,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetConnectionLockResponse)

	response.List.Each(func(info *serialize.LockInfo) {
		info.ClusterID = cluster
	})

	return response.List, err

}
