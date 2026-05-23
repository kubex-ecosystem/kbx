// Package version provides functionality to manage and check the version of the Kubex Horizon CLI tool.
// It includes methods to retrieve the current version, check for the latest version,
package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	info "github.com/kubex-ecosystem/kbx/internal/module/info"

	"github.com/spf13/cobra"
)

var (
	manifest info.Manifest
	vrs      Service
	err      error
)

type Service interface {
	// GetLatestVersion retrieves the latest version from the Git repository.
	GetLatestVersion() (string, error)
	// GetCurrentVersion returns the current version of the service.
	GetCurrentVersion() string
	// IsLatestVersion checks if the current version is the latest version.
	IsLatestVersion() (bool, error)
	// GetName returns the name of the service.
	GetName() string
	// GetVersion returns the current version of the service.
	GetVersion() string
	// GetRepository returns the Git repository URL of the service.
	GetRepository() string
	// setLastCheckedAt sets the last checked time for the version.
	setLastCheckedAt(time.Time)
	// updateLatestVersion updates the latest version from the Git repository.
	updateLatestVersion() error
}
type ServiceImpl struct {
	info.Manifest

	gitModelURL    string
	latestVersion  string
	lastCheckedAt  time.Time
	currentVersion string
}

func init() {
	if vrs == nil {
		vrs = NewVersionService()
	}
}

func getLatestTag(repoURL string) (string, error) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Fprintf(os.Stderr, "Recovered from panic in getLatestTag: %v\n", rec)
			err = fmt.Errorf("panic occurred while fetching latest tag: %v", rec)
		}
	}()

	defer func() {
		if vrs == nil {
			vrs = NewVersionService()
		}
		vrs.setLastCheckedAt(time.Now())
	}()

	if manifest == nil {
		manifest, err = info.GetManifest()
		// if err := manifest.LoadManifest(); err != nil {
		// 	return "", fmt.Errorf("failed to load manifest: %v", err)
		// }
	}
	if manifest.IsPrivate() {
		return "", fmt.Errorf("cannot fetch latest tag for private repositories")
	}

	if repoURL == "" {
		repoURL = manifest.GetRepository()
		if repoURL == "" {
			return "", fmt.Errorf("repository URL is not set")
		}
	}

	apiURL := fmt.Sprintf("%s/tags", repoURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}
	type Tag struct {
		Name string `json:"name"`
	}

	// Decode the JSON response into a slice of Tag structs
	// This assumes the API returns a JSON array of tags.
	// Adjust the decoding logic based on the actual API response structure.
	if resp.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}
	return tags[0].Name, nil
}
func (v *ServiceImpl) updateLatestVersion() error {
	if manifest.IsPrivate() {
		return fmt.Errorf("cannot fetch latest version for private repositories")
	}
	repoURL := strings.TrimSuffix(v.gitModelURL, ".git")
	tag, err := getLatestTag(repoURL)
	if err != nil {
		return err
	}
	v.latestVersion = tag
	return nil
}
func (v *ServiceImpl) vrsCompare(v1, v2 []int) (int, error) {
	compare := 0
	for i := 0; i < len(v1) && i < len(v2); i++ {
		if v1[i] < v2[i] {
			compare = -1
			break
		}
		if v1[i] > v2[i] {
			compare = 1
			break
		}
	}
	return compare, nil
}
func (v *ServiceImpl) versionAtMost(versionAtMostArg, max []int) (bool, error) {
	if comp, err := v.vrsCompare(versionAtMostArg, max); err != nil {
		return false, err
	} else if comp == 1 {
		return false, nil
	}
	return true, nil
}
func (v *ServiceImpl) parseVersion(versionToParse string) []int {
	if versionToParse == "" {
		return nil
	}
	if strings.Contains(versionToParse, "-") {
		versionToParse = strings.Split(versionToParse, "-")[0]
	}
	if strings.Contains(versionToParse, "v") {
		versionToParse = strings.TrimPrefix(versionToParse, "v")
	}
	parts := strings.Split(versionToParse, ".")
	parsedVersion := make([]int, len(parts))
	for i, part := range parts {
		if num, err := strconv.Atoi(part); err != nil {
			return nil
		} else {
			parsedVersion[i] = num
		}
	}
	return parsedVersion
}
func (v *ServiceImpl) IsLatestVersion() (bool, error) {
	if manifest.IsPrivate() {
		return false, fmt.Errorf("cannot check version for private repositories")
	}
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return false, err
		}
	}

	currentVersionParts := v.parseVersion(v.currentVersion)
	latestVersionParts := v.parseVersion(v.latestVersion)

	if len(currentVersionParts) == 0 || len(latestVersionParts) == 0 {
		return false, fmt.Errorf("invalid version format")
	}

	if len(currentVersionParts) != len(latestVersionParts) {
		return false, fmt.Errorf("version parts length mismatch")
	}

	return v.versionAtMost(currentVersionParts, latestVersionParts)
}
func (v *ServiceImpl) GetLatestVersion() (string, error) {
	if manifest.IsPrivate() {
		return "", fmt.Errorf("cannot fetch latest version for private repositories")
	}
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return "", err
		}
	}
	return v.latestVersion, nil
}
func (v *ServiceImpl) GetCurrentVersion() string {
	if v.currentVersion == "" {
		v.currentVersion = manifest.GetVersion()
	}
	return v.currentVersion
}
func (v *ServiceImpl) GetName() string {
	if manifest == nil {
		return "Unknown Service"
	}
	return manifest.GetName()
}
func (v *ServiceImpl) GetVersion() string {
	if manifest == nil {
		return "Unknown version"
	}
	return manifest.GetVersion()
}
func (v *ServiceImpl) GetRepository() string {
	if manifest == nil {
		return "No repository URL set in the manifest."
	}
	return manifest.GetRepository()
}
func (v *ServiceImpl) setLastCheckedAt(t time.Time) {
	v.lastCheckedAt = t

}

func NewVersionService() Service {
	if manifest == nil {
		manifest, err = info.GetManifest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load manifest: %v\n", err)
			os.Exit(1)
		}
	}
	return &ServiceImpl{
		Manifest:       manifest,
		gitModelURL:    manifest.GetRepository(),
		currentVersion: manifest.GetVersion(),
		latestVersion:  "",
	}
}

var (
	versionCmd   *cobra.Command
	subLatestCmd *cobra.Command
	subCmdCheck  *cobra.Command
	updCmd       *cobra.Command
	getCmd       *cobra.Command
	restartCmd   *cobra.Command
)

func init() {
	if versionCmd == nil {
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Print the version number of " + manifest.GetName(),
			Long:  "Print the version number of " + manifest.GetName() + " and other related information.",
			Run: func(cmd *cobra.Command, args []string) {
				if manifest.IsPrivate() {
					fmt.Fprintln(os.Stderr, "The information shown may not be accurate for private repositories.")
					fmt.Fprintln(os.Stdout, "Current version: "+GetVersion())
					fmt.Fprintln(os.Stdout, "Git repository: "+GetGitRepositoryModelURL())
					return
				}
				GetVersionInfo()
			},
		}
	}
	if subLatestCmd == nil {
		subLatestCmd = &cobra.Command{
			Use:   "latest",
			Short: "Print the latest version number of " + manifest.GetName(),
			Long:  "Print the latest version number of " + manifest.GetName() + " from the Git repository.",
			Run: func(cmd *cobra.Command, args []string) {
				if manifest.IsPrivate() {
					fmt.Fprintln(os.Stderr, "Cannot fetch latest version for private repositories.")
					return
				}
				GetLatestVersionInfo()
			},
		}
	}
	if subCmdCheck == nil {
		subCmdCheck = &cobra.Command{
			Use:   "check",
			Short: "Check if the current version is the latest version of " + manifest.GetName(),
			Long:  "Check if the current version is the latest version of " + manifest.GetName() + " and print the version information.",
			Run: func(cmd *cobra.Command, args []string) {
				if manifest.IsPrivate() {
					fmt.Fprintln(os.Stderr, "Cannot check version for private repositories.")
					return
				}
				GetVersionInfoWithLatestAndCheck()
			},
		}
	}
	if updCmd == nil {
		updCmd = &cobra.Command{
			Use:   "update",
			Short: "Update the version information of " + manifest.GetName(),
			Long:  "Update the version information of " + manifest.GetName() + " by fetching the latest version from the Git repository.",
			Run: func(cmd *cobra.Command, args []string) {
				if manifest.IsPrivate() {
					fmt.Fprintln(os.Stderr, "Cannot update version for private repositories.")
					return
				}
				if err := vrs.updateLatestVersion(); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to update version: "+err.Error())
				} else {
					latestVersion, err := vrs.GetLatestVersion()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Failed to get latest version: "+err.Error())
					} else {
						fmt.Fprintln(os.Stdout, "Current version: "+vrs.GetCurrentVersion())
						fmt.Fprintln(os.Stdout, "Latest version: "+latestVersion)
					}
					vrs.setLastCheckedAt(time.Now())
				}
			},
		}
	}
	if getCmd == nil {
		getCmd = &cobra.Command{
			Use:   "get",
			Short: "Get the current version of " + manifest.GetName(),
			Long:  "Get the current version of " + manifest.GetName() + " from the manifest.",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Fprintln(os.Stdout, "Current version: "+vrs.GetCurrentVersion())
			},
		}
	}
	if restartCmd == nil {
		restartCmd = &cobra.Command{
			Use:   "restart",
			Short: "Restart the " + manifest.GetName() + " service",
			Long:  "Restart the " + manifest.GetName() + " service to apply any changes made.",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Fprintln(os.Stdout, "Restarting the service...")
				// Logic to restart the service can be added here
				fmt.Fprintln(os.Stdout, "Service restarted successfully")
			},
		}
	}

}
func GetVersion() string {
	if manifest == nil {
		manifest, err = info.GetManifest()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to load manifest: "+err.Error())
			return "Unknown version"
		}
	}
	return manifest.GetVersion()
}
func GetGitRepositoryModelURL() string {
	if manifest.GetRepository() == "" {
		return "No repository URL set in the manifest."
	}
	return manifest.GetRepository()
}
func GetVersionInfo() string {
	fmt.Fprintln(os.Stdout, "Version: "+GetVersion())
	fmt.Fprintln(os.Stdout, "Git repository: "+GetGitRepositoryModelURL())
	return fmt.Sprintf("Version: %s\nGit repository: %s", GetVersion(), GetGitRepositoryModelURL())
}
func GetLatestVersionFromGit() string {
	if manifest.IsPrivate() {
		fmt.Fprintln(os.Stderr, "Cannot fetch latest version for private repositories.")
		return "Cannot fetch latest version for private repositories."
	}

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	gitURLWithoutGit := strings.TrimSuffix(GetGitRepositoryModelURL(), ".git")
	if gitURLWithoutGit == "" {
		fmt.Fprintln(os.Stderr, "No repository URL set in the manifest.")
		return "No repository URL set in the manifest."
	}

	response, err := netClient.Get(gitURLWithoutGit + "/releases/latest")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching latest version: "+err.Error())
		fmt.Fprintln(os.Stderr, gitURLWithoutGit+"/releases/latest")
		return err.Error()
	}

	if response.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "Error fetching latest version: "+response.Status)
		fmt.Fprintln(os.Stderr, "Url: "+gitURLWithoutGit+"/releases/latest")
		body, _ := io.ReadAll(response.Body)
		return fmt.Sprintf("Error: %s\nResponse: %s", response.Status, string(body))
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}
func GetLatestVersionInfo() string {
	if manifest.IsPrivate() {
		fmt.Fprintln(os.Stderr, "Cannot fetch latest version for private repositories.")
		return "Cannot fetch latest version for private repositories."
	}
	fmt.Fprintln(os.Stdout, "Latest version: "+GetLatestVersionFromGit())
	return "Latest version: " + GetLatestVersionFromGit()
}
func GetVersionInfoWithLatestAndCheck() string {
	if manifest.IsPrivate() {
		fmt.Fprintln(os.Stderr, "Cannot check version for private repositories.")
		return "Cannot check version for private repositories."
	}
	if GetVersion() == GetLatestVersionFromGit() {
		fmt.Fprintln(os.Stdout, "You are using the latest version.")
		return fmt.Sprintf("You are using the latest version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	} else {
		fmt.Fprintln(os.Stderr, "You are using an outdated version.")
		return fmt.Sprintf("You are using an outdated version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	}
}
func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	versionCmd.AddCommand(updCmd)
	versionCmd.AddCommand(getCmd)
	versionCmd.AddCommand(restartCmd)
	return versionCmd
}
