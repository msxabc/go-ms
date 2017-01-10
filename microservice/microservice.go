package microservice

import (
	"strings"
	"strconv"
	//"gt-go-ms/config"
	//"gt-go-ms/log"
	"gt-go-ms/endpoint"
	"net/http"
	"golang.org/x/net/context"
	httptransport "github.com/go-kit/kit/transport/http"
	gkEndpoint "github.com/go-kit/kit/endpoint"

)

type MSOperation interface {
	AddEndpoint(epHandler endpoint.EndpointHandler, epName string, req endpoint.RequestUnmarshaler, resp endpoint.ResponseMarshaller)
	Start(port int)
	Stop()
}


type msHandler struct {
	requestHandler httptransport.DecodeRequestFunc
	responseHandler httptransport.EncodeResponseFunc
	epHandler gkEndpoint.Endpoint
}

type Microservice struct {
	eps map[string]*msHandler
	ctx context.Context
}

func New() (*Microservice, error){

	if err:= msInit(); err != nil {
		return nil, err
	}

	ms := &Microservice{}
	ms.ctx = context.Background()
	ms.eps = make(map[string]*msHandler)
	return ms, nil
}

func (ms *Microservice) AddEndpoint(epHandler endpoint.EndpointHandler, epName string, req endpoint.RequestUnmarshaler, resp endpoint.ResponseMarshaller){
	reqFunc := buildRequestUnmarshalFunc(req)
	respFunc := buildResponseMarshalFunc(resp)
	endpoint := buildEndpoint(epHandler)
	msh := &msHandler{reqFunc, respFunc, endpoint}
	ms.eps[epName] = msh
}

//TODO: start transport based on transport configuration
func (ms *Microservice) Start(port int){
	startHTTPService(ms, port)
}


func (ms *Microservice) Stop(){

}

//initialize everything enabled by the framework
func msInit() error{

	return nil
}

func buildRequestUnmarshalFunc(reqFunc endpoint.RequestUnmarshaler) httptransport.DecodeRequestFunc{
	return func (_ context.Context, r *http.Request) (interface{}, error) {
		return reqFunc(r.Body)
	}
}

func buildResponseMarshalFunc(respFunc endpoint.ResponseMarshaller) httptransport.EncodeResponseFunc {
	return func (_ context.Context, w http.ResponseWriter, resp interface{}) error {
		return respFunc(w, resp)
	}
}

func buildEndpoint(epHandler endpoint.EndpointHandler) gkEndpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
			return epHandler(request)
	}
}

func startHTTPService(ms *Microservice, port int) (error){
	for ep, msh := range ms.eps {
		httpServer := httptransport.NewServer(
			ms.ctx,
			msh.epHandler,
			msh.requestHandler,
			msh.responseHandler,
		)
		if (strings.HasPrefix(ep, "/")){
			http.Handle(ep, httpServer)
		}else{
			http.Handle("/"+ep, httpServer)
		}
	}
	return http.ListenAndServe(":" + strconv.Itoa(port), nil)
}










