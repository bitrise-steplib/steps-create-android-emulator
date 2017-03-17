package avdmanager

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/bitrise-tools/go-android/sdkmanager"
)

// Model ...
type Model struct {
	legacy bool
	binPth string
}

// New ...
func New(androidHome string) (*Model, error) {
	binPth := filepath.Join(androidHome, "tools", "bin", "avdmanager")

	legacy, err := sdkmanager.IsLegacySDKVersion(androidHome)
	if err != nil {
		return nil, err
	} else if legacy {
		binPth = filepath.Join(androidHome, "tools", "android")
	}

	if exist, err := pathutil.IsPathExists(binPth); err != nil {
		return nil, err
	} else if !exist {
		return nil, fmt.Errorf("no sdk manager tool found at: %s", binPth)
	}

	return &Model{
		legacy: legacy,
		binPth: binPth,
	}, nil
}

// CreateAVDCommand ...
func (model Model) CreateAVDCommand(name string, systemImage sdkcomponent.SystemImage, options ...string) *command.Model {
	if model.legacy {
		args := append([]string{"create", "avd", "--force", "--name", name, "--target", systemImage.Platform, "--abi", systemImage.ABI}, options...)
		return command.New(model.binPth, args...)
	}

	args := append([]string{"create", "avd", "--force", "--package", systemImage.GetSDKStylePath(), "--name", name, "--abi", systemImage.ABI}, options...)
	return command.New(model.binPth, args...)
}
