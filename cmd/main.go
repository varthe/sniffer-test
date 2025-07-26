package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"sniffer/internal/logger"
	"time"
)

func main() {
	plexURL := os.Getenv("PLEX_URL")
	if plexURL == "" {
		logger.Fatal("missing PLEX_URL")
	}

	target, err := url.Parse(plexURL)
	if err != nil {
		logger.Fatal("parse target: %v", err)
	}

	if err := ping(target); err != nil {
		logger.Fatal("plex unreachable: %v", err)
	}

	proxy := newProxy(target)

	log.Printf("proxy listening :80 → %s\n", plexURL)
	log.Fatal(http.ListenAndServe(":80", proxy))
}

func newProxy(target *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(target)
	orig := rp.Director
	rp.Director = func(req *http.Request) {
		orig(req)
		logger.File(req.RemoteAddr, req.URL.RequestURI(), "")
		logger.Console("%s %s", req.RemoteAddr, req.URL.RequestURI())
	}
	return rp
}

func ping(u *url.URL) error {
	v := *u
	v.Path = path.Join(v.Path, "web/index.html")

	c := &http.Client{Timeout: 10 * time.Second}
	r, err := c.Get(v.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", r.StatusCode)
	}
	return nil
}