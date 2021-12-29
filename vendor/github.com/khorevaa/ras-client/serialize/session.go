package serialize

import (
	"context"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
	"io"
	"sync"
	"time"
)

type SessionCloser interface {
	TerminateSession(ctx context.Context, cluster uuid.UUID, session uuid.UUID, msg string) error
}

type SessionnSig interface {
	Sig() (cluster uuid.UUID, session uuid.UUID)
}

type SessionInfoList []*SessionInfo

func (l SessionInfoList) ByID(id uuid.UUID) (*SessionInfo, bool) {

	if id == uuid.Nil {
		return nil, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.UUID, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return nil, false

}

func (l SessionInfoList) ByProcess(id uuid.UUID) (SessionInfoList, bool) {

	if id == uuid.Nil {
		return SessionInfoList{}, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.ProcessID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) ByInfobase(id uuid.UUID) (SessionInfoList, bool) {

	if id == uuid.Nil {
		return SessionInfoList{}, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.InfobaseID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) Find(fn func(info *SessionInfo) bool) (SessionInfoList, bool) {

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) First(fn func(info *SessionInfo) bool) (*SessionInfo, bool) {

	val := l.filter(fn, 1)

	if len(val) == 0 {
		return nil, false
	}

	return val[0], true

}

func (l SessionInfoList) Filter(fn func(info *SessionInfo) bool) SessionInfoList {

	return l.filter(fn, 0)

}

func (l SessionInfoList) Each(fn func(info *SessionInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l SessionInfoList) TerminateSessions(ctx context.Context, closer SessionCloser, msg string) error {

	var mErr *multierror.Error
	var muErr sync.Mutex
	var wg sync.WaitGroup
	l.Each(func(info *SessionInfo) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			errDisconnect := closer.TerminateSession(ctx, info.ClusterID, info.UUID, msg)

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

func (l SessionInfoList) filter(fn func(info *SessionInfo) bool, count int) (val SessionInfoList) {

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

func (l *SessionInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls SessionInfoList

	for i := 0; i < count; i++ {

		info := &SessionInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type SessionInfo struct {
	UUID                          uuid.UUID        `rac:"session" json:"uuid" example:"1fb5f037-99e8-4924-a99d-a9e687522d32"`
	ID                            int              `rac:"session-id" json:"id" example:"12"`
	InfobaseID                    uuid.UUID        `json:"infobase_id" example:"aea71760-15b3-485a-9a35-506eb8a0b04a"`
	ConnectionID                  uuid.UUID        `json:"connection_id" example:"8adf4514-0379-4333-a153-0b2689edc415"`
	ProcessID                     uuid.UUID        `json:"process_id" example:"1af2e54f-d95a-4370-9b45-8277280cad23"`
	UserName                      string           `json:"user_name" example:"АКузнецов"`
	Host                          string           `json:"host" example:"host"`
	AppId                         string           `json:"app_id" example:"Designer"`
	Locale                        string           `json:"locale" example:"ru_RU"`
	StartedAt                     time.Time        `json:"started_at" example:"2018-04-09T14:51:31"`
	LastActiveAt                  time.Time        `json:"last_active_at" example:"2018-04-09T14:51:31"`
	Hibernate                     bool             `json:"hibernate" example:"true"`
	PassiveSessionHibernateTime   int              `json:"passive_session_hibernate_time" example:"1200"`
	HibernateDessionTerminateTime int              `json:"hibernate_dession_terminate_time" example:"86400"`
	BlockedByDbms                 int              `json:"blocked_by_dbms" example:"0"`
	BlockedByLs                   int              `json:"blocked_by_ls" example:"0"`
	BytesAll                      int64            `json:"bytes_all" example:"105972550"`
	BytesLast5min                 int64            `rac:"bytes-last-5min" json:"bytes_last_5_min" example:"0"`
	CallsAll                      int              `json:"calls_all" example:"119052"`
	CallsLast5min                 int64            `rac:"calls-last-5min" json:"calls_last_5_min" example:"0"`
	DbmsBytesAll                  int64            `json:"dbms_bytes_all" example:"317824922"`
	DbmsBytesLast5min             int64            `rac:"dbms-bytes-last-5min" json:"dbms_bytes_last_5_min" example:"0"`
	DbProcInfo                    string           `json:"db_proc_info" example:"DbProcInfo"`
	DbProcTook                    int              `json:"db_proc_took" example:"0"`
	DbProcTookAt                  time.Time        `json:"db_proc_took_at" example:"2018-04-09T14:51:31"`
	DurationAll                   int              `json:"duration_all" example:"66184"`
	DurationAllDbms               int              `json:"duration_all_dbms" example:"43242"`
	DurationCurrent               int              `json:"duration_current" example:"0"`
	DurationCurrentDbms           int              `json:"duration_current_dbms" example:"0"`
	DurationLast5Min              int64            `rac:"duration-last-5min" json:"duration_last_5_min" example:"0"`
	DurationLast5MinDbms          int64            `rac:"duration-last-5min-dbms" json:"duration_last_5_min_dbms" example:"0"`
	MemoryCurrent                 int64            `json:"memory_current" example:"0"`
	MemoryLast5min                int64            `json:"memory_last_5_min" example:"416379"`
	MemoryTotal                   int64            `json:"memory_total" example:"23178863"`
	ReadCurrent                   int64            `json:"read_current" example:"156162"`
	ReadLast5min                  int64            `json:"read_last_5_min" example:"156162"`
	ReadTotal                     int64            `json:"read_total" example:"15616"`
	WriteCurrent                  int64            `json:"write_current" example:"0"`
	WriteLast5min                 int64            `json:"write_last_5_min" example:"123"`
	WriteTotal                    int64            `json:"write_total" example:"1071457"`
	DurationCurrentService        int              `json:"duration_current_service" example:"0"`
	DurationLast5minService       int64            `json:"duration_last_5_min_service" example:"30"`
	DurationAllService            int              `json:"duration_all_service" example:"515"`
	CurrentServiceName            string           `json:"current_service_name" example:"name"`
	CpuTimeCurrent                int64            `json:"cpu_time_current" example:"0"`
	CpuTimeLast5min               int64            `json:"cpu_time_last_5_min" example:"280"`
	CpuTimeTotal                  int64            `json:"cpu_time_total" example:"6832"`
	DataSeparation                string           `json:"data_separation" example:"sep=1"`
	ClientIPAddress               string           `json:"client_ip_address" example:"127.0.0.1"`
	Licenses                      *LicenseInfoList `json:"licenses"`
	ClusterID                     uuid.UUID        `json:"cluster_id" example:"0e588a25-8354-4344-b935-53442312aa30"`
}

func (i SessionInfo) Sig() (cluster uuid.UUID, session uuid.UUID) {
	return i.ClusterID, i.UUID
}

func (i *SessionInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.AppId, r)
	decoder.IntPtr(&i.BlockedByDbms, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.Int64Ptr(&i.BytesAll, r)
	decoder.Int64Ptr(&i.BytesLast5min, r)
	decoder.IntPtr(&i.CallsAll, r)
	decoder.Int64Ptr(&i.CallsLast5min, r)
	decoder.UuidPtr(&i.ConnectionID, r)
	decoder.Int64Ptr(&i.DbmsBytesAll, r)
	decoder.Int64Ptr(&i.DbmsBytesLast5min, r)
	decoder.StringPtr(&i.DbProcInfo, r)
	decoder.IntPtr(&i.DbProcTook, r)
	decoder.TimePtr(&i.DbProcTookAt, r)
	decoder.IntPtr(&i.DurationAll, r)
	decoder.IntPtr(&i.DurationAllDbms, r)
	decoder.IntPtr(&i.DurationCurrent, r)
	decoder.IntPtr(&i.DurationCurrentDbms, r)
	decoder.Int64Ptr(&i.DurationLast5Min, r)
	decoder.Int64Ptr(&i.DurationLast5MinDbms, r)

	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.InfobaseID, r)
	decoder.TimePtr(&i.LastActiveAt, r)
	decoder.BoolPtr(&i.Hibernate, r)
	decoder.IntPtr(&i.PassiveSessionHibernateTime, r)
	decoder.IntPtr(&i.HibernateDessionTerminateTime, r)

	licenseList := &LicenseInfoList{}
	licenseList.Parse(decoder, version, r)
	i.Licenses = licenseList

	decoder.StringPtr(&i.Locale, r)
	decoder.UuidPtr(&i.ProcessID, r)
	decoder.IntPtr(&i.ID, r)
	decoder.TimePtr(&i.StartedAt, r)
	decoder.StringPtr(&i.UserName, r)

	if version >= 4 {
		decoder.Int64Ptr(&i.MemoryCurrent, r)
		decoder.Int64Ptr(&i.MemoryLast5min, r)
		decoder.Int64Ptr(&i.MemoryTotal, r)
		decoder.Int64Ptr(&i.ReadCurrent, r)
		decoder.Int64Ptr(&i.ReadLast5min, r)
		decoder.Int64Ptr(&i.ReadTotal, r)
		decoder.Int64Ptr(&i.WriteCurrent, r)
		decoder.Int64Ptr(&i.WriteLast5min, r)
		decoder.Int64Ptr(&i.WriteTotal, r)
	}

	if version >= 5 {
		decoder.IntPtr(&i.DurationCurrentService, r)
		decoder.Int64Ptr(&i.DurationLast5minService, r)
		decoder.IntPtr(&i.DurationAllService, r)
		decoder.StringPtr(&i.CurrentServiceName, r)
	}

	if version >= 6 {
		decoder.Int64Ptr(&i.CpuTimeCurrent, r)
		decoder.Int64Ptr(&i.CpuTimeLast5min, r)
		decoder.Int64Ptr(&i.CpuTimeTotal, r)
	}

	if version >= 7 {
		decoder.StringPtr(&i.DataSeparation, r)
	}

	if version >= 10 {
		decoder.StringPtr(&i.ClientIPAddress, r)
	}

	i.Licenses.Each(func(info *LicenseInfo) {
		info.SessionID = i.UUID
		info.ProcessID = i.ProcessID
		info.Host = i.Host
		info.AppId = i.AppId
		info.UserName = i.UserName
	})

}
