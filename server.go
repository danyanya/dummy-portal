// dummy-portal
// Adds MAC and IP to URL and make 307 redirect

package main

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/mostlygeek/arp"
)

// environment config
type environment struct {
	// serve address (host:port)
	ServeAddr string `envconfig:"SERVE_ADDR" default:":9998"`
	// redirect URL
	RedirectURL string `envconfig:"REDIRECT_URL" default:"https://ya.ru/"`
}

func findMAC(ip string) string {
	arp.CacheUpdate()
	return arp.Search(ip)
}

func main() {
	var env environment
	var err = envconfig.Process("", &env)
	if err != nil {
		panic(err)
	}

	// router initialization
	var e = echo.New()
	e.HideBanner = true

	// include standard recover and logger for router
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	var ur, _ = url.Parse(env.RedirectURL)

	e.GET("*",
		func(c echo.Context) error {
			var ip = c.RealIP()
			var u = ur
			var mac = findMAC(ip)
			q := u.Query()
			q.Set("mac", mac)
			q.Set("client_ip", ip)
			u.RawQuery = q.Encode()
			return c.Redirect(307, u.String())
		},
	)

	// start server
	if err = e.Start(env.ServeAddr); err != nil {
		panic(err.Error())
	}
}
