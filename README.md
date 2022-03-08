# This is a Clash For Flutter Service.

<!-- https://github.com/mozey/run-as-admin -->

## windows run as admin

```bash
mshta vbscript:createobject("shell.application").shellexecute("absolutePath...\clash-for-flutter-service-windows-amd64.exe","install start","","runas",1)(window.close)
```

## macos run as admin

```bash
osascript -e 'do shell script "absolutePath.../clash-for-flutter-service-darwin-amd64.exe install start"'
```
