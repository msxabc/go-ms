package microservice

import (
	"strings"
	"strconv"
	//"gt-go-ms/config"
	//"gt-go-ms/route"
	//"gt-go-ms/log"
	"gt-go-ms/endpoint"
	"net/http"
	"golang.org/x/net/context"
	httptransport "github.com/go-kit/kit/transport/http"
	gkEndpoint "github.com/go-kit/kit/endpoint"
)

type MSOperation interface {
	Send()
	Receive()
	Close()
	decodeRequest(_ context.Context, r *http.Request) (interface{}, error)
	encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error
	makeEndpoint() gkEndpoint.Endpoint

}

type Microservice struct {
	requestHandler endpoint.RequestUnmarshaler
	responseHandler endpoint.ResponseMarshaller
	epHandler endpoint.EndpointHandler
	epName string
	port int
	ctx context.Context
}

//TODO: right now we can only create one endpoint per port, should use a map structure to allow for multiple definitions 
func New (epHandler endpoint.EndpointHandler, epName string, req endpoint.RequestUnmarshaler, resp endpoint.ResponseMarshaller, port int) (*Microservice, error){
	ms := &Microservice{}
	ms.ctx = context.Background()
	ms.requestHandler = req
	ms.responseHandler = resp
	ms.epHandler = epHandler

	if (strings.HasPrefix(epName, "/")){
		ms.epName = epName
	}else{
		ms.epName = "/" + epName
	}
	
	ms.port = port

	return ms, newHTTPService(ms)
}

func newHTTPService(ms *Microservice) (error){

	httpServer := httptransport.NewServer(
		ms.ctx,
		ms.makeEndpoint(),
		ms.decodeRequest,
		ms.encodeResponse,
	)


	http.Handle(ms.epName, httpServer)
	return http.ListenAndServe(":" + strconv.Itoa(ms.port), nil)
}

func (ms *Microservice) makeEndpoint() gkEndpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
			return ms.epHandler(request)
	}
}


func (ms *Microservice) decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ms.requestHandler(r.Body)
}


func (ms *Microservice) encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return ms.responseHandler (w, resp)
}

func (ms *Microservice) Send(){
	
}

func (ms *Microservice) Receive(){
	
}

func (ms *Microservice) Close(){

}
