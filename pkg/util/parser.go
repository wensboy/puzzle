package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"go.yaml.in/yaml/v3"
)

func ParseJsonString(sr string, dest any) error {
	if err := json.Unmarshal([]byte(sr), dest); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	return nil
}

func ParseJsonFile(path string, dest any) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(dest); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	return nil
}

func ParseYamlString(sr string, dest any) error {
	if err := yaml.Unmarshal([]byte(sr), dest); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	return nil
}

func ParseYamlFile(path string, dest any) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	defer f.Close()
	if err = yaml.NewDecoder(bufio.NewReader(f)).Decode(dest); err != nil {
		clog.Error(fmt.Sprintf("<pkg.util.parse> %s", err.Error()))
		return err
	}
	clog.Info(fmt.Sprintf("<pkg.util.parse> %s parsed successfully!", palette.SkyBlue(absPath)))
	return nil
}
