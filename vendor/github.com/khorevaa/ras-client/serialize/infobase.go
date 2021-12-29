package serialize

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"io"
	"strings"
	"time"
)

type InfobaseSig interface {
	Sig() (cluster uuid.UUID, infobase uuid.UUID)
}

type InfobaseInfoGetter interface {
	GetInfobaseInfo(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (InfobaseInfo, error)
}

type InfobaseConnectionsGetter interface {
	GetInfobaseConnections(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (ConnectionShortInfoList, error)
}

type InfobaseSessionsGetter interface {
	GetInfobaseSessions(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (SessionInfoList, error)
}

type InfobaseLocksGetter interface {
	GetInfobaseLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (LocksList, error)
}

type InfobaseDropper interface {
	DropInfobase(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, mode int) error
}

type InfobaseSummaryUpdater interface {
	UpdateSummaryInfobase(ctx context.Context, cluster uuid.UUID, info InfobaseSummaryInfo) error
}

type InfobaseUpdater interface {
	UpdateInfobase(ctx context.Context, cluster uuid.UUID, info InfobaseInfo) error
}

type InfobaseBlocker interface {
	InfobaseInfoGetter
	InfobaseUpdater
}

type InfobaseSummaryList []*InfobaseSummaryInfo

func (l InfobaseSummaryList) ByID(id uuid.UUID) (*InfobaseSummaryInfo, bool) {

	if id == uuid.Nil {
		return nil, false
	}

	fn := func(info *InfobaseSummaryInfo) bool {
		return uuid.Equal(info.UUID, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return nil, false

}

func (l InfobaseSummaryList) ByName(name string) (*InfobaseSummaryInfo, bool) {

	if len(name) == 0 {
		return nil, false
	}

	fn := func(info *InfobaseSummaryInfo) bool {
		return strings.EqualFold(info.Name, name)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return nil, false

}

func (l InfobaseSummaryList) Find(fn func(info *InfobaseSummaryInfo) bool) (*InfobaseSummaryInfo, bool) {

	val := l.filter(fn, 1)

	if len(val) == 0 {
		return nil, false
	}

	return val[0], true

}

func (l InfobaseSummaryList) Filter(fn func(info *InfobaseSummaryInfo) bool) InfobaseSummaryList {

	return l.filter(fn, 0)

}

func (l InfobaseSummaryList) Each(fn func(info *InfobaseSummaryInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l InfobaseSummaryList) filter(fn func(info *InfobaseSummaryInfo) bool, count int) (val InfobaseSummaryList) {

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

func (l *InfobaseSummaryList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls InfobaseSummaryList
	for i := 0; i < count; i++ {

		info := &InfobaseSummaryInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type InfobaseSummaryInfo struct {
	ClusterID   uuid.UUID `rac:"-" json:"cluster_id" example:""`
	UUID        uuid.UUID `rac:"infobase" json:"uuid" example:"efa3672f-947a-4d84-bd58-b21997b83561"`
	Name        string    `json:"name" example:"УППБоеваяБаза"`
	Description string    `rac:"descr" json:"descr" example:"УППБоеваяБаза"`
}

func (i InfobaseSummaryInfo) Sig() (uuid.UUID, uuid.UUID) {
	return i.ClusterID, i.UUID
}

func (i InfobaseSummaryInfo) FullInfo(ctx context.Context, runner InfobaseInfoGetter) (InfobaseInfo, error) {
	cluster, infobase := i.Sig()
	return runner.GetInfobaseInfo(ctx, cluster, infobase)
}

func (i InfobaseSummaryInfo) Drop(ctx context.Context, runner InfobaseDropper, mode int) error {
	cluster, infobase := i.Sig()
	return runner.DropInfobase(ctx, cluster, infobase, mode)
}

func (i InfobaseSummaryInfo) Update(ctx context.Context, runner InfobaseSummaryUpdater) error {
	return runner.UpdateSummaryInfobase(ctx, i.ClusterID, i)
}

func (i *InfobaseSummaryInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Name, r)

}

func (i InfobaseSummaryInfo) Format(encoder Encoder, _ int, w io.Writer) {

	encoder.Uuid(i.UUID, w)
	encoder.String(i.Description, w)
	encoder.String(i.Name, w)

}

type InfobaseInfo struct {
	UUID                                   uuid.UUID `rac:"infobase" json:"uuid" example:"efa3672f-947a-4d84-bd58-b21997b83561"`
	Name                                   string    `json:"name" example:"УППБоеваяБаза"`
	Description                            string    `rac:"descr" json:"descr" example:"Это очень хорошая база"`
	Dbms                                   string    `json:"dbms" example:"MSSQLServer"`
	DbServer                               string    `json:"db_server" example:"sql"`
	DbName                                 string    `json:"db_name" example:"base"`
	DbUser                                 string    `json:"db_user" example:"user"`
	DbPwd                                  string    `rac:"-" json:"db_pwd" example:"password"`
	SecurityLevel                          int       `json:"security_level" example:"0"`
	LicenseDistribution                    int       `json:"license_distribution" example:"0"`
	ScheduledJobsDeny                      bool      `json:"scheduled_jobs_deny" example:"false"`
	SessionsDeny                           bool      `json:"sessions_deny" example:"false"`
	DeniedFrom                             time.Time `json:"denied_from" example:"2020-10-01T08:30:00Z"`
	DeniedMessage                          string    `json:"denied_message" example:"Выполняется обновление базы"`
	DeniedParameter                        string    `json:"denied_parameter" example:"123"`
	DeniedTo                               time.Time `json:"denied_to" example:"2020-10-01T08:30:00Z"`
	PermissionCode                         string    `json:"permission_code" example:"123"`
	ExternalSessionManagerConnectionString string    `json:"external_session_manager_connection_string" example:"http://auth2.com"`
	ExternalSessionManagerRequired         bool      `json:"external_session_manager_required" example:"false"`
	SecurityProfileName                    string    `json:"security_profile_name"  example:"sec_profile1"`
	SafeModeSecurityProfileName            string    `json:"safe_mode_security_profile_name" example:"profile1"`
	ReserveWorkingProcesses                bool      `json:"reserve_working_processes" example:"false"`
	DateOffset                             int       `json:"date_offset" example:"0"`
	Locale                                 string    `json:"locale" example:"ru_RU"`
	ClusterID                              uuid.UUID `rac:"-" json:"cluster_id" example:"efa3672f-947a-4d84-bd58-b21997b83561"`

	Cluster     *ClusterInfo             `json:"-" swaggerignore:"true"`
	Connections *ConnectionShortInfoList `json:"-" swaggerignore:"true"`
	Sessions    *SessionInfoList         `json:"-" swaggerignore:"true"`
	Locks       *LocksList               `json:"-" swaggerignore:"true"`
}

func (i *InfobaseInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.IntPtr(&i.DateOffset, r)
	decoder.StringPtr(&i.Dbms, r)
	decoder.StringPtr(&i.DbName, r)
	decoder.StringPtr(&i.DbPwd, r)
	decoder.StringPtr(&i.DbServer, r)
	decoder.StringPtr(&i.DbUser, r)
	decoder.TimePtr(&i.DeniedFrom, r)
	decoder.StringPtr(&i.DeniedMessage, r)
	decoder.StringPtr(&i.DeniedParameter, r)
	decoder.TimePtr(&i.DeniedTo, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Locale, r)
	decoder.StringPtr(&i.Name, r)
	decoder.StringPtr(&i.PermissionCode, r)
	decoder.BoolPtr(&i.ScheduledJobsDeny, r)
	decoder.IntPtr(&i.SecurityLevel, r)
	decoder.BoolPtr(&i.SessionsDeny, r)
	decoder.IntPtr(&i.LicenseDistribution, r)
	decoder.StringPtr(&i.ExternalSessionManagerConnectionString, r)
	decoder.BoolPtr(&i.ExternalSessionManagerRequired, r)
	decoder.StringPtr(&i.SecurityProfileName, r)
	decoder.StringPtr(&i.SafeModeSecurityProfileName, r)
	if version >= 9 {
		decoder.BoolPtr(&i.ReserveWorkingProcesses, r)
	}

}

func (i InfobaseInfo) Format(encoder Encoder, version int, w io.Writer) {

	encoder.Uuid(i.UUID, w)
	encoder.Int(i.DateOffset, w)
	encoder.String(i.Dbms, w)
	encoder.String(i.DbName, w)
	encoder.String(i.DbPwd, w)
	encoder.String(i.DbServer, w)
	encoder.String(i.DbUser, w)
	encoder.Time(i.DeniedFrom, w)
	encoder.String(i.DeniedMessage, w)
	encoder.String(i.DeniedParameter, w)
	encoder.Time(i.DeniedTo, w)
	encoder.String(i.Description, w)
	encoder.String(i.Locale, w)
	encoder.String(i.Name, w)
	encoder.String(i.PermissionCode, w)
	encoder.Bool(i.ScheduledJobsDeny, w)
	encoder.Int(i.SecurityLevel, w)
	encoder.Bool(i.SessionsDeny, w)
	encoder.Int(i.LicenseDistribution, w)
	encoder.String(i.ExternalSessionManagerConnectionString, w)
	encoder.Bool(i.ExternalSessionManagerRequired, w)
	encoder.String(i.SecurityProfileName, w)
	encoder.String(i.SafeModeSecurityProfileName, w)
	if version >= 9 {
		encoder.Bool(i.ReserveWorkingProcesses, w)
	}

}

func (i InfobaseInfo) Summary() InfobaseSummaryInfo {
	return InfobaseSummaryInfo{
		ClusterID:   i.ClusterID,
		UUID:        i.UUID,
		Name:        i.Name,
		Description: i.Description,
	}

}

func (i InfobaseInfo) Sig() (uuid.UUID, uuid.UUID) {
	return i.ClusterID, i.UUID
}

func (i InfobaseInfo) Drop(ctx context.Context, runner InfobaseDropper, mode int) error {
	cluster, infobase := i.Sig()

	return runner.DropInfobase(ctx, cluster, infobase, mode)
}

func (i *InfobaseInfo) Update(ctx context.Context, runner InfobaseUpdater) error {

	return runner.UpdateInfobase(ctx, i.ClusterID, *i)

}

func (i *InfobaseInfo) Reload(ctx context.Context, runner InfobaseInfoGetter) error {
	cluster, infobase := i.Sig()
	newInfo, err := runner.GetInfobaseInfo(ctx, cluster, infobase)
	if err != nil {
		return err
	}

	*i = newInfo

	return nil

}

func (i *InfobaseInfo) GetConnections(ctx context.Context, runner InfobaseConnectionsGetter) (*ConnectionShortInfoList, error) {
	cluster, infobase := i.Sig()
	list, err := runner.GetInfobaseConnections(ctx, cluster, infobase)
	if err != nil {
		return nil, err
	}
	i.Connections = &list

	return i.Connections, err
}

func (i *InfobaseInfo) GetSessions(ctx context.Context, runner InfobaseSessionsGetter) (*SessionInfoList, error) {

	cluster, infobase := i.Sig()
	list, err := runner.GetInfobaseSessions(ctx, cluster, infobase)
	if err != nil {
		return nil, err
	}
	i.Sessions = &list

	return i.Sessions, err
}

func (i *InfobaseInfo) GetLocks(ctx context.Context, runner InfobaseLocksGetter) (*LocksList, error) {

	cluster, infobase := i.Sig()
	list, err := runner.GetInfobaseLocks(ctx, cluster, infobase)
	if err != nil {
		return nil, err
	}
	i.Locks = &list

	return i.Locks, err
}

func (i *InfobaseInfo) Blocker(reload bool) BlockerInfobase {

	return BlockerInfobase{
		infobase: i,
		Reload:   reload,
	}

}

type BlockerInfobase struct {
	DeniedFrom        time.Time
	DeniedTo          time.Time
	Message           string
	PermissionCode    string
	ScheduledJobsDeny bool
	Reload            bool

	runner   InfobaseBlocker
	infobase *InfobaseInfo
}

func (b *BlockerInfobase) Msg(msg string) *BlockerInfobase {

	b.Message = msg
	return b
}

func (b *BlockerInfobase) Code(code string) *BlockerInfobase {

	b.PermissionCode = code
	return b
}

func (b *BlockerInfobase) From(from time.Time) *BlockerInfobase {

	b.DeniedFrom = from
	return b
}

func (b *BlockerInfobase) To(to time.Time) *BlockerInfobase {

	b.DeniedTo = to
	return b
}

func (b *BlockerInfobase) ScheduledJobs(deny bool) *BlockerInfobase {

	b.ScheduledJobsDeny = deny
	return b
}

func (b *BlockerInfobase) Block(ctx context.Context, runner InfobaseBlocker) error {

	b.runner = runner

	blockInfo := *b.infobase
	blockInfo.DeniedTo = b.DeniedTo
	blockInfo.DeniedFrom = b.DeniedFrom

	if len(b.Message) > 0 {
		blockInfo.DeniedMessage = b.Message
	}

	if !blockInfo.ScheduledJobsDeny && b.ScheduledJobsDeny {
		blockInfo.ScheduledJobsDeny = b.ScheduledJobsDeny
	}
	blockInfo.SessionsDeny = true
	blockInfo.PermissionCode = b.PermissionCode

	return runner.UpdateInfobase(ctx, blockInfo.ClusterID, blockInfo)

}

func (b BlockerInfobase) UnblockWithRunner(ctx context.Context, runner InfobaseBlocker) error {

	unblockInfo := *b.infobase
	unblockInfo.SessionsDeny = false
	err := runner.UpdateInfobase(ctx, unblockInfo.ClusterID, unblockInfo)

	if err != nil {
		return err
	}

	if b.Reload {
		err = b.infobase.Reload(ctx, runner)
	}

	return err

}

func (b BlockerInfobase) Unblock(ctx context.Context) error {

	return b.UnblockWithRunner(ctx, b.runner)

}
