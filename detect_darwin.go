//go:build darwin
// +build darwin

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const DOCKER_DEPS_VERSION = `>=18362`

type OSInfo struct {
	CurrentVersion            string
	ProductName               string
	CurrentMajorVersionNumber uint64
	CurrentMinorVersionNumber uint64
	ReleaseVersion            string
	BuildVersion              string
}

func GetOSInfo() (*OSInfo, error) {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return nil, err
	}

	sysname := unix.ByteSliceToString(uts.Sysname[:])
	release := unix.ByteSliceToString(uts.Release[:])
	dotPos := strings.Index(release, ".")
	if dotPos == -1 {
		return nil, errors.New(fmt.Sprintf("nvalid release value: %s", release))
	}

	major := release[:dotPos]
	majorVersion, err := strconv.Atoi(major)
	if err != nil {
		return nil, err
	}
	strings.SplitN(release, ".", 2)

	return &OSInfo{
		CurrentVersion:            release,
		ProductName:               sysname,
		CurrentMajorVersionNumber: uint64(majorVersion),
		CurrentMinorVersionNumber: uint64(majorVersion),
		BuildVersion:              release,
		ReleaseVersion:            release,
	}, nil
}

func (this *OSInfo) IsWindows10() bool {
	return strings.Contains(this.ProductName, `Windows 10`)
}

func (this *OSInfo) IsMacOS() bool {
	return strings.Contains(this.ProductName, "Darwin")
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
	return false, ""
}

func IsDockerInstalled() error {
	_, err := os.Stat("/Applications/Docker.app")
	if err != nil {
		log.Error(err)
		return NewLauncherError(ERROR_DOCKER_DESKTOP_NOT_INSTALLED,
			"请先安装 Docker Desktop, https://doc.weixin.qq.com/doc/w2_AGUAAwb2AKY9YO5hUlhSjODSUvwi6?scode=AEwAtAeZAAkacJBfyk")
	}
	output, err := exec.Command("docker", "info").Output()
	if err != nil {
		log.Error(err)
		return NewLauncherError(ERROR_DOCKER_DESKTOP_NOT_RUNNING,
			"DockerDesktop 还未运行")
	}

	if !strings.Contains(string(output), "Server Version:") {
		return NewLauncherError(ERROR_DOCKER_DESKTOP_NOT_RUNNING,
			"DockerDesktop 还未运行")
	}

	return nil
}

func IsWindowTerminalInstalled() error {
	return errors.New("OS not support")
}

const WINSERVICE_STATUS_STARTED = 4
const WINSERVICE_STATUS_STOPPED = 1

type WinService struct {
	Status      int    `json:"Status"`
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

func IsMacOS() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

func DetectService(name string) (error, *WinService) {
	return errors.New("OS not support"), nil
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

func IsWordPressProject(dir string) bool {
	projectRoot := filepath.Join(dir, "wp-load.php")
	if _, err := os.Stat(projectRoot); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GetPHPDependFromWPDir(dir string) (string, error) {
	includedPaths := filepath.Join(dir, "wp-includes", "version.php")
	versionContent, err := ioutil.ReadFile(includedPaths)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`\$required_php_version = '([^']+)'`)
	found := r.FindStringSubmatch(string(versionContent))
	if len(found) > 1 {
		return found[1], nil
	}

	return "", errors.New("php version can not be found")
}

type WPPHPVerion struct {
	WPVersion  string `json:"wp"`
	PHPVersion string `json:"php"`
}

func MatchWordPressPHPVersion(wpVersion string) (string, error) {
	baseVersion := "7.0"
	wpPHPVersionMapping := []WPPHPVerion{
		WPPHPVerion{WPVersion: "~5.6", PHPVersion: "8"},
		WPPHPVerion{WPVersion: "~5.3", PHPVersion: "7.2"},
		WPPHPVerion{WPVersion: "~5.2", PHPVersion: "7.2"},
		WPPHPVerion{WPVersion: "~5.0", PHPVersion: "7.2"},
		WPPHPVerion{WPVersion: "~4.9", PHPVersion: "7.2"},
		WPPHPVerion{WPVersion: "~4.7", PHPVersion: "7.1"},
		WPPHPVerion{WPVersion: "~4.4", PHPVersion: "7.0"},
		WPPHPVerion{WPVersion: "~4.1", PHPVersion: "7.0"},
	}

	ver := semver.MustParse(wpVersion)

	for _, targetConstraint := range wpPHPVersionMapping {
		constraint, err := semver.NewConstraint(targetConstraint.WPVersion)
		if err != nil {
			return "", err
		}

		if constraint.Check(ver) {
			return targetConstraint.PHPVersion, nil
		}
	}

	return baseVersion, nil
}

func FindPublicDir(baseDir string) (string, error) {
	if IsWordPressProject(baseDir) {
		return baseDir, nil
	}
	baseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", errors.New("不是一个合法的项目地址")
	}

	composerConfigPath := filepath.Join(baseDir, "composer.json")
	depth := 1
tryMatch:
	if _, err := os.Stat(composerConfigPath); err != nil && os.IsNotExist(err) {
		baseDir = filepath.Dir(baseDir)
		composerConfigPath = filepath.Join(baseDir, "composer.json")
		depth++
		if depth < 5 {
			goto tryMatch
		}
	}

	publicDir := filepath.Join(baseDir, "public")
	if _, err := os.Stat(publicDir); err != nil && os.IsNotExist(err) {
		return "", errors.New("不是一个合法的项目地址")
	}

	return publicDir, nil
}

func DetectPHPVersion(dir string) (string, error) {
	composerConfigPath := filepath.Join(dir, "../composer.json")
	if IsWordPressProject(dir) {
		composerConfigPath = filepath.Join(dir, "composer.json")
		wpVersion, err := GetPHPDependFromWPDir(dir)
		if err != nil {
			log.Error(err)
			return "", err
		}
		phpVersion, err := MatchWordPressPHPVersion(wpVersion)
		if err != nil {
			log.Error(err)
			return "", err
		}

		return phpVersion, nil
	}
	if _, err := os.Stat(composerConfigPath); os.IsNotExist(err) {
		log.Error("未发现 composer.json")
		return "", err
	}
	conf, err := LoadComposerJSON(composerConfigPath)
	if err != nil {
		log.Error(err)
		return "", err
	}
	phpver, err := conf.MatchVersion("7.0", "7.2.99", "8.0.999", "8.1.999", "8.2.999")
	if err != nil {
		log.Error(err)
		return "", err
	}
	if strings.EqualFold(phpver, "7.2.99") {
		phpver = "7.2"
	}
	if strings.EqualFold(phpver, "8.0.999") {
		phpver = "8.0"
	}
	if strings.EqualFold(phpver, "8.1.999") {
		phpver = "8.1"
	}
	if strings.EqualFold(phpver, "8.2.999") {
		phpver = "8.2"
	}

	log.Infof("配到PHP 运行版本为 %s \n", phpver)
	return phpver, nil
}
