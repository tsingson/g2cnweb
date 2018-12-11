package main

import (
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)



 //
func getToken(log *zap.Logger) (string, error) {
	url := "http://50.7.118.186/rpc/auth"

	// post payload
	var requestLogin ReqAuth

	requestLogin.UserID = "440e41f7-a706-4e41-afa0-70ad208e6639"
	requestLogin.MacAddress1 = "9c:f8:db:05:c5:f9"
	requestLogin.MacAddress2 = "9c:f8:db:05:c5:f0"
	requestLogin.ReleaseSn = "vk-v1.0.0-20180401"

	postBodyByte, _ := jsoniter.Marshal(requestLogin)

	client := &fasthttp.Client{}

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetConnectionClose()
	request.SetRequestURI(url)
	request.SetBody(postBodyByte)
	request.Header.SetMethod("POST")
	request.Header.Set("Content-Type", "application/json")

	if err := client.Do(request, response); err != nil {
		return "", err
	}
	code := response.StatusCode()
	if code == fasthttp.StatusOK {

		kk := make([]RespItem, 1)
		err1 := jsoniter.Unmarshal(response.Body(), &kk)
		if err1 != nil {
			return "", err1
		}
		result := kk[0]

		return result.Token, nil

	}

	return "", errors.New("result Error ")

}
