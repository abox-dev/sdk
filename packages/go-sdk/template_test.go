package agentbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestTemplateRunUserArguments(t *testing.T) {
	tests := []struct {
		name  string
		apply func(*TemplateBuilder) *TemplateBuilder
		args  []string
	}{
		{"run-as-current", func(b *TemplateBuilder) *TemplateBuilder { return b.RunAs("", "echo one", "echo two") }, []string{"echo one && echo two"}},
		{"run-as-explicit", func(b *TemplateBuilder) *TemplateBuilder { return b.RunAs("root", "id -un") }, []string{"id -un", "root"}},
		{"npm-default", func(b *TemplateBuilder) *TemplateBuilder { return b.NPMInstall(PackageInstallOptions{}, "lodash") }, []string{"npm install lodash"}},
		{"npm-dev", func(b *TemplateBuilder) *TemplateBuilder {
			return b.NPMInstall(PackageInstallOptions{Dev: true}, "typescript")
		}, []string{"npm install --save-dev typescript"}},
		{"npm-global", func(b *TemplateBuilder) *TemplateBuilder {
			return b.NPMInstall(PackageInstallOptions{Global: true}, "tsx")
		}, []string{"npm install -g tsx", "root"}},
		{"bun-default", func(b *TemplateBuilder) *TemplateBuilder { return b.BunInstall(PackageInstallOptions{}, "lodash") }, []string{"bun install lodash"}},
		{"bun-dev", func(b *TemplateBuilder) *TemplateBuilder {
			return b.BunInstall(PackageInstallOptions{Dev: true}, "typescript")
		}, []string{"bun install --dev typescript"}},
		{"bun-global", func(b *TemplateBuilder) *TemplateBuilder {
			return b.BunInstall(PackageInstallOptions{Global: true}, "tsx")
		}, []string{"bun install -g tsx", "root"}},
		{"git-default", func(b *TemplateBuilder) *TemplateBuilder { return b.GitClone("https://example.test/repo.git", nil) }, []string{"git clone 'https://example.test/repo.git'"}},
		{"git-options", func(b *TemplateBuilder) *TemplateBuilder {
			return b.GitClone("https://example.test/repo.git", &GitCloneOptions{Path: "repo", Depth: 1})
		}, []string{"git clone 'https://example.test/repo.git' --depth 1 'repo'"}},
		{"git-explicit", func(b *TemplateBuilder) *TemplateBuilder {
			return b.GitClone("https://example.test/repo.git", &GitCloneOptions{User: "root"})
		}, []string{"git clone 'https://example.test/repo.git'", "root"}},
	}
	for _, currentUser := range []string{"", "app"} {
		for _, test := range tests {
			t.Run(currentUser+"/"+test.name, func(t *testing.T) {
				builder := NewTemplate("").FromBase()
				if currentUser != "" {
					builder.User(currentUser)
				}
				data, err := test.apply(builder).JSON()
				if err != nil {
					t.Fatal(err)
				}
				var request struct {
					Steps []struct {
						Type string   `json:"type"`
						Args []string `json:"args"`
					} `json:"steps"`
				}
				if err := json.Unmarshal(data, &request); err != nil {
					t.Fatal(err)
				}
				if len(request.Steps) == 0 {
					t.Fatal("missing build steps")
				}
				step := request.Steps[len(request.Steps)-1]
				if step.Type != "RUN" || !slices.Equal(step.Args, test.args) {
					t.Fatalf("RUN arguments = %#v, want %#v (type %q)", step.Args, test.args, step.Type)
				}
				if currentUser != "" && (request.Steps[0].Type != "USER" || !slices.Equal(request.Steps[0].Args, []string{currentUser})) {
					t.Fatal("configured build user was not preserved")
				}
			})
		}
	}
}

func TestTemplateBuilder(t *testing.T) {
	contextPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(contextPath, "hello.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	gzipDisabled := false
	builder := NewTemplate(contextPath, "*.ignore").FromPython("3.13").Copy("hello.txt", "/app/", &CopyOptions{User: "user", Mode: 0o640, ForceUpload: true, Gzip: &gzipDisabled}).Workdir("/app").User("user").Env(map[string]string{"B": "2", "A": "1"}).Run("echo one", "echo two").MakeDir("/tmp/x").Rename("a", "b").Symlink("b", "c").Remove("c").PipInstall("requests").NPMInstall(PackageInstallOptions{Dev: true}, "typescript").BunInstall(PackageInstallOptions{Global: true}, "tsx").AptInstall(AptInstallOptions{NoInstallRecommends: true}, "git").GitClone("https://example.test/repo.git", &GitCloneOptions{Path: "repo", Branch: "main", Depth: 1}).Start("python app.py", WaitForPort(8000)).Ready(WaitForURL("http://localhost", 200)).SkipCache()
	data, err := builder.JSON()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"fromImage":"python:3.13"`) || strings.Contains(string(data), `"filesHash"`) {
		t.Fatalf("unexpected JSON: %s", data)
	}
	dockerfile := builder.Dockerfile()
	for _, expected := range []string{"FROM python:3.13", "COPY hello.txt /app/", "WORKDIR /app", "USER user", "RUN echo one && echo two", "ENV A=1"} {
		if !strings.Contains(dockerfile, expected) {
			t.Fatalf("missing %q in %s", expected, dockerfile)
		}
	}
	if WaitForFile("/tmp/ready") == "" || WaitForProcess("nginx") == "" {
		t.Fatal("ready helpers empty")
	}
	if WaitForTimeout(0) != "sleep 1" {
		t.Fatal("timeout helper")
	}
	if _, err := NewTemplate(contextPath).Copy("../escape", "/", nil).JSON(); err == nil {
		t.Fatal("expected escaping copy error")
	}
	if _, err := NewTemplate(contextPath).FromImage("").JSON(); err == nil {
		t.Fatal("expected image error")
	}
	parsed := NewTemplate(contextPath).FromDockerfile("FROM alpine:3.24\nENV A=1 B 2\nWORKDIR /app\nUSER app\nRUN echo ok\nCOPY hello.txt /app/")
	if _, err := parsed.JSON(); err != nil || !strings.Contains(parsed.Dockerfile(), "FROM alpine:3.24") {
		t.Fatalf("dockerfile parse: %v", err)
	}
	quoted := NewTemplate(contextPath).FromDockerfile("FROM alpine:3.24\nENV A=\"hello world\"\nCOPY \"hello world.txt\" /app/")
	if _, err := quoted.JSON(); err != nil || quoted.steps[2].Args[1] != "hello world" || quoted.steps[3].Args[0] != "hello world.txt" {
		t.Fatalf("quoted Dockerfile parse: %#v %v", quoted.steps, err)
	}
	if _, err := NewTemplate(contextPath).FromDockerfile("FROM alpine\nENV A='unterminated").JSON(); err == nil {
		t.Fatal("expected unterminated quote error")
	}
	if _, err := NewTemplate(contextPath).FromDockerfile("EXPOSE 80").JSON(); err == nil {
		t.Fatal("expected missing FROM")
	}
	if _, err := NewTemplate(contextPath).FromDockerfile("FROM a\nFROM b").JSON(); err == nil {
		t.Fatal("expected multi-stage error")
	}
	if _, err := NewTemplate(contextPath).FromDockerfile("FROM alpine\nARG FLAG\nCOPY --chown=user hello.txt other.txt /app/\nEXPOSE 80\nVOLUME /data\nCMD [\"echo\",\"ok\"]\nLABEL x=y").JSON(); err != nil {
		t.Fatal(err)
	}
	NewTemplate("").FromBase().FromDebian("").FromUbuntu("").FromFedora("").FromAlpine("").FromArch("").FromNode("").FromBun("").FromTemplate("base").FromRegistry("image", "u", "p").FromAWSRegistry("image", "a", "s", "r").FromGCPRegistry("image", `{}`)
	NewTemplate(contextPath).FromPython("").SkipCache().FromTemplate("base").PipInstall().NPMInstall(PackageInstallOptions{Global: true}).BunInstall(PackageInstallOptions{Dev: true})
	if err := os.WriteFile(filepath.Join(contextPath, "Dockerfile"), []byte("FROM alpine:3.24\\\n\nRUN echo ok\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewTemplate(contextPath).FromDockerfile("Dockerfile").JSON(); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateArchiveDeterministic(t *testing.T) {
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "dir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "dir", "a.txt"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "skip.ignore"), []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, hash1, err := archiveCopy(directory, ".", []string{"*.ignore"}, false, true)
	if err != nil {
		t.Fatal(err)
	}
	second, hash2, err := archiveCopy(directory, ".", []string{"*.ignore"}, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if hash1 != hash2 || string(first) != string(second) {
		t.Fatal("archive is not deterministic")
	}
	if _, _, err := archiveCopy(directory, "../escape", nil, false, false); err == nil {
		t.Fatal("expected escape error")
	}
	if !ignored("anything", []string{"", "# comment", "any*"}) {
		t.Fatal("ignore pattern not applied")
	}
	if err := os.Symlink("dir/a.txt", filepath.Join(directory, "link")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := archiveCopy(directory, "link", nil, true, false); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateService(t *testing.T) {
	contextPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(contextPath, "hello.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	statusCalls := 0
	uploaded := false
	cachePresent := false
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch {
		case request.URL.Path == "/v3/templates":
			writer.WriteHeader(http.StatusAccepted)
			fmt.Fprint(writer, `{"templateID":"tpl","buildID":"build","names":["name"],"public":false,"tags":["v1"]}`)
		case strings.HasPrefix(request.URL.Path, "/templates/tpl/files/"):
			writer.WriteHeader(http.StatusCreated)
			fmt.Fprintf(writer, `{"present":%t,"url":%q}`, cachePresent, server.URL+"/upload")
		case request.URL.Path == "/upload":
			uploaded = request.ContentLength > 0
			writer.WriteHeader(http.StatusOK)
		case request.URL.Path == "/v2/templates/tpl/builds/build":
			writer.WriteHeader(http.StatusAccepted)
		case request.URL.Path == "/templates/tpl/builds/build/status":
			statusCalls++
			status := "building"
			if statusCalls > 1 {
				status = "ready"
			}
			fmt.Fprintf(writer, `{"templateID":"tpl","buildID":"build","status":%q,"logs":[],"logEntries":[]}`, status)
		case request.URL.Path == "/templates/aliases/missing":
			writer.WriteHeader(http.StatusNotFound)
			fmt.Fprint(writer, `{"message":"missing"}`)
		case request.URL.Path == "/templates/aliases/forbidden":
			writer.WriteHeader(http.StatusForbidden)
			fmt.Fprint(writer, `{"message":"forbidden"}`)
		case request.URL.Path == "/templates/aliases/name":
			fmt.Fprint(writer, `{"templateID":"tpl","public":false}`)
		case request.URL.Path == "/v2/templates":
			writer.Header().Set("X-Next-Token", "next")
			fmt.Fprint(writer, `[]`)
		case request.URL.Path == "/templates/tpl" && request.Method == http.MethodGet:
			fmt.Fprint(writer, `{"templateID":"tpl","names":["name"],"builds":[],"public":false,"spawnCount":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`)
		case request.URL.Path == "/templates/tpl" && request.Method == http.MethodDelete:
			writer.WriteHeader(http.StatusNoContent)
		case request.URL.Path == "/v2/templates/tpl":
			fmt.Fprint(writer, `{"names":["name"]}`)
		case request.URL.Path == "/templates/tags" && request.Method == http.MethodPost:
			writer.WriteHeader(http.StatusCreated)
			fmt.Fprint(writer, `{"buildID":"00000000-0000-0000-0000-000000000001","tags":["v1"]}`)
		case request.URL.Path == "/templates/tags" && request.Method == http.MethodDelete:
			writer.WriteHeader(http.StatusNoContent)
		case request.URL.Path == "/templates/tpl/tags":
			fmt.Fprint(writer, `[{"tag":"v1","buildID":"00000000-0000-0000-0000-000000000001","createdAt":"2024-01-01T00:00:00Z"}]`)
		case request.URL.Path == "/templates/tpl/builds/build/logs":
			fmt.Fprint(writer, `{"logs":[],"nextCursor":"cursor"}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, err := NewClient(WithAPIURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	builder := NewTemplate(contextPath).FromAlpine("").Copy("hello.txt", "/app/", nil)
	ref, err := client.Templates.Build(context.Background(), builder, "name", &TemplateBuildOptions{Tags: []string{"v1"}, PollInterval: time.Millisecond})
	if err != nil || ref.TemplateID != "tpl" || !uploaded {
		t.Fatalf("build: %#v %v uploaded=%v", ref, err, uploaded)
	}
	if exists, err := client.Templates.Exists(context.Background(), "name"); err != nil || !exists {
		t.Fatalf("exists: %v %v", exists, err)
	}
	if exists, err := client.Templates.Exists(context.Background(), "missing"); err != nil || exists {
		t.Fatalf("missing exists: %v %v", exists, err)
	}
	if exists, err := client.Templates.Exists(context.Background(), "forbidden"); err != nil || !exists {
		t.Fatalf("forbidden exists: %v %v", exists, err)
	}
	cachePresent = true
	if _, err := client.Templates.BuildInBackground(context.Background(), NewTemplate(contextPath).Copy("hello.txt", "/cached/", nil), "cached", nil); err != nil {
		t.Fatal(err)
	}
	if page, err := client.Templates.List(context.Background(), &TemplateListOptions{Limit: 10}); err != nil || page.NextToken != "next" {
		t.Fatalf("list: %#v %v", page, err)
	}
	if _, err := client.Templates.Info(context.Background(), "tpl", &TemplateInfoOptions{Limit: 10}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Templates.SetPublic(context.Background(), "tpl", true); err != nil {
		t.Fatal(err)
	}
	if id, err := client.Templates.AssignTags(context.Background(), "name", []string{"v1"}); err != nil || id == "" {
		t.Fatalf("tags: %s %v", id, err)
	}
	if tags, err := client.Templates.Tags(context.Background(), "tpl"); err != nil || len(tags) != 1 {
		t.Fatalf("tags list: %v", err)
	}
	if err := client.Templates.RemoveTags(context.Background(), "name", []string{"v1"}); err != nil {
		t.Fatal(err)
	}
	if logs, err := client.Templates.BuildLogs(context.Background(), "tpl", "build", &TemplateLogOptions{Limit: 10, Direction: "forward", Level: "info", Source: "persistent"}); err != nil || logs.NextToken != "cursor" {
		t.Fatalf("logs: %#v %v", logs, err)
	}
	if err := client.Templates.Delete(context.Background(), "tpl"); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateServiceFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, `{"message":"failure"}`)
	}))
	defer server.Close()
	client, _ := NewClient(WithAPIURL(server.URL))
	ctx := context.Background()
	builder := NewTemplate(t.TempDir())
	calls := []func() error{
		func() error { _, err := client.Templates.BuildInBackground(ctx, builder, "name", nil); return err },
		func() error { _, err := client.Templates.BuildStatus(ctx, "tpl", "build", 0); return err },
		func() error { _, err := client.Templates.Exists(ctx, "name"); return err },
		func() error { _, err := client.Templates.List(ctx, nil); return err },
		func() error { _, err := client.Templates.Info(ctx, "tpl", nil); return err },
		func() error { return client.Templates.Delete(ctx, "tpl") },
		func() error { _, err := client.Templates.SetPublic(ctx, "tpl", true); return err },
		func() error { _, err := client.Templates.AssignTags(ctx, "name", nil); return err },
		func() error { return client.Templates.RemoveTags(ctx, "name", nil) },
		func() error { _, err := client.Templates.Tags(ctx, "tpl"); return err },
		func() error { _, err := client.Templates.BuildLogs(ctx, "tpl", "build", nil); return err },
	}
	for index, call := range calls {
		if err := call(); err == nil {
			t.Fatalf("call %d unexpectedly succeeded", index)
		}
	}
	if _, err := client.Templates.BuildInBackground(ctx, nil, "", nil); err == nil {
		t.Fatal("expected builder validation")
	}
}

func TestTemplateBuildTerminalFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch {
		case request.URL.Path == "/v3/templates":
			writer.WriteHeader(http.StatusAccepted)
			fmt.Fprint(writer, `{"templateID":"tpl","buildID":"build","names":[],"public":false,"tags":[]}`)
		case request.URL.Path == "/v2/templates/tpl/builds/build":
			writer.WriteHeader(http.StatusAccepted)
		case request.URL.Path == "/templates/tpl/builds/build/status":
			fmt.Fprint(writer, `{"templateID":"tpl","buildID":"build","status":"error","logs":[],"logEntries":[],"reason":{"message":"broken"}}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(WithAPIURL(server.URL))
	if _, err := client.Templates.Build(context.Background(), NewTemplate(t.TempDir()), "name", nil); err == nil {
		t.Fatal("expected build failure")
	} else {
		var build *BuildError
		if !errors.As(err, &build) {
			t.Fatalf("error type %T", err)
		}
	}
}
