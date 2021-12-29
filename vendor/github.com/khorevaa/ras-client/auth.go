package rclient

import uuid "github.com/satori/go.uuid"

var _ authApi = (*Client)(nil)

func (c *Client) AuthenticateAgent(user, password string) {

	c.base.SetAgentAuth(user, password)

}

func (c *Client) AuthenticateCluster(cluster uuid.UUID, user, password string) {

	c.base.SetClusterAuth(cluster, user, password)

}

func (c *Client) AuthenticateInfobase(cluster uuid.UUID, user, password string) {

	c.base.SetInfobaseAuth(cluster, user, password)

}
