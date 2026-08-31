<p align="center"><a href="https://agentbox.ru"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></a></p>

# `@abox-dev/cli`

The AgentBox command-line client for sandboxes and templates.

```bash
npm install --global @abox-dev/cli
agentbox configure ab_... --project-id PROJECT_ID
agentbox sandbox list
```

Configuration is stored in `~/.agentbox/config.json` with private directory and file permissions. `AGENTBOX_API_KEY` and `AGENTBOX_PROJECT_ID` take precedence over the file.

Create an API key using the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/). See the [CLI guide](https://docs.agentbox.ru/en/cli/) and [configuration reference](https://docs.agentbox.ru/en/cli/configuration/).
