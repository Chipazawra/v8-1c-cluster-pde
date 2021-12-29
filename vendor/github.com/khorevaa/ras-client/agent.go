package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
)

var _ agentApi = (*Client)(nil)

func (c *Client) GetAgentVersion(ctx context.Context) (string, error) {

	switch c.serviceVersion {
	case "4.0":
		return "8.3.11", nil
	case "5.0":
		return "8.3.12", nil
	case "6.0":
		return "8.3.13", nil
	case "7.0":
		return "8.3.14", nil
	case "8.0":
		return "8.3.15", nil
	}

	req := &messages.GetAgentVersionRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return "", err
	}

	response := resp.(*messages.GetAgentVersionResponse)

	return response.Version, err
}

func (c *Client) GetAgentAdmins(ctx context.Context) (serialize.UsersList, error) {

	req := &messages.GetAgentAdminsRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetAgentAdminsResponse)

	return response.Users, err
}

func (c *Client) RegAgentAdmin(ctx context.Context, user serialize.UserInfo) error {

	req := &messages.RegAgentAdminRequest{
		User: user,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UnregAgentAdmin(ctx context.Context, user string) error {

	req := &messages.UnregAgentAdminRequest{
		User: user,
	}

	_, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return err
	}
	return nil
}
