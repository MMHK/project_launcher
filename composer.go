package main

import (
	"encoding/json"
	"errors"
	semver "github.com/Masterminds/semver/v3"
	"os"
	"strings"
)

type PHPComposerConfig struct {
	Require *RequireDeps    `json:"require"`
	Config  *ComposerConfig `json:"config"`
}

type RequireDeps struct {
	PHPCondition string `json:"php"`
}

type RuntimePlatform struct {
	PHPCondition string `json:"php"`
}

type ComposerConfig struct {
	Platform *RuntimePlatform `json:"platform"`
}

func LoadComposerJSON(filePath string) (*PHPComposerConfig, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decode := json.NewDecoder(f)
	config := new(PHPComposerConfig)
	err = decode.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (this *PHPComposerConfig) GetPHPCondition() string {
	if this.Config != nil && this.Config.Platform != nil {
		return this.Config.Platform.PHPCondition
	}

	if this.Require != nil {
		return this.Require.PHPCondition
	}

	return "0"
}

func (this *PHPComposerConfig) MatchVersion(versionList ...string) (string, error) {
	formatedCondition := strings.Replace(this.GetPHPCondition(), `|`, `||`, 1)
	//log.Debug(formatedCondition)
	constraint, err := semver.NewConstraint(formatedCondition)
	if err != nil {
		return "", err
	}

	for _, v := range versionList {
		ver := semver.MustParse(v)
		if constraint.Check(ver) {
			return v, nil
		}
	}

	return "", errors.New(`no version match`)
}
