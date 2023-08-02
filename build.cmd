go clean
del .\projectLauncher.syso
rsrc -manifest projectLauncher.exe.manifest -o projectLauncher.syso
go build -ldflags "-s -w" -tags windows -o projectLauncher.exe .