package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ infobaseApi = (*Client)(nil)

func (c *Client) CreateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error) {

	req := &messages.CreateInfobaseRequest{
		ClusterID: cluster,
		Infobase:  &infobase,
		Mode:      mode,
	}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return serialize.InfobaseInfo{}, err
	}

	response := resp.(*messages.CreateInfobaseResponse)

	return c.GetInfobaseInfo(ctx, cluster, response.InfobaseID)
}

func (c *Client) DropInfobase(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, mode int) error {

	req := &messages.DropInfobaseRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
		Mode:       mode,
	}
	_, err := c.sendEndpointRequest(ctx, req)

	return err

}

func (c *Client) UpdateSummaryInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error {

	req := &messages.UpdateInfobaseShortRequest{ClusterID: cluster, Infobase: infobase}
	_, err := c.sendEndpointRequest(ctx, req)

	return err
}

func (c *Client) UpdateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo) error {

	req := &messages.UpdateInfobaseRequest{ClusterID: cluster, Infobase: infobase}
	_, err := c.sendEndpointRequest(ctx, req)
	return err

}

func (c *Client) GetInfobaseInfo(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error) {

	req := &messages.GetInfobaseInfoRequest{ClusterID: cluster, InfobaseID: infobase}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return serialize.InfobaseInfo{}, err
	}

	response := resp.(*messages.GetInfobaseInfoResponse)
	response.Infobase.ClusterID = cluster

	return response.Infobase, err
}
