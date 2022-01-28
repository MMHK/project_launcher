package main

import "testing"

func TestBootStrapConfig_BuildConfig(t *testing.T) {
	conf := &BootStrapConfig{
		Frp: &FRPConfig{
			ServiceHost: "192.168.33.6",
			SubDomain:   "client03",
		},
		Docker: &DockerComposeConfig{
			ImageVersion: "8",
		},
	}

	err := conf.BuildConfig(getLocalPath("./tests"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadFrpcConfig(t *testing.T) {
	path := "F:\\Projects\\speedyagency\\code\\public\\frpc.ini"
	config, err := LoadFrpcConfig(path)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Logf("%v+", config)
}
