package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KnicKnic/go-powershell/pkg/powershell"
	"github.com/Masterminds/semver/v3"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const DOCKER_DEPS_VERSION = `>=18362`

type OSInfo struct {
	CurrentVersion string
	ProductName string
	CurrentMajorVersionNumber uint64
	CurrentMinorVersionNumber uint64
	ReleaseVersion string
	BuildVersion string
}

func GetOSInfo() (*OSInfo, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer k.Close()

	cv, _, err := k.GetStringValue("CurrentVersion")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	pn , _, err := k.GetStringValue("ProductName")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rv, _, err := k.GetStringValue("ReleaseId")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &OSInfo{
		CurrentVersion: cv,
		ProductName: pn,
		CurrentMajorVersionNumber: maj,
		CurrentMinorVersionNumber: min,
		BuildVersion: cb,
		ReleaseVersion: rv,
	}, nil
}

func (this *OSInfo) IsWindows10() bool {
	return strings.Contains(this.ProductName, `Windows 10`)
}

func (this *OSInfo) MatchBuildVersion(Condition string) bool {
	constraint, err := semver.NewConstraint(Condition)
	if err != nil {
		log.Error(err)
		return false
	}

	ver := semver.MustParse(this.BuildVersion)
	return constraint.Check(ver)
}

func IsScoopInstalled() (bool, string) {
	PathEnv := os.Getenv("Path")
	if len(PathEnv) > 0 {
		pathList := strings.Split(PathEnv, `;`)
		scoopShimPath := ""
		for _, row := range pathList {
			if strings.Contains(row, `scoop.ps1`) {
				scoopShimPath = row
				break
			}
		}
		if len(scoopShimPath) > 0 {
			return true, filepath.Join(scoopShimPath, `scoop`)
		}
	}

	userDir, err := os.UserHomeDir()
	if err == nil {
		scoopDefault := filepath.Join(userDir, "scoop", "shims", "scoop.ps1")
		if _, err := os.Stat(scoopDefault); err == nil {
			return true, scoopDefault
		}
	}

	return false, ""
}

func IsWinGetInstalled() (bool, string) {
	output := ""
	err := RunScript(func(runner powershell.Runspace) error {
		cmd := `winget -v`
		log.Debug(cmd)
		
		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			for _, ele := range res.Objects {
				output += ele.ToString()
			}
			return nil
		}
		
		return errors.New(res.Exception.ToString())
	})
	
	if err != nil {
		log.Error(err)
		return false, ""
	}

	output = strings.TrimSpace(output)
	if len(output) <= 0 {
		return false, ""
	}

	return true, output
}

func IsDockerInstalled() (error) {
	err, _ := DetectService(`com.docker.service`)
	if err != nil {
		log.Error(err)
		return errors.New("install docker desktop first plz, https://docs.docker.com/docker-for-windows/install/")
	}
	
	cmd := exec.Command("docker", "info")
	err = cmd.Start()
	err = cmd.Wait()

	if err != nil {
		log.Error(err)
		
		return errors.New("start docker desktop first plz")
	}

	return nil
}


const WINSERVICE_STATUS_STARTED = 4;
const WINSERVICE_STATUS_STOPPED = 1;

type WinService struct {
	Status      int `json:"Status"`
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
}

func parseService(raw string) *WinService {
	item := new(WinService)
	
	
	decoder := json.NewDecoder(strings.NewReader(raw))
	err := decoder.Decode(item)
	if err != nil {
		log.Error(err)
	}
	return item
}

func DetectService(name string) (error, *WinService) {
	var service *WinService;
	err := RunScript(func(runner powershell.Runspace) error {
		cmd := fmt.Sprintf(`Get-Service "%s" | ConvertTo-Json -Compress`, name)
		//log.Debug(cmd)
		res := runner.ExecScript(cmd, true, nil)
		defer res.Close()
		if res.Success() {
			for _, ele := range res.Objects {
				service = parseService(ele.ToString())
				if len(service.Name) > 0 && !strings.EqualFold(service.Name, "Name") {
					return nil;
				}
 			}
			
			return errors.New("Service not found")
		}
		
		return errors.New(res.Exception.ToString())
	})
	
	return err, service
}

func LoadExistFrpcConfig(dir string) (string, error) {
	filePath := filepath.Join(dir, "frpc.ini")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("frpc.ini not found")
	}
	bin, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return "", err
	}
	r := regexp.MustCompile(`subdomain = ([^=]+)`)
	found := r.FindStringSubmatch(string(bin))
	if len(found) > 1 {
		return strings.TrimSpace(found[1]), nil
	}
	
	return "", errors.New("sub domain not found")
}