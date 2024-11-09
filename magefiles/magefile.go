//go:build mage

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
)

const (
	appName            = "wows-fast-stats"
	semver             = "0.16.2-alpha1"
	wgAppID            = "e25e1a2af190880c9e33d3be7cc5313d"
	binaryName         = appName + "-" + semver + ".exe"
	biscordWebhookJSON = "discord_webhook.json"
	npm = "npm --prefix ./frontend"
)

var commonLdflags = map[string]string{
	"main.AppName": appName,
	"main.Semver":  semver,
	"main.WGAppID": wgAppID,
}

type DiscordWebhook struct {
	Dev  Webhook `json:"dev"`
	Prod Webhook `json:"prod"`
}

type Webhook struct {
	Alert string `json:"alert"`
	Info  string `json:"info"`
}

func getFormattedLDFlags(isDev bool) string {
	var discordWebhook DiscordWebhook

	file, err := os.ReadFile(biscordWebhookJSON)
	if err != nil {
		fmt.Printf("[WARN] Failed to parse %s: %v\n", biscordWebhookJSON, err)
	} else {
		json.Unmarshal(file, &discordWebhook)
	}

	flags := make(map[string]string)
	for k, v := range commonLdflags {
		flags[k] = v
	}

	if isDev {
		flags["main.IsDev"] = "true"
		flags["main.AlertDiscordWebhookURL"] = discordWebhook.Dev.Alert
		flags["main.InfoDiscordWebhookURL"] = discordWebhook.Dev.Info
	} else {
		if discordWebhook.Prod.Alert == "" || discordWebhook.Prod.Info == "" {
			panic("Discord webhook URL not defined")
		}
		flags["main.AlertDiscordWebhookURL"] = discordWebhook.Prod.Alert
		flags["main.InfoDiscordWebhookURL"] = discordWebhook.Prod.Info
	}

	var result []string
	for k, v := range flags {
		result = append(result, fmt.Sprintf("-X %s=%s", k, v))
	}
	return strings.Join(result, " ")
}

func run(command string) {
	fmt.Printf("$ %s\n", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println()
}

func Setup() {
	pkgs := []string{
		"github.com/wailsapp/wails/v2/cmd/wails@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3",
		"go.uber.org/mock/mockgen@latest",
		"github.com/fdaines/arch-go@latest",
	}
	for _, pkg := range pkgs {
		run("go install " + pkg)
	}
	run(npm + " ci")
}

func GenMock() {
	run("go generate ./...")
}

func Lint() {
	run("golangci-lint run")
	run(npm + " run check")
	run(npm + " run lint")
}

func Fmt() {
	run("golangci-lint run --fix")
	run("go fmt")
	run(npm + " run fmt")
}

func Test() {
	run("arch-go")
	run("go test -cover ./... -count=1")
	run(npm + " run test")
}

func Dev() {
	run(fmt.Sprintf("wails dev -ldflags \"%s\"", getFormattedLDFlags(true)))
}

func Chbtl() {
	testReplayPath := "./test_install_dir/replays"
	tempArenaInfoName := "tempArenaInfo.json"

	files, err := os.ReadDir(testReplayPath)
	if err != nil {
		panic(err)
	}

	var jsonFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") && file.Name() != tempArenaInfoName {
			jsonFiles = append(jsonFiles, filepath.Join(testReplayPath, file.Name()))
		}
	}

	for i, file := range jsonFiles {
		fmt.Println(i, file)
	}

	fmt.Print("index? > ")
	var index int
	fmt.Scan(&index)

	if index >= 0 && index < len(jsonFiles) {
		err := os.Rename(jsonFiles[index], filepath.Join(testReplayPath, tempArenaInfoName))
		if err != nil {
			panic(err)
		}
	}
}

func Build() {
	mg.Deps(Test)
	run(fmt.Sprintf("wails build -ldflags \"%s\" -platform windows/amd64 -o %s -trimpath -webview2 embed", getFormattedLDFlags(false), binaryName))
}

func Pkg() {
	mg.Deps(Build)

	tempDir, err := os.MkdirTemp("", appName)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	dstPath := filepath.Join(tempDir, appName)
	os.Mkdir(dstPath, os.ModePerm)

	err = os.Rename(filepath.Join("./build/bin/", binaryName), filepath.Join(dstPath, binaryName))
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(appName + ".zip")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	archive := zip.NewWriter(outFile)
	defer archive.Close()

	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, tempDir+"/")
		if !info.IsDir() {
			f, err := archive.Create(relPath)
			if err != nil {
				return err
			}
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()
			_, err = io.Copy(f, src)
		}
		return err
	})
	if err != nil {
		panic(err)
	}
}

func Clean() {
	targets := []string{"config/", "cache/", "temp_arena_info/", "screenshot/", "persistent_data/"}
	for _, target := range targets {
		os.RemoveAll(target)
	}
}

func Uddep() {
	run("go get -u")
	run(npm + " update")
}
