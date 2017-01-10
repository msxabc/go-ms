package comm

import (
	"gt-go-ms/route"
	"net/http"
	"errors"
	"bytes"
	"io/ioutil"
)

type Response struct {
	Data []byte
	Err error
}


func Send(serviceName string, input []byte ) Response {
	return sendHttpMessage(serviceName, input)
}

func SendAsync(serviceName string, input []byte, c chan Response) {
	go func (c chan Response){
		c <- sendHttpMessage(serviceName, input)
	}(c)
}	


func sendHttpMessage (serviceName string, input []byte ) Response {
	path,err := route.Get(serviceName)

	if (err != nil) {return Response{nil, err}}

	if req, err := http.NewRequest("POST", path, bytes.NewBuffer(input)); err == nil {

		client := &http.Client{}
	    resp, err := client.Do(req)
		defer resp.Body.Close()
		if (err != nil){ return Response{nil, err}}


		//some sort of weird issue has occurred
		if (resp.StatusCode >= 300){
			return Response{nil, errors.New("Received non 2xx http status code from " + serviceName + ", " + path + " :" + resp.Status)}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil){ return Response{nil, err}}

		return Response{body, nil}
		
	}else {
			return Response{nil, err}
	}

}
