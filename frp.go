package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FrpApi struct {
	EndPoint string
	AuthUser string
	AuthPwd  string
	Auth bool
}

func NewFrpApi(endpoint string, adminUser string, adminPwd string) *FrpApi {
	return &FrpApi{
		EndPoint: endpoint,
		AuthUser: adminUser,
		AuthPwd:  adminPwd,
		Auth: true,
	}
}

type ProxyConfig struct {
	LocalIP   string `json:"local_ip"`
	LocalPort int32  `json:"local_port"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	SubDomain string `json:"subdomain"`
}

type ProxyEntry struct {
	ConnectCount int          `json:"cur_conns"`
	Name         string       `json:"name"`
	Status       string       `json:"status"`
	Conf         *ProxyConfig `json:"conf"`
}

type ProxyResp struct {
	Proxies []*ProxyEntry `json:"proxies"`
}

type InfoResp struct {
	SubDomainHost string `json:"subdomain_host"`
}

func (this *FrpApi) DisableAuth() {
	this.Auth = false
}

func (this *FrpApi) EnableAuth() {
	this.Auth = true
}

func (this *FrpApi) GetServiceInfo() (*InfoResp, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(`%s/serverinfo`, this.EndPoint), nil)
	if err != nil {
		return nil, err
	}
	if this.Auth {
		req.SetBasicAuth(this.AuthUser, this.AuthPwd)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	decode := json.NewDecoder(resp.Body)
	result := new(InfoResp)
	err = decode.Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *FrpApi) GetHTTPProxyList() (*ProxyResp, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(`%s/proxy/http`, this.EndPoint), nil)
	if err != nil {
		return nil, err
	}
	if this.Auth {
		req.SetBasicAuth(this.AuthUser, this.AuthPwd)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	decode := json.NewDecoder(resp.Body)
	result := new(ProxyResp)
	err = decode.Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 判定 使用的子域名是否在使用中
func (this *FrpApi) SubDomainExist(subDomain string) bool {
	list, err := this.GetHTTPProxyList()
	if err != nil {
		return false
	}

	if len(list.Proxies) == 0 {
		return false
	}

	for _, item := range list.Proxies {
		if strings.EqualFold(item.Name, subDomain) && strings.EqualFold(item.Status, "online") {
			return true
		}
	}

	return false
}
