package agentbox

import "fmt"

// APIError describes a non-successful HTTP or Connect response.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Cause      error
}

// Error formats the AgentBox API failure.
func (e *APIError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("agentbox: %d: %s", e.StatusCode, e.Message)
	}
	return "agentbox: " + e.Message
}

// Unwrap returns the underlying request error, if any.
func (e *APIError) Unwrap() error { return e.Cause }

// SandboxError is the base error for sandbox operations.
type SandboxError struct{ APIError }

// Error formats the sandbox operation failure.
func (e *SandboxError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying sandbox operation error, if any.
func (e *SandboxError) Unwrap() error { return e.APIError.Unwrap() }

// AuthenticationError reports missing or invalid credentials.
type AuthenticationError struct{ APIError }

// Error formats the authentication failure.
func (e *AuthenticationError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying authentication error, if any.
func (e *AuthenticationError) Unwrap() error { return e.APIError.Unwrap() }

// RateLimitError reports an exhausted API quota.
type RateLimitError struct{ APIError }

// Error formats the rate-limit failure.
func (e *RateLimitError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying rate-limit error, if any.
func (e *RateLimitError) Unwrap() error { return e.APIError.Unwrap() }

// TimeoutError reports a request, execution, or sandbox timeout.
type TimeoutError struct{ APIError }

// Error formats the timeout failure.
func (e *TimeoutError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying timeout error, if any.
func (e *TimeoutError) Unwrap() error { return e.APIError.Unwrap() }

// InvalidArgumentError reports invalid SDK input.
type InvalidArgumentError struct {
	Message string
	Cause   error
}

// Error formats the invalid argument failure.
func (e *InvalidArgumentError) Error() string { return "agentbox: " + e.Message }

// Unwrap returns the underlying validation error, if any.
func (e *InvalidArgumentError) Unwrap() error { return e.Cause }

// FileNotFoundError reports a missing sandbox file.
type FileNotFoundError struct{ APIError }

// Error formats the missing-file failure.
func (e *FileNotFoundError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying missing-file error, if any.
func (e *FileNotFoundError) Unwrap() error { return e.APIError.Unwrap() }

// NotEnoughSpaceError reports exhausted sandbox storage.
type NotEnoughSpaceError struct{ APIError }

// Error formats the storage-capacity failure.
func (e *NotEnoughSpaceError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying storage error, if any.
func (e *NotEnoughSpaceError) Unwrap() error { return e.APIError.Unwrap() }

// SandboxNotFoundError reports a missing or expired sandbox.
type SandboxNotFoundError struct{ APIError }

// Error formats the missing-sandbox failure.
func (e *SandboxNotFoundError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying missing-sandbox error, if any.
func (e *SandboxNotFoundError) Unwrap() error { return e.APIError.Unwrap() }

// TemplateError reports an invalid or incompatible template.
type TemplateError struct{ APIError }

// Error formats the template operation failure.
func (e *TemplateError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying template operation error, if any.
func (e *TemplateError) Unwrap() error { return e.APIError.Unwrap() }

// BuildError reports a failed template build.
type BuildError struct{ APIError }

// Error formats the template build failure.
func (e *BuildError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying template build error, if any.
func (e *BuildError) Unwrap() error { return e.APIError.Unwrap() }

// FileUploadError reports a failed template file upload.
type FileUploadError struct{ APIError }

// Error formats the file upload failure.
func (e *FileUploadError) Error() string { return e.APIError.Error() }

// Unwrap returns the underlying file upload error, if any.
func (e *FileUploadError) Unwrap() error { return e.APIError.Unwrap() }
