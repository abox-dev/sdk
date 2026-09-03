package agentbox

import "fmt"

// APIError describes a non-successful HTTP or Connect response.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Cause      error
}

func (e *APIError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("agentbox: %d: %s", e.StatusCode, e.Message)
	}
	return "agentbox: " + e.Message
}

func (e *APIError) Unwrap() error { return e.Cause }

// SandboxError is the base error for sandbox operations.
type SandboxError struct{ APIError }

func (e *SandboxError) Error() string { return e.APIError.Error() }
func (e *SandboxError) Unwrap() error { return e.APIError.Unwrap() }

// AuthenticationError reports missing or invalid credentials.
type AuthenticationError struct{ APIError }

func (e *AuthenticationError) Error() string { return e.APIError.Error() }
func (e *AuthenticationError) Unwrap() error { return e.APIError.Unwrap() }

// RateLimitError reports an exhausted API quota.
type RateLimitError struct{ APIError }

func (e *RateLimitError) Error() string { return e.APIError.Error() }
func (e *RateLimitError) Unwrap() error { return e.APIError.Unwrap() }

// TimeoutError reports a request, execution, or sandbox timeout.
type TimeoutError struct{ APIError }

func (e *TimeoutError) Error() string { return e.APIError.Error() }
func (e *TimeoutError) Unwrap() error { return e.APIError.Unwrap() }

// InvalidArgumentError reports invalid SDK input.
type InvalidArgumentError struct {
	Message string
	Cause   error
}

func (e *InvalidArgumentError) Error() string { return "agentbox: " + e.Message }
func (e *InvalidArgumentError) Unwrap() error { return e.Cause }

// FileNotFoundError reports a missing sandbox file.
type FileNotFoundError struct{ APIError }

func (e *FileNotFoundError) Error() string { return e.APIError.Error() }
func (e *FileNotFoundError) Unwrap() error { return e.APIError.Unwrap() }

// NotEnoughSpaceError reports exhausted sandbox storage.
type NotEnoughSpaceError struct{ APIError }

func (e *NotEnoughSpaceError) Error() string { return e.APIError.Error() }
func (e *NotEnoughSpaceError) Unwrap() error { return e.APIError.Unwrap() }

// SandboxNotFoundError reports a missing or expired sandbox.
type SandboxNotFoundError struct{ APIError }

func (e *SandboxNotFoundError) Error() string { return e.APIError.Error() }
func (e *SandboxNotFoundError) Unwrap() error { return e.APIError.Unwrap() }

// TemplateError reports an invalid or incompatible template.
type TemplateError struct{ APIError }

func (e *TemplateError) Error() string { return e.APIError.Error() }
func (e *TemplateError) Unwrap() error { return e.APIError.Unwrap() }

// BuildError reports a failed template build.
type BuildError struct{ APIError }

func (e *BuildError) Error() string { return e.APIError.Error() }
func (e *BuildError) Unwrap() error { return e.APIError.Unwrap() }

// FileUploadError reports a failed template file upload.
type FileUploadError struct{ APIError }

func (e *FileUploadError) Error() string { return e.APIError.Error() }
func (e *FileUploadError) Unwrap() error { return e.APIError.Unwrap() }
