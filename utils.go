package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type ErrorCode struct {
	ErrorId          	string `json:"ErrorId"`
	Description       string `json:"Description"`
}

func readConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"port":     "80",
		"url": "localhost",
		"sslport": "443",
	}
	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}



func getErrorString(e ErrorCode)(j string){
	jsonData, err := json.Marshal(e)
	if err != nil{
		fmt.Printf("Invalid error code provided, %s", err.Error())
	}
	return string(jsonData)
}

func truncateExtraInvertedCommas(s string)(r string)  {
	return s[1:len(s)-1]
}


func decodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, _ :=base64.StdEncoding.Decode(base64Text, []byte(message))
	return string(base64Text[:l])
}