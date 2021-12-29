package serialize

import (
	"context"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
	"io"
	"sync"
	"time"
)

type ConnectionCloser interface {
	DisconnectConnection(ctx context.Context, cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error
}

type ConnectionSig interface {
	Sig() (cluster uuid.UUID, process uuid.UUID, connection uuid.UUID)
}

type ConnectionShortInfoList []*ConnectionShortInfo

func (l ConnectionShortInfoList) ByID(id uuid.UUID) (*ConnectionShortInfo, bool) {

	if id == uuid.Nil {
		return nil, false
	}

	fn := func(info *ConnectionShortInfo) bool {
		return uuid.Equal(info.UUID, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return nil, false

}

func (l ConnectionShortInfoList) ByProcess(id uuid.UUID) (ConnectionShortInfoList, bool) {

	if id == uuid.Nil {
		return ConnectionShortInfoList{}, false
	}

	fn := func(info *ConnectionShortInfo) bool {
		return uuid.Equal(info.Process, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l ConnectionShortInfoList) ByInfobase(id uuid.UUID) (ConnectionShortInfoList, bool) {

	if id == uuid.Nil {
		return ConnectionShortInfoList{}, false
	}

	fn := func(info *ConnectionShortInfo) bool {
		return uuid.Equal(info.InfobaseID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l ConnectionShortInfoList) Find(fn func(info *ConnectionShortInfo) bool) (ConnectionShortInfoList, bool) {

	val := l.filter(fn, 0)

	return val, true

}

func (l ConnectionShortInfoList) Filter(fn func(info *ConnectionShortInfo) bool) ConnectionShortInfoList {

	return l.filter(fn, 0)

}

func (l ConnectionShortInfoList) Each(fn func(info *ConnectionShortInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l ConnectionShortInfoList) Disconnect(ctx context.Context, closer ConnectionCloser) (err error) {

	var mErr *multierror.Error
	var muErr sync.Mutex
	var wg sync.WaitGroup
	l.Each(func(info *ConnectionShortInfo) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			cluster, process, connection := info.Sig()
			errDisconnect := closer.DisconnectConnection(ctx, cluster, process, connection)

			if errDisconnect != nil {
				muErr.Lock()
				_ = multierror.Append(mErr, errDisconnect)
				muErr.Unlock()
			}
		}()

	})

	wg.Wait()

	return mErr.ErrorOrNil()

}

func (l ConnectionShortInfoList) filter(fn func(info *ConnectionShortInfo) bool, count int) (val ConnectionShortInfoList) {

	n := 0

	for _, info := range l {

		if n == count && count > 0 {
			break
		}

		result := fn(info)

		if result {
			n += 1
			val = append(val, info)
		}

	}

	return

}

func (l *ConnectionShortInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls ConnectionShortInfoList

	for i := 0; i < count; i++ {

		info := &ConnectionShortInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type ConnectionShortInfo struct {
	UUID        uuid.UUID `rac:"connection" json:"uuid" example:"cd16cde9-6372-4664-ac61-b0ae5cb24478"`
	ID          int       `rac:"conn-id" json:"conn_id" example:"8714"`
	Host        string    `json:"host" example:"host"`
	Process     uuid.UUID `json:"process" example:"94232f94-be78-4acd-a11e-09911bd4f4ed"`
	ClusterID   uuid.UUID `rac:"-" json:"cluster_id" example:"94232f94-be78-4acd-a11e-09911bd4f4ed"`
	InfobaseID  uuid.UUID `json:"infobase_id" example:"efa3672f-947a-4d84-bd58-b21997b83561"`
	Application string    `json:"application" example:"1CV8"`
	ConnectedAt time.Time `json:"connected_at" example:"2020-10-01T07:29:57"`
	SessionID   int       `rac:"session-number" json:"session_id" example:"148542"`
	BlockedByLs int       `json:"blocked_by_ls" example:"0"`
}

func (i *ConnectionShortInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Application, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.TimePtr(&i.ConnectedAt, r)
	decoder.IntPtr(&i.ID, r)
	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.InfobaseID, r)
	decoder.UuidPtr(&i.Process, r)
	decoder.IntPtr(&i.SessionID, r)

}

func (i ConnectionShortInfo) Sig() (cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) {
	return i.ClusterID, i.Process, i.UUID
}

type ConnectionInfoList []*ConnectionInfo

type ConnectionInfo struct {
	UUID        uuid.UUID `rac:"connection" json:"uuid" example:"cd16cde9-6372-4664-ac61-b0ae5cb24478"`
	ID          int       `rac:"conn-id" json:"conn_id" example:"8714"`
	Host        string    `json:"host" example:"host"`
	Process     uuid.UUID `json:"process" example:"94232f94-be78-4acd-a11e-09911bd4f4ed"`
	ClusterID   uuid.UUID `rac:"-" json:"cluster_id" example:"94232f94-be78-4acd-a11e-09911bd4f4ed"`
	InfobaseID  uuid.UUID `json:"infobase_id" example:"efa3672f-947a-4d84-bd58-b21997b83561"`
	Application string    `json:"application" example:"1CV8"`
	ConnectedAt time.Time `json:"connected_at" example:"2020-10-01T07:29:57"`
	SessionID   int       `rac:"session-number" json:"session_id" example:"148542"`
	BlockedByLs int       `json:"blocked_by_ls" example:"0"`
}

func (i *ConnectionInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Application, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.TimePtr(&i.ConnectedAt, r)
	decoder.IntPtr(&i.ID, r)
	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.InfobaseID, r)
	decoder.UuidPtr(&i.Process, r)
	decoder.IntPtr(&i.SessionID, r)

	//builder.appId(decoder.decodeString(buffer)).
	//blockedByDbms(decoder.decodeInt(buffer)).
	//bytesAll(decoder.decodeLong(buffer)).
	//bytesLast5Min(decoder.decodeLong(buffer)).
	//callsAll((int)decoder.decodeLong(buffer)).
	//callsLast5Min(decoder.decodeLong(buffer)).
	//connectedAt(dateFromTicks(decoder.decodeLong(buffer))).
	//connId(decoder.decodeInt(buffer)).
	//dbConnMode(decoder.decodeInt(buffer))
	//.dbmsBytesAll(decoder.decodeLong(buffer)).
	//dbmsBytesLast5Min(decoder.decodeLong(buffer)).
	//dbProcInfo(decoder.decodeString(buffer)).
	//dbProcTook(decoder.decodeInt(buffer)).
	//dbProcTookAt(dateFromTicks(decoder.decodeLong(buffer))).
	//durationAll(decoder.decodeInt(buffer)).
	//durationAllDbms(decoder.decodeInt(buffer)).
	//durationCurrent(decoder.decodeInt(buffer)).
	//durationCurrentDbms(decoder.decodeInt(buffer)).
	//durationLast5Min(decoder.decodeLong(buffer)).
	//durationLast5MinDbms(decoder.decodeLong(buffer)).
	//hostName(decoder.decodeString(buffer)).
	//ibConnMode(decoder.decodeInt(buffer)).
	//threadMode(decoder.decodeInt(buffer)).
	//userName(decoder.decodeString(buffer));
	//if (version >= 4) {
	//	builder.memoryCurrent(decoder.decodeLong(buffer)).
	//	memoryLast5Min(decoder.decodeLong(buffer)).
	//	memoryTotal(decoder.decodeLong(buffer)).
	//	readBytesCurrent(decoder.decodeLong(buffer)).
	//	readBytesLast5Min(decoder.decodeLong(buffer)).
	//	readBytesTotal(decoder.decodeLong(buffer)).
	//	writeBytesCurrent(decoder.decodeLong(buffer)).
	//	writeBytesLast5Min(decoder.decodeLong(buffer)).
	//	writeBytesTotal(decoder.decodeLong(buffer));
	//}
	//return builder.build();

}

func (i ConnectionInfo) Sig() (cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) {
	return i.ClusterID, i.Process, i.UUID
}
