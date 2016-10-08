package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/GeertJohan/yubigo"
	"github.com/kisom/whitelist"
	"github.com/qbit/puber/backend"
)

type puberReq struct {
	User   string `json:"user"`
	PubKey string `json:"pubKey"`
	YKey   string `json:"yKey"`
}

var yubiServer string
var yubiSKey string
var yubiCID string
var err error
var debug bool
var store = backend.MemStore{}
var be backend.Backend

func auth(ykey string) (*yubigo.YubiResponse, bool, error) {
	yubiAuth, err := yubigo.NewYubiAuth(yubiCID, yubiSKey)
	if yubiServer != "" {
		log.Printf("using '%s'", yubiServer)
		yubiAuth.SetApiServerList(yubiServer)
	}
	if err != nil {
		log.Println(err)
	}

	return yubiAuth.Verify(ykey)
}

func handleIdx(w http.ResponseWriter, r *http.Request) {
	data, err := be.GetCount()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "I currently have %d key(s)", data)
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	keys, err := be.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data string
	for i := range keys {
		for j := range keys[i] {
			data += keys[i][j] + "\n"
		}
	}

	fmt.Fprintf(w, data)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	parts := strings.Split(u.String(), "/")
	user := parts[len(parts)-1]

	match, err := regexp.Match(",", []byte(user))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var data string

	if match {
		users := strings.Split(user, ",")
		var d []string
		for i := range users {
			d, err = be.Get(users[i])
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			for j := range d {
				data += d[j] + "\n"
			}

		}
	} else {
		var d []string
		d, err = be.Get(user)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data += strings.Join(d, "\n")

	}
	fmt.Fprintf(w, data)
}

func readPost(r *http.Request) (*puberReq, error) {
	var data puberReq

	if r.Body == nil {
		return nil, errors.New("empty request")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	data, err := readPost(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, ok, err := auth(data.YKey)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if debug {
		log.Println(result.GetRequestQuery())
	}

	if ok {
		be.Add(data.User, data.PubKey)
		fmt.Fprintf(w, "Added\n")
	} else {
		log.Println("Auth Failed!")
		http.Error(w, "Auth Failed!", http.StatusUnauthorized)
	}
}

func handleRm(w http.ResponseWriter, r *http.Request) {
	data, err := readPost(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok, err := auth(data.YKey)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok {
		if data.PubKey != "" {
			_, err := be.RM(data.User, data.PubKey)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(w, fmt.Sprintf("Removed '%s' from '%s'\n", data.PubKey, data.User))
		} else {
			_, err := be.RMAll(data.User)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "Removed\n")
		}
	} else {
		log.Println("Auth Failed!")
		http.Error(w, "Auth Failed!", http.StatusUnauthorized)
	}
}

func init() {
	store.Init()
	be = backend.Backend(&store)
}

func main() {
	addWL := whitelist.NewBasic()
	rmWL := whitelist.NewBasic()
	wlAddHandler, err := whitelist.NewHandlerFunc(handleAdd, nil, addWL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	wlRMHandler, err := whitelist.NewHandlerFunc(handleRm, nil, rmWL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	listen := flag.String("listen", ":8081", "listen string")
	whiteList := flag.String("wl", "localhost", "comma seperated list of hosts to allow access to /add")

	flag.StringVar(&yubiServer, "yserver", "api.yubico.com/wsapi/2.0/verify", "YubiAuth server to authenticate against")
	flag.StringVar(&yubiSKey, "yskey", "", "Yubi API Key")
	flag.StringVar(&yubiCID, "ycid", "", "Yubi Client ID")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")

	flag.Parse()

	if debug {
		be.Add("debug", "debug key")
		be.Add("debug", "debug key two")
		be.Add("debug", "debug key three")
	}

	for _, host := range strings.Split(*whiteList, ",") {
		ip, err := net.LookupHost(host)
		if err != nil {
			log.Fatal(err)
		}
		for i := range ip {
			log.Printf("whitelisting: %s (%s)\n", host, ip[i])

			// TODO: It's ghetto to have two whitelists.
			// Reduce this down to use one.
			addWL.Add(net.ParseIP(ip[i]))
			rmWL.Add(net.ParseIP(ip[i]))
		}
	}

	http.HandleFunc("/", handleIdx)
	http.HandleFunc("/all", handleAll)
	http.HandleFunc("/user/", handleUser)
	http.Handle("/add", wlAddHandler)
	http.Handle("/rm", wlRMHandler)

	log.Printf("Listening on '%s'\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
