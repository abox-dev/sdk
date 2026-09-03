package agentbox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	pathpkg "path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	filesystem "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/filesystem"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/filesystem/filesystemconnect"
)

// FileType identifies a filesystem entry kind.
type FileType string

const (
	FileTypeFile      FileType = "file"
	FileTypeDirectory FileType = "dir"
	FileTypeSymlink   FileType = "symlink"
)

// EntryInfo describes a sandbox filesystem entry.
type EntryInfo struct {
	Name          string
	Type          FileType
	Path          string
	Size          int64
	Mode          uint32
	Permissions   string
	Owner         string
	Group         string
	ModifiedAt    time.Time
	SymlinkTarget string
	Metadata      map[string]string
}

// WriteFile describes one batch upload.
type WriteFile struct {
	Path     string
	Data     io.Reader
	Metadata map[string]string
}

// WriteFileOptions configures file ownership, metadata, and upload timeout.
type WriteFileOptions struct {
	User     string
	Metadata map[string]string
	// RequestTimeout limits the complete streaming upload. Zero leaves the
	// upload bounded only by ctx and a custom HTTP client timeout.
	RequestTimeout time.Duration
}

// WatchOptions configures recursive and enriched filesystem events.
type WatchOptions struct{ Recursive, IncludeEntry, AllowNetworkMounts bool }

// FileEvent describes a filesystem change.
type FileEvent struct {
	Name  string
	Type  string
	Entry *EntryInfo
}

// WatchHandle owns a directory watch stream.
type WatchHandle struct {
	Events <-chan FileEvent
	done   <-chan struct{}
	cancel context.CancelFunc
	mu     sync.Mutex
	err    error
}

// Close stops the watcher.
func (handle *WatchHandle) Close() error {
	handle.cancel()
	<-handle.done
	handle.mu.Lock()
	defer handle.mu.Unlock()
	return handle.err
}

// FileService reads and mutates sandbox files.
type FileService struct {
	sandbox *Sandbox
	client  filesystemconnect.FilesystemClient
}

func newFileService(sandbox *Sandbox) *FileService {
	return &FileService{sandbox: sandbox, client: filesystemconnect.NewFilesystemClient(sandbox.client.envdClient, sandbox.envdURL(envdPort, false), connect.WithCodec(tolerantJSONCodec{}), connect.WithAcceptCompression("gzip", nil, nil))}
}

// Read opens a streaming file response. The caller must close it.
func (service *FileService) Read(ctx context.Context, path, user string) (io.ReadCloser, error) {
	if path == "" {
		return nil, &InvalidArgumentError{Message: "file path cannot be empty"}
	}
	user = service.sandbox.resolveUser(user)
	endpoint, _ := url.Parse(service.sandbox.envdURL(envdPort, false) + "/files")
	query := endpoint.Query()
	query.Set("path", path)
	if user != "" {
		query.Set("username", user)
	}
	endpoint.RawQuery = query.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header = service.sandbox.envdHeaders(envdPort)
	response, err := service.sandbox.client.httpClient.Do(request)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		defer response.Body.Close()
		return nil, decodeFileHTTPError(response)
	}
	return response.Body, nil
}

// ReadBytes reads a complete file.
func (service *FileService) ReadBytes(ctx context.Context, path, user string) ([]byte, error) {
	reader, err := service.Read(ctx, path, user)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

// ReadText reads a UTF-8 file as a string.
func (service *FileService) ReadText(ctx context.Context, path, user string) (string, error) {
	data, err := service.ReadBytes(ctx, path, user)
	return string(data), err
}

// ReadTo streams a file into writer.
func (service *FileService) ReadTo(ctx context.Context, path, user string, writer io.Writer) (int64, error) {
	reader, err := service.Read(ctx, path, user)
	if err != nil {
		return 0, err
	}
	defer reader.Close()
	return io.Copy(writer, reader)
}

// Write uploads a file from reader.
func (service *FileService) Write(ctx context.Context, path string, reader io.Reader, options *WriteFileOptions) (*EntryInfo, error) {
	if path == "" || reader == nil {
		return nil, &InvalidArgumentError{Message: "file path and reader are required"}
	}
	if options == nil {
		options = &WriteFileOptions{}
	}
	if options.RequestTimeout < 0 {
		return nil, &InvalidArgumentError{Message: "file upload timeout cannot be negative"}
	}
	if len(options.Metadata) > 0 && !envdAtLeast(service.sandbox.EnvdVersion, 0, 6, 2) {
		return nil, &TemplateError{APIError: APIError{Message: "file metadata requires envd 0.6.2 or later"}}
	}
	user := service.sandbox.resolveUser(options.User)
	endpoint, _ := url.Parse(service.sandbox.envdURL(envdPort, false) + "/files")
	query := endpoint.Query()
	query.Set("path", path)
	if user != "" {
		query.Set("username", user)
	}
	endpoint.RawQuery = query.Encode()
	headers := service.sandbox.envdHeaders(envdPort)
	if err := addMetadataHeaders(headers, options.Metadata); err != nil {
		return nil, err
	}
	requestCtx, cancel := withRequestTimeout(ctx, options.RequestTimeout)
	defer cancel()
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodPost, endpoint.String(), pipeReader)
	if err != nil {
		pipeReader.Close()
		pipeWriter.Close()
		return nil, err
	}
	go writeMultipartFile(pipeWriter, multipartWriter, pathpkg.Base(path), reader)
	request.Header = headers
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	response, err := service.sandbox.client.httpClient.Do(request)
	if err != nil {
		pipeReader.CloseWithError(err)
		return nil, normalizeRequestError(err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, decodeFileHTTPError(response)
	}
	var entries []struct {
		Name     string            `json:"name"`
		Type     string            `json:"type"`
		Path     string            `json:"path"`
		Metadata map[string]string `json:"metadata"`
	}
	if err := json.NewDecoder(response.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("agentbox: decode upload response: %w", err)
	}
	if len(entries) == 0 {
		return nil, &FileUploadError{APIError: APIError{Message: "upload response did not contain a file"}}
	}
	return &EntryInfo{Name: entries[0].Name, Type: FileType(entries[0].Type), Path: entries[0].Path, Metadata: entries[0].Metadata}, nil
}

// WriteText uploads a string.
func (service *FileService) WriteText(ctx context.Context, path, text string, options *WriteFileOptions) (*EntryInfo, error) {
	return service.Write(ctx, path, strings.NewReader(text), options)
}

// WriteBytes uploads bytes.
func (service *FileService) WriteBytes(ctx context.Context, path string, data []byte, options *WriteFileOptions) (*EntryInfo, error) {
	return service.Write(ctx, path, bytes.NewReader(data), options)
}

// WriteBatch writes files in order and stops at the first failure.
func (service *FileService) WriteBatch(ctx context.Context, files []WriteFile, user string) ([]EntryInfo, error) {
	result := make([]EntryInfo, 0, len(files))
	for _, file := range files {
		entry, err := service.Write(ctx, file.Path, file.Data, &WriteFileOptions{User: user, Metadata: file.Metadata})
		if err != nil {
			return result, err
		}
		result = append(result, *entry)
	}
	return result, nil
}

// Stat returns information about a path.
func (service *FileService) Stat(ctx context.Context, path string) (*EntryInfo, error) {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&filesystem.StatRequest{Path: path})
	service.addHeaders(request.Header())
	response, err := service.client.Stat(requestCtx, request)
	if err != nil {
		return nil, fileConnectError(err)
	}
	return mapEntry(response.Msg.GetEntry()), nil
}

// Exists reports whether path exists.
func (service *FileService) Exists(ctx context.Context, path string) (bool, error) {
	_, err := service.Stat(ctx, path)
	var notFound *FileNotFoundError
	if errors.As(err, &notFound) {
		return false, nil
	}
	return err == nil, err
}

// List lists path recursively up to depth.
func (service *FileService) List(ctx context.Context, path string, depth uint32) ([]EntryInfo, error) {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&filesystem.ListDirRequest{Path: path, Depth: depth})
	service.addHeaders(request.Header())
	response, err := service.client.ListDir(requestCtx, request)
	if err != nil {
		return nil, fileConnectError(err)
	}
	result := make([]EntryInfo, 0, len(response.Msg.GetEntries()))
	for _, entry := range response.Msg.GetEntries() {
		result = append(result, *mapEntry(entry))
	}
	return result, nil
}

// MakeDir creates a directory.
func (service *FileService) MakeDir(ctx context.Context, path string) (*EntryInfo, error) {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&filesystem.MakeDirRequest{Path: path})
	service.addHeaders(request.Header())
	response, err := service.client.MakeDir(requestCtx, request)
	if err != nil {
		return nil, fileConnectError(err)
	}
	return mapEntry(response.Msg.GetEntry()), nil
}

// Rename moves a filesystem entry.
func (service *FileService) Rename(ctx context.Context, source, destination string) (*EntryInfo, error) {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&filesystem.MoveRequest{Source: source, Destination: destination})
	service.addHeaders(request.Header())
	response, err := service.client.Move(requestCtx, request)
	if err != nil {
		return nil, fileConnectError(err)
	}
	return mapEntry(response.Msg.GetEntry()), nil
}

// Remove recursively removes a filesystem entry.
func (service *FileService) Remove(ctx context.Context, path string) error {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&filesystem.RemoveRequest{Path: path})
	service.addHeaders(request.Header())
	_, err := service.client.Remove(requestCtx, request)
	return fileConnectError(err)
}

func writeMultipartFile(pipe *io.PipeWriter, writer *multipart.Writer, name string, reader io.Reader) {
	part, err := writer.CreateFormFile("file", name)
	if err == nil {
		_, err = io.Copy(part, reader)
	}
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		_ = pipe.CloseWithError(err)
		return
	}
	_ = pipe.Close()
}

// Watch watches a directory until context cancellation or Close.
func (service *FileService) Watch(ctx context.Context, path string, options *WatchOptions) (*WatchHandle, error) {
	if options == nil {
		options = &WatchOptions{}
	}
	if options.Recursive && !envdAtLeast(service.sandbox.EnvdVersion, 0, 1, 4) {
		return nil, &SandboxError{APIError: APIError{Message: "recursive watch requires envd 0.1.4 or later"}}
	}
	if options.IncludeEntry && !envdAtLeast(service.sandbox.EnvdVersion, 0, 6, 3) {
		return nil, &SandboxError{APIError: APIError{Message: "watch entry details require envd 0.6.3 or later"}}
	}
	if options.AllowNetworkMounts && !envdAtLeast(service.sandbox.EnvdVersion, 0, 6, 4) {
		return nil, &SandboxError{APIError: APIError{Message: "watching network mounts requires envd 0.6.4 or later"}}
	}
	watchCtx, cancel := context.WithCancel(ctx)
	request := connect.NewRequest(&filesystem.WatchDirRequest{Path: path, Recursive: options.Recursive, IncludeEntry: options.IncludeEntry, AllowNetworkMounts: options.AllowNetworkMounts})
	service.addHeaders(request.Header())
	stream, err := service.client.WatchDir(watchCtx, request)
	if err != nil {
		cancel()
		return nil, fileConnectError(err)
	}
	events := make(chan FileEvent, 16)
	done := make(chan struct{})
	handle := &WatchHandle{Events: events, done: done, cancel: cancel}
	go func() {
		defer close(done)
		defer close(events)
		defer func() { _ = stream.Close() }()
		for stream.Receive() {
			event := stream.Msg().GetFilesystem()
			if event != nil {
				select {
				case events <- FileEvent{Name: event.GetName(), Type: mapEventType(event.GetType()), Entry: mapEntry(event.GetEntry())}:
				case <-watchCtx.Done():
					return
				}
			}
		}
		if err := stream.Err(); err != nil && !errors.Is(err, context.Canceled) {
			handle.mu.Lock()
			handle.err = fileConnectError(err)
			handle.mu.Unlock()
		}
	}()
	return handle, nil
}

// SignedReadURL creates a directly usable download URL.
func (service *FileService) SignedReadURL(path, user string, expiration time.Time) (string, error) {
	return service.signedURL(path, user, "read", expiration)
}

// SignedWriteURL creates a directly usable upload URL.
func (service *FileService) SignedWriteURL(path, user string, expiration time.Time) (string, error) {
	return service.signedURL(path, user, "write", expiration)
}
func (service *FileService) signedURL(path, user, operation string, expiration time.Time) (string, error) {
	user = service.sandbox.resolveUser(user)
	signature, unix, err := fileSignature(path, operation, user, service.sandbox.envdAccessToken, expiration)
	if err != nil {
		return "", err
	}
	endpoint, _ := url.Parse(service.sandbox.envdURL(envdPort, true) + "/files")
	query := endpoint.Query()
	query.Set("path", path)
	query.Set("signature", signature)
	if user != "" {
		query.Set("username", user)
	}
	if unix != 0 {
		query.Set("signature_expiration", strconv.FormatInt(unix, 10))
	}
	endpoint.RawQuery = query.Encode()
	return endpoint.String(), nil
}

func (service *FileService) addHeaders(header http.Header) {
	for key, values := range service.sandbox.envdHeaders(envdPort) {
		header[key] = slices.Clone(values)
	}
	header.Set("Keepalive-Ping-Interval", "50")
}
func mapEntry(entry *filesystem.EntryInfo) *EntryInfo {
	if entry == nil {
		return nil
	}
	result := &EntryInfo{Name: entry.GetName(), Path: entry.GetPath(), Size: entry.GetSize(), Mode: entry.GetMode(), Permissions: entry.GetPermissions(), Owner: entry.GetOwner(), Group: entry.GetGroup(), SymlinkTarget: entry.GetSymlinkTarget(), Metadata: maps.Clone(entry.GetMetadata())}
	if entry.GetModifiedTime() != nil {
		result.ModifiedAt = entry.GetModifiedTime().AsTime()
	}
	switch entry.GetType() {
	case filesystem.FileType_FILE_TYPE_FILE:
		result.Type = FileTypeFile
	case filesystem.FileType_FILE_TYPE_DIRECTORY:
		result.Type = FileTypeDirectory
	case filesystem.FileType_FILE_TYPE_SYMLINK:
		result.Type = FileTypeSymlink
	}
	return result
}
func mapEventType(value filesystem.EventType) string {
	switch value {
	case filesystem.EventType_EVENT_TYPE_CREATE:
		return "create"
	case filesystem.EventType_EVENT_TYPE_WRITE:
		return "write"
	case filesystem.EventType_EVENT_TYPE_REMOVE:
		return "remove"
	case filesystem.EventType_EVENT_TYPE_RENAME:
		return "rename"
	case filesystem.EventType_EVENT_TYPE_CHMOD:
		return "chmod"
	default:
		return "unknown"
	}
}
func fileConnectError(err error) error {
	if err == nil {
		return nil
	}
	var value *connect.Error
	if errors.As(err, &value) && value.Code() == connect.CodeNotFound {
		return &FileNotFoundError{APIError: APIError{Code: value.Code().String(), Message: value.Message(), Cause: err}}
	}
	return connectError(err)
}

var metadataKey = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)

func addMetadataHeaders(header http.Header, metadata map[string]string) error {
	for key, value := range metadata {
		if !metadataKey.MatchString(key) {
			return &InvalidArgumentError{Message: "invalid file metadata key: " + key}
		}
		for _, char := range []byte(value) {
			if char < 0x20 || char > 0x7e {
				return &InvalidArgumentError{Message: "file metadata values must be printable ASCII"}
			}
		}
		header.Set("X-Metadata-"+key, value)
	}
	return nil
}
