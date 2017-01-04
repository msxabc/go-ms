package endpoint

import (
	"io"
)

type EndpointHandler func(request interface{}) (response interface{}, err error)
type RequestUnmarshaler func (r io.Reader) (interface{}, error) 
type ResponseMarshaller func (w io.Writer, response interface{}) error