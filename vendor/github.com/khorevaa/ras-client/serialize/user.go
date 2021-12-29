package serialize

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"io"
)

type UsersList []*UserInfo

func (l UsersList) Each(fn func(info *UserInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l *UsersList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls UsersList

	for i := 0; i < count; i++ {

		info := &UserInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type UserInfo struct {
	Name                string `json:"name" example:"admin"`
	Description         string `rac:"descr" json:"descr" example:"дополнительное описание"`
	Password            string `json:"password"`
	PasswordAuthAllowed bool   `json:"password_auth_allowed"`
	SysAuthAllowed      bool   `json:"sys_auth_allowed"`
	SysUserName         string `json:"sys_user_name"`
}

func (i *UserInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&i.Name, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Password, r)
	decoder.BoolPtr(&i.PasswordAuthAllowed, r)
	decoder.BoolPtr(&i.SysAuthAllowed, r)
	decoder.StringPtr(&i.SysUserName, r)

}

func (i *UserInfo) Format(encoder codec.Encoder, v int, w io.Writer) {
	encoder.String(i.Name, w)
	encoder.String(i.Description, w)
	encoder.String(i.Password, w)
	encoder.Bool(i.PasswordAuthAllowed, w)
	encoder.Bool(i.SysAuthAllowed, w)
	encoder.String(i.SysUserName, w)
}
