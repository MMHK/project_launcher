package main

import "testing"

func Test_SubDomainExist(t *testing.T) {
	api := NewFrpApi(`http://192.168.33.6:7001/api`, "admin", "admin")

	resp, err := api.GetHTTPProxyList()
	if err != nil {
		t.Error(err)
		return
	}

	for _, i := range resp.Proxies {
		t.Log(i.Name)
	}

	exist := api.SubDomainExist("client02")

	t.Log(exist)
}

func TestFrpApi_GetServiceInfo(t *testing.T) {
	api := NewFrpApi(`http://127.0.0.1:4000/api`, "admin", "admin")
	resp, err := api.GetServiceInfo()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", resp)
}
