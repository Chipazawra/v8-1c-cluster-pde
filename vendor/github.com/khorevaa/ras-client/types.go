package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

type Api interface {
	Version() string

	Close() error

	agentApi
	authApi
	clusterApi
	sessionApi
	locksApi
	connectionsApi
	infobaseApi
	processesApi
	serversApi
}

type authApi interface {
	AuthenticateAgent(user, password string)
	AuthenticateCluster(cluster uuid.UUID, user, password string)
	AuthenticateInfobase(cluster uuid.UUID, user, password string)
}

type agentApi interface {
	GetAgentVersion(ctx context.Context) (string, error)
	GetAgentAdmins(ctx context.Context) (serialize.UsersList, error)
	RegAgentAdmin(ctx context.Context, user serialize.UserInfo) error
	UnregAgentAdmin(ctx context.Context, user string) error
}

type clusterApi interface {
	RegCluster(ctx context.Context, info serialize.ClusterInfo) (uuid.UUID, error)
	UnregCluster(ctx context.Context, clusterId uuid.UUID) error
	GetClusterAdmins(ctx context.Context, clusterID uuid.UUID) (serialize.UsersList, error)
	RegClusterAdmin(ctx context.Context, clusterID uuid.UUID, user serialize.UserInfo) error
	UnregClusterAdmin(ctx context.Context, clusterID uuid.UUID, user string) error
	GetClusters(ctx context.Context) ([]*serialize.ClusterInfo, error)
	GetClusterInfo(ctx context.Context, cluster uuid.UUID) (serialize.ClusterInfo, error)
	GetClusterInfobases(ctx context.Context, id uuid.UUID) (serialize.InfobaseSummaryList, error)
	GetClusterServices(ctx context.Context, id uuid.UUID) ([]*serialize.ServiceInfo, error)
	GetClusterManagers(ctx context.Context, id uuid.UUID) ([]*serialize.ManagerInfo, error)
}

type sessionApi interface {
	GetClusterSessions(ctx context.Context, cluster uuid.UUID) (serialize.SessionInfoList, error)
	GetInfobaseSessions(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error)
	TerminateSession(ctx context.Context, cluster uuid.UUID, session uuid.UUID, msg string) error
}

type connectionsApi interface {
	GetClusterConnections(ctx context.Context, uuid uuid.UUID) (serialize.ConnectionShortInfoList, error)
	GetInfobaseConnections(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error)
	DisconnectConnection(ctx context.Context, cluster uuid.UUID, process uuid.UUID, connection uuid.UUID, infobase uuid.UUID) error
}

type locksApi interface {
	GetClusterLocks(ctx context.Context, cluster uuid.UUID) (serialize.LocksList, error)
	GetInfobaseLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error)
	GetSessionLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error)
	GetConnectionLocks(ctx context.Context, cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error)
}

type processesApi interface {
	GetWorkingProcesses(ctx context.Context, clusterID uuid.UUID) (serialize.ProcessInfoList, error)
	GetWorkingProcessInfo(ctx context.Context, clusterID, processID uuid.UUID) (*serialize.ProcessInfo, error)
}

type serversApi interface {
	GetWorkingServers(ctx context.Context, clusterID uuid.UUID) ([]*serialize.ServerInfo, error)
	GetWorkingServerInfo(ctx context.Context, clusterID, serverID uuid.UUID) (*serialize.ServerInfo, error)
	RegWorkingServer(ctx context.Context, clusterID uuid.UUID, info *serialize.ServerInfo) (*serialize.ServerInfo, error)
	UnRegWorkingServer(ctx context.Context, clusterID, serverID uuid.UUID) error
}

type infobaseApi interface {
	CreateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error)
	UpdateSummaryInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error
	UpdateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo) error
	DropInfobase(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, mode int) error
	GetInfobaseInfo(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error)
}
