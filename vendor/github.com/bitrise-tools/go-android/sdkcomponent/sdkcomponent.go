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
	ABI      string
	Tag      string

	SDKStylePath       string
	LegacySDKStylePath string
}

// GetSDKStylePath ...
func (component SystemImage) GetSDKStylePath() string {
	if component.SDKStylePath != "" {
		return component.SDKStylePath
	}

	tag := "default"
	if component.Tag != "" {
		tag = component.Tag
	}

	return fmt.Sprintf("system-images;%s;%s;%s", component.Platform, tag, component.ABI)
}

// GetLegacySDKStylePath ...
func (component SystemImage) GetLegacySDKStylePath() string {
	if component.LegacySDKStylePath != "" {
		return component.LegacySDKStylePath
	}

	platform := component.Platform
	if component.Tag != "" && component.Tag != "default" {
		split := strings.Split(component.Platform, "-")
		if len(split) == 2 {
			platform = component.Tag + "-" + split[1]
		}
	}

	return fmt.Sprintf("sys-img-%s-%s", component.ABI, platform)
}

// InstallPathInAndroidHome ...
func (component SystemImage) InstallPathInAndroidHome() string {
	componentTag := "default"
	if component.Tag != "" {
		componentTag = component.Tag
	}

	return filepath.Join("system-images", component.Platform, componentTag, component.ABI)
}
