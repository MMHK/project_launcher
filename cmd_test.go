package main

import (
	"testing"
)

func TestSearchAppPackage(t *testing.T) {
	appName := "docker"
	err, list := SearchAppPackage(appName)
	if err != nil {
		t.Error(err)
		return
	}

	for _, row := range list {
		t.Logf("%+v\n", row)
	}
}

func TestInstallAppPackage(t *testing.T) {
	appID := `Docker.DockerDesktop`
	err := InstallAppPackage(appID)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestParseAppItem(t *testing.T) {
	raw := `Docker Desktop      Docker.DockerDesktop           3.5.2       Moniker: docker`

	item := ParseAppItem(raw)

	t.Logf("%+v\n", item)
}

func TestEnableWSL(t *testing.T) {
	err := EnableWSL()
	if err != nil {
		t.Error(err)
	}
}

func TestEnableHyperV(t *testing.T) {
	err := EnableHyperV()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableVM(t *testing.T) {
	err := EnableVM()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestReloadPathEnv(t *testing.T) {
	err := ReloadPathEnv()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestStartContainer(t *testing.T) {
	path := "D:/projects/wechat-coupon/code/php/public"
	err := StartContainer(path, "sam-wechat")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestStartDockerDesktop(t *testing.T) {
	err := StartDockerDesktop()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestRunPHPConsole(t *testing.T) {
	err := RunPHPConsole("speedyagency")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestPHPComposerInit(t *testing.T) {
	path := "F:/Projects/speedyagency/code/public"
	err := PHPComposerInit(path)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestStopContainer(t *testing.T) {
	path := "F:/Projects/speedyagency/insurance/public"
	err := StopContainer(path, "sam-sa-ins")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}