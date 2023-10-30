//go:build mage

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	APP_NAME             = "wows-fast-stats"
	SEMVER               = "0.11.1"
	TEST_DIR             = "test_install_dir/replays/"
	TEMP_ARENA_INFO_FILE = "tempArenaInfo.json"
)

var (
	Default       = Dev
	AppExeName    = fmt.Sprintf("%s-%s.exe", APP_NAME, SEMVER)
	LDFlagsCommon = map[string]string{
		"main.AppName": APP_NAME,
		"main.Semver":  SEMVER,
	}
	LDFlagsDev = map[string]string{
		"main.DiscordWebhookURL": readFile("discord_webhook_url_dev"),
		"main.IsDev":             "true",
	}
	LDFlagsProd = map[string]string{
		"main.DiscordWebhookURL": readFile("discord_webhook_url_prod"),
	}
)

func ldflags(isDev bool) string {
	var flags []string

	for k, v := range LDFlagsCommon {
		flags = append(flags, fmt.Sprintf("-X %s=%s", k, v))
	}

	if isDev {
		for k, v := range LDFlagsDev {
			flags = append(flags, fmt.Sprintf("-X %s=%s", k, v))
		}
	} else {
		for k, v := range LDFlagsProd {
			flags = append(flags, fmt.Sprintf("-X %s=%s", k, v))
		}
	}

	return strings.Join(flags, " ")
}

func execNPM(additionals ...string) error {
	args := []string{
		"--prefix",
		"./frontend",
	}

	for _, additional := range additionals {
		args = append(args, additional)
	}

	return sh.RunV("npm", args...)
}

func readFile(path string) string {
	ret, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("path: %s, error: %s", path, err.Error()))
	}

	return strings.Trim(string(ret), "\n")
}

func createZip(dst string, src string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	err = filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		zipPath, err := filepath.Rel(src, filePath)
		if err != nil {
			return err
		}

		zipEntry, err := w.Create(zipPath)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func Gen() error {
	return sh.RunV("go", "generate", "./...")
}

func Lint() error {
	if err := sh.RunV("golangci-lint", "run"); err != nil {
		return err
	}

	if err := execNPM("run", "check"); err != nil {
		return err
	}

	if err := execNPM("run", "lint"); err != nil {
		return err
	}

	return nil
}

func Fmt() error {
	if err := sh.RunV("golangci-lint", "run", "--fix"); err != nil {
		return err
	}

	if err := sh.RunV("go", "fmt"); err != nil {
		return err
	}

	if err := execNPM("run", "fmt"); err != nil {
		return err
	}

	return nil
}

func Test() error {
	if err := sh.RunV("go", "test", "-cover", "./..."); err != nil {
		return err
	}

	if err := execNPM("run", "test"); err != nil {
		return err
	}

	return nil
}

func Dev() error {
	mg.Deps(Gen)

	return sh.RunV("wails", "dev", "-ldflags", ldflags(true))
}

func Setup() error {
	pkgs := []string{
		"github.com/wailsapp/wails/v2/cmd/wails@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"golang.org/x/tools/cmd/stringer@latest",
	}

	for _, pkg := range pkgs {
		if err := sh.RunV("go", "install", pkg); err != nil {
			return err
		}
	}

	if err := execNPM("ci", "./frontend"); err != nil {
		return err
	}

	return nil
}

func Build() error {
	mg.SerialDeps(Gen, Test)

	return sh.RunV("wails", "build", "-ldflags", ldflags(false), "-platform", "windows/amd64", "-o", AppExeName, "-trimpath")
}

func Pkg() error {
	mg.Deps(Build)

	tempDir := "tmp"
	appDir := filepath.Join(tempDir, APP_NAME)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	dst := filepath.Join(appDir, AppExeName)
	src := filepath.Join("build/bin", AppExeName)
	if err := sh.Copy(dst, src); err != nil {
		return err
	}

	var appZip = fmt.Sprintf("%s.zip", APP_NAME)
	if err := createZip(appZip, tempDir); err != nil {
		return err
	}

	return nil
}

func PutTempArenaInfo() error {
	jsons := []string{}
	err := filepath.Walk(TEST_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if !strings.HasPrefix(fileName, TEMP_ARENA_INFO_FILE) {
			return nil
		}

		jsons = append(jsons, fileName)

		return nil
	})

	if err != nil {
		return err
	}

	for i, json := range jsons {
		fmt.Println(i, json)
	}

	var input string
	fmt.Print("index? > ")
	if _, err := fmt.Scan(&input); err != nil {
		return err
	}

	index, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	return sh.Copy(filepath.Join(TEST_DIR, TEMP_ARENA_INFO_FILE), filepath.Join(TEST_DIR, jsons[index]))
}

func Clean() error {
	targets := []string{
		"config/",
		"cache/",
		"temp_arena_info/",
		"screenshot/",
	}

	for _, target := range targets {
		if err := os.RemoveAll(target); err != nil {
			return err
		}
	}

	return nil
}
