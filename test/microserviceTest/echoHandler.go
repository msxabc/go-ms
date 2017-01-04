package mstest

import (
	"strings"
	"errors"
	"gt-go-ms/endpoint"
	"io"
	"encoding/json"
)


//Business logic handler, internal to the handler
func echo (data string, command string) (string, error){
	switch (strings.ToUpper(command)){
		case "ECHO":
			return data, nil
		case "DOUBLE":
			return data+data, nil 
		default:
			return "", errors.New("Unknown command")
	}
}

//Request JSON object, internal to the hanlder
//Make sure mapping rules are exportable 
type echoRequest struct {
	Data string `json:"data"`
	Command string `json:"command"`
}

//Response JSON object, internal to the hanlder
//Make sure mapping rules are exportable 
type echoResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

//Endpoint handler
func MakeEchoEndpoint() endpoint.EndpointHandler {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(echoRequest)
		if (!ok){
			return nil, errors.New("Unknown request format")
		}
		v, err := echo(req.Data, req.Command)
		if err != nil {
			return echoResponse{v, err.Error()}, nil
		}
		return echoResponse{v, ""}, nil
	}
}

//Service controled marshaling and unmarshaling rules
func DecodeRequest() endpoint.RequestUnmarshaler{
	return func (r io.Reader)  (interface{}, error) {
		var request echoRequest
		if err := json.NewDecoder(r).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
}
func EncodeResponse() endpoint.ResponseMarshaller {
	return func (w io.Writer, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
}

