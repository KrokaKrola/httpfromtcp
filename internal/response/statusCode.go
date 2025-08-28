package response

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func (s StatusCode) String() string {
	switch s {
	case OK:
		return "OK"
	case BadRequest:
		return "Bad Request"
	case InternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Status Code"
	}
}
