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
	APP_NAME          = "wows-fast-stats"
	SEMVER            = "0.11.1"
	TEST_DIR          = "./test_install_dir/replays/"
	TEMP_ARENA_INFO   = "tempArenaInfo.json"
	CMD_GO            = "go"
	CMD_WAILS         = "wails"
	CMD_GOLANGCI_LINT = "golangci-lint"
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
	add := func(slice []string, elems map[string]string) []string {
		for k, v := range elems {
			slice = append(slice, fmt.Sprintf("-X %s=%s", k, v))
		}

		return slice
	}

	flags := make([]string, 0)
	flags = add(flags, LDFlagsCommon)

	if isDev {
		flags = add(flags, LDFlagsDev)
	} else {
		flags = add(flags, LDFlagsProd)
	}

	return strings.Join(flags, " ")
}

func npm(additionals ...string) error {
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
	return sh.RunV(CMD_GO, "generate", "./...")
}

func Lint() error {
	if err := sh.RunV(CMD_GOLANGCI_LINT, "run"); err != nil {
		return err
	}

	if err := npm("run", "check"); err != nil {
		return err
	}

	if err := npm("run", "lint"); err != nil {
		return err
	}

	return nil
}

func Fmt() error {
	if err := sh.RunV(CMD_GOLANGCI_LINT, "run", "--fix"); err != nil {
		return err
	}

	if err := sh.RunV(CMD_GO, "fmt"); err != nil {
		return err
	}

	if err := npm("run", "fmt"); err != nil {
		return err
	}

	return nil
}

func Test() error {
	if err := sh.RunV(CMD_GO, "test", "-cover", "./..."); err != nil {
		return err
	}

	if err := npm("run", "test"); err != nil {
		return err
	}

	return nil
}

func Dev() error {
	mg.Deps(Gen)

	return sh.RunV(
		CMD_WAILS, "dev",
		"-ldflags", ldflags(true),
	)
}

func Setup() error {
	pkgs := []string{
		"github.com/wailsapp/wails/v2/cmd/wails@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"golang.org/x/tools/cmd/stringer@latest",
	}

	for _, pkg := range pkgs {
		if err := sh.RunV(CMD_GO, "install", pkg); err != nil {
			return err
		}
	}

	if err := npm("ci", "./frontend"); err != nil {
		return err
	}

	return nil
}

func Build() error {
	mg.SerialDeps(Gen, Test)

	return sh.RunV(
		CMD_WAILS, "build",
		"-ldflags", ldflags(false),
		"-platform", "windows/amd64",
		"-o", AppExeName,
		"-trimpath",
	)
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

func PutTestData() error {
	testDataFiles := []string{}
	err := filepath.Walk(TEST_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if !strings.HasPrefix(fileName, TEMP_ARENA_INFO) {
			return nil
		}

		testDataFiles = append(testDataFiles, fileName)

		return nil
	})

	if err != nil {
		return err
	}

	for i, file := range testDataFiles {
		fmt.Println(i, file)
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

	return sh.Copy(
		filepath.Join(TEST_DIR, TEMP_ARENA_INFO),
		filepath.Join(TEST_DIR, testDataFiles[index]),
	)
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
