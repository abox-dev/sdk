package agentbox

import (
	"encoding/json"
	"fmt"
	"time"
)

// SandboxState is the lifecycle state of a sandbox.
type SandboxState string

const (
	// SandboxRunning indicates that a sandbox accepts requests.
	SandboxRunning SandboxState = "running"
	// SandboxPaused indicates that a sandbox is suspended.
	SandboxPaused SandboxState = "paused"
)

// SandboxMetadata contains user-defined sandbox metadata.
type SandboxMetadata map[string]string

// SandboxNetworkConfig configures sandbox egress and public traffic.
type SandboxNetworkConfig struct {
	AllowOut           *[]string                        `json:"allowOut,omitzero"`
	AllowPublicTraffic *bool                            `json:"allowPublicTraffic,omitzero"`
	DenyOut            *[]string                        `json:"denyOut,omitzero"`
	MaskRequestHost    *string                          `json:"maskRequestHost,omitzero"`
	Rules              *map[string][]SandboxNetworkRule `json:"rules,omitzero"`
}

// SandboxNetworkRule applies request transformations to matching traffic.
type SandboxNetworkRule struct {
	Transform *SandboxNetworkTransform `json:"transform,omitzero"`
}

// SandboxNetworkTransform describes headers injected into matching requests.
type SandboxNetworkTransform struct {
	Headers *map[string]string `json:"headers,omitzero"`
}

// SandboxIAM configures workload identity tokens.
type SandboxIAM struct {
	Tokens *SandboxIAMTokens `json:"tokens,omitzero"`
}

// SandboxIAMTokens contains named workload identity token definitions.
type SandboxIAMTokens map[string]SandboxIAMToken

// SandboxIAMToken configures one workload identity token.
type SandboxIAMToken struct {
	Audience  string `json:"audience"`
	TokenType string `json:"tokenType"`
}

// SandboxLifecycle describes timeout and auto-resume behavior.
type SandboxLifecycle struct {
	AutoResume bool   `json:"autoResume"`
	OnTimeout  string `json:"onTimeout"`
}

// SandboxInfo contains current sandbox state and configuration.
type SandboxInfo struct {
	Alias               *string               `json:"alias,omitzero"`
	AllowInternetAccess *bool                 `json:"allowInternetAccess,omitzero"`
	CPUCount            int32                 `json:"cpuCount"`
	DiskSizeMB          int32                 `json:"diskSizeMB"`
	Domain              *string               `json:"domain,omitzero"`
	EndAt               time.Time             `json:"endAt"`
	EnvdVersion         string                `json:"envdVersion"`
	Lifecycle           *SandboxLifecycle     `json:"lifecycle,omitzero"`
	MemoryMB            int32                 `json:"memoryMB"`
	Metadata            *SandboxMetadata      `json:"metadata,omitzero"`
	Network             *SandboxNetworkConfig `json:"network,omitzero"`
	SandboxID           string                `json:"sandboxID"`
	StartedAt           time.Time             `json:"startedAt"`
	State               SandboxState          `json:"state"`
	TemplateID          string                `json:"templateID"`
}

// ListedSandbox is a compact sandbox list entry.
type ListedSandbox struct {
	Alias       *string          `json:"alias,omitzero"`
	CPUCount    int32            `json:"cpuCount"`
	DiskSizeMB  int32            `json:"diskSizeMB"`
	EndAt       time.Time        `json:"endAt"`
	EnvdVersion string           `json:"envdVersion"`
	MemoryMB    int32            `json:"memoryMB"`
	Metadata    *SandboxMetadata `json:"metadata,omitzero"`
	SandboxID   string           `json:"sandboxID"`
	StartedAt   time.Time        `json:"startedAt"`
	State       SandboxState     `json:"state"`
	TemplateID  string           `json:"templateID"`
}

// SandboxMetric is one timestamped resource-usage sample.
type SandboxMetric struct {
	CPUCount      int32   `json:"cpuCount"`
	CPUUsedPct    float32 `json:"cpuUsedPct"`
	DiskTotal     int64   `json:"diskTotal"`
	DiskUsed      int64   `json:"diskUsed"`
	MemCache      int64   `json:"memCache"`
	MemTotal      int64   `json:"memTotal"`
	MemUsed       int64   `json:"memUsed"`
	TimestampUnix int64   `json:"timestampUnix"`
}

// SnapshotInfo identifies a saved sandbox snapshot.
type SnapshotInfo struct {
	Names      []string `json:"names"`
	SnapshotID string   `json:"snapshotID"`
}

// SandboxLogEntry is one structured sandbox log record.
type SandboxLogEntry struct {
	Fields    map[string]string `json:"fields"`
	ID        *string           `json:"id,omitzero"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
}

// TemplateBuildStatus is a template build state.
type TemplateBuildStatus string

const (
	// BuildWaiting indicates that a template build is queued.
	BuildWaiting TemplateBuildStatus = "waiting"
	// BuildBuilding indicates that a template build is in progress.
	BuildBuilding TemplateBuildStatus = "building"
	// BuildReady indicates that a template build completed successfully.
	BuildReady TemplateBuildStatus = "ready"
	// BuildFailed indicates that a template build failed.
	BuildFailed TemplateBuildStatus = "error"
)

// BuildLogEntry is one structured template build log record.
type BuildLogEntry struct {
	ID        *string   `json:"id,omitzero"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Step      *string   `json:"step,omitzero"`
	Timestamp time.Time `json:"timestamp"`
}

// BuildStatusReason explains a terminal template build status.
type BuildStatusReason struct {
	LogEntries *[]BuildLogEntry `json:"logEntries,omitzero"`
	Message    string           `json:"message"`
	Step       *string          `json:"step,omitzero"`
}

// TemplateBuildInfo contains the latest status and logs for a build.
type TemplateBuildInfo struct {
	BuildID    string              `json:"buildID"`
	LogEntries []BuildLogEntry     `json:"logEntries"`
	Logs       []string            `json:"logs"`
	Reason     *BuildStatusReason  `json:"reason,omitzero"`
	Status     TemplateBuildStatus `json:"status"`
	TemplateID string              `json:"templateID"`
}

// TemplateBuild describes one historical template build.
type TemplateBuild struct {
	BuildID     string              `json:"buildID"`
	CPUCount    int32               `json:"cpuCount"`
	CreatedAt   time.Time           `json:"createdAt"`
	DiskSizeMB  *int32              `json:"diskSizeMB,omitzero"`
	EnvdVersion *string             `json:"envdVersion,omitzero"`
	FinishedAt  *time.Time          `json:"finishedAt,omitzero"`
	MemoryMB    int32               `json:"memoryMB"`
	Status      TemplateBuildStatus `json:"status"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// TeamUser identifies the user who created a template.
type TeamUser struct {
	ID string `json:"id"`
}

// TemplateInfo is a compact template list entry.
type TemplateInfo struct {
	BuildCount    int32               `json:"buildCount"`
	BuildID       string              `json:"buildID"`
	BuildStatus   TemplateBuildStatus `json:"buildStatus"`
	CPUCount      int32               `json:"cpuCount"`
	CreatedAt     time.Time           `json:"createdAt"`
	CreatedBy     *TeamUser           `json:"createdBy,omitzero"`
	DiskSizeMB    int32               `json:"diskSizeMB"`
	EnvdVersion   string              `json:"envdVersion"`
	LastSpawnedAt *time.Time          `json:"lastSpawnedAt,omitzero"`
	MemoryMB      int32               `json:"memoryMB"`
	Names         []string            `json:"names"`
	Public        bool                `json:"public"`
	SpawnCount    int64               `json:"spawnCount"`
	TemplateID    string              `json:"templateID"`
	UpdatedAt     time.Time           `json:"updatedAt"`
}

// TemplateWithBuilds contains a template and its build history.
type TemplateWithBuilds struct {
	Builds        []TemplateBuild `json:"builds"`
	CreatedAt     time.Time       `json:"createdAt"`
	LastSpawnedAt *time.Time      `json:"lastSpawnedAt,omitzero"`
	Names         []string        `json:"names"`
	Public        bool            `json:"public"`
	SpawnCount    int64           `json:"spawnCount"`
	TemplateID    string          `json:"templateID"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

// Page contains one page and an optional opaque continuation token.
type Page[T any] struct {
	Items     []T
	NextToken string
}

// convertModel deliberately copies only the fields in the public model. This
// keeps generated API changes internal until the public SDK model is updated.
func convertModel[Target any](source any) (Target, error) {
	var target Target
	data, err := json.Marshal(source)
	if err != nil {
		return target, fmt.Errorf("agentbox: encode API model: %w", err)
	}
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("agentbox: decode API model: %w", err)
	}
	return target, nil
}
