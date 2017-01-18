package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gokyle/goconfig"
)

type puberReq struct {
	kind   string
	User   string `json:"user"`
	PubKey string `json:"pubKey"`
	YKey   string `json:"yKey"`
}

//RES=$(curl -sH "Content-Type: application/json" -X POST -d "{\"user\": \"${KEY}\", \"pubKey\": \"asdfasdfasdfasdf\", \"yKey\": \"${YK}\"}" http://localhost:8081/add)

func readUser(m string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(m)
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n")
}

func getYKey() string {
	return readUser("Press your yubikey: ")
}

func getUser() string {
	return readUser("Enter the user to associate pubkey with: ")
}

func getPubKey() string {
	return readUser("Enter the pubkey: ")
}

func request(c goconfig.ConfigMap, r *puberReq) {
	server := c["server"]["url"] + "/" + r.kind

	encR, err := json.Marshal(r)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	res, err := http.Post(server, "application/json", bytes.NewBuffer(encR))
	if err != nil {
		log.Printf("%v\n", encR)
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("%v\n", encR)
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
}

func findConfig() string {
	path := os.Getenv("HOME")
	path += "/.puberrc"

	_, err := os.Stat(path)
	if err != nil {
		path = "/etc/puberrc"
		_, err := os.Stat(path)
		if err != nil {
			fmt.Println("No config file found in /etc/puberrc or ~/.puberrc!")
			os.Exit(1)
		}
	}

	return path
}

func usage() {
	fmt.Printf("Usage: %s [add,rm]\n", os.Args[0])
	os.Exit(0)
}

func main() {
	cfile := findConfig()
	conf, err := goconfig.ParseFile(cfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]

	var r = &puberReq{}

	switch cmd {
	case "add":
		r.kind = "add"
	case "rm":
		r.kind = "rm"
	default:
		usage()

	}

	r.User = getUser()
	r.PubKey = getPubKey()
	r.YKey = getYKey()

	request(conf, r)
}
