package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	nexmo "github.com/judy2k/nexmo-go"
	cache "github.com/patrickmn/go-cache"
)

var cacheDB *cache.Cache

var (
	API_KEY    = os.Getenv("API_KEY")
	API_SECRET = os.Getenv("API_SECRET")
	APP_ID     = os.Getenv("APP_ID")
	NCCO_URL   = os.Getenv("NCCO_URL")
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello fren!")
}

func klaus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/index.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("phonenumber:", r.Form["phonenumber"])
		phone := strings.Join(r.Form["phonenumber"], "")
		when, found := cacheDB.Get(phone)
		if !found {
			sendKlaus(phone)
			fmt.Fprintf(w, "Ok fren, Klausing %s!", phone)
		} else {
			fmt.Fprintf(w, "No fren, Klaused %s %s! try again later", phone, when)
		}
	}
}

func nexmoResp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Nexmo callback")
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	fmt.Println("----------")
	w.Write([]byte("OK"))
}

func sendKlaus(numberK string) {

	privateKey, err := ioutil.ReadFile("private.key")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	auth := nexmo.NewAuthSet()

	auth.SetAPISecret(API_KEY, API_SECRET)
	auth.SetApplicationAuth(APP_ID, privateKey)

	client := nexmo.NewClient(http.DefaultClient, auth)

	callReq := nexmo.CreateCallRequest{
		To: []interface{}{
			nexmo.PhoneCallEndpoint{
				Type:   "phone",
				Number: numberK,
			},
		},
		From: nexmo.PhoneCallEndpoint{
			Type:   "phone",
			Number: "1234567890",
		},
		AnswerURL: []string{NCCO_URL},
	}
	callResp, resp, err := client.Call.CreateCall(callReq)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(callResp, resp)
	now := time.Now()
	cacheDB.Set(numberK, now.String(), cache.DefaultExpiration)
}

func main() {
	cacheDB = cache.New(59*time.Minute, 1*time.Minute)
	http.HandleFunc("/", index)
	http.HandleFunc("/klaus", klaus)
	http.HandleFunc("/nexmo", nexmoResp)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
