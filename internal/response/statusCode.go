package response

type StatusCode int

const (
	OK StatusCode = iota
	BadRequest
	InternalServerError
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

func (s StatusCode) Number() int {
	switch s {
	case OK:
		return 200
	case BadRequest:
		return 400
	case InternalServerError:
		return 500
	default:
		return 0
	}
}
