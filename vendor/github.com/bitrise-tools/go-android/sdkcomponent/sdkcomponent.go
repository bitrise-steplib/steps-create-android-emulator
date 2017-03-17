package sdkcomponent

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Model ...
type Model interface {
	GetSDKStylePath() string
	GetLegacySDKStylePath() string
	InstallPathInAndroidHome() string
}

// BuildTool ...
type BuildTool struct {
	Version string

	SDKStylePath       string
	LegacySDKStylePath string
}

// GetSDKStylePath ...
func (component BuildTool) GetSDKStylePath() string {
	if component.SDKStylePath != "" {
		return component.SDKStylePath
	}
	return fmt.Sprintf("build-tools;%s", component.Version)
}

// GetLegacySDKStylePath ...
func (component BuildTool) GetLegacySDKStylePath() string {
	if component.LegacySDKStylePath != "" {
		return component.LegacySDKStylePath
	}
	return fmt.Sprintf("build-tools-%s", component.Version)
}

// InstallPathInAndroidHome ...
func (component BuildTool) InstallPathInAndroidHome() string {
	return filepath.Join("build-tools", component.Version)
}

// Platform ...
type Platform struct {
	Version string

	SDKStylePath       string
	LegacySDKStylePath string
}

// GetSDKStylePath ...
func (component Platform) GetSDKStylePath() string {
	if component.SDKStylePath != "" {
		return component.SDKStylePath
	}
	return fmt.Sprintf("platforms;%s", component.Version)
}

// GetLegacySDKStylePath ...
func (component Platform) GetLegacySDKStylePath() string {
	if component.LegacySDKStylePath != "" {
		return component.LegacySDKStylePath
	}
	return component.Version
}

// InstallPathInAndroidHome ...
func (component Platform) InstallPathInAndroidHome() string {
	return filepath.Join("platforms", component.Version)
}

// SystemImage ...
type SystemImage struct {
	Platform string
	Type     string
	ABI      string

	SDKStylePath       string
	LegacySDKStylePath string
}

// GetSDKStylePath ...
func (component SystemImage) GetSDKStylePath() string {
	if component.SDKStylePath != "" {
		return component.SDKStylePath
	}

	componentType := "default"
	if component.Type != "" {
		componentType = component.Type
	}

	return fmt.Sprintf("system-images;%s;%s;%s", component.Platform, componentType, component.ABI)
}

// GetLegacySDKStylePath ...
func (component SystemImage) GetLegacySDKStylePath() string {
	if component.LegacySDKStylePath != "" {
		return component.LegacySDKStylePath
	}

	platform := component.Platform
	if component.Type != "" && component.Type != "default" {
		split := strings.Split(component.Platform, "-")
		if len(split) == 2 {
			platform = component.Type + "-" + split[1]
		}
	}

	return fmt.Sprintf("sys-img-%s-%s", component.ABI, platform)
}

// InstallPathInAndroidHome ...
func (component SystemImage) InstallPathInAndroidHome() string {
	componentType := "default"
	if component.Type != "" {
		componentType = component.Type
	}

	return filepath.Join("system-images", component.Platform, componentType, component.ABI)
}
