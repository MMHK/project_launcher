package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DEMO_SERVICE_HOST = `192.168.33.6`
const FRPS_API = `http://192.168.33.6:7001/api`

func main() {
	ReloadPathEnv()
	
	info, err := GetOSInfo()
	if err != nil {
		log.Error(err)
		return
	}
	
	if !info.IsWindows10() {
		log.Error("支持windows 10，其他系统请自己解决")
		fmt.Scanln()
		return
	}
	if !info.MatchBuildVersion(DOCKER_DEPS_VERSION) {
		log.Error("支持windows 10 版本过低，装不了Docker，请自己解决")
		fmt.Scanln()
		return
	}
	
	
	err = IsDockerInstalled()
	
	if err != nil {
		log.Error(err)
		fmt.Scanln()
		return
	}

	api := NewFrpApi(FRPS_API, "", "")
	api.DisableAuth()
	
	
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		return
	}
	
	input := ""
	subDomain, err := LoadExistFrpcConfig(dir)
	if err == nil {
		input = subDomain
	}


	user_input:
	fmt.Println("请输入三级域名，最终将会能通过 {Your Name}.dev.mixmedia.com 访问")
	if len(subDomain) > 0 {
		fmt.Printf("已经检测到存在配置域名为 %s.dev.mixmedia.com， 请重新输入 %s\n", subDomain, subDomain)
	}
	_, err = fmt.Scanln(&input)
	if err != nil {
		log.Error(err)
		return
	}
	
	input = strings.TrimSpace(input)
	if len(input) <= 0 && len(subDomain) > 0 {
		fmt.Printf("使用现有域名 %s.dev.mixmedia.com\n", subDomain)
		input = subDomain
	}
	
	exist := api.SubDomainExist(input)
	if exist {
		log.Errorf("三级域名 %s 已经存在，请重新选择另外一个名字\n", input)
		goto user_input
	}
	
	composerConfigPath := filepath.Join(dir, "../composer.json")
	if _, err := os.Stat(composerConfigPath); os.IsNotExist(err) {
		log.Error("composer.json not found")
		return
	}
	conf, err := LoadComposerJOSN(composerConfigPath)
	if err != nil {
		log.Error(err)
		return
	}
	
	phpver, err := conf.MatchVersion("7.0", "7.2.99", "8")
	if strings.EqualFold(phpver, "7.2.99") {
		phpver = "7.2"
	}
	log.Infof("配到PHP 运行版本为 %s \n", phpver)
	
	builder := &BootStrapConfig{
		Frp: &FRPConfig{
			ServiceHost: DEMO_SERVICE_HOST,
			SubDomain: input,
		},
		Docker: &DockerComposeConfig{
			ImageVersion: phpver,
		},
	}
	
	err = builder.BuildConfig(dir)
	if err != nil {
		log.Error(err)
		return
	}
	
	go func() {
		log.Infof("请访问 http://%s.dev.mixmedia.com \n", input)
		OpenBrowser(fmt.Sprintf(`http://%s.dev.mixmedia.com`, input))
	}()
	
	log.Info("准备启动容器")
	err = StartContainer(dir)
	if err != nil {
		log.Error(err)
		return
	}
}


