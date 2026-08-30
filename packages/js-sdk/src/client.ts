import { ConnectionOpts } from './connectionConfig'
import { Sandbox } from './sandbox'
import { Template, TemplateBase } from './template'
import { callableTemplate } from './template/callable'

/**
 * Connection options bound to an {@link AgentBox} client.
 *
 * Same as {@link ConnectionOpts} without `signal`, which cancels a single
 * request and therefore can only be passed per call.
 */
export type AgentBoxClientOpts = Omit<ConnectionOpts, 'signal'>

/**
 * AgentBox client with an explicitly bound connection configuration.
 *
 * The resources exposed by the client ({@link AgentBox.Sandbox} and
 * {@link AgentBox.Template}) behave exactly like the top-level `Sandbox` and `Template` exports,
 * except the options passed to the client are used as the defaults instead of
 * the environment variables.
 * Per-call options still take precedence over the client's options.
 *
 * Multiple clients are fully isolated from each other and from the top-level
 * env-configured exports.
 *
 * @example
 * ```ts
 * import { AgentBox } from '@abox-dev/sdk'
 *
 * const client = new AgentBox({ apiKey: 'ab_...' })
 *
 * const sandbox = await client.Sandbox.create()
 * await client.Template.build(client.Template().fromPythonImage('3'), 'my-env')
 * ```
 */
export class AgentBox {
  /**
   * `Sandbox` class bound to this client's connection configuration.
   */
  readonly Sandbox: typeof Sandbox

  /**
   * `Template` bound to this client's connection configuration. Both the
   * builder (`client.Template()`) and the statics
   * (`client.Template.build(...)`, `client.Template.exists(...)`, …) work like
   * the top-level `Template`.
   */
  readonly Template: typeof Template

  /**
   * Create a new client with the connection options bound to it.
   *
   * @param opts connection options used as the defaults for every call made
   *   through this client's resource classes.
   */
  constructor(opts?: AgentBoxClientOpts) {
    // Options are copied so later mutations of the caller's object cannot
    // change the bound configuration. `signal` is dropped rather than only
    // typed away, since it cancels a single request and a caller passing a
    // wider-typed object (or plain JS) would otherwise bind it to every call.
    const boundOpts: AgentBoxClientOpts = { ...(opts ?? {}) }
    delete (boundOpts as ConnectionOpts).signal

    this.Sandbox = class extends Sandbox {
      protected static override readonly boundOpts = boundOpts
    }

    this.Template = callableTemplate(
      class extends TemplateBase {
        protected static override readonly boundOpts = boundOpts
      }
    )
  }
}
