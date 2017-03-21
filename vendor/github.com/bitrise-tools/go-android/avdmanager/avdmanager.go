package avdmanager

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/sdkcomponent"
)

// Model ...
type Model struct {
	legacy bool
	binPth string
}

// New ...
func New(androidHome string) (*Model, error) {
	if exist, err := pathutil.IsDirExists(androidHome); err != nil {
		return nil, err
	} else if !exist {
		return nil, fmt.Errorf("android home not exists at: %s", androidHome)
	}

	binPth := filepath.Join(androidHome, "tools", "bin", "avdmanager")
	avdManagerExists, err := pathutil.IsPathExists(binPth)
	if err != nil {
		return nil, err
	} else if !avdManagerExists {
		binPth = filepath.Join(androidHome, "tools", "android")
	}

	if exist, err := pathutil.IsPathExists(binPth); err != nil {
		return nil, err
	} else if !exist {
		return nil, fmt.Errorf("no avd manager tool found at: %s", binPth)
	}

	return &Model{
		legacy: !avdManagerExists,
		binPth: binPth,
	}, nil
}

// CreateAVDCommand ...
func (model Model) CreateAVDCommand(name string, systemImage sdkcomponent.SystemImage, options ...string) *command.Model {
	if model.legacy {
		args := append([]string{"create", "avd", "--force", "--name", name, "--target", systemImage.Platform, "--abi", systemImage.ABI}, options...)
		return command.New(model.binPth, args...)
	}

	args := []string{"create", "avd", "--force", "--package", systemImage.GetSDKStylePath(), "--name", name, "--abi", systemImage.ABI}
	if systemImage.Tag != "" && systemImage.Tag != "default" {
		args = append(args, "--tag", systemImage.Tag)
	}
	args = append(args, options...)
	return command.New(model.binPth, args...)
}
