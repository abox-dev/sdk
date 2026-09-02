export { ApiClient } from './api'
export type { components, paths } from './api'

export { ConnectionConfig } from './connectionConfig'
export type { ConnectionOpts, Username } from './connectionConfig'
export {
  AuthenticationError,
  FileNotFoundError,
  InvalidArgumentError,
  NotEnoughSpaceError,
  SandboxError,
  SandboxNotFoundError,
  TemplateError,
  TimeoutError,
  RateLimitError,
  BuildError,
  FileUploadError,
} from './errors'
export type { Logger } from './logs'

export { getSignature } from './sandbox/signature'

export { FileType } from './sandbox/filesystem'
export type {
  WriteInfo,
  EntryInfo,
  Filesystem,
  FilesystemWriteOpts,
  FilesystemReadOpts,
} from './sandbox/filesystem'
export { FilesystemEventType } from './sandbox/filesystem/watchHandle'
export type {
  FilesystemEvent,
  WatchHandle,
} from './sandbox/filesystem/watchHandle'

export { CommandExitError } from './sandbox/commands/commandHandle'
export type {
  CommandResult,
  Stdout,
  Stderr,
  PtyOutput,
  CommandHandle,
} from './sandbox/commands/commandHandle'
export type {
  SandboxInfo,
  SandboxMetrics,
  SandboxOpts,
  SandboxApiOpts,
  SandboxConnectOpts,
  SandboxForkOpts,
  SandboxMetricsOpts,
  SandboxPauseOpts,
  SandboxState,
  SandboxListOpts,
  SandboxPaginator,
  SandboxIamOpts,
  SandboxIamToken,
  SandboxIamTokenType,
  SandboxNetworkOpts,
  SandboxNetworkInfo,
  SandboxNetworkSelector,
  SandboxNetworkSelectorContext,
  SandboxNetworkRule,
  SandboxNetworkRuleInfo,
  SandboxNetworkRules,
  SandboxNetworkTransform,
  SandboxNetworkTransformContext,
  SandboxNetworkTransformResolver,
  SandboxNetworkUpdate,
  SandboxOnTimeout,
  SandboxLifecycle,
  SandboxInfoLifecycle,
  SnapshotInfo,
  SnapshotListOpts,
  SnapshotPaginator,
  CreateSnapshotOpts,
} from './sandbox/sandboxApi'

export { ALL_TRAFFIC } from './sandbox/network'

export type {
  ProcessInfo,
  CommandRequestOpts,
  CommandConnectOpts,
  CommandStartOpts,
  Commands,
  Pty,
} from './sandbox/commands'

export { Sandbox }
import { Sandbox } from './sandbox'

export { AgentBox, type AgentBoxClientOpts } from './client'

export default Sandbox

export * from './template'

export {
  ReadyCmd,
  waitForPort,
  waitForURL,
  waitForProcess,
  waitForFile,
  waitForTimeout,
} from './template/readycmd'

export {
  LogEntry,
  LogEntryStart,
  LogEntryEnd,
  type LogEntryLevel,
  defaultBuildLogger,
} from './template/logger'
