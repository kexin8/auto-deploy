package deploy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	releaseUrl = "https://api.github.com/repos/kexin8/auto-deploy/releases/latest"
	url        = "https://github.com/kexin8/auto-deploy/releases/download/%s/deploy-%s-%s.tgz"
)

// GetLatestVersion 获取最新版本号
func GetLatestVersion() (version string, err error) {
	// 请求github api 获取最新版本号
	resp, err := http.Get(releaseUrl)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 解析json
	var latestReleases map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&latestReleases); err != nil {
		return "", err
	}

	// 获取tag_name or name
	if latestReleases["tag_name"] != nil {
		return latestReleases["tag_name"].(string), nil
	} else if latestReleases["name"] != nil {
		return latestReleases["name"].(string), nil
	} else {
		return "", fmt.Errorf("can't find the latest version")
	}
}
