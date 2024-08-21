package harbor

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	caCerts          = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----"
	get_confirations = "/api/v2.0/configurations"
	get_projects     = "/api/v2.0/projects"
	get_repositories = "/api/v2.0/projects/%s/repositories"
	get_artifacts    = "/api/v2.0/projects/%s/repositories/%s/artifacts"
)

type Page struct {
	//当前页
	CurrentPage int `json:"currentPage"`
	//分页数量
	PageSize int `json:"pageSize"`
	//总条数
	TotalCount int `json:"totalCount"`
	//总页数
	TotalPage int `json:"totalPage"`
	//数据
	Data any `json:"data"`
}

// @Author: morris
type Harbor struct {
	username string
	password string
	host     string
	scheme   string
	client   http.Client
}
type Artifact struct {
	Accessories   interface{} `json:"accessories"`
	Host          string      `json:"host"`
	AdditionLinks struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"build_history"`
	} `json:"addition_links"`
	Digest     string `json:"digest"`
	ExtraAttrs struct {
		Architecture string `json:"architecture"`
		Author       string `json:"author"`
		Config       struct {
			Entrypoint   []string `json:"Entrypoint"`
			Env          []string `json:"Env"`
			ExposedPorts struct {
				Tcp struct {
				} `json:"8082/tcp"`
			} `json:"ExposedPorts"`
			Labels struct {
				MAINTAINER                     string    `json:"MAINTAINER"`
				OrgOpencontainersImageCreated  time.Time `json:"org.opencontainers.image.created"`
				OrgOpencontainersImageRevision string    `json:"org.opencontainers.image.revision"`
				OrgOpencontainersImageSource   string    `json:"org.opencontainers.image.source"`
				OrgOpencontainersImageUrl      string    `json:"org.opencontainers.image.url"`
			} `json:"Labels"`
			WorkingDir string `json:"WorkingDir"`
		} `json:"config"`
		Created time.Time `json:"created"`
		Os      string    `json:"os"`
	} `json:"extra_attrs"`
	Icon              string      `json:"icon"`
	Id                int         `json:"id"`
	Labels            interface{} `json:"labels"`
	ManifestMediaType string      `json:"manifest_media_type"`
	MediaType         string      `json:"media_type"`
	ProjectId         int         `json:"project_id"`
	PullTime          time.Time   `json:"pull_time"`
	PushTime          time.Time   `json:"push_time"`
	References        interface{} `json:"references"`
	RepositoryId      int         `json:"repository_id"`
	Size              int         `json:"size"`
	Tags              []struct {
		ArtifactId   int       `json:"artifact_id"`
		Id           int       `json:"id"`
		Immutable    bool      `json:"immutable"`
		Name         string    `json:"name"`
		PullTime     time.Time `json:"pull_time"`
		PushTime     time.Time `json:"push_time"`
		RepositoryId int       `json:"repository_id"`
		Signed       bool      `json:"signed"`
	} `json:"tags"`
	Type string `json:"type"`
}
type Repository struct {
	ArtifactCount int       `json:"artifact_count"`
	CreationTime  time.Time `json:"creation_time"`
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	ProjectId     int       `json:"project_id"`
	PullCount     int       `json:"pull_count"`
	UpdateTime    time.Time `json:"update_time"`
}
type Project struct {
	CreationTime       time.Time `json:"creation_time"`
	CurrentUserRoleId  int       `json:"current_user_role_id"`
	CurrentUserRoleIds []int     `json:"current_user_role_ids"`
	CveAllowlist       struct {
		CreationTime time.Time     `json:"creation_time"`
		Id           int           `json:"id"`
		Items        []interface{} `json:"items"`
		ProjectId    int           `json:"project_id"`
		UpdateTime   time.Time     `json:"update_time"`
	} `json:"cve_allowlist"`
	Metadata struct {
		Public string `json:"public"`
	} `json:"metadata"`
	Name       string    `json:"name"`
	OwnerId    int       `json:"owner_id"`
	OwnerName  string    `json:"owner_name"`
	ProjectId  int       `json:"project_id"`
	RepoCount  int       `json:"repo_count"`
	UpdateTime time.Time `json:"update_time"`
}

func InitHarbor(scheme, host, username, password string) (harbor *Harbor, err error) {
	url := fmt.Sprintf("%s://%s%s", scheme, host, get_confirations)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(username, password)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCerts))
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resMap := make(map[string]any)
	err = json.Unmarshal(respBytes, &resMap)
	if err != nil {
		return
	}
	if _, ok := resMap["errors"]; ok {
		err = fmt.Errorf("认证失败：%s", string(respBytes))
		return
	}
	//赋值harbor
	harbor = &Harbor{
		username, password, host, scheme, client,
	}
	return
}

// 查询项目列表
func (h *Harbor) GetProjects(curPage, pageSize int, keyword string) Page {
	page := Page{
		CurrentPage: curPage,
		PageSize:    pageSize,
	}
	url := fmt.Sprintf("%s://%s%s", h.scheme, h.host, get_projects)
	url = fmt.Sprintf("%s?page=%d&page_size=%d", url, curPage, pageSize)
	if keyword != "" {
		url = fmt.Sprintf("%s&q=name=~%s", url, keyword)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return page
	}
	req.SetBasicAuth(h.username, h.password)
	resp, err := h.client.Do(req)
	if err != nil {
		return page
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return page
	}
	xTotalCount := resp.Header.Get("x-total-count")
	newXTotalCount, _ := strconv.Atoi(xTotalCount)
	projects := make([]Project, newXTotalCount)
	err = json.Unmarshal(respBody, &projects)
	if err != nil {
		return page
	}
	page.Data = projects
	page.TotalPage = int(math.Ceil(float64(newXTotalCount) / float64(pageSize)))
	page.TotalCount = newXTotalCount
	return page
}
func (h *Harbor) GetRepositories(projectName string, curPage, pageSize int, keyword string) Page {
	page := Page{
		CurrentPage: curPage,
		PageSize:    pageSize,
	}
	path := fmt.Sprintf(get_repositories, projectName)
	url := fmt.Sprintf("%s://%s%s", h.scheme, h.host, path)
	url = fmt.Sprintf("%s?page=%d&page_size=%d", url, curPage, pageSize)
	if keyword != "" {
		url = fmt.Sprintf("%s&q=name=~%s", url, keyword)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return page
	}
	req.SetBasicAuth(h.username, h.password)
	resp, err := h.client.Do(req)
	if err != nil {
		return page
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return page
	}
	xTotalCount := resp.Header.Get("x-total-count")
	newXTotalCount, _ := strconv.Atoi(xTotalCount)
	repositories := make([]Repository, newXTotalCount)
	err = json.Unmarshal(respBody, &repositories)
	if err != nil {
		return page
	}
	page.Data = repositories
	page.TotalPage = int(math.Ceil(float64(newXTotalCount) / float64(pageSize)))
	page.TotalCount = newXTotalCount
	return page
}
func (h *Harbor) GetArtifacts(projectName, repositoryName string, curPage, pageSize int, keyword string) Page {
	page := Page{
		CurrentPage: curPage,
		PageSize:    pageSize,
	}
	path := fmt.Sprintf(get_artifacts, projectName, repositoryName)
	url := fmt.Sprintf("%s://%s%s", h.scheme, h.host, path)
	url = fmt.Sprintf("%s?page=%d&page_size=%d", url, curPage, pageSize)
	if keyword != "" {
		url = fmt.Sprintf("%s&q=tags=~%s", url, keyword)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return page
	}
	req.SetBasicAuth(h.username, h.password)
	resp, err := h.client.Do(req)
	if err != nil {
		return page
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return page
	}
	xTotalCount := resp.Header.Get("x-total-count")
	newXTotalCount, _ := strconv.Atoi(xTotalCount)
	artifacts := make([]Artifact, newXTotalCount)
	err = json.Unmarshal(respBody, &artifacts)
	if err != nil {
		return page
	}
	for index, _ := range artifacts {
		artifacts[index].Host = h.host
	}
	page.Data = artifacts
	page.TotalPage = int(math.Ceil(float64(newXTotalCount) / float64(pageSize)))
	page.TotalCount = newXTotalCount
	return page
}
func (h *Harbor) MatchImage(keyword string) []string {
	//kubeimooc-server:v1
	keywordArr := strings.Split(keyword, ":")
	image := ""
	tag := ""
	if len(keywordArr) == 1 {
		image = keywordArr[0]
	} else {
		image = keywordArr[0]
		tag = keywordArr[1]
	}
	images := make([]string, 0)
	//循环projects
	projectsPage := h.GetProjects(1, 20, "")
	if projectsPage.TotalCount == 0 {
		return images
	}
	var countFlag int
	projects := projectsPage.Data.([]Project)
	for _, project := range projects {
		//匹配镜像仓库
		repositoriesPage := h.GetRepositories(project.Name, 1, 10, image)
		if repositoriesPage.TotalCount == 0 {
			continue
		}
		repositories := repositoriesPage.Data.([]Repository)
		for _, repository := range repositories {
			//预期：kubeimooc-server 实际：kubeimoc/kubeimooc-server
			repositoryName := filepath.Base(repository.Name)
			//匹配tag
			artifactsPage := h.GetArtifacts(project.Name, repositoryName, 1, 10, tag)
			if artifactsPage.TotalCount == 0 {
				continue
			}
			artifacts := artifactsPage.Data.([]Artifact)
			for _, artifact := range artifacts {
				for _, s := range artifact.Tags {
					imageinfo := fmt.Sprintf("%s/%s/%s:%s", artifact.Host, project.Name, repositoryName, s.Name)
					//当匹配到10个结束
					images = append(images, imageinfo)
					countFlag++
					if countFlag == 10 {
						return images
					}
				}
			}
		}
	}
	return images
}
