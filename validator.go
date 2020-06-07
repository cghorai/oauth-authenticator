package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type EventRequest struct {
	Token          string `json:"Token"`
	Config       string `json:"Config"`
}

type KeyPair struct {
	Kid 	string 	`json:"Kid"`
	RSA256Key 	string  `json:"SigningKey"`
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	var newEvent EventRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		invalidDataError := ErrorCode{ ErrorId: "1", Description: "The data is not correct."}
		json.NewEncoder(w).Encode(getErrorString(invalidDataError))
	}
	error := json.Unmarshal(reqBody, &newEvent)
	if error != nil{
		invalidDataError := ErrorCode{ ErrorId: "7", Description: "Error while unmarshalling dataset."}
		json.NewEncoder(w).Encode(getErrorString(invalidDataError))
	}
	if newEvent.Token == "" || newEvent.Config == "" {
		blankDataError := ErrorCode{ ErrorId: "2", Description: "The data is empty." }
		json.NewEncoder(w).Encode(getErrorString(blankDataError))
	} else {
		validateEvent(w,newEvent)
	}
}

func validateEvent(w http.ResponseWriter, e EventRequest){
	fmt.Println("Starting the application...")
	response, err := http.Get(e.Config)
	if err != nil {
		invalidDataError := ErrorCode{ ErrorId: "3", Description: "The config url is not valid."}
		json.NewEncoder(w).Encode(getErrorString(invalidDataError))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(data, &result)
		kurl,_ := json.Marshal(result["jwks_uri"])
		if len(kurl) == 0 {
			invalidDataError := ErrorCode{ ErrorId: "3", Description: "The config url is not valid."}
			json.NewEncoder(w).Encode(getErrorString(invalidDataError))
		} else {
			url := string(kurl)
			response, err := http.Get(url[1:len(url)-1])
			if err != nil {
				invalidDataError := ErrorCode{ ErrorId: "4", Description: "The key URL is not valid."}
				json.NewEncoder(w).Encode(getErrorString(invalidDataError))
			} else {
				data, err := ioutil.ReadAll(response.Body)
				if err != nil{
					invalidDataError := ErrorCode{ ErrorId: "5", Description: "The Key URL response body is empty."}
					json.NewEncoder(w).Encode(getErrorString(invalidDataError))
				}
				keypair := findKeys(w, data)
				//fmt.Println(keypair)
				keyObject := validate(e.Token,keypair)
				json.NewEncoder(w).Encode(keyObject)
			}
		}
	}
}



func validate(token string, keypair []KeyPair)(key KeyPair ) {
	header:=decodeB64(strings.Split(token, ".")[0])

	var result map[string]interface{}
	json.Unmarshal([]byte(header), &result)
	kid ,err := json.Marshal(result["kid"])
	if err != nil{
		log.Fatal("error:",err)
	}
	for _, result := range keypair {
		if result.Kid == truncateExtraInvertedCommas(string(kid)) {
			return result
		}
	}
	return
}

func findKeys(w http.ResponseWriter, s []byte)(keypair []KeyPair){
	var result map[string]interface{}
	json.Unmarshal(s, &result)
	var results []map[string]interface{}
	jsonResult, _ := json.Marshal(result["keys"])
	json.Unmarshal(jsonResult, &results)
	var keypairs []KeyPair
	for _, result := range results {
		kid ,err := json.Marshal(result["kid"])
		if err != nil{
			invalidDataError := ErrorCode{ ErrorId: "6", Description: "Invalid JSON data found."}
			json.NewEncoder(w).Encode(getErrorString(invalidDataError))
		}
		n ,err := json.Marshal(result["n"])
		if err != nil{
			invalidDataError := ErrorCode{ ErrorId: "6", Description: "Invalid JSON data found."}
			json.NewEncoder(w).Encode(getErrorString(invalidDataError))
		}
		var keypair = &KeyPair{truncateExtraInvertedCommas(string(kid)), truncateExtraInvertedCommas(string(n))}
		keypairs = append(keypairs, *keypair)
	}
	return keypairs
}
