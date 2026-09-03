package agentbox

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
)

func (service *TemplateService) prepareCopySteps(ctx context.Context, builder *TemplateBuilder, templateID string) ([]api.TemplateStep, error) {
	steps := make([]api.TemplateStep, len(builder.steps))
	for index, step := range builder.steps {
		args := slices.Clone(step.Args)
		force := step.Force
		steps[index] = api.TemplateStep{Type: step.Type, Args: &args, Force: &force}
		if step.Type != "COPY" {
			continue
		}
		archive, hash, err := archiveCopy(builder.contextPath, step.Args[0], builder.ignore, step.ResolveSymlinks, step.Gzip)
		if err != nil {
			return nil, &FileUploadError{APIError: APIError{Message: "prepare " + step.Args[0], Cause: err}}
		}
		steps[index].FilesHash = &hash
		requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
		response, err := service.client.api.GetTemplatesTemplateIDFilesHashWithResponse(requestCtx, templateID, hash)
		cancel()
		if err != nil {
			return nil, normalizeRequestError(err)
		}
		if response.JSON201 == nil {
			return nil, &FileUploadError{APIError: APIError{StatusCode: response.StatusCode(), Message: string(response.Body)}}
		}
		if response.JSON201.Present && !step.ForceUpload {
			continue
		}
		if response.JSON201.URL == nil || *response.JSON201.URL == "" {
			return nil, &FileUploadError{APIError: APIError{Message: "upload URL is missing"}}
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPut, *response.JSON201.URL, bytes.NewReader(archive))
		if err != nil {
			return nil, err
		}
		request.ContentLength = int64(len(archive))
		upload, err := service.client.httpClient.Do(request)
		if err != nil {
			return nil, &FileUploadError{APIError: APIError{Message: "upload failed", Cause: err}}
		}
		_, _ = io.Copy(io.Discard, upload.Body)
		upload.Body.Close()
		if upload.StatusCode < 200 || upload.StatusCode >= 300 {
			return nil, &FileUploadError{APIError: APIError{StatusCode: upload.StatusCode, Message: upload.Status}}
		}
	}
	return steps, nil
}

func archiveCopy(contextPath, source string, ignore []string, resolveSymlinks, useGzip bool) ([]byte, string, error) {
	root, err := filepath.Abs(contextPath)
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(root, filepath.Clean(source))
	if !strings.HasPrefix(path, root+string(filepath.Separator)) && path != root {
		return nil, "", fmt.Errorf("source escapes context: %s", source)
	}
	var buffer bytes.Buffer
	var output io.Writer = &buffer
	var zipper *gzip.Writer
	if useGzip {
		zipper, _ = gzip.NewWriterLevel(&buffer, gzip.BestSpeed)
		zipper.Header.ModTime = time.Unix(0, 0)
		zipper.Header.OS = 255
		output = zipper
	}
	archive := tar.NewWriter(output)
	paths := make([]string, 0)
	err = filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		if ignored(relative, ignore) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		paths = append(paths, current)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	slices.Sort(paths)
	for _, current := range paths {
		info, err := os.Lstat(current)
		if err != nil {
			return nil, "", err
		}
		link := ""
		if info.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(current)
			if err != nil {
				return nil, "", err
			}
			if resolveSymlinks {
				info, err = os.Stat(current)
				if err != nil {
					return nil, "", err
				}
				link = ""
			}
		}
		header, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return nil, "", err
		}
		header.Name, _ = filepath.Rel(root, current)
		header.Name = filepath.ToSlash(header.Name)
		header.ModTime = time.Unix(0, 0)
		header.AccessTime = time.Time{}
		header.ChangeTime = time.Time{}
		header.Uid, header.Gid, header.Uname, header.Gname = 0, 0, "", ""
		if err := archive.WriteHeader(header); err != nil {
			return nil, "", err
		}
		if info.Mode().IsRegular() {
			file, err := os.Open(current)
			if err != nil {
				return nil, "", err
			}
			_, copyErr := io.Copy(archive, file)
			closeErr := file.Close()
			if copyErr != nil {
				return nil, "", copyErr
			}
			if closeErr != nil {
				return nil, "", closeErr
			}
		}
	}
	if err := archive.Close(); err != nil {
		return nil, "", err
	}
	if zipper != nil {
		if err := zipper.Close(); err != nil {
			return nil, "", err
		}
	}
	digest := sha256.Sum256(buffer.Bytes())
	return buffer.Bytes(), fmt.Sprintf("%x", digest[:]), nil
}

func ignored(path string, patterns []string) bool {
	path = filepath.ToSlash(path)
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}
