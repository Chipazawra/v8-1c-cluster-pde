package messages

type UnknownMessageError struct {
	Type       byte
	Data       []byte
	EndpointID int
	Err        error
	ServiceID  string
}

func (m *UnknownMessageError) Error() string {

	return m.Err.Error()

}
