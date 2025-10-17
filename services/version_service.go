package services

import (
	"os/exec"
	"strings"
	"time"
)

type VersionInfo struct {
	Version     string    `json:"version"`
	BuildTime   time.Time `json:"build_time"`
	GitCommit   string    `json:"git_commit"`
	GitBranch   string    `json:"git_branch"`
	GoVersion   string    `json:"go_version"`
}

type VersionService struct{
	GitCommit string
	GitBranch string
	BuildTime string
}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func NewVersionServiceWithBuildInfo(gitCommit, gitBranch, buildTime string) *VersionService {
	return &VersionService{
		GitCommit: gitCommit,
		GitBranch: gitBranch,
		BuildTime: buildTime,
	}
}

// GetVersionInfo 獲取版本資訊
func (vs *VersionService) GetVersionInfo() *VersionInfo {
	version := &VersionInfo{
		Version:   "2.0.0",
		BuildTime: time.Now(),
	}

	// 優先使用構建時注入的 Git 資訊
	if vs.GitCommit != "" {
		version.GitCommit = vs.GitCommit
		if len(version.GitCommit) > 7 {
			version.GitCommit = version.GitCommit[:7] // 只取前7位
		}
	} else {
		// 如果沒有構建時資訊，嘗試從 Git 獲取
		if commit, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			version.GitCommit = strings.TrimSpace(string(commit))
			if len(version.GitCommit) > 7 {
				version.GitCommit = version.GitCommit[:7] // 只取前7位
			}
		}
	}

	// 優先使用構建時注入的 Git branch
	if vs.GitBranch != "" {
		version.GitBranch = vs.GitBranch
	} else {
		// 如果沒有構建時資訊，嘗試從 Git 獲取
		if branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
			version.GitBranch = strings.TrimSpace(string(branch))
		}
	}

	// 使用構建時間
	if vs.BuildTime != "" {
		if buildTime, err := time.Parse(time.RFC3339, vs.BuildTime); err == nil {
			version.BuildTime = buildTime
		}
	}

	// 獲取 Go 版本
	if goVersion, err := exec.Command("go", "version").Output(); err == nil {
		version.GoVersion = strings.TrimSpace(string(goVersion))
	}

	return version
}

// GetShortVersion 獲取簡短版本號（用於顯示）
func (vs *VersionService) GetShortVersion() string {
	info := vs.GetVersionInfo()
	
	// 如果有 Git commit，使用 commit hash 作為版本號
	if info.GitCommit != "" {
		return info.Version + "-" + info.GitCommit
	}
	
	// 如果有 Git branch 且不是 main/master，包含 branch 資訊
	if info.GitBranch != "" && info.GitBranch != "main" && info.GitBranch != "master" {
		return info.Version + "-" + info.GitBranch
	}
	
	return info.Version
}
