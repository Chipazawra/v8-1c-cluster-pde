package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

func (c *Client) GetWorkingProcesses(ctx context.Context, clusterID uuid.UUID) (serialize.ProcessInfoList, error) {

	req := &messages.GetWorkingProcessesRequest{
		ClusterID: clusterID,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetWorkingProcessesResponse)
	response.Processes.Each(func(info *serialize.ProcessInfo) {
		info.ClusterID = clusterID
	})

	return response.Processes, err
}

func (c *Client) GetWorkingProcessInfo(ctx context.Context, clusterID, processID uuid.UUID) (*serialize.ProcessInfo, error) {

	req := &messages.GetWorkingProcessInfoRequest{
		ClusterID: clusterID,
		ProcessID: processID,
	}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetWorkingProcessInfoResponse)
	info := response.Info
	info.ClusterID = clusterID
	return info, err
}
