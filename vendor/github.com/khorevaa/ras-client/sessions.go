package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ sessionApi = (*Client)(nil)

func (c *Client) TerminateSession(ctx context.Context, cluster uuid.UUID, session uuid.UUID, msg string) error {

	req := &messages.TerminateSessionRequest{
		ClusterID: cluster,
		SessionID: session,
		Message:   msg,
	}
	_, err := c.sendEndpointRequest(ctx, req)

	return err
}

func (c *Client) GetInfobaseSessions(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error) {

	req := &messages.GetInfobaseSessionsRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobaseSessionsResponse)

	response.Sessions.Each(func(info *serialize.SessionInfo) {
		info.ClusterID = cluster
		info.InfobaseID = infobase
	})
	return response.Sessions, err

}

func (c *Client) GetClusterSessions(ctx context.Context, cluster uuid.UUID) (serialize.SessionInfoList, error) {

	req := &messages.GetSessionsRequest{
		ClusterID: cluster,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetSessionsResponse)

	response.Sessions.Each(func(info *serialize.SessionInfo) {
		info.ClusterID = cluster
	})
	return response.Sessions, err

}
