package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

const DEMO_SERVICE_HOST = `192.168.33.6`
const LOCAL_SERVICE_HOST = `host.docker.internal`
const FRPS_API = `http://192.168.33.6:7001/api`
const FRPS_LOCAL_API = `http://host.docker.internal:7001/api`

func prepareRuntime() error {
	info, err := GetOSInfo()
	if err != nil {
		log.Error(err)
		return err
	}
	if !info.IsWindows10() {
		log.Error("支持windows 10，其他系统请自己解决")
		return err
	}
	if !info.MatchBuildVersion(DOCKER_DEPS_VERSION) {
		log.Error("支持windows 10 版本过低，装不了Docker，请自己解决")
		return err
	}

	err = IsDockerInstalled()
	if err != nil {
		log.Error(err)
		if MatchLauncherError(err, ERROR_DOCKER_DESKTOP_NOT_RUNNING) {
			startError := StartDockerDesktop()
			if startError != nil {
				log.Error(startError)
				return startError
			}
		}
		return err
	}

	return nil
}

func getProjectDomain(defaultValue string, asLocalService bool) (string, error) {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("必须提供一个项目名")
		}

		frpsEndPoint := FRPS_API
		if asLocalService {
			frpsEndPoint = FRPS_LOCAL_API
		}

		api := NewFrpApi(frpsEndPoint, "", "")
		api.DisableAuth()
		exist := api.SubDomainExist(input)
		if exist {
			return errors.New(fmt.Sprintf("三级域名 %s 已经存在并使用中，请重新选择另外一个名字", input))
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "请起一个霸气的项目名称, 最终将会能通过 {项目名}.dev.mixmedia.com 访问",
		Validate: validate,
		Default:  defaultValue,
	}

	return prompt.Run()
}

func StartPHPWebProject(asLocalService bool) error {
	if asLocalService {
		log.Info(`启动本地 FPRS service`)
		err := StartLocalFRPS()
		if err != nil {
			log.Error(err)
			return err
		}
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		return err
	}

	dir, err = FindPublicDir(dir)
	if err != nil {
		log.Error(err)
		return err
	}

	frpcPath := filepath.Join(dir, "frpc.ini")
	frpcCfg, err := LoadFrpcConfig(frpcPath)
	defaultProjectName := ""
	if err == nil && len(frpcCfg.SubDomain) > 0 {
		defaultProjectName = frpcCfg.SubDomain
	}

inputProjectName:
	projectName, err := getProjectDomain(defaultProjectName, asLocalService)
	if err != nil {
		log.Error(err)
		goto inputProjectName
	}

	// 匹配项目 composer.json
	// 并根据composer 分析出php 版本
	// 选择对应的 docker images
	phpver, err := DetectPHPVersion(dir)
	if err != nil {
		log.Error(err)
		return err
	}

	serviceHost := DEMO_SERVICE_HOST
	if asLocalService {
		serviceHost = LOCAL_SERVICE_HOST
	}

	builder := &BootStrapConfig{
		Frp: &FRPConfig{
			ServiceHost: serviceHost,
			SubDomain:   projectName,
		},
		Docker: &DockerComposeConfig{
			ImageVersion: phpver,
		},
	}
	// 覆盖项目 frpc.ini 及 docker-compose.yml
	err = builder.BuildConfig(dir)
	if err != nil {
		log.Error(err)
		return err
	}

	if asLocalService {
		vHostName := fmt.Sprintf(`%s.localhost`, projectName)
		err = AddLocalHostName(vHostName)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// 启动 docker compose 服务组
	log.Info("准备启动容器")
	err = StartContainer(dir, projectName)
	if err != nil {
		log.Error(err)
		return err
	}
	// 打开browser 访问 外网绑定网址
	go func() {
		vhostURL := fmt.Sprintf(`http://%s.dev.mixmedia.com`, projectName)
		if asLocalService {
			vhostURL = fmt.Sprintf(`http://%s.localhost`, projectName)
		}

		log.Infof("请访问 %s \n", vhostURL)
		OpenBrowser(vhostURL)
	}()

	return nil
}

func PHPConsole() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		return err
	}

	frpcPath := filepath.Join(dir, "frpc.ini")
	frpcCfg, err := LoadFrpcConfig(frpcPath)
	if err == nil && len(frpcCfg.SubDomain) > 0 {
		return RunPHPConsole(frpcCfg.SubDomain)
	}

	return nil
}

func ComposerInit() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		return err
	}

	return PHPComposerInit(dir)
}

func StopPHPWebProject() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
		return err
	}

	frpcPath := filepath.Join(dir, "frpc.ini")
	frpcCfg, err := LoadFrpcConfig(frpcPath)
	if err == nil && len(frpcCfg.SubDomain) > 0 {
		return StopContainer(dir, frpcCfg.SubDomain)
	}

	return nil
}

func StartMySQLServer() error {
	err := StartLocalMySQLServer()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Warning("\n-------------\nMySQL 服务启动成功，即将打开 adminer 管理入 \nMySQL 默认账密 root/mysql50\nMySQL endpoint: host.docker.internal\n-------------\n")
	OpenBrowser(`http://localhost:8088`)
	return nil
}

func StartLocalRedisServer() error {
	err := StartRedisService()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Warning("\n-------------\nRedis 服务启动成功 \nRedis endpoint: host.docker.internal\n-------------\n")
	return nil
}

func SelectMethods() {
	prompt := promptui.Select{
		Label: "请选择操作",
		Items: []string{
			"1 启动PHP web项目 (外网地址)",
			"2 启动PHP web项目 (本地地址)",
			"3 进入PHP 项目 console",
			"4 初始化 PHP项目 (composer update)",
			"5 停止PHP web项目",
			"6 启动 MySQL 服务",
			"7 启动 Redis 服务",
		},
	}

	index, _, err := prompt.Run()

	if err != nil {
		log.Error(err)
		return
	}

	switch index {
	case 0:
		StartPHPWebProject(false)
		return
	case 1:
		StartPHPWebProject(true)
		return
	case 2:
		PHPConsole()
		return
	case 3:
		ComposerInit()
		return
	case 4:
		StopPHPWebProject()
		return
	case 5:
		StartMySQLServer()
		return
	case 6:
		StartLocalRedisServer()
		return
	}
}

func main() {
	ReloadPathEnv()

	runtimeError := prepareRuntime()
	if runtimeError != nil {
		log.Error(`环境检查有问题，请解决后重试， 输入任意键退出`)
		fmt.Scanf("h")
		os.Exit(1)
		return
	}

start:

	SelectMethods()

	goto start
}
