package mstest

import (
	"gt-go-ms/endpoint"
	"encoding/json"
	"io"
)

//Request JSON object, internal to the hanlder
//Make sure mapping rules are exportable for unmarshler
type addRequest struct {
	Numbers []numberData `json:"numbers"`
}

type numberData struct {
	Number int `json:"number"`
}

//Response JSON object, internal to the hanlder
//Make sure mapping rules are exportable for marshaller
type addResponse struct {
	Sum   int `json:"sum"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}


//Business logic handler, internal to the service handler
func add (a int, b int) int{
	return a + b
}


//Endpoint handler, rules of process request, engage business logic, and produce response
func MakeAddEndpoint() endpoint.EndpointHandler {
	return func(request interface{}) (interface{}, error) {
		s := 0
		req, ok := request.(addRequest)
		if (!ok){
			return addResponse{0, "Unknown request format"}, nil
		}
		
		for _, n := range req.Numbers {
			s = add (s, n.Number)
		}
		return addResponse{s, ""}, nil
	}
}

//Request unmarshaling rules
func DecodeAddRequest() endpoint.RequestUnmarshaler{
	return func (r io.Reader)  (interface{}, error) {
		var request addRequest
		if err := json.NewDecoder(r).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
}


//Response marshaling rules
func EncodeAddResponse() endpoint.ResponseMarshaller {
	return func (w io.Writer, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
}