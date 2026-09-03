package agentbox

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
)

const defaultBaseImage = "docker.io/retailcrm/agentbox-base:v0.0.1@sha256:493b753044eb82f90b99afb3974033591ccbd7037b39e7963dfe49f38790376e"

// TemplateStep is one layer in a template build.
type TemplateStep struct {
	Type            string   `json:"type"`
	Args            []string `json:"args"`
	Force           bool     `json:"force"`
	FilesHash       string   `json:"filesHash,omitempty"`
	ForceUpload     bool     `json:"forceUpload,omitempty"`
	ResolveSymlinks bool     `json:"resolveSymlinks,omitempty"`
	Gzip            bool     `json:"gzip,omitempty"`
}

// CopyOptions configures a COPY template layer.
type CopyOptions struct {
	User            string
	Mode            os.FileMode
	ForceUpload     bool
	ResolveSymlinks bool
	Gzip            *bool
}

// PackageInstallOptions configures npm/bun installs.
type PackageInstallOptions struct{ Global, Dev bool }

// AptInstallOptions configures apt-get.
type AptInstallOptions struct{ NoInstallRecommends, FixMissing bool }

// GitCloneOptions configures a git clone layer.
type GitCloneOptions struct {
	Path, Branch, User string
	Depth              int
}

// TemplateBuilder builds a declarative template definition.
type TemplateBuilder struct {
	contextPath  string
	ignore       []string
	baseImage    string
	baseTemplate string
	registry     *api.FromImageRegistry
	startCmd     string
	readyCmd     string
	force        bool
	forceNext    bool
	steps        []TemplateStep
	err          error
}

// NewTemplate starts a template definition. Empty contextPath uses the current directory.
func NewTemplate(contextPath string, ignore ...string) *TemplateBuilder {
	if contextPath == "" {
		contextPath = "."
	}
	return &TemplateBuilder{contextPath: contextPath, ignore: slices.Clone(ignore), baseImage: defaultBaseImage}
}

func (builder *TemplateBuilder) fromImage(image string) *TemplateBuilder {
	if image == "" {
		builder.err = &InvalidArgumentError{Message: "base image cannot be empty"}
		return builder
	}
	builder.baseImage, builder.baseTemplate = image, ""
	if builder.forceNext {
		builder.force = true
	}
	return builder
}

// FromImage selects an OCI image.
func (builder *TemplateBuilder) FromImage(image string) *TemplateBuilder {
	return builder.fromImage(image)
}
func (builder *TemplateBuilder) FromDebian(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "stable"
	}
	return builder.fromImage("debian:" + variant)
}
func (builder *TemplateBuilder) FromUbuntu(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "latest"
	}
	return builder.fromImage("ubuntu:" + variant)
}
func (builder *TemplateBuilder) FromFedora(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "44"
	}
	return builder.fromImage("fedora:" + variant)
}
func (builder *TemplateBuilder) FromAlpine(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "3.24"
	}
	return builder.fromImage("alpine:" + variant)
}
func (builder *TemplateBuilder) FromArch(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "latest"
	}
	return builder.fromImage("archlinux:" + variant)
}
func (builder *TemplateBuilder) FromPython(version string) *TemplateBuilder {
	if version == "" {
		version = "3"
	}
	return builder.fromImage("python:" + version)
}
func (builder *TemplateBuilder) FromNode(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "lts"
	}
	return builder.fromImage("node:" + variant)
}
func (builder *TemplateBuilder) FromBun(variant string) *TemplateBuilder {
	if variant == "" {
		variant = "latest"
	}
	return builder.fromImage("oven/bun:" + variant)
}
func (builder *TemplateBuilder) FromBase() *TemplateBuilder {
	return builder.fromImage(defaultBaseImage)
}

// FromTemplate selects another AgentBox template.
func (builder *TemplateBuilder) FromTemplate(template string) *TemplateBuilder {
	builder.baseTemplate, builder.baseImage = template, ""
	if builder.forceNext {
		builder.force = true
	}
	return builder
}

// FromRegistry selects a password-authenticated OCI registry image.
func (builder *TemplateBuilder) FromRegistry(image, username, password string) *TemplateBuilder {
	builder.fromImage(image)
	var registry api.FromImageRegistry
	if err := registry.FromGeneralRegistry(api.GeneralRegistry{Type: api.Registry, Username: username, Password: password}); err != nil {
		builder.err = err
	}
	builder.registry = &registry
	return builder
}

// FromAWSRegistry selects an AWS ECR image.
func (builder *TemplateBuilder) FromAWSRegistry(image, accessKeyID, secretAccessKey, region string) *TemplateBuilder {
	builder.fromImage(image)
	var registry api.FromImageRegistry
	if err := registry.FromAWSRegistry(api.AWSRegistry{Type: api.Aws, AwsAccessKeyID: accessKeyID, AwsSecretAccessKey: secretAccessKey, AwsRegion: region}); err != nil {
		builder.err = err
	}
	builder.registry = &registry
	return builder
}

// FromGCPRegistry selects a GCP Artifact Registry image.
func (builder *TemplateBuilder) FromGCPRegistry(image, serviceAccountJSON string) *TemplateBuilder {
	builder.fromImage(image)
	var registry api.FromImageRegistry
	if err := registry.FromGCPRegistry(api.GCPRegistry{Type: api.Gcp, ServiceAccountJSON: serviceAccountJSON}); err != nil {
		builder.err = err
	}
	builder.registry = &registry
	return builder
}

// FromDockerfile parses common FROM/RUN/COPY/ENV/WORKDIR/USER directives.
func (builder *TemplateBuilder) FromDockerfile(contentOrPath string) *TemplateBuilder {
	content := contentOrPath
	if data, err := os.ReadFile(filepath.Join(builder.contextPath, contentOrPath)); err == nil {
		content = string(data)
	}
	lines := joinDockerfileLines(content)
	fromCount := 0
	for _, line := range lines {
		if fields := strings.Fields(line); len(fields) > 0 && strings.EqualFold(fields[0], "FROM") {
			fromCount++
		}
	}
	if fromCount != 1 {
		builder.err = &InvalidArgumentError{Message: "Dockerfile must contain exactly one FROM instruction"}
		return builder
	}
	builder.User("root").Workdir("/")
	userChanged, workdirChanged := false, false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		directive := strings.ToUpper(fields[0])
		rest := strings.TrimSpace(line[len(fields[0]):])
		switch directive {
		case "FROM":
			values, err := splitDockerWords(rest)
			if err != nil {
				builder.err = err
				continue
			}
			if len(values) == 0 {
				builder.err = &InvalidArgumentError{Message: "invalid Dockerfile FROM"}
			} else {
				builder.FromImage(values[0])
			}
		case "RUN":
			builder.Run(rest)
		case "WORKDIR":
			builder.Workdir(rest)
			workdirChanged = true
		case "USER":
			builder.User(rest)
			userChanged = true
		case "ENV", "ARG":
			values, err := splitDockerWords(rest)
			if err != nil {
				builder.err = err
				continue
			}
			env := map[string]string{}
			for index := 0; index < len(values); index++ {
				key, value, found := strings.Cut(values[index], "=")
				if found {
					env[key] = value
				} else if index+1 < len(values) {
					env[key] = values[index+1]
					index++
				} else if directive == "ARG" {
					env[key] = ""
				}
			}
			builder.Env(env)
		case "COPY", "ADD":
			values, err := splitDockerWords(rest)
			if err != nil {
				builder.err = err
				continue
			}
			copyOptions := &CopyOptions{}
			for len(values) > 0 && strings.HasPrefix(values[0], "--") {
				if owner, ok := strings.CutPrefix(values[0], "--chown="); ok {
					copyOptions.User = owner
				}
				values = values[1:]
			}
			if len(values) < 2 {
				builder.err = &InvalidArgumentError{Message: "invalid Dockerfile " + directive}
			} else {
				for _, source := range values[:len(values)-1] {
					builder.Copy(source, values[len(values)-1], copyOptions)
				}
			}
		case "CMD", "ENTRYPOINT":
			command := rest
			var execArgs []string
			if json.Unmarshal([]byte(rest), &execArgs) == nil {
				command = strings.Join(execArgs, " ")
			}
			builder.Start(command, WaitForTimeout(20*time.Second))
		case "EXPOSE", "VOLUME":
			// These directives do not require a template build operation.
		default:
			// Match the other SDKs: ignore unknown Dockerfile directives.
		}
	}
	if !userChanged {
		builder.User("user")
	}
	if !workdirChanged {
		builder.Workdir("/home/user")
	}
	return builder
}

func (builder *TemplateBuilder) add(kind string, args ...string) *TemplateBuilder {
	builder.steps = append(builder.steps, TemplateStep{Type: kind, Args: args, Force: builder.forceNext})
	return builder
}

// Copy adds a file or directory from the context.
func (builder *TemplateBuilder) Copy(source, destination string, options *CopyOptions) *TemplateBuilder {
	if filepath.IsAbs(source) || strings.HasPrefix(filepath.Clean(source), "..") {
		builder.err = &InvalidArgumentError{Message: "copy source must stay inside template context"}
		return builder
	}
	step := TemplateStep{Type: "COPY", Args: []string{source, destination, "", ""}, Force: builder.forceNext, Gzip: true}
	if options != nil {
		step.Args[2] = options.User
		if options.Mode != 0 {
			step.Args[3] = fmt.Sprintf("%04o", options.Mode.Perm())
		}
		step.ForceUpload, step.ResolveSymlinks = options.ForceUpload, options.ResolveSymlinks
		step.Force = step.Force || options.ForceUpload
		if options.Gzip != nil {
			step.Gzip = *options.Gzip
		}
	}
	builder.steps = append(builder.steps, step)
	return builder
}

func (builder *TemplateBuilder) Run(commands ...string) *TemplateBuilder {
	return builder.add("RUN", strings.Join(commands, " && "))
}
func (builder *TemplateBuilder) RunAs(user string, commands ...string) *TemplateBuilder {
	return builder.add("RUN", strings.Join(commands, " && "), user)
}
func (builder *TemplateBuilder) Workdir(path string) *TemplateBuilder {
	return builder.add("WORKDIR", path)
}
func (builder *TemplateBuilder) User(user string) *TemplateBuilder { return builder.add("USER", user) }
func (builder *TemplateBuilder) Env(values map[string]string) *TemplateBuilder {
	args := make([]string, 0, len(values)*2)
	keys := slices.Sorted(maps.Keys(values))
	for _, key := range keys {
		args = append(args, key, values[key])
	}
	return builder.add("ENV", args...)
}
func (builder *TemplateBuilder) Remove(paths ...string) *TemplateBuilder {
	quoted := make([]string, len(paths))
	for index, path := range paths {
		quoted[index] = shellQuote(path)
	}
	return builder.Run("rm -rf " + strings.Join(quoted, " "))
}
func (builder *TemplateBuilder) Rename(source, destination string) *TemplateBuilder {
	return builder.Run("mv " + shellQuote(source) + " " + shellQuote(destination))
}
func (builder *TemplateBuilder) MakeDir(paths ...string) *TemplateBuilder {
	quoted := make([]string, len(paths))
	for index, path := range paths {
		quoted[index] = shellQuote(path)
	}
	return builder.Run("mkdir -p " + strings.Join(quoted, " "))
}
func (builder *TemplateBuilder) Symlink(source, destination string) *TemplateBuilder {
	return builder.Run("ln -s " + shellQuote(source) + " " + shellQuote(destination))
}
func (builder *TemplateBuilder) PipInstall(packages ...string) *TemplateBuilder {
	if len(packages) == 0 {
		packages = []string{"."}
	}
	return builder.RunAs("root", "pip install "+strings.Join(packages, " "))
}
func (builder *TemplateBuilder) NPMInstall(options PackageInstallOptions, packages ...string) *TemplateBuilder {
	flags := ""
	if options.Global {
		flags += " -g"
	}
	if options.Dev {
		flags += " --save-dev"
	}
	user := ""
	if options.Global {
		user = "root"
	}
	return builder.RunAs(user, strings.TrimSpace("npm install"+flags+" "+strings.Join(packages, " ")))
}
func (builder *TemplateBuilder) BunInstall(options PackageInstallOptions, packages ...string) *TemplateBuilder {
	flags := ""
	if options.Global {
		flags += " -g"
	}
	if options.Dev {
		flags += " --dev"
	}
	user := ""
	if options.Global {
		user = "root"
	}
	return builder.RunAs(user, strings.TrimSpace("bun install"+flags+" "+strings.Join(packages, " ")))
}
func (builder *TemplateBuilder) AptInstall(options AptInstallOptions, packages ...string) *TemplateBuilder {
	flags := ""
	if options.NoInstallRecommends {
		flags += " --no-install-recommends"
	}
	if options.FixMissing {
		flags += " --fix-missing"
	}
	return builder.RunAs("root", "apt-get update", "DEBIAN_FRONTEND=noninteractive DEBCONF_NOWARNINGS=yes apt-get install -y"+flags+" "+strings.Join(packages, " "))
}
func (builder *TemplateBuilder) GitClone(repository string, options *GitCloneOptions) *TemplateBuilder {
	args := []string{"git clone", shellQuote(repository)}
	user := ""
	if options != nil {
		user = options.User
		if options.Branch != "" {
			args = append(args, "--branch", shellQuote(options.Branch), "--single-branch")
		}
		if options.Depth > 0 {
			args = append(args, "--depth", strconv.Itoa(options.Depth))
		}
		if options.Path != "" {
			args = append(args, shellQuote(options.Path))
		}
	}
	return builder.RunAs(user, strings.Join(args, " "))
}

// SkipCache forces this and all subsequent layers.
func (builder *TemplateBuilder) SkipCache() *TemplateBuilder {
	builder.forceNext = true
	return builder
}
func (builder *TemplateBuilder) Start(command, readyCommand string) *TemplateBuilder {
	builder.startCmd, builder.readyCmd = command, readyCommand
	return builder
}
func (builder *TemplateBuilder) Ready(command string) *TemplateBuilder {
	builder.readyCmd = command
	return builder
}
func WaitForPort(port int) string {
	return fmt.Sprintf(`[ -n "$(ss -Htuln sport = :%d)" ]`, port)
}
func WaitForURL(value string, status int) string {
	return fmt.Sprintf(`curl -s -o /dev/null -w "%%{http_code}" %s | grep -q "%d"`, shellQuote(value), status)
}
func WaitForFile(path string) string       { return "test -e " + shellQuote(path) }
func WaitForProcess(process string) string { return "pgrep " + shellQuote(process) + " >/dev/null" }

// WaitForTimeout waits a fixed duration before marking a service ready.
func WaitForTimeout(timeout time.Duration) string {
	seconds := max(1, int(timeout/time.Second))
	return "sleep " + strconv.Itoa(seconds)
}

// JSON returns the build request representation without computed copy hashes.
func (builder *TemplateBuilder) JSON() ([]byte, error) {
	if builder.err != nil {
		return nil, builder.err
	}
	return json.Marshal(builder.request(nil))
}

// Dockerfile returns a human-readable equivalent definition.
func (builder *TemplateBuilder) Dockerfile() string {
	var output strings.Builder
	if builder.baseImage != "" {
		fmt.Fprintf(&output, "FROM %s\n", builder.baseImage)
	}
	for _, step := range builder.steps {
		switch step.Type {
		case "RUN":
			fmt.Fprintf(&output, "RUN %s\n", step.Args[0])
		case "COPY":
			fmt.Fprintf(&output, "COPY %s %s\n", step.Args[0], step.Args[1])
		case "WORKDIR", "USER":
			fmt.Fprintf(&output, "%s %s\n", step.Type, step.Args[0])
		case "ENV":
			for index := 0; index+1 < len(step.Args); index += 2 {
				fmt.Fprintf(&output, "ENV %s=%s\n", step.Args[index], step.Args[index+1])
			}
		}
	}
	return output.String()
}

func (builder *TemplateBuilder) request(steps []api.TemplateStep) api.TemplateBuildStartV2 {
	if steps == nil {
		steps = make([]api.TemplateStep, len(builder.steps))
		for index, step := range builder.steps {
			args := slices.Clone(step.Args)
			force := step.Force
			steps[index] = api.TemplateStep{Type: step.Type, Args: &args, Force: &force}
			if step.FilesHash != "" {
				steps[index].FilesHash = &step.FilesHash
			}
		}
	}
	force := builder.force
	request := api.TemplateBuildStartV2{Steps: &steps, Force: &force}
	if builder.baseImage != "" {
		request.FromImage = &builder.baseImage
	}
	if builder.baseTemplate != "" {
		request.FromTemplate = &builder.baseTemplate
	}
	if builder.registry != nil {
		request.FromImageRegistry = builder.registry
	}
	if builder.startCmd != "" {
		request.StartCmd = &builder.startCmd
	}
	if builder.readyCmd != "" {
		request.ReadyCmd = &builder.readyCmd
	}
	return request
}

// TemplateService manages AgentBox templates.
type TemplateService struct{ client *Client }
type TemplateBuildOptions struct {
	Tags               []string
	CPUCount, MemoryMB int
	SkipCache          bool
	PollInterval       time.Duration
	OnLog              func(BuildLogEntry)
}
type TemplateBuildRef struct {
	Name                string
	Tags                []string
	TemplateID, BuildID string
}
type TemplateListOptions struct {
	TeamID, NextToken string
	Limit             int
}

// TemplateInfoOptions paginates a template's build history.
type TemplateInfoOptions struct {
	NextToken string
	Limit     int
}
type TemplateLogOptions struct {
	Cursor                   string
	Timestamp                int64
	Limit                    int
	Direction, Level, Source string
}
type TemplateTag struct {
	Tag, BuildID string
	CreatedAt    time.Time
}

// BuildInBackground uploads COPY contexts and starts a build.
func (service *TemplateService) BuildInBackground(ctx context.Context, builder *TemplateBuilder, name string, options *TemplateBuildOptions) (*TemplateBuildRef, error) {
	if builder == nil || name == "" {
		return nil, &InvalidArgumentError{Message: "template and name are required"}
	}
	if builder.err != nil {
		return nil, builder.err
	}
	if options == nil {
		options = &TemplateBuildOptions{}
	}
	cpu, memory := int32(options.CPUCount), int32(options.MemoryMB)
	if cpu == 0 {
		cpu = 2
	}
	if memory == 0 {
		memory = 1024
	}
	builder.force = builder.force || options.SkipCache
	body := api.TemplateBuildRequestV3{Name: &name, CPUCount: &cpu, MemoryMB: &memory}
	if len(options.Tags) > 0 {
		tags := slices.Clone(options.Tags)
		body.Tags = &tags
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.PostV3TemplatesWithResponse(requestCtx, body)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON202 == nil {
		return nil, &BuildError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	value := response.JSON202
	steps, err := service.prepareCopySteps(ctx, builder, value.TemplateID)
	if err != nil {
		return nil, err
	}
	triggerCtx, triggerCancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer triggerCancel()
	trigger, err := service.client.api.PostV2TemplatesTemplateIDBuildsBuildIDWithResponse(triggerCtx, value.TemplateID, value.BuildID, builder.request(steps))
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if trigger.StatusCode() < 200 || trigger.StatusCode() >= 300 {
		return nil, &BuildError{APIError: APIError{StatusCode: trigger.StatusCode(), Message: string(trigger.Body)}}
	}
	return &TemplateBuildRef{Name: name, Tags: value.Tags, TemplateID: value.TemplateID, BuildID: value.BuildID}, nil
}

// Build starts a build and waits for a terminal state.
func (service *TemplateService) Build(ctx context.Context, builder *TemplateBuilder, name string, options *TemplateBuildOptions) (*TemplateBuildRef, error) {
	reference, err := service.BuildInBackground(ctx, builder, name, options)
	if err != nil {
		return nil, err
	}
	interval := 200 * time.Millisecond
	if options != nil && options.PollInterval > 0 {
		interval = options.PollInterval
	}
	offset := 0
	for {
		status, err := service.BuildStatus(ctx, reference.TemplateID, reference.BuildID, offset)
		if err != nil {
			return nil, err
		}
		offset += len(status.LogEntries)
		if options != nil && options.OnLog != nil {
			for _, entry := range status.LogEntries {
				options.OnLog(entry)
			}
		}
		switch status.Status {
		case BuildReady:
			return reference, nil
		case BuildFailed:
			message := "template build failed"
			if status.Reason != nil {
				message = status.Reason.Message
			}
			return nil, &BuildError{APIError: APIError{Message: message}}
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func (service *TemplateService) BuildStatus(ctx context.Context, templateID, buildID string, logsOffset int) (*TemplateBuildInfo, error) {
	offset := int32(logsOffset)
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetTemplatesTemplateIDBuildsBuildIDStatusWithResponse(requestCtx, templateID, buildID, &api.GetTemplatesTemplateIDBuildsBuildIDStatusParams{LogsOffset: &offset})
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, &BuildError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	result, err := convertModel[TemplateBuildInfo](*response.JSON200)
	return &result, err
}
func (service *TemplateService) Exists(ctx context.Context, alias string) (bool, error) {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetTemplatesAliasesAliasWithResponse(requestCtx, alias)
	if err != nil {
		return false, normalizeRequestError(err)
	}
	if response.StatusCode() == http.StatusNotFound {
		return false, nil
	}
	if response.StatusCode() == http.StatusForbidden {
		return true, nil
	}
	if response.JSON200 == nil {
		return false, &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	return true, nil
}
func (service *TemplateService) AssignTags(ctx context.Context, target string, tags []string) (string, error) {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.PostTemplatesTagsWithResponse(requestCtx, api.AssignTemplateTagsRequest{Target: target, Tags: tags})
	if err != nil {
		return "", normalizeRequestError(err)
	}
	if response.JSON201 == nil {
		return "", &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	return response.JSON201.BuildID.String(), nil
}
func (service *TemplateService) RemoveTags(ctx context.Context, name string, tags []string) error {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.DeleteTemplatesTagsWithResponse(requestCtx, api.DeleteTemplateTagsRequest{Name: name, Tags: tags})
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	return nil
}
func (service *TemplateService) Tags(ctx context.Context, templateID string) ([]TemplateTag, error) {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetTemplatesTemplateIDTagsWithResponse(requestCtx, templateID)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	result := make([]TemplateTag, 0, len(*response.JSON200))
	for _, tag := range *response.JSON200 {
		result = append(result, TemplateTag{Tag: tag.Tag, BuildID: tag.BuildID.String(), CreatedAt: tag.CreatedAt})
	}
	return result, nil
}

// List returns one template page.
func (service *TemplateService) List(ctx context.Context, options *TemplateListOptions) (Page[TemplateInfo], error) {
	params := &api.GetV2TemplatesParams{}
	if options != nil {
		if options.TeamID != "" {
			params.TeamID = &options.TeamID
		}
		if options.NextToken != "" {
			params.NextToken = &options.NextToken
		}
		if options.Limit > 0 {
			limit := int32(options.Limit)
			params.Limit = &limit
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetV2TemplatesWithResponse(requestCtx, params)
	if err != nil {
		return Page[TemplateInfo]{}, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return Page[TemplateInfo]{}, &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	items, err := convertModel[[]TemplateInfo](*response.JSON200)
	return Page[TemplateInfo]{Items: items, NextToken: response.HTTPResponse.Header.Get("X-Next-Token")}, err
}

// Info returns a template and its build history.
func (service *TemplateService) Info(ctx context.Context, templateID string, options *TemplateInfoOptions) (*TemplateWithBuilds, error) {
	params := &api.GetTemplatesTemplateIDParams{}
	if options != nil {
		if options.NextToken != "" {
			params.NextToken = &options.NextToken
		}
		if options.Limit > 0 {
			value := int32(options.Limit)
			params.Limit = &value
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetTemplatesTemplateIDWithResponse(requestCtx, templateID, params)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	result, err := convertModel[TemplateWithBuilds](*response.JSON200)
	return &result, err
}

// Delete removes a template.
func (service *TemplateService) Delete(ctx context.Context, templateID string) error {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.DeleteTemplatesTemplateIDWithResponse(requestCtx, templateID)
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	return nil
}

// SetPublic changes template visibility and returns its names.
func (service *TemplateService) SetPublic(ctx context.Context, templateID string, public bool) ([]string, error) {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.PatchV2TemplatesTemplateIDWithResponse(requestCtx, templateID, api.TemplateUpdateRequest{Public: &public})
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, &TemplateError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	return response.JSON200.Names, nil
}

// BuildLogs returns one page of structured build logs.
func (service *TemplateService) BuildLogs(ctx context.Context, templateID, buildID string, options *TemplateLogOptions) (Page[BuildLogEntry], error) {
	params := &api.GetTemplatesTemplateIDBuildsBuildIDLogsParams{}
	if options != nil {
		if options.Cursor != "" {
			params.PageCursor = &options.Cursor
		}
		if options.Timestamp != 0 {
			params.Cursor = &options.Timestamp
		}
		if options.Limit > 0 {
			limit := int32(options.Limit)
			params.Limit = &limit
		}
		if options.Direction != "" {
			value := api.LogsDirection(options.Direction)
			params.Direction = &value
		}
		if options.Level != "" {
			value := api.LogLevel(options.Level)
			params.Level = &value
		}
		if options.Source != "" {
			value := api.LogsSource(options.Source)
			params.Source = &value
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetTemplatesTemplateIDBuildsBuildIDLogsWithResponse(requestCtx, templateID, buildID, params)
	if err != nil {
		return Page[BuildLogEntry]{}, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return Page[BuildLogEntry]{}, &BuildError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
	}
	items, err := convertModel[[]BuildLogEntry](response.JSON200.Logs)
	if err != nil {
		return Page[BuildLogEntry]{}, err
	}
	page := Page[BuildLogEntry]{Items: items}
	if response.JSON200.NextCursor != nil {
		page.NextToken = *response.JSON200.NextCursor
	}
	return page, nil
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }

func splitDockerWords(value string) ([]string, error) {
	words := make([]string, 0)
	var current strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}
	for _, char := range value {
		if escaped {
			current.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' && quote != '\'' {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
			continue
		}
		switch char {
		case '\'', '"':
			quote = char
		case ' ', '\t', '\r', '\n':
			flush()
		default:
			current.WriteRune(char)
		}
	}
	if escaped || quote != 0 {
		return nil, &InvalidArgumentError{Message: "unterminated Dockerfile escape or quote"}
	}
	flush()
	return words, nil
}

func joinDockerfileLines(content string) []string {
	raw := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	result := make([]string, 0, len(raw))
	current := ""
	for _, line := range raw {
		line = strings.TrimSpace(line)
		current += line
		if strings.HasSuffix(current, "\\") {
			current = strings.TrimSuffix(current, "\\") + " "
			continue
		}
		if current != "" {
			result = append(result, current)
		}
		current = ""
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
