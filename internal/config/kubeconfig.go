package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetDefaultKubeconfigPath() string {
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	switch runtime.GOOS {
	case "windows":
		if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
			return filepath.Join(userProfile, ".kube", "config")
		}
	default:
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".kube", "config")
		}
	}

	return filepath.Join(".kube", "config")
}

func expandPath(path string) (string, error) {
	if isWSL() && isWindowsPath(path) {
		path = convertWindowsPathToWSL(path)
	}

	if strings.HasPrefix(path, "~/") {
		var home string
		var err error

		if homeEnv := os.Getenv("HOME"); homeEnv != "" {
			home = homeEnv
		} else {
			home, err = os.UserHomeDir()
			if err != nil {
				return path, err
			}
		}
		return filepath.Join(home, path[2:]), nil
	}

	return os.ExpandEnv(path), nil
}

func ValidateKubeconfigPath(path string) (string, error) {
	if path == "" {
		return GetDefaultKubeconfigPath(), nil
	}

	expanded, err := expandPath(path)
	if err != nil {
		return "", err
	}

	return expanded, nil
}

func isWSL() bool {
	if _, err := os.Stat("/proc/version"); err == nil {
		if content, err := os.ReadFile("/proc/version"); err == nil {
			return strings.Contains(strings.ToLower(string(content)), "microsoft") ||
				strings.Contains(strings.ToLower(string(content)), "wsl")
		}
	}

	return os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSLENV") != ""
}

func isWindowsPath(path string) bool {
	if len(path) >= 2 && path[1] == ':' &&
		((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z')) {
		return true
	}

	return strings.HasPrefix(path, "\\\\")
}

func convertWindowsPathToWSL(path string) string {
	if len(path) >= 2 && path[1] == ':' {
		drive := strings.ToLower(string(path[0]))
		rest := path[2:]
		rest = strings.ReplaceAll(rest, "\\", "/")

		if strings.Contains(rest, "/.kube/") {
			kubeIndex := strings.Index(rest, "/.kube/")
			if kubeIndex != -1 {
				kubePath := rest[kubeIndex:]

				if homeEnv := os.Getenv("HOME"); homeEnv != "" {
					wslPath := homeEnv + kubePath
					if _, err := os.Stat(wslPath); err == nil {
						return wslPath
					}
				}
			}
		}

		return "/mnt/" + drive + rest
	}

	return strings.ReplaceAll(path, "\\", "/")
}
