package http

import (
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)


func SendPost(url string,Param interface{}) (response []byte) {
	method := "POST"
	p,_:=json.MarshalIndent(Param,"", "	")
	payload := strings.NewReader(string(p))
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return response
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return response
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return response
	}
	return body
}
