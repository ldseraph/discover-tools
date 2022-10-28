package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wifx/gonetworkmanager"
	"github.com/google/uuid"
)

type eth0 struct {
	nm               gonetworkmanager.NetworkManager
	device           gonetworkmanager.Device
	activeConnection gonetworkmanager.ActiveConnection
	staticConnection gonetworkmanager.Connection
	dhcpConnection   gonetworkmanager.Connection
	settings         gonetworkmanager.Settings
}

func (e *eth0) GetAddress() (addrs []gonetworkmanager.IP4Address) {
	ip4config, err := e.device.GetPropertyIP4Config()
	if ip4config == nil || err != nil {
		return
	}
	addrs, _ = ip4config.GetPropertyAddresses()
	return
}

func (e *eth0) IsDHCP() bool {
	c, _ := e.activeConnection.GetPropertyConnection()
	setting, _ := c.GetSettings()
	return setting["ipv4"]["method"] == "auto"
}

type StaticIP struct {
	IPv4    string `json:"ip"`
	NetMask uint32 `json:"mask"`
	Gateway string `json:"gateway"`
}

func vaildIP(ip string) bool {
	strs := strings.Split(ip, ".")
	if len(strs) != 4 {
		return false
	}

	for _, s := range strs {
		if len(s) == 0 || (len(s) > 1 && s[0] == '0') {
			return false
		}
		if s[0] < '0' || s[0] > '9' {
			return false
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if n < 0 || n > 255 {
			return false
		}
	}
	return true
}

func (s *StaticIP) Vaild() bool {
	return vaildIP(s.IPv4) && vaildIP(s.Gateway)
}

func ipToInt(ip string) (value uint32) {
	for i, fields := range strings.Split(ip, ".") {
		a, _ := strconv.Atoi(fields)
		value = value + uint32(a*(1<<(8*i)))
	}
	return
}

func (e *eth0) SetIP(dhcp bool, IP4s []gonetworkmanager.IP4Address) error {
	var connection gonetworkmanager.Connection
	if dhcp {
		connection = e.dhcpConnection
	} else {
		connection = e.staticConnection
		setting, err := connection.GetSettings()
		if err != nil {
			return err
		}

		addressArray := make([][]uint32, len(IP4s))
		for _, ip4 := range IP4s {
			addresses := make([]uint32, 3)
			addresses[0] = ipToInt(ip4.Address)
			addresses[1] = uint32(ip4.Prefix)
			addresses[2] = ipToInt(ip4.Gateway)
			addressArray = append(addressArray, addresses)
		}

		setting["ipv4"]["addresses"] = addressArray
		delete(setting["ipv6"], "addresses")
		delete(setting["ipv6"], "routes")

		err = connection.Update(setting)
		if err != nil {
			return err
		}
	}

	activConnection, err := e.nm.ActivateConnection(connection, e.device, nil)
	if err != nil {
		return err
	}
	e.activeConnection = activConnection

	return nil
}

// func nmSetManaged(device string, managed bool) error {
// 	nm, err := gonetworkmanager.NewNetworkManager()
// 	if err != nil {
// 		return nil
// 	}

// 	d, err := nm.GetDeviceByIpIface(device)
// 	if err != nil {
// 		return err
// 	}
// 	current, err := d.GetPropertyManaged()
// 	if err != nil {
// 		return err
// 	}
// 	if current == managed {
// 		return nil
// 	}

// 	if err = d.SetPropertyManaged(managed); err != nil {
// 		return err
// 	}

// 	return nil
// }

func createConnection() (map[string]map[string]interface{}, error) {
	connection := make(map[string]map[string]interface{})
	connection["802-3-ethernet"] = make(map[string]interface{})
	connection["802-3-ethernet"]["auto-negotiate"] = true

	connection["connection"] = make(map[string]interface{})
	connection["connection"]["type"] = "802-3-ethernet"
	connectionUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	connection["connection"]["uuid"] = connectionUUID.String()
	connection["connection"]["interface-name"] = "eth0"
	connection["connection"]["autoconnect"] = true
	connection["ipv4"] = make(map[string]interface{})
	connection["ipv4"]["never-default"] = false
	connection["ipv6"] = make(map[string]interface{})
	connection["ipv6"]["method"] = "auto"
	return connection, nil
}

func createStaticConnection(settings gonetworkmanager.Settings) (gonetworkmanager.Connection, error) {
	connection, err := createConnection()
	if err != nil {
		return nil, err
	}
	connection["connection"]["id"] = "eth0_static"
	addresses := make([]uint32, 3)
	addresses[0] = ipToInt("192.168.101.3")
	addresses[1] = 24
	addresses[2] = ipToInt("192.168.101.1")
	addressArray := make([][]uint32, 1)
	addressArray[0] = addresses
	connection["ipv4"]["addresses"] = addressArray
	connection["ipv4"]["method"] = "manual"

	return settings.AddConnection(connection)
}

func createDHCPConnection(settings gonetworkmanager.Settings) (gonetworkmanager.Connection, error) {

	connection, err := createConnection()
	if err != nil {
		return nil, err
	}
	connection["connection"]["id"] = "eth0_dhcp"
	connection["ipv4"]["method"] = "auto"

	return settings.AddConnection(connection)
}

func NewEth() (*eth0, error) {

	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return nil, err
	}

	settings, err := gonetworkmanager.NewSettings()
	if err != nil {
		return nil, err
	}

	connections, err := settings.ListConnections()
	if err != nil {
		return nil, err
	}

	var staticConnection gonetworkmanager.Connection
	var dhcpConnection gonetworkmanager.Connection
	for _, c := range connections {
		s, err := c.GetSettings()
		if err != nil {
			return nil, err
		}
		id := s["connection"]["id"]
		switch id {
		case "eth0_static":
			staticConnection = c
		case "eth0_dhcp":
			dhcpConnection = c
		}
	}

	if staticConnection == nil {
		staticConnection, err = createStaticConnection(settings)
		if err != nil {
			return nil, err
		}
	}

	if dhcpConnection == nil {
		dhcpConnection, err = createDHCPConnection(settings)
		if err != nil {
			return nil, err
		}
	}

	eth0Device, err := nm.GetDeviceByIpIface("eth0")
	if err != nil {
		return nil, err
	}

	activeConnection, err := eth0Device.GetPropertyActiveConnection()
	if err != nil {
		return nil, err
	}

	id, err := activeConnection.GetPropertyID()
	if err != nil {
		return nil, err
	}

	if (id != "eth0_static") && (id != "eth0_dhcp") {
		activeConnection, err = nm.ActivateConnection(dhcpConnection, eth0Device, nil)
		if err != nil {
			return nil, err
		}
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			ip4config, _ := eth0Device.GetPropertyIP4Config()
			if ip4config != nil {
				break
			}
		}
	}

	return &eth0{
		nm:               nm,
		device:           eth0Device,
		staticConnection: staticConnection,
		dhcpConnection:   dhcpConnection,
		activeConnection: activeConnection,
		settings:         settings,
	}, nil
}
