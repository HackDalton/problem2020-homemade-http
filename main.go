package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/HackDalton/homemade-http/router"
	"github.com/mbndr/figlet4go"
)

var asciiRenderer *figlet4go.AsciiRender
var asciiOptions *figlet4go.RenderOptions
var rtr router.Router

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	asciiRenderer = figlet4go.NewAsciiRender()
	asciiOptions = figlet4go.NewRenderOptions()
	asciiOptions.FontColor = []figlet4go.Color{
		figlet4go.ColorCyan,
		figlet4go.ColorRed,
	}

	rtr = router.CreateRouter()
	rtr.AddRoute(router.POST("/login", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type", "text/plain")
		if req.FormValue("username") == "" || req.FormValue("password") == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing parameters."))
			return
		} else if req.FormValue("username") != "admin" || req.FormValue("password") != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid login."))
			return
		}
		flag, err := ioutil.ReadFile("./flag.txt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Add("content-type", "text/plain")
			w.Write([]byte("Couldn't read flag.txt."))
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(flag)
	})))
	rtr.AddRoute(router.GET("/", http.FileServer(http.Dir("./static"))))

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	captcha1 := rand.Intn(10)
	captcha2 := rand.Intn(10)
	captchaAnswer := captcha1 + captcha2
	captcha, _ := asciiRenderer.RenderOpts(strconv.Itoa(captcha1)+"+"+strconv.Itoa(captcha2), asciiOptions)
	conn.Write([]byte("Solve this captcha to connect:\n"))
	conn.Write([]byte(captcha))
	buff := bufio.NewReader(conn)
	captchaAnswerBytes, _, err := buff.ReadLine()
	if err != nil {
		conn.Write([]byte("An error occured reading your response\n"))
		conn.Close()
		return
	}

	submittedAnswer, err := strconv.Atoi(string(captchaAnswerBytes))
	if err != nil {
		conn.Write([]byte("You must submit a positive integer as your captcha answer.\n"))
		conn.Close()
		return
	}

	if submittedAnswer != captchaAnswer {
		conn.Write([]byte("That answer was not correct.\n"))
		conn.Close()
		return
	}

	conn.Write([]byte("Now accepting HTTP/1.1 requests.\n\n"))

	for true {
		req, err := http.ReadRequest(buff)
		if err != nil {
			conn.Write([]byte("An error occured processing your request:\n"))
			conn.Write([]byte(err.Error()))
			conn.Close()
			return
		} else if req.ProtoAtLeast(2, 0) {
			conn.Write([]byte("This server only supports HTTP/1.1 and below."))
			conn.Close()
			return
		} else if req.Host != "problems.hackdalton.com" {
			conn.Write([]byte("This server can only serve connections to problems.hackdalton.com."))
			conn.Close()
			return
		}

		resp := rtr.Handle(*req)
		resp.Write(conn)
		conn.Write([]byte("\n\n"))
	}
}
