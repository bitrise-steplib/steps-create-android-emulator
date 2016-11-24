package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/kballard/go-shellquote"
)

const (
	bitriseEmulatorName = "BITRISE_EMULATOR_NAME"
)

// ConfigsModel ...
type ConfigsModel struct {
	Name        string
	Platform    string
	Abi         string
	Options     string
	AndroidHome string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Name:        os.Getenv("name"),
		Platform:    os.Getenv("platform"),
		Abi:         os.Getenv("abi"),
		Options:     os.Getenv("options"),
		AndroidHome: os.Getenv("ANDROID_HOME"),
	}
}

func (configs ConfigsModel) print() {
	log.Info("Configs:")
	log.Detail("- Name: %s", configs.Name)
	log.Detail("- Platform: %s", configs.Platform)
	log.Detail("- Abi: %s", configs.Abi)
	log.Detail("- Options: %s", configs.Options)
	log.Detail("- AndroidHome: %s", configs.AndroidHome)
}

func (configs ConfigsModel) validate() error {
	if configs.Name == "" {
		return errors.New("no Name parameter specified")
	}
	if configs.Platform == "" {
		return errors.New("no Platform parameter specified")
	}
	if configs.Abi == "" {
		return errors.New("no Abi parameter specified")
	}
	if configs.AndroidHome == "" {
		return errors.New("no ANDROID_HOME env set")
	}
	return nil
}

func fail(format string, v ...interface{}) {
	log.Error(format, v...)
	os.Exit(1)
}

func exportEnvironmentWithEnvman(keyStr, valueStr string) error {
	cmd := cmdex.NewCommand("envman", "add", "--key", keyStr)
	cmd.SetStdin(strings.NewReader(valueStr))
	return cmd.Run()
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
	log.Info("Check if platform installed")

	platformInstalled := false

	platformPth := filepath.Join(configs.AndroidHome, "platforms", configs.Platform)
	if exist, err := pathutil.IsPathExists(platformPth); err != nil {
		fail("Failed to check if path (%s) exists, error: %s", platformPth, err)
	} else {
		platformInstalled = exist
	}

	log.Done("installed: %v", platformInstalled)
	// ---

	//
	// Install platform
	if !platformInstalled {
		fmt.Println()
		log.Info("Installing: %s", configs.Platform)

		args := []string{"android", "update", "sdk", "--no-ui", "--all", "--filter", configs.Platform}
		cmd := cmdex.NewCommand(args[0], args[1:]...)
		cmd.SetStdin(strings.NewReader("y"))
		cmd.SetStdout(os.Stdout)
		cmd.SetStdout(os.Stderr)

		fmt.Println()
		log.Done("$ %s", cmdex.PrintableCommandArgs(false, args))
		fmt.Println()

		if err := cmd.Run(); err != nil {
			fail("Failed to install platform, error: %s", err)
		}

		log.Done("Installed")
	}
	// ---

	//
	// Check if system image installed
	fmt.Println()
	log.Info("Check if system image installed")

	systemImageInstalled := false

	systemImagePth := filepath.Join(configs.AndroidHome, "system-images", configs.Platform, configs.Abi)
	log.Detail("checking path: %s", systemImagePth)

	if exist, err := pathutil.IsPathExists(systemImagePth); err != nil {
		fail("Failed to check if path (%s) exists, error: %s", systemImagePth, err)
	} else if !exist {
		systemImagePth = filepath.Join(configs.AndroidHome, "system-images", configs.Platform, "default", configs.Abi)
		log.Detail("checking path: %s", systemImagePth)

		if exist, err := pathutil.IsPathExists(systemImagePth); err != nil {
			fail("Failed to check if path (%s) exists, error: %s", systemImagePth, err)
		} else {
			systemImageInstalled = exist
		}
	} else {
		systemImageInstalled = true
	}

	log.Done("installed: %v", systemImageInstalled)
	// ---

	//
	// Install system image
	if !systemImageInstalled {
		systemImage := fmt.Sprintf("sys-img-%s-%s", configs.Abi, configs.Platform)

		fmt.Println()
		log.Info("Installing: %s", systemImage)

		args := []string{"android", "update", "sdk", "--no-ui", "--all", "--filter", systemImage}
		cmd := cmdex.NewCommand(args[0], args[1:]...)
		cmd.SetStdin(strings.NewReader("y"))
		cmd.SetStdout(os.Stdout)
		cmd.SetStdout(os.Stderr)

		fmt.Println()
		log.Done("$ %s", cmdex.PrintableCommandArgs(false, args))
		fmt.Println()

		if err := cmd.Run(); err != nil {
			fail("Failed to install system image, error: %s", err)
		}

		// Check if install succed
		systemImagePth := filepath.Join(configs.AndroidHome, "system-images", configs.Platform, configs.Abi)
		log.Detail("checking if system image created at path: %s", systemImagePth)

		if exist, err := pathutil.IsPathExists(systemImagePth); err != nil {
			fail("Failed to check if path (%s) exists, error: %s", systemImagePth, err)
		} else if !exist {
			systemImagePth = filepath.Join(configs.AndroidHome, "system-images", configs.Platform, "default", configs.Abi)
			log.Detail("checking if system image created at path: %s", systemImagePth)

			if exist, err := pathutil.IsPathExists(systemImagePth); err != nil {
				fail("Failed to check if path (%s) exists, error: %s", systemImagePth, err)
			} else if !exist {
				fail("system image: %s not installed", systemImage)
			}
		}

		log.Done("Installed")
	}
	// ---

	//
	// Create AVD image
	fmt.Println()
	log.Info("Creating AVD image")

	options := []string{}
	if configs.Options != "" {
		opts, err := shellquote.Split(configs.Options)
		if err != nil {
			fail("Failed to split custom options: %v", configs.Options)
		}
		options = opts
	}

	args := []string{"android", "create", "avd", "--force", "--name", configs.Name, "--target", configs.Platform, "--abi", configs.Abi}
	args = append(args, options...)

	cmd := cmdex.NewCommand(args[0], args[1:]...)
	cmd.SetStdin(strings.NewReader("n"))
	cmd.SetStdout(os.Stdout)
	cmd.SetStdout(os.Stderr)

	fmt.Println()
	log.Done("$ %s", cmdex.PrintableCommandArgs(false, args))
	fmt.Println()

	if err := cmd.Run(); err != nil {
		fail("Failed to create image, error: %s", err)
	}
	// ---

	if err := exportEnvironmentWithEnvman(bitriseEmulatorName, configs.Name); err != nil {
		fail("Failed to export %s, error: %s", bitriseEmulatorName, err)
	}
	log.Done("Emaultor name is exported in environment variable: %s (value: %s)", bitriseEmulatorName, configs.Name)
}
