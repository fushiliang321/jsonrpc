package common

const (
	WithoutError   = 0
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
	InternalPanic  = -32604
)

type InternalErr struct {
	Data any
	Text string
}

func (e InternalErr) Error() string {
	return e.Text
}

func NewInternalErr(text string, data any) *InternalErr {
	return &InternalErr{
		Text: text,
		Data: data,
	}
}

var CodeMap = map[int]string{
	ParseError:     "Parse error",
	InvalidRequest: "Invalid request",
	MethodNotFound: "Method not found",
	InvalidParams:  "Invalid params",
	InternalError:  "Internal error",
	InternalPanic:  "Internal panic",
}
