package main

import "github.com/op/go-logging"

var log = logging.MustGetLogger("ipa2s3")

func init() {
	format := logging.MustStringFormatter(
		`ProjectLauncher %{color} %{shortfunc} %{level:.4s} %{shortfile}
%{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}
