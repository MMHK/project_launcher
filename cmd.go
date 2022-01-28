package main

import (
	"errors"
	"fmt"
	"github.com/KnicKnic/go-powershell/pkg/powershell"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
	return RunScript(func(runner powershell.Runspace) error {
		cmd := fmt.Sprintf(`winget install --id %s`, AppID)
		log.Debug(cmd)

		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			return nil
		}

		return errors.New(res.Exception.ToString())
	})
}

func SearchAppPackage(appName string) (error, []*SearchAppItem) {
	resultList := make([]*SearchAppItem, 0)
	err := RunScript(func(runner powershell.Runspace) error {
		cmd := fmt.Sprintf(`winget search -q %s`, appName)
		log.Debug(cmd)

		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			for _, ele := range res.Objects {
				app := ParseAppItem(ele.ToString())
				if len(app.ID) > 0 {
					resultList = append(resultList, app)
				}
			}
			return nil
		}

		return errors.New(res.Exception.ToString())
	})

	return err, resultList
}

func RunScript(callback func(powershell.Runspace) error) error {
	runSpace := powershell.CreateRunspace(new(PowerShellLogger), nil)
	defer runSpace.Close()

	return callback(runSpace)
}

func ReloadPathEnv() error {
	commandLine := `$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")`
	cmd := exec.Command("powershell", "-Command", commandLine)

	return cmd.Run()
}

func EnableWSL() error {
	return RunScript(func(runner powershell.Runspace) error {
		cmd := `dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart`
		log.Debug(cmd)

		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			return nil
		}

		return errors.New(res.Exception.ToString())
	})
}

func EnableHyperV() error {
	return RunScript(func(runner powershell.Runspace) error {
		cmd := `dism.exe /online /enable-feature /featurename:Microsoft-Hyper-V /all /norestart`
		log.Debug(cmd)

		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			return nil
		}

		return errors.New(res.Exception.ToString())
	})
}

func EnableVM() error {
	return RunScript(func(runner powershell.Runspace) error {
		cmd := `dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart`
		log.Debug(cmd)

		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			return nil
		}

		return errors.New(res.Exception.ToString())
	})
}

func StartContainer(dir string, containerName string) error {
	wtCmd := ""
	if err := IsWindowTerminalInstalled(); err == nil{
		wtCmd = "wt"
	}
	cmd := exec.Command("cmd", "/C", "start", wtCmd, "docker-compose",
		"--project-directory", filepath.FromSlash(dir),
		"--file", fmt.Sprintf(`%s/docker-compose.yml`, dir),
		"--project-name", containerName,
		"--", "up", "--detach", "--force-recreate")
	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}
	err := cmd.Wait()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func RunPHPConsole(containerName string) error {
	wtCmd := ""
	if err := IsWindowTerminalInstalled(); err == nil{
		wtCmd = "wt"
	}

	cmd := exec.Command("cmd", "/C", "start", wtCmd,
		"docker", "exec", "--workdir=/var/www", "-it", fmt.Sprintf(`%s_php_1`, containerName), "/bin/sh")

	//log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func PHPComposerInit(dir string) error {
	wtCmd := ""
	if err := IsWindowTerminalInstalled(); err == nil{
		wtCmd = "wt"
	}

	cmd := exec.Command("cmd", "/C", "start", "/wait", "/D", filepath.FromSlash(dir), wtCmd,
		"docker-compose", "run", "--no-deps", "--rm", "--workdir=/var/www", "--", "php", "composer", "update")

	log.Debugf("%s\n", cmd)
	if err := cmd.Start(); err != nil {
		log.Error(err)
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
	binPath, err := exec.LookPath(`docker.exe`)
	if err != nil {
		return err
	}

	execPath := filepath.Join(filepath.Dir(binPath), "../../Docker Desktop.exe")

	cmd := exec.Command(execPath)
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
