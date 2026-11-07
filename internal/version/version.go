package version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Version - основная версия приложения
var Version = "0.4.1-mvp.1"

// BuildInfo - информация о сборке
type BuildInfo struct {
	Version   string
	Commit    string
	DateTime  string
	Modified  bool
}

// GetVersion - возвращает версию приложения
// Если версия была встроена через ldflags, используется она
// Иначе используется версия по умолчанию
func GetVersion() string {
	if info, ok := getBuildInfo(); ok {
		return info.Version
	}
	return Version
}

// GetBuildInfo - возвращает полную информацию о сборке
func GetBuildInfo() BuildInfo {
	if info, ok := getBuildInfo(); ok {
		return info
	}
	return BuildInfo{
		Version:  Version,
		Commit:   "unknown",
		DateTime: "unknown",
		Modified: false,
	}
}

// getBuildInfo - извлекает информацию о сборке из runtime/debug
func getBuildInfo() (BuildInfo, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{}, false
	}

	// Извлекаем дополнительную информацию из main.Version
	// которая была установлена через ldflags
	var embeddedVersion string
	var embeddedCommit string
	var embeddedDateTime string
	var modified bool

	// Ищем main.Version в build info
	for _, setting := range info.Settings {
		switch setting.Key {
		case "main.Version":
			embeddedVersion = setting.Value
		case "vcs.revision":
			embeddedCommit = setting.Value
		case "vcs.time":
			embeddedDateTime = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	// Если версия не встроена через ldflags, используем из build info
	if embeddedVersion == "" {
		if info.Main.Version != "(devel)" {
			embeddedVersion = info.Main.Version
		} else {
			return BuildInfo{}, false
		}
	}

	return BuildInfo{
		Version:   embeddedVersion,
		Commit:    embeddedCommit,
		DateTime:  embeddedDateTime,
		Modified:  modified,
	}, true
}

// FormatVersion - форматирует версию для отображения
func FormatVersion() string {
	info := GetBuildInfo()
	
	if info.Commit != "" {
		shortCommit := info.Commit
		if len(shortCommit) > 7 {
			shortCommit = shortCommit[:7]
		}
		
		if info.Modified {
			return fmt.Sprintf("%s (dirty:%s)", info.Version, shortCommit)
		}
		return fmt.Sprintf("%s (%s)", info.Version, shortCommit)
	}
	
	return info.Version
}

// FormatFullVersion - возвращает полную информацию о версии
func FormatFullVersion() string {
	info := GetBuildInfo()
	
	var parts []string
	parts = append(parts, fmt.Sprintf("Version: %s", info.Version))
	
	if info.Commit != "" {
		parts = append(parts, fmt.Sprintf("Commit: %s", info.Commit))
	}
	
	if info.DateTime != "" {
		parts = append(parts, fmt.Sprintf("Build: %s", info.DateTime))
	}
	
	if info.Modified {
		parts = append(parts, "Modified")
	}
	
	return strings.Join(parts, " | ")
}