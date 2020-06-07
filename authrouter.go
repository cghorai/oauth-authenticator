package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Config map[string]string

var configFile = "./resources/config"

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home! Please use /validateToken path to validate your token signing key")
}

func startRouter(configFileArgs string)  {
	if configFileArgs != "" {
		configFile = configFileArgs
	}
	config, err := readConfig(configFile)

	if err != nil {
		fmt.Println(err)
	}

	// virtualhost, ok := config["url"]
	// if !ok {
	// 	error
	// }
	virtualhost := config["url"]
	httpPort := config["httpPort"]
	httpsPort := config["httpsPort"]
	username := config["username"]
	password :=config["password"]
	banner:=config["banner"]
	
	
	
	fmt.Println()
	fmt.Printf("Application is starting up at URL:%s and Port:%s(http) & %s(https)",virtualhost,httpPort,httpsPort)

	router := mux.NewRouter().StrictSlash(true)
	// this endpoint heartbeat checking
	router.HandleFunc(virtualhost+"/", homeLink)
	router.HandleFunc("/", homeLink)
	// basic authentication on this url.
	router.HandleFunc(virtualhost+"/validateToken",basicAuthentication(validateToken, username, password, banner))
	router.HandleFunc("/validateToken",basicAuthentication(validateToken, username, password, banner))
	log.Fatal(http.ListenAndServe(":80", router))
	//log.Fatal(http.ListenAndServeTLS(":443",))

}

