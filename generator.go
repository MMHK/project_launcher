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

type FRPConfig struct {
	ServiceHost string
	SubDomain   string
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
