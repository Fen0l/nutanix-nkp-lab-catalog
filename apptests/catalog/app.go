// Package catalog implements a generic install scenario for applications laid
// out in the NKP catalog format (applications/<name>/<version>/helmrelease).
//
// It reuses the public test harness from mesosphere/kommander-applications
// (kind cluster + Flux) and applies the application's helmrelease kustomization
// with the same substitution variables Kommander injects at deploy time.
package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mesosphere/kommander-applications/apptests/environment"
	"golang.org/x/mod/semver"
)

const (
	// DefaultNamespace is used as ${releaseNamespace} and ${workspaceNamespace}
	// in tests. It exists on every cluster.
	DefaultNamespace = "default"
	PollInterval     = 5 * time.Second
)

// App is a catalog application under test.
type App struct {
	Name    string
	Version string
}

// New returns an App for the given name and version. When version is empty the
// latest version directory (semver order) is used.
func New(name, version string) (*App, error) {
	if version == "" {
		v, err := LatestVersion(name)
		if err != nil {
			return nil, err
		}
		version = v
	}
	return &App{Name: name, Version: version}, nil
}

// Dir returns applications/<name>/<version>.
func (a *App) Dir() (string, error) {
	root, err := ApplicationsDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, a.Name, a.Version)
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("application %s version %s not found at %s", a.Name, a.Version, dir)
	}
	return dir, nil
}

// Substitutions mirrors the variables Kommander sets on the Flux Kustomization
// postBuild.substitute for an AppDeployment.
func (a *App) Substitutions() map[string]string {
	return map[string]string{
		"releaseName":        a.Name,
		"releaseNamespace":   DefaultNamespace,
		"workspaceNamespace": DefaultNamespace,
		"appName":            a.Name,
		"appVersion":         a.Version,
		// Multi-instance apps resolve ${projectNamespace:=<app>} into appNamespace in the
		// outer (Kommander) substitution; the kind harness has no outer layer, so set it here.
		"appNamespace": a.Name,
	}
}

// Install applies the helmrelease kustomization of the application.
func (a *App) Install(ctx context.Context, env *environment.Env) error {
	dir, err := a.Dir()
	if err != nil {
		return err
	}
	return env.ApplyKustomizations(ctx, filepath.Join(dir, "helmrelease"), a.Substitutions())
}

// ApplicationsDir walks up from the working directory to find applications/.
func ApplicationsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "applications")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no applications/ directory found above %s", dir)
		}
		dir = parent
	}
}

// Versions lists the version directories of an application, oldest first.
func Versions(name string) ([]string, error) {
	root, err := ApplicationsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(root, name))
	if err != nil {
		return nil, err
	}
	var versions []string
	for _, e := range entries {
		if e.IsDir() && e.Name()[0] != '.' {
			versions = append(versions, e.Name())
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare("v"+versions[i], "v"+versions[j]) < 0
	})
	return versions, nil
}

// LatestVersion returns the highest version directory of an application.
func LatestVersion(name string) (string, error) {
	versions, err := Versions(name)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found for application %s", name)
	}
	return versions[len(versions)-1], nil
}
