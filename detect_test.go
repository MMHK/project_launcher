package main

import (
	"net"
	"testing"
)

func TestGetOSInfo(t *testing.T) {
	info, err := GetOSInfo()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("%+v\n", info)
}

func TestOSInfo_IsWindows10(t *testing.T) {
	info, err := GetOSInfo()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(info.IsWindows10())
}

func TestOSInfo_MatchBuildVersion(t *testing.T) {
	info, err := GetOSInfo()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(info.MatchBuildVersion(DOCKER_DEPS_VERSION))
}

func TestIsScoopInstalled(t *testing.T) {
	t.Log(IsScoopInstalled())
}

func TestIsDockerInstalled(t *testing.T) {
	ReloadPathEnv()
	err := IsDockerInstalled()
	if err != nil {
		t.Error(err)
		t.Fail()

		if MatchLauncherError(err, ERROR_DOCKER_DESKTOP_NOT_RUNNING) {
			err = StartDockerDesktop()
			if err != nil {
				t.Error(err)
				t.Fail()
			}
		}

		return
	}
}

func TestIsWinGetInstalled(t *testing.T) {
	ok, p := IsWinGetInstalled()
	if !ok {
		t.Error(`winget is not installed, Goto https://www.microsoft.com/p/app-installer/9nblggh4nns1 install`)
		return
	}

	t.Log(p)
}

func TestDetectService(t *testing.T) {
	serviceName := `com.docker.service`
	err, service := DetectService(serviceName)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v\n", service)
}

func TestLoadExistFrpcConfig(t *testing.T) {
	dir := getLocalPath(".")
	subDomain, err := LoadExistFrpcConfig(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(subDomain)
}

func TestIsWindowTerminalInstalled(t *testing.T) {
	err := IsWindowTerminalInstalled()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestIsWordPressProject(t *testing.T) {
	dir := "/Users/mixmedia/Downloads/wordpress"
	flag := IsWordPressProject(dir)
	if flag {
		return
	}

	t.Error("not a wordpress project")
	t.Fail()

}

func TestGetPHPDependFromWPDir(t *testing.T) {
	dir := "/Users/mixmedia/Downloads/wordpress"
	version, err := GetPHPDependFromWPDir(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(version)
}

func TestMatchWordPressPHPVersion(t *testing.T) {
	dir := "/Users/mixmedia/Downloads/wordpress"
	wpVersion, err := GetPHPDependFromWPDir(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(wpVersion)

	version, err := MatchWordPressPHPVersion(wpVersion)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(version)
}

func TestFindPublicDir(t *testing.T) {
	dir := "/Users/mixmedia/Desktop/Projects/mm-chat/motors-chat/public/static/bootstrap/css"
	publicDir, err := FindPublicDir(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(publicDir)
}

func TestDetectPHPVersion(t *testing.T) {
	dir := "/Users/mixmedia/Desktop/Projects/mm-chat/motors-chat/public/static/bootstrap/css"

	baseDir, err := FindPublicDir(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	phpVersion, err := DetectPHPVersion(baseDir)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(phpVersion)
}

func TestLocalIPRange(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		t.Error(err)
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			t.Log(ipNet.IP)
		}
	}
}
