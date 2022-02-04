package main

import "testing"

func getInstance() (*PHPComposerConfig, error) {
	return LoadComposerJSON(getLocalPath(`composer.json`))
}

func TestLoadComposerJOSN(t *testing.T) {
	conf, err := getInstance()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("%v", conf)
}

func TestPHPComposerConfig_GetPHPCondition(t *testing.T) {
	conf, err := getInstance()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(conf.GetPHPCondition())
}

func TestPHPComposerConfig_MatchVersion(t *testing.T) {
	conf, err := getInstance()
	if err != nil {
		t.Fatal(err)
		return
	}

	ver, err := conf.MatchVersion("5.6", "7.2.99", "8", "8.1")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(ver)
}
