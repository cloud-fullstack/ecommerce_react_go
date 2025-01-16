package handler

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

func CorsAnywhere(c *gin.Context) {
	r := c.Request
	w := c.Writer
	target := c.Query("url")
	toCall := target

	// always allow access origin
	//w.Header().Add("Access-Control-Allow-Origin", origin)
	// https://www.flickr.com/services/oembed/?format=json&url=https://flic.kr/p/2mZnNKo&maxheight=600&maxwidth=600
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, HEAD, TRACE, DELETE, PATCH, COPY, HEAD, LINK, OPTIONS")

	if r.Method == "OPTIONS" {
		for n, h := range r.Header {
			if strings.Contains(n, "Access-Control-Request") {
				for _, h := range h {
					k := strings.Replace(n, "Request", "Allow", 1)
					w.Header().Add(k, h)
				}
			}
		}
		return
	}

	// create the request to server
	req, err := http.NewRequest(r.Method, toCall, r.Body)

	// add ALL headers to the connection
	for n, h := range r.Header {
		for _, h := range h {
			req.Header.Add(n, h)
		}
	}

	// create a basic client to send the request
	client := http.Client{}
	if r.TLS != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for h, v := range resp.Header {
		for _, v := range v {
			w.Header().Add(h, v)
		}
	}
	// copy the response from the server to the connected client request
	w.WriteHeader(resp.StatusCode)

	wr, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Println(wr, err)
	}
}
