package basc


const (
	NetWorkErr int = iota
	NoItemErr
	ItemPolluted
	OtherErr
)

type BascErr struct {
	Msg string `json:"msg"`
	Code int   `json:"code"`
}

func (be *BascErr)Error() string {
	return be.Msg
}

func NewError(msg string, code int) *BascErr {
	return &BascErr{Msg: msg,Code: code}
}

