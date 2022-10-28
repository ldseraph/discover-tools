package main

import (
	"bytes"
	"context"
	"discover-tools/pkg/mdns"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func init() {
	RegistryModule(NewManager())
}

type IP4Address struct {
	Address string `json:"address"`
	Prefix  uint8  `json:"prefix"`
	Gateway string `json:"gateway"`
}

type ethRespone struct {
	IP4  []IP4Address `json:"ip4"`
	DHCP bool         `json:"dhcp"`
}

type Node struct {
	ID string `json:"uuid"`
	ethRespone
	url string
}

func (n *Node) fetch() error {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	res, err := client.Get(n.url + "/eth")
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	respone := &ethRespone{}
	err = json.Unmarshal(data, respone)
	if err != nil {
		return err
	}
	n.DHCP = respone.DHCP
	n.IP4 = respone.IP4
	return nil
}

type NodeManager struct {
	entriesCh       chan *mdns.ServiceEntry
	interfaces      []net.Interface
	nodes           []*Node
	service         string
	queryPeriod     time.Duration
	heartbeatPeriod time.Duration
	nodeEmitPeriod  time.Duration
	ctx             context.Context
}

func NewManager() *NodeManager {
	n := &NodeManager{
		entriesCh:       make(chan *mdns.ServiceEntry),
		service:         "._innovation._tcp.local.",
		heartbeatPeriod: 2 * time.Second,
		queryPeriod:     1 * time.Second,
		nodeEmitPeriod:  1 * time.Second,
	}

	inter, _ := net.Interfaces()

	for _, i := range inter {
		if (i.Flags&net.FlagUp) != 0 && (i.Flags&net.FlagLoopback) == 0 {
			n.interfaces = append(n.interfaces, i)
		}
	}
	return n
}

func (n *NodeManager) SetEth(newNode *Node) error {
	for i, oldNode := range n.nodes {
		if oldNode.ID == newNode.ID {
			json_data, err := json.Marshal(newNode)
			if err != nil {
				return err
			}
			_, err = http.Post(oldNode.url+"/eth", "application/json",
				bytes.NewBuffer(json_data))
			n.nodes = append(n.nodes[:i], n.nodes[i+1:]...)
			return err
		}
	}
	return errors.New("no find")
}

func (n *NodeManager) GetEth(ID string) (*Node, error) {
	for _, node := range n.nodes {
		if node.ID == ID {
			return node, nil
		}
	}

	return nil, errors.New("no found")
}

func (n *NodeManager) OnStartup(ctx context.Context) {
	n.ctx = ctx
	go n.heartbeat()
	go n.query()
	go n.nodeEmit()
}

func (n *NodeManager) findNode(nodeID string) *Node {
	for _, node := range n.nodes {
		if node.ID == nodeID {
			return node
		}
	}
	return nil
}

func (n *NodeManager) parserName(name string) (string, bool) {
	return strings.TrimSuffix(name, n.service), !strings.HasSuffix(name, n.service)
}

func (n *NodeManager) handle() {
	for entry := range n.entriesCh {
		nodeID, vaild := n.parserName(entry.Name)
		if vaild {
			continue
		}

		node := n.findNode(nodeID)
		if node == nil {
			node := &Node{
				ID:  nodeID,
				url: "http://" + entry.AddrV4.String() + ":1323",
			}
			err := node.fetch()
			if err != nil {
				continue
			}
			n.nodes = append(n.nodes, node)
		}
	}
}

func (n *NodeManager) query() {
	go n.handle()

	ticker := time.NewTicker(n.queryPeriod)
	for {
		for _, i := range n.interfaces {
			params := mdns.DefaultParams(strings.TrimSuffix(n.service, ".local."))
			params.Entries = n.entriesCh
			params.Interface = &i
			params.DisableIPv6 = true
			err := mdns.Query(params)
			if err != nil {
				runtime.EventsEmit(n.ctx, "query_error", err)
			}
		}
		<-ticker.C
	}
}

func (n *NodeManager) nodeEmit() {
	ticker := time.NewTicker(n.nodeEmitPeriod)
	for {
		<-ticker.C
		runtime.EventsEmit(n.ctx, "node", n.nodes)
	}
}

func (n *NodeManager) heartbeat() {
	ticker := time.NewTicker(n.heartbeatPeriod)
	for {
		<-ticker.C
		var updateNode []*Node
		for _, node := range n.nodes {
			err := node.fetch()
			if err != nil {
				continue
			}
			updateNode = append(updateNode, node)
		}
		n.nodes = updateNode
	}
}
