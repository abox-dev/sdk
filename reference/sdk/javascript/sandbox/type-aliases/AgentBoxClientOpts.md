[agentbox-sdk-monorepo](../README.md) / AgentBoxClientOpts

# Type Alias: AgentBoxClientOpts

> **AgentBoxClientOpts** = `Omit`\<`ConnectionOpts`, `"signal"`\>

Connection options bound to an [AgentBox](../classes/AgentBox.md) client.

Same as ConnectionOpts without `signal`, which cancels a single
request and therefore can only be passed per call.
