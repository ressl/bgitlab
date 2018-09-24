package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	_ "strings"
	//"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	_ "reflect"
	//	"sync"
	"flag"
	"os/user"
	"time"
)

type GitLabProjects []struct {
	ID            int           `json:"id"`
	Description   string        `json:"description"`
	DefaultBranch string        `json:"default_branch"`
	Public        bool          `json:"public"`
	Visibility    string        `json:"visibility"`
	SSHURLToRepo  string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo string        `json:"http_url_to_repo"`
	WebURL        string        `json:"web_url"`
	TagList       []interface{} `json:"tag_list"`
	Owner         struct {
		ID               int         `json:"id"`
		Username         string      `json:"username"`
		Email            string      `json:"email"`
		Name             string      `json:"name"`
		State            string      `json:"state"`
		CreatedAt        interface{} `json:"created_at"`
		Bio              string      `json:"bio"`
		Location         string      `json:"location"`
		Skype            string      `json:"skype"`
		Linkedin         string      `json:"linkedin"`
		Twitter          string      `json:"twitter"`
		WebsiteURL       string      `json:"website_url"`
		Organization     string      `json:"organization"`
		ExternUID        string      `json:"extern_uid"`
		Provider         string      `json:"provider"`
		ThemeID          int         `json:"theme_id"`
		LastActivityOn   interface{} `json:"last_activity_on"`
		ColorSchemeID    int         `json:"color_scheme_id"`
		IsAdmin          bool        `json:"is_admin"`
		AvatarURL        string      `json:"avatar_url"`
		CanCreateGroup   bool        `json:"can_create_group"`
		CanCreateProject bool        `json:"can_create_project"`
		ProjectsLimit    int         `json:"projects_limit"`
		CurrentSignInAt  interface{} `json:"current_sign_in_at"`
		LastSignInAt     interface{} `json:"last_sign_in_at"`
		TwoFactorEnabled bool        `json:"two_factor_enabled"`
		Identities       interface{} `json:"identities"`
		External         bool        `json:"external"`
	} `json:"owner"`
	Name                     string    `json:"name"`
	NameWithNamespace        string    `json:"name_with_namespace"`
	Path                     string    `json:"path"`
	PathWithNamespace        string    `json:"path_with_namespace"`
	IssuesEnabled            bool      `json:"issues_enabled"`
	OpenIssuesCount          int       `json:"open_issues_count"`
	MergeRequestsEnabled     bool      `json:"merge_requests_enabled"`
	ApprovalsBeforeMerge     int       `json:"approvals_before_merge"`
	JobsEnabled              bool      `json:"jobs_enabled"`
	WikiEnabled              bool      `json:"wiki_enabled"`
	SnippetsEnabled          bool      `json:"snippets_enabled"`
	ContainerRegistryEnabled bool      `json:"container_registry_enabled"`
	CreatedAt                time.Time `json:"created_at"`
	LastActivityAt           time.Time `json:"last_activity_at"`
	CreatorID                int       `json:"creator_id"`
	Namespace                struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Kind     string `json:"kind"`
		FullPath string `json:"full_path"`
	} `json:"namespace"`
	ImportStatus string `json:"import_status"`
	ImportError  string `json:"import_error"`
	Permissions  struct {
		ProjectAccess interface{} `json:"project_access"`
		GroupAccess   interface{} `json:"group_access"`
	} `json:"permissions"`
	Archived                                  bool          `json:"archived"`
	AvatarURL                                 string        `json:"avatar_url"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	ForksCount                                int           `json:"forks_count"`
	StarCount                                 int           `json:"star_count"`
	RunnersToken                              string        `json:"runners_token"`
	PublicJobs                                bool          `json:"public_jobs"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	LfsEnabled                                bool          `json:"lfs_enabled"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	MergeMethod                               string        `json:"merge_method"`
	ForkedFromProject                         interface{}   `json:"forked_from_project"`
	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
	Statistics                                interface{}   `json:"statistics"`
	Links                                     struct {
		Self          string `json:"self"`
		Issues        string `json:"issues"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Labels        string `json:"labels"`
		Events        string `json:"events"`
		Members       string `json:"members"`
	} `json:"_links"`
	CiConfigPath interface{} `json:"ci_config_path"`
}

func main() {

	currentUser, _ := user.Current()
	tokenPtr := flag.String("token", "Yahnax1aeSaiCheireth", "GitLab api Token")
	urlPtr := flag.String("url", "https://git.example.ch", "GitLab URL")
	pagesPtr := flag.Int("pages", 50, "Number of items to list per page (max: 100)")
	cDirPtr := flag.String("cdir", currentUser.HomeDir+"/git", "Define the directory where the projets should be cloned")
	flag.Parse()

	url := fmt.Sprintf(*urlPtr+"/api/v4/projects?per_page=%d", *pagesPtr)

	var projectName []string
	var projectNamespace []string
	var projectSSHURL []string

	gitlabClient := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "bgitlab")
	req.Header.Set("PRIVATE-TOKEN", *tokenPtr)

	res, getErr := gitlabClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var test GitLabProjects
	err = json.Unmarshal(body, &test)
	if err != nil {
		fmt.Println(err)
	}

	for element := range test {
		projectNamespace = append(projectNamespace, test[element].PathWithNamespace)
		projectName = append(projectName, test[element].Name)
		projectSSHURL = append(projectSSHURL, test[element].SSHURLToRepo)
	}

	defer res.Body.Close()

	fmt.Printf("The len of projects: %d\n", len(test))
	cmds := "git"

	numberof, _ := strconv.Atoi(res.Header["X-Total-Pages"][0])
	fmt.Printf("The total number of pages: %d\n", numberof)

	for i := 2; i <= numberof; i++ {
		url = fmt.Sprintf(*urlPtr+"/api/v4/projects?per_page=%d", *pagesPtr, i)

		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "bgitlab")
		req.Header.Set("PRIVATE-TOKEN", "")

		res, getErr := gitlabClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		bodyfor, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()
		err = json.Unmarshal(bodyfor, &test)
		if err != nil {
			fmt.Println(err)
		}

		for element := range test {
			projectNamespace = append(projectNamespace, test[element].PathWithNamespace)
			projectName = append(projectName, test[element].Name)
			projectSSHURL = append(projectSSHURL, test[element].SSHURLToRepo)
		}
	}

	for element := range projectName {
		if _, err := os.Stat(*cDirPtr + "/" + projectNamespace[element] + "/.git"); os.IsNotExist(err) {
			fmt.Printf("Create folder %s\n", projectNamespace[element])
			os.MkdirAll(*cDirPtr+"/"+projectNamespace[element], os.ModePerm)

			fmt.Printf("Clone project %s into folder %s/%s\n", projectName[element], *cDirPtr, projectNamespace[element])
			args := []string{"clone", projectSSHURL[element], *cDirPtr + projectNamespace[element]}
			cmd := exec.Command(cmds, args...)
			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Fatalf("Failed to install %s: %s", cmds, err)
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatalf("Failed to install %s: %s", cmds, err)
			}

			if err := cmd.Start(); err != nil {
				log.Fatalf("Failed to install %s: %s", cmds, err)
			}
			stdoutp, _ := ioutil.ReadAll(stdout)
			fmt.Printf("%s\n", stdoutp)
			stderrp, _ := ioutil.ReadAll(stderr)
			fmt.Printf("%s\n", stderrp)

			if err := cmd.Wait(); err != nil {
				log.Fatalf("Failed to install %s: %s", cmds, err)
			}
		}
	}
}
