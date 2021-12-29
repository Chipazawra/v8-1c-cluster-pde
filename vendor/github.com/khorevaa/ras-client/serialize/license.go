package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"strings"
)

type LicenseInfoList []*LicenseInfo

func (l LicenseInfoList) BySeries(series string) (LicenseInfoList, bool) {

	if len(series) == 0 {
		return LicenseInfoList{}, false
	}

	fn := func(info *LicenseInfo) bool {
		return strings.EqualFold(info.Series, series)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l LicenseInfoList) ByProcess(id uuid.UUID) (LicenseInfoList, bool) {

	if id == uuid.Nil {
		return LicenseInfoList{}, false
	}

	fn := func(info *LicenseInfo) bool {
		return uuid.Equal(info.ProcessID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l LicenseInfoList) BySession(id uuid.UUID) (LicenseInfoList, bool) {

	if id == uuid.Nil {
		return LicenseInfoList{}, false
	}

	fn := func(info *LicenseInfo) bool {
		return uuid.Equal(info.SessionID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l LicenseInfoList) Filter(fn func(info *LicenseInfo) bool) LicenseInfoList {

	return l.filter(fn, 0)

}

func (l LicenseInfoList) Each(fn func(info *LicenseInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l LicenseInfoList) filter(fn func(info *LicenseInfo) bool, count int) (val LicenseInfoList) {

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

func (l *LicenseInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls LicenseInfoList

	for i := 0; i < count; i++ {

		info := &LicenseInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type LicenseInfo struct {
	ProcessID         uuid.UUID `json:"process_id" example:"94232f94-be78-4acd-a11e-09911bd4f4ed"`
	SessionID         uuid.UUID `json:"session_id" example:"e45c1c2b-b3ac-4fea-9f0c-0583ad65d117"`
	UserName          string    `json:"user_name" example:"User"`
	Host              string    `json:"host" example:"host"`
	AppId             string    `json:"app_id" example:"1CV8"`
	FullName          string    `json:"full_name" example:"FullName"`
	Series            string    `json:"series" example:"ORG8A"`
	IssuedByServer    bool      `json:"issued_by_server" example:"true"`
	LicenseType       int       `json:"license_type" example:"0"` // license-type       : HASP
	Net               bool      `json:"net" example:"true"`
	MaxUsersAll       int       `json:"max_users_all" example:"300"`
	MaxUsersCur       int       `json:"max_users_cur" example:"300"`
	RmngrAddress      string    `json:"rmngr_address" example:"app"`
	RmngrPort         int       `json:"rmngr_port" example:"1541"`
	RmngrPid          string    `json:"rmngr_pid" example:"2300"`
	ShortPresentation string    `json:"short_presentation" example:"Сервер, ORG8A Сет 300"`
	FullPresentation  string    `json:"full_presentation" example:"Сервер, 2300, app, 1541, ORG8A Сетевой 300"`
}

func (i *LicenseInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&i.FullName, r)
	decoder.StringPtr(&i.FullPresentation, r)
	decoder.BoolPtr(&i.IssuedByServer, r)
	decoder.IntPtr(&i.LicenseType, r)
	decoder.IntPtr(&i.MaxUsersAll, r)
	decoder.IntPtr(&i.MaxUsersCur, r)
	decoder.BoolPtr(&i.Net, r)
	decoder.StringPtr(&i.RmngrAddress, r)
	decoder.StringPtr(&i.RmngrPid, r)
	decoder.IntPtr(&i.RmngrPort, r)
	decoder.StringPtr(&i.Series, r)
	decoder.StringPtr(&i.ShortPresentation, r)

}
