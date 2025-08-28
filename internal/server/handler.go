package server

import (
	"httpfromtcp.krokakrola.com/internal/request"
	"httpfromtcp.krokakrola.com/internal/response"
)

type Handler func(req *request.Request, res *response.Writer)
