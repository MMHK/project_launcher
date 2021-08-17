package main

import (
	_ "embed"
	"html/template"
	"os"
	"path/filepath"
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

	dcokerFile, err := os.Create(filepath.Join(savePath, `docker-compose.yml`))
	if err != nil {
		return err
	}
	defer dcokerFile.Close()


	err = dockerBuilder.Execute(dcokerFile, this.Docker)
	if err != nil {
		return err
	}

	return nil
}
