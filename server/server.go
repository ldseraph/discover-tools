package main

import (
	"discover-service/mdns"
	"fmt"
	"net/http"

	"github.com/Wifx/gonetworkmanager"
	"github.com/denisbrodbeck/machineid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Serve struct {
	eth  *eth0
	e    *echo.Echo
	port string
	mdns *mdns.Server
}

func (s *Serve) runMDNS() {
	if s.mdns != nil {
		s.mdns.Shutdown()
	}
	uuid, _ := machineid.ID()
	service, _ := mdns.NewMDNSService(uuid, "._innovation._tcp", "", "", 5353, s.eth, nil)
	s.mdns, _ = mdns.NewServer(&mdns.Config{Zone: service})
}

func (s *Serve) Run() {
	s.runMDNS()
	s.e.Logger.Fatal(s.e.Start(s.port))
}

func NewServer(eth *eth0, port string) *Serve {
	// Echo instance
	e := echo.New()
	s := &Serve{
		eth:  eth,
		e:    e,
		port: port,
	}

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("eth0", eth)
			return next(c)
		}
	})
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/eth", getEth)
	e.POST("/eth", postEth)
	return s
}

type getEthRespone struct {
	IP4  []gonetworkmanager.IP4Address `json:"ip4"`
	DHCP bool                          `json:"dhcp"`
}

func getEth(c echo.Context) error {
	eth := c.Get("eth0").(*eth0)
	return c.JSON(http.StatusOK, getEthRespone{
		IP4:  eth.GetAddress(),
		DHCP: eth.IsDHCP(),
	})
}

type postEthRequest struct {
	IP4  []gonetworkmanager.IP4Address `json:"ip4"`
	DHCP bool                          `json:"dhcp"`
}

type postEthRespone struct {
	Error string `json:"error"`
}

func postEth(c echo.Context) error {
	eth := c.Get("eth0").(*eth0)
	requset := &postEthRequest{}
	c.Bind(requset)
	if len(requset.IP4) == 0 {
		fmt.Println("ip4 len is 0")
		return c.JSON(http.StatusInternalServerError, postEthRespone{
			Error: "ip4 len is 0",
		})
	}
	err := eth.SetIP(requset.DHCP, requset.IP4)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, postEthRespone{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, postEthRespone{})
}
