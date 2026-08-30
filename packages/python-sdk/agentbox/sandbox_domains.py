SUPPORTED_SANDBOX_DOMAINS = ("agentbox-runtime.ru",)


def is_supported_sandbox_domain(sandbox_domain: str) -> bool:
    return sandbox_domain in SUPPORTED_SANDBOX_DOMAINS
