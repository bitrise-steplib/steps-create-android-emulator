package avdmanager

import (
	"fmt"
	"path/filepath"

	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/sdk"
	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/bitrise-tools/go-android/sdkmanager"
)

// Model ...
type Model struct {
	legacy bool
	binPth string
}

// IsLegacyAVDManager ...
func IsLegacyAVDManager(androidHome string) (bool, error) {
	exist, err := pathutil.IsPathExists(filepath.Join(androidHome, "tools", "bin", "avdmanager"))
	return !exist, err
}

// New ...
func New(sdk sdk.AndroidSdkInterface) (*Model, error) {
	binPth := filepath.Join(sdk.GetAndroidHome(), "tools", "bin", "avdmanager")

	legacySdk, err := sdkmanager.IsLegacySDKManager(sdk.GetAndroidHome())
	if err != nil {
		return nil, err
	}

	legacyAvd, err := IsLegacyAVDManager(sdk.GetAndroidHome())
	if err != nil {
		return nil, err
	} else if legacyAvd && legacySdk {
		binPth = filepath.Join(sdk.GetAndroidHome(), "tools", "android")
	} else if legacyAvd && !legacySdk {
		binPth = filepath.Join(sdk.GetAndroidHome(), "tools", "android")
		sdkManager, err := sdkmanager.New(sdk)
		if err == nil {
			updateCmd := sdkManager.UpdateToolsCommand()
			updateCmd.SetStderr(os.Stderr)
			updateCmd.SetStdout(os.Stdout)
			if err := updateCmd.Run(); err == nil {
				legacyAvd, err = IsLegacyAVDManager(sdk.GetAndroidHome())
				if err == nil && !legacyAvd {
					binPth = filepath.Join(sdk.GetAndroidHome(), "tools", "bin", "avdmanager")
				}
			}
		}
	}

	if exist, err := pathutil.IsPathExists(binPth); err != nil {
		return nil, err
	} else if !exist {
		return nil, fmt.Errorf("no avd manager tool found at: %s", binPth)
	}

	return &Model{
		legacy: legacyAvd,
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
