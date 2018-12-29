package main

import (
	"github.com/integrii/flaggy"
	"github.com/sanity-io/litter"
	"github.com/tsingson/fastx/utils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func _main() {
	/**
	// Requests will be spread among these servers.
	servers := []string{

		"198.245.49.120:80",
	}

	// Prepare clients for each server

	var lbc fasthttp.LBClient
	for _, addr := range servers {
		c := &fasthttp.HostClient{
			Addr: addr,
		}
		lbc.Clients = append(lbc.Clients, c)
	}
	*/

	// Declare variables and their defaults
	var url = "http://198.245.49.120/category"

	var g *flaggy.Subcommand
	g = flaggy.NewSubcommand("g")
	g.Description = "My great subcommand!"

	// Add a flag
	g.String(&url, "u", "url", "A test string flag")

	flaggy.AttachSubcommand(g, 1)
	// Parse the flag
	flaggy.Parse()

	if g.Used {

		c := &fasthttp.Client{}

		// Send requests to load-balanced servers

		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()

		// 	url := "http://198.245.49.120/category"
		// url := "/category"
		req.SetRequestURI(url)
		// 	req.SetHost()
		// req.Header.SetMethod("POST")
		req.Header.SetMethod("GET")
		jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIiA6ICJ0ZXJtaW5hbF9zdGIiLCAidXNlcl9pZCIgOiAiNDQwZTQxZjctYTcwNi00ZTQxLWFmYTAtNzBhZDIwOGU2NjM5IiwgImV4cCIgOiAxNTQyOTgxOTk2fQ.kMQ2PGc3jxpzrHDkqtK591MBtwrOF61UO3s5t0vLX8I"
		req.Header.Set("Authorization", utils.StrBuilder("Bearer ", jwtToken))
		req.Header.Add("Content-Type", "application/json")

		if err := c.Do(req, resp); err != nil {
			log.Fatal("Error when sending request", zap.Error(err))
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			log.Fatal("unexpected status code", zap.Int("status code", resp.StatusCode()))

		}
		litter.Dump(utils.BytesToStringUnsafe(resp.Body()))
	}

	// useResponseBody(resp.Body())

}
func useResponseBody(body []byte) {
	litter.Dump(body)
}
