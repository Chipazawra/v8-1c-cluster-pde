package rclient

/*

  locks api
  file: locks.go

	GetClusterLocks(ctx context.Context, cluster uuid.UUID) (serialize.LocksList, error)
	GetInfobaseLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error)
	GetSessionLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error)
	GetConnectionLocks(ctx context.Context, cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error)


*/
