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

func TestBuildMySQLConfig(t *testing.T) {
	mysqlCfgPath, err := BuildMySQLConfig()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Logf(`mysql compose config file path: %s`, mysqlCfgPath)
}

func TestBuildFrpsConfig(t *testing.T) {
	confPath, err := BuildFrpsConfig()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}


	t.Logf(`frps compose config file path: %s`, confPath)
}

func TestBuildRedisConfig(t *testing.T) {
	confPath, err := BuildRedisConfig()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(confPath)
}