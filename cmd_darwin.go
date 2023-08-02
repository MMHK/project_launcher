//go:build darwin
// +build darwin

package main

import (
	"errors"
	"fmt"
	"github.com/goodhosts/hostsfile"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type PowerShellLogger struct {
}

type SearchAppItem struct {
	Name     string
	Version  string
	ID       string
	Category string
}

func (*PowerShellLogger) Write(arg string) {
	log.Debug(arg)
}

var r = regexp.MustCompile(`([0-9a-zA-Z\ \:\.]+)[ ]+([0-9a-zA-Z\.]+)[ ]+([0-9a-zA-Z\.]+)[ ]+([0-9a-zA-Z\ \:\.]+)`)

func ParseAppItem(raw string) *SearchAppItem {
	found := r.FindStringSubmatch(raw)
	item := new(SearchAppItem)

	if len(found) > 1 {
		item.Name = found[1]
	}
	if len(found) > 2 {
		item.ID = found[2]
	}
	if len(found) > 3 {
		item.Version = found[3]
	}
	if len(found) > 4 {
		item.Category = found[4]
	}
	return item
}

func InstallAppPackage(AppID string) error {
	return errors.New("OS not support")
}

func SearchAppPackage(appName string) (error, []*SearchAppItem) {
	return errors.New("OS not support"), nil
}

func ReloadPathEnv() error {
	oldPath := os.Getenv("PATH")
	newPath := "/usr/local/bin:" + oldPath
	os.Setenv("PATH", newPath)
	return nil
}

func EnableWSL() error {
	return errors.New("OS not support")
}

func EnableHyperV() error {
	return errors.New("OS not support")
}

func EnableVM() error {
	return errors.New("OS not support")
}

func StartContainer(dir string, containerName string) error {
	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker-compose",
		"--project-directory", filepath.FromSlash(dir),
		"--file", fmt.Sprintf(`%s/docker-compose.yml`, dir),
		"--project-name", containerName,
		"up", "--detach", "--force-recreate"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	if err := cmd.Start(); err != nil {
		log.Error(err)
		log.Error(os.Stderr)
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func StopContainer(dir string, containerName string) error {
	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker-compose",
		"--project-directory", filepath.FromSlash(dir),
		"--file", fmt.Sprintf(`%s/docker-compose.yml`, dir),
		"--project-name", containerName,
		"down"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	if err := cmd.Start(); err != nil {
		log.Error(err)
		log.Error(os.Stderr)
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func RunPHPConsole(containerName string) error {
	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker", "exec", "--workdir=/var/www",
		"-it", fmt.Sprintf(`%s-php-1`, containerName), "/bin/sh"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func PHPComposerInit(dir string) error {
	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerComposerYMLPath := filepath.FromSlash(filepath.Join(dir, "docker-compose.yml"))
	dockerScript := []string{"docker-compose", "--file", dockerComposerYMLPath, "run", "--no-deps", "--rm",
		"--workdir=/var/www", "php", "composer", "update"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func StartDockerDesktop() error {
	scriptWrap := `tell application "Docker" to activate`
	cmd := exec.Command("osascript", "-s", "h", "-e", scriptWrap)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func StartLocalMySQLServer() error {
	mysqlCfgPath, err := BuildMySQLConfig()
	if err != nil {
		return err
	}

	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker-compose",
		"--file", filepath.FromSlash(mysqlCfgPath),
		"--project-name",
		"mysql", "up", "-d"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func StartLocalFRPS() error {
	frpsConfPath, err := BuildFrpsConfig()
	if err != nil {
		return err
	}

	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker-compose",
		"--file", filepath.FromSlash(frpsConfPath),
		"--project-name", "frps", "up", "-d"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}

func AddLocalHostName(hostname string) error {
	hfile, err := hostsfile.NewHosts()
	if err != nil {
		return err
	}

	if hfile.Has(`127.0.0.1`, hostname) {
		return nil
	}

	err = hfile.Add(`127.0.0.1`, hostname)
	if err != nil {
		return err
	}

	return hfile.Flush()
}

func StartRedisService() error {
	redisConfPath, err := BuildRedisConfig()
	if err != nil {
		return err
	}

	frontEndScript := `activate application "Terminal"`
	scriptWrap := `tell application "Terminal" to do script "%s" `
	dockerScript := []string{"docker-compose",
		"--file", filepath.FromSlash(redisConfPath),
		"--project-name", "redis", "up", "-d"}
	scriptWrap = fmt.Sprintf(scriptWrap, strings.Join(dockerScript, " "))
	cmd := exec.Command("osascript", "-s", "h", "-e", frontEndScript, "-e", scriptWrap)

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Error(err)
		log.Error(cmd.Output())
		return err
	}

	return nil
}
