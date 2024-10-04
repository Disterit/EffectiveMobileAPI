package request

const (
	OkReq             = "Ok"
	badReq            = "Bad Request"
	InternalServerReq = "Internal Server Error"
)

type Request struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}

func Ok() *Request {
	return &Request{Description: OkReq, Error: ""}
}

func BadRequest(Err string) *Request {
	return &Request{Description: badReq, Error: Err}
}

func InternalServer(Err string) *Request {
	return &Request{Description: InternalServerReq, Error: Err}
}
