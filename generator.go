package main

import (
	_ "embed"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

//go:embed frpc.ini
var frpConfigFile string

//go:embed docker-compose.yml
var ComposeFile string

//go:embed mysql-compose.yml
var MySQLComposeFile string

//go:embed frps-compose.yml
var FrpsComposeFile string
//go:embed frps.ini
var FrpsIniFile string
//go:embed redis-compose.yml
var RedisComposeFile string

type FRPConfig struct {
	ServiceHost string
	SubDomain   string
}

type FRPSComposeConfig struct {
	FrpsConfPath string
}

type MySQLComposeConfig struct {
	MySQLDATAPath string
	AdminerPort int
}

type RedisComposeConfig struct {
	RedisDataPath string
}

type DockerComposeConfig struct {
	ImageVersion string
}

type BootStrapConfig struct {
	Frp    *FRPConfig
	Docker *DockerComposeConfig
}

func (this *BootStrapConfig) BuildConfig(savePath string) error {
	frpBuilder, err := template.New("frp").Parse(frpConfigFile)
	if err != nil {
		return err
	}
	frpFile, err := os.Create(filepath.Join(savePath, `frpc.ini`))
	if err != nil {
		return err
	}
	defer frpFile.Close()

	err = frpBuilder.Execute(frpFile, this.Frp)
	if err != nil {
		return err
	}

	dockerBuilder, err := template.New("docker").Parse(ComposeFile)
	if err != nil {
		return err
	}

	dockerComposePath := filepath.Join(savePath, `docker-compose.yml`)
	if _, err := os.Stat(dockerComposePath); os.IsNotExist(err) {
		dcokerFile, err := os.Create(dockerComposePath)
		if err != nil {
			return err
		}
		defer dcokerFile.Close()

		err = dockerBuilder.Execute(dcokerFile, this.Docker)
		if err != nil {
			return err
		}
	}


	return nil
}

func LoadFrpcConfig(path string) (*FRPConfig, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	domainRule := regexp.MustCompile(`subdomain = ([^=]+)`)
	hostRule := regexp.MustCompile(`server_addr = ([^=]+)`)
	domainFound := domainRule.FindStringSubmatch(string(raw))
	hostFound := hostRule.FindStringSubmatch(string(raw))
	ServiceHost := ""
	SubDomain := ""
	if len(domainFound) > 1 {
		SubDomain = domainFound[1]
	}
	if len(domainFound) > 1 {
		ServiceHost = hostFound[1]
	}

	return &FRPConfig{
		ServiceHost: ServiceHost,
		SubDomain:   SubDomain,
	}, nil
}

func GetDefaultMySQLDataPath() string {
	dataPath, err := os.UserHomeDir()
	if err != nil {
		log.Error(err)
		return ""
	}

	return dataPath
}

func BuildMySQLConfig() (string, error) {
	savePath := os.TempDir()

	mysqlDataPath := filepath.Join(GetDefaultMySQLDataPath(), "mysql", "data")
	if _, err := os.Stat(mysqlDataPath); os.IsNotExist(err) {
		os.MkdirAll(mysqlDataPath, 0777)
	}

	mysqlCfgBuilder, err := template.New("frp").Parse(MySQLComposeFile)
	if err != nil {
		return "", err
	}
	tempCfgFilePath := filepath.Join(savePath, `mysql-compose.yml`)
	//log.Debugf(`mysql compose file path = %s`, tempCfgFilePath)
	mysqlCfgFile, err := os.Create(tempCfgFilePath)
	if err != nil {
		return "", err
	}
	defer mysqlCfgFile.Close()

	err = mysqlCfgBuilder.Execute(mysqlCfgFile, &MySQLComposeConfig{
		MySQLDATAPath: mysqlDataPath,
		AdminerPort: 8088,
	})
	if err != nil {
		return "", err
	}

	return tempCfgFilePath, nil
}

func BuildRedisConfig() (string, error) {
	savePath := os.TempDir()

	redisDataPath := filepath.Join(GetDefaultMySQLDataPath(), "redis", "data")
	if _, err := os.Stat(redisDataPath); os.IsNotExist(err) {
		os.MkdirAll(redisDataPath, 0777)
	}

	tmpRedisConfPath := filepath.Join(savePath, `redis-compose.yml`)
	redisConfFile, err := os.Create(tmpRedisConfPath)
	if err != nil {
		return "", err
	}
	defer redisConfFile.Close()

	redisBuilder, err := template.New("redis").Parse(RedisComposeFile)
	err = redisBuilder.Execute(redisConfFile, &RedisComposeConfig{
		RedisDataPath: redisDataPath,
	})
	if err != nil {
		return "", err
	}

	return tmpRedisConfPath, nil
}

func BuildFrpsConfig() (string, error) {
	savePath := os.TempDir()

	FrpsIniPath := filepath.Join(savePath, `frps.ini`)
	FrpsConfPath := filepath.Join(savePath, `frps-compose.yml`)

	err := ioutil.WriteFile(FrpsIniPath, []byte(FrpsIniFile), 0777)
	if err != nil {
		return "", err
	}

	frpsConfBuilder, err := template.New("frps").Parse(FrpsComposeFile)
	if err != nil {
		return "", err
	}
	log.Debugf(`mysql compose file path = %s`, FrpsConfPath)
	frpsConf, err := os.Create(FrpsConfPath)
	if err != nil {
		return "", err
	}
	defer frpsConf.Close()

	err = frpsConfBuilder.Execute(frpsConf, &FRPSComposeConfig{
		FrpsConfPath: filepath.ToSlash(FrpsIniPath),
	})
	if err != nil {
		return "", err
	}

	return FrpsConfPath, nil
}
