package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/avdmanager"
	"github.com/bitrise-tools/go-android/sdk"
	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/bitrise-tools/go-android/sdkmanager"
	"github.com/bitrise-tools/go-steputils/tools"
	"github.com/kballard/go-shellquote"
)

const (
	bitriseEmulatorName = "BITRISE_EMULATOR_NAME"
)

// ConfigsModel ...
type ConfigsModel struct {
	Name                         string
	Platform                     string
	Abi                          string
	Tag                          string
	Options                      string
	CustomHardwareProfileContent string
	AndroidHome                  string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Name:     os.Getenv("name"),
		Platform: os.Getenv("platform"),
		Abi:      os.Getenv("abi"),
		Tag:      os.Getenv("tag"),
		Options:  os.Getenv("options"),
		CustomHardwareProfileContent: os.Getenv("custom_hardware_profile_content"),
		AndroidHome:                  os.Getenv("ANDROID_HOME"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Name: %s", configs.Name)
	log.Printf("- Platform: %s", configs.Platform)
	log.Printf("- Abi: %s", configs.Abi)
	log.Printf("- Tag: %s", configs.Tag)
	log.Printf("- Options: %s", configs.Options)
	log.Printf("- AndroidHome: %s", configs.AndroidHome)
	log.Printf("- CustomHardwareProfileContent:")
	log.Printf(configs.CustomHardwareProfileContent)
}

func (configs ConfigsModel) validate() error {
	if configs.Name == "" {
		return errors.New("no Name parameter specified")
	}

	if configs.Platform == "" {
		return errors.New("no Platform parameter specified")
	}

	validAbis := []string{"armeabi-v7a", "arm64-v8a", "mips", "x86", "x86_64"}
	if configs.Abi == "" {
		return errors.New("no Abi parameter specified")
	} else if !isValueValid(configs.Abi, validAbis) {
		return fmt.Errorf("invalid Abi parameter specified (%s), valid options: %s", configs.Abi, validAbis)
	}

	validTags := []string{"default", "google_apis", "google_apis_playstore", "android-tv", "android-wear"}
	if configs.Tag == "" {
		return errors.New("no Tag parameter specified")
	} else if !isValueValid(configs.Tag, validTags) {
		return fmt.Errorf("invalid Tag parameter specified (%s), valid options: %s", configs.Tag, validTags)
	}

	if configs.AndroidHome == "" {
		return errors.New("no ANDROID_HOME env set")
	}

	return nil
}

func isValueValid(value string, validValues []string) bool {
	for _, v := range validValues {
		if v == value {
			return true
		}
	}
	return false
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	//
	// Check is platform installed
	fmt.Println()
	log.Infof("Check if platform installed")

	androidSdk, err := sdk.New(configs.AndroidHome)
	if err != nil {
		fail("Failed to create sdk, error: %s", err)
	}

	manager, err := sdkmanager.New(androidSdk)
	if err != nil {
		fail("Failed to create sdk manager, error: %s", err)
	}

	platformComponent := sdkcomponent.Platform{
		Version: configs.Platform,
	}

	platformInstalled, err := manager.IsInstalled(platformComponent)
	if err != nil {
		fail("Failed to check if platform (%s) installed, error: %s", platformComponent.Version)
	}

	log.Donef("installed: %v", platformInstalled)
	// ---

	//
	// Install platform
	if !platformInstalled {
		fmt.Println()
		log.Infof("Installing: %s", configs.Platform)

		installCmd := manager.InstallCommand(platformComponent)
		installCmd.SetStdin(strings.NewReader("y"))
		installCmd.SetStdout(os.Stdout)
		installCmd.SetStderr(os.Stderr)

		fmt.Println()
		log.Donef("$ %s", installCmd.PrintableCommandArgs())
		fmt.Println()

		if err := installCmd.Run(); err != nil {
			fail("Failed to install platform, error: %s", err)
		}

		if installed, err := manager.IsInstalled(platformComponent); err != nil {
			fail("Failed to check if platform (%s) installed, error: %s", platformComponent.Version)
		} else if !installed {
			fail("Failed to install platform")
		}

		log.Donef("Installed")
	}
	// ---

	//
	// Check if system image installed
	fmt.Println()
	log.Infof("Check if system image installed")

	systemImageComponent := sdkcomponent.SystemImage{
		Platform: configs.Platform,
		Tag:      configs.Tag,
		ABI:      configs.Abi,
	}

	log.Printf("Checking path: %s", systemImageComponent.InstallPathInAndroidHome())

	systemImageInstalled, err := manager.IsInstalled(systemImageComponent)
	if err != nil {
		fail("Failed to check if system image (platform: %s abi: %s tag: %s) installed, error: %s", systemImageComponent.Platform, systemImageComponent.ABI, systemImageComponent.Tag, err)
	}

	log.Donef("installed: %v", systemImageInstalled)
	// ---

	//
	// Install system image
	if !systemImageInstalled {
		fmt.Println()
		log.Infof("Installing system image (platform: %s abi: %s tag: %s)", systemImageComponent.Platform, systemImageComponent.ABI, systemImageComponent.Tag)

		installCmd := manager.InstallCommand(systemImageComponent)
		installCmd.SetStdin(strings.NewReader("y"))
		installCmd.SetStdout(os.Stdout)
		installCmd.SetStderr(os.Stderr)

		fmt.Println()
		log.Donef("$ %s", installCmd.PrintableCommandArgs())
		fmt.Println()

		if err := installCmd.Run(); err != nil {
			fail("Failed to install platform, error: %s", err)
		}

		if installed, err := manager.IsInstalled(systemImageComponent); err != nil {
			fail("Failed to check if system image (platform: %s abi: %s tag: %s) installed, error: %s", systemImageComponent.Platform, systemImageComponent.ABI, systemImageComponent.Tag, err)
		} else if !installed {
			fail("Failed to install system image (platform: %s abi: %s tag: %s)", systemImageComponent.Platform, systemImageComponent.ABI, systemImageComponent.Tag)
		}

		log.Donef("Installed")
	}
	// ---

	//
	// Create AVD image
	fmt.Println()
	log.Infof("Creating AVD image")

	options := []string{}
	if configs.Options != "" {
		opts, err := shellquote.Split(configs.Options)
		if err != nil {
			fail("Failed to split custom options: %v", configs.Options)
		}
		options = opts
	}

	avdManager, err := avdmanager.New(androidSdk)
	if err != nil {
		fail("Failed to create avd manager, error: %s", err)
	}

	cmd := avdManager.CreateAVDCommand(configs.Name, systemImageComponent, options...)
	cmd.SetStdin(strings.NewReader("n"))
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	fmt.Println()
	log.Donef("$ %s", cmd.PrintableCommandArgs())
	fmt.Println()

	if err := cmd.Run(); err != nil {
		fail("Failed to create image, error: %s", err)
	}
	// ---

	//
	// Write custom hardware profile
	if configs.CustomHardwareProfileContent != "" {
		fmt.Println()
		log.Infof("Applying custom hardware profile")

		homeDir := pathutil.UserHomeDir()
		avdImageDir := filepath.Join(homeDir, ".android/avd", configs.Name+".avd")

		if exist, err := pathutil.IsDirExists(avdImageDir); err != nil {
			fail("Failed to check if avd image dir (%s) exists, error: %s", avdImageDir, err)
		} else if !exist {
			fail("The avd image (%s) created but not found at: %s", configs.Name, avdImageDir)
		}

		configPth := filepath.Join(avdImageDir, "config.ini")
		if err := fileutil.WriteStringToFile(configPth, configs.CustomHardwareProfileContent); err != nil {
			fail("Failed to write custom hardware profile, error: %s", err)
		}

		log.Donef("config.ini path: %s", configPth)
		fmt.Println()
	}
	// ---

	if err := tools.ExportEnvironmentWithEnvman(bitriseEmulatorName, configs.Name); err != nil {
		fail("Failed to export %s, error: %s", bitriseEmulatorName, err)
	}
	log.Donef("Emulator name is exported in environment variable: %s (value: %s)", bitriseEmulatorName, configs.Name)
}
