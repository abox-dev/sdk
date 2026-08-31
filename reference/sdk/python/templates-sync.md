## `agentbox.template_sync.main`

**Classes:**

- [**Template**](#agentbox.template_sync.main.Template) – Synchronous template builder for AgentBox sandboxes.

### `agentbox.template_sync.main.Template`

Bases: <code>[TemplateBase](#agentbox.template.main.TemplateBase)</code>

Synchronous template builder for AgentBox sandboxes.

**Functions:**

- [**assign_tags**](#agentbox.template_sync.main.Template.assign_tags) – Assign tag(s) to an existing template build.
- [**build**](#agentbox.template_sync.main.Template.build) – Build and deploy a template to AgentBox infrastructure.
- [**build_in_background**](#agentbox.template_sync.main.Template.build_in_background) – Build and deploy a template to AgentBox infrastructure without waiting for completion.
- [**exists**](#agentbox.template_sync.main.Template.exists) – Check if a template with the given name exists.
- [**from_alpine_image**](#agentbox.template_sync.main.Template.from_alpine_image) – Start template from an Alpine base image.
- [**from_arch_image**](#agentbox.template_sync.main.Template.from_arch_image) – Start template from an Arch Linux base image.
- [**from_aws_registry**](#agentbox.template_sync.main.Template.from_aws_registry) – Start template from an AWS ECR registry image.
- [**from_base_image**](#agentbox.template_sync.main.Template.from_base_image) – Start template from the pinned AgentBox base image.
- [**from_bun_image**](#agentbox.template_sync.main.Template.from_bun_image) – Start template from a Bun base image.
- [**from_debian_image**](#agentbox.template_sync.main.Template.from_debian_image) – Start template from a Debian base image.
- [**from_dockerfile**](#agentbox.template_sync.main.Template.from_dockerfile) – Parse a Dockerfile and convert it to Template SDK format.
- [**from_fedora_image**](#agentbox.template_sync.main.Template.from_fedora_image) – Start template from a Fedora base image.
- [**from_gcp_registry**](#agentbox.template_sync.main.Template.from_gcp_registry) – Start template from a GCP Artifact Registry or Container Registry image.
- [**from_image**](#agentbox.template_sync.main.Template.from_image) – Start template from a Docker image.
- [**from_node_image**](#agentbox.template_sync.main.Template.from_node_image) – Start template from a Node.js base image.
- [**from_python_image**](#agentbox.template_sync.main.Template.from_python_image) – Start template from a Python base image.
- [**from_template**](#agentbox.template_sync.main.Template.from_template) – Start template from an existing AgentBox template.
- [**from_ubuntu_image**](#agentbox.template_sync.main.Template.from_ubuntu_image) – Start template from an Ubuntu base image.
- [**get_build_status**](#agentbox.template_sync.main.Template.get_build_status) – Get the status of a build.
- [**get_tags**](#agentbox.template_sync.main.Template.get_tags) – Get all tags for a template.
- [**remove_tags**](#agentbox.template_sync.main.Template.remove_tags) – Remove tag(s) from a template.
- [**skip_cache**](#agentbox.template_sync.main.Template.skip_cache) – Skip cache for all subsequent build instructions from this point.
- [**to_dockerfile**](#agentbox.template_sync.main.Template.to_dockerfile) – Convert a template to Dockerfile format.
- [**to_json**](#agentbox.template_sync.main.Template.to_json) – Convert a template to JSON representation.

#### `agentbox.template_sync.main.Template.assign_tags`

```python
assign_tags(target_name, tags, **opts)
```

Assign tag(s) to an existing template build.

:param target_name: Template name in 'name:tag' format (the source build to tag from)
:param tags: Tag or tags to assign
:return: TemplateTagInfo with build_id and assigned tags

Example
```python
from agentbox import Template

# Assign a single tag
result = Template.assign_tags('my-template:v1.0', 'production')

# Assign multiple tags
result = Template.assign_tags('my-template:v1.0', ['production', 'stable'])
```

#### `agentbox.template_sync.main.Template.build`

```python
build(
    template,
    name,
    *,
    tags=None,
    cpu_count=2,
    memory_mb=1024,
    skip_cache=False,
    on_build_logs=None,
    **opts
)
```

Build and deploy a template to AgentBox infrastructure.

:param template: The template to build
:param name: Template name in 'name' or 'name:tag' format
:param tags: Optional additional tags to assign to the template
:param cpu_count: Number of CPUs allocated to the sandbox
:param memory_mb: Amount of memory in MB allocated to the sandbox
:param skip_cache: If True, forces a complete rebuild ignoring cache
:param on_build_logs: Callback function to receive build logs during the build process

Example
```python
from agentbox import Template

template = (
    Template()
    .from_python_image('3')
    .copy('requirements.txt', '/home/user/')
    .run_cmd('pip install -r /home/user/requirements.txt')
)

# Build with single tag
Template.build(template, 'my-python-env:v1.0')

# Build with multiple tags
Template.build(template, 'my-python-env', tags=['v1.1.0', 'stable'])
```

#### `agentbox.template_sync.main.Template.build_in_background`

```python
build_in_background(
    template,
    name,
    *,
    tags=None,
    cpu_count=2,
    memory_mb=1024,
    skip_cache=False,
    on_build_logs=None,
    **opts
)
```

Build and deploy a template to AgentBox infrastructure without waiting for completion.

:param template: The template to build
:param name: Template name in 'name' or 'name:tag' format
:param tags: Optional additional tags to assign to the template
:param cpu_count: Number of CPUs allocated to the sandbox
:param memory_mb: Amount of memory in MB allocated to the sandbox
:param skip_cache: If True, forces a complete rebuild ignoring cache
:return: BuildInfo containing the template ID and build ID

Example
```python
from agentbox import Template

template = (
    Template()
    .from_python_image('3')
    .run_cmd('echo "test"')
    .set_start_cmd('echo "Hello"', 'sleep 1')
)

# Build with single tag
build_info = Template.build_in_background(template, 'my-python-env:v1.0')

# Build with multiple tags
build_info = Template.build_in_background(template, 'my-python-env', tags=['v1.1.0', 'stable'])
```

#### `agentbox.template_sync.main.Template.exists`

```python
exists(name, **opts)
```

Check if a template with the given name exists.

:param name: Template name to check
:return: True if the name exists, False otherwise

Example
```python
from agentbox import Template

exists = Template.exists('my-python-env')
if exists:
    print('Template exists!')
```

#### `agentbox.template_sync.main.Template.from_alpine_image`

```python
from_alpine_image(variant='3.24')
```

Start template from an Alpine base image.

:param variant: Alpine image variant (default: '3.24')

:return: `TemplateBuilder` class

Example
```python
Template().from_alpine_image('3.24')
```

#### `agentbox.template_sync.main.Template.from_arch_image`

```python
from_arch_image(variant='latest')
```

Start template from an Arch Linux base image.

Defaults to `latest`: Arch is a rolling release and template
provisioning runs `pacman -Syu`, so pinning a tag would not change
the built result.

:param variant: Arch Linux image variant (default: 'latest')

:return: `TemplateBuilder` class

Example
```python
Template().from_arch_image('base-devel')
```

#### `agentbox.template_sync.main.Template.from_aws_registry`

```python
from_aws_registry(image, access_key_id, secret_access_key, region)
```

Start template from an AWS ECR registry image.

:param image: Docker image name from AWS ECR
:param access_key_id: AWS access key ID
:param secret_access_key: AWS secret access key
:param region: AWS region

:return: `TemplateBuilder` class

Example
```python
Template().from_aws_registry(
    '123456789.dkr.ecr.us-west-2.amazonaws.com/myimage:latest',
    access_key_id='AKIA...',
    secret_access_key='...',
    region='us-west-2'
)
```

#### `agentbox.template_sync.main.Template.from_base_image`

```python
from_base_image()
```

Start template from the pinned AgentBox base image.

:return: `TemplateBuilder` class

Example
```python
Template().from_base_image()
```

#### `agentbox.template_sync.main.Template.from_bun_image`

```python
from_bun_image(variant='latest')
```

Start template from a Bun base image.

:param variant: Bun image variant (default: 'latest')

:return: `TemplateBuilder` class

#### `agentbox.template_sync.main.Template.from_debian_image`

```python
from_debian_image(variant='stable')
```

Start template from a Debian base image.

:param variant: Debian image variant

:return: `TemplateBuilder` class

Example
```python
Template().from_debian_image('bookworm')
```

#### `agentbox.template_sync.main.Template.from_dockerfile`

```python
from_dockerfile(dockerfile_content_or_path)
```

Parse a Dockerfile and convert it to Template SDK format.

:param dockerfile_content_or_path: Either the Dockerfile content as a string, or a path to a Dockerfile file

:return: `TemplateBuilder` class

Example
```python
Template().from_dockerfile('Dockerfile')
Template().from_dockerfile('FROM python:3\nRUN pip install numpy')
```

#### `agentbox.template_sync.main.Template.from_fedora_image`

```python
from_fedora_image(variant='44')
```

Start template from a Fedora base image.

:param variant: Fedora image variant (default: '44')

:return: `TemplateBuilder` class

Example
```python
Template().from_fedora_image('44')
```

#### `agentbox.template_sync.main.Template.from_gcp_registry`

```python
from_gcp_registry(image, service_account_json)
```

Start template from a GCP Artifact Registry or Container Registry image.

:param image: Docker image name from GCP registry
:param service_account_json: Service account JSON string, dict, or path to JSON file

:return: `TemplateBuilder` class

Example
```python
Template().from_gcp_registry(
    'gcr.io/myproject/myimage:latest',
    service_account_json='path/to/service-account.json'
)
```

#### `agentbox.template_sync.main.Template.from_image`

```python
from_image(image, username=None, password=None)
```

Start template from a Docker image.

:param image: Docker image name (e.g., 'ubuntu:24.04')
:param username: Username for private registry authentication
:param password: Password for private registry authentication

:return: `TemplateBuilder` class

Example
```python
Template().from_image('python:3')

# With credentials (optional)
Template().from_image('myregistry.com/myimage:latest', username='user', password='pass')
```

#### `agentbox.template_sync.main.Template.from_node_image`

```python
from_node_image(variant='lts')
```

Start template from a Node.js base image.

:param variant: Node.js image variant (default: 'lts')

:return: `TemplateBuilder` class

Example
```python
Template().from_node_image('24')
```

#### `agentbox.template_sync.main.Template.from_python_image`

```python
from_python_image(version='3')
```

Start template from a Python base image.

:param version: Python version (default: '3')

:return: `TemplateBuilder` class

Example
```python
Template().from_python_image('3')
```

#### `agentbox.template_sync.main.Template.from_template`

```python
from_template(template)
```

Start template from an existing AgentBox template.

:param template: AgentBox template ID or alias

:return: `TemplateBuilder` class

Example
```python
Template().from_template('my-base-template')
```

#### `agentbox.template_sync.main.Template.from_ubuntu_image`

```python
from_ubuntu_image(variant='latest')
```

Start template from an Ubuntu base image.

:param variant: Ubuntu image variant (default: 'latest')

:return: `TemplateBuilder` class

Example
```python
Template().from_ubuntu_image('24.04')
```

#### `agentbox.template_sync.main.Template.get_build_status`

```python
get_build_status(build_info, logs_offset=0, **opts)
```

Get the status of a build.

:param build_info: Build identifiers returned from build_in_background
:param logs_offset: Offset for fetching logs
:return: TemplateBuild containing the build status and logs

Example
```python
from agentbox import Template

build_info = Template.build_in_background(template, 'my-template')
status = Template.get_build_status(build_info, logs_offset=0)
```

#### `agentbox.template_sync.main.Template.get_tags`

```python
get_tags(template_id, **opts)
```

Get all tags for a template.

:param template_id: Template ID or name
:return: List of TemplateTag with tag name, build_id, and created_at

Example
```python
from agentbox import Template

tags = Template.get_tags('my-template')
for tag in tags:
    print(f"Tag: {tag.tag}, Build: {tag.build_id}, Created: {tag.created_at}")
```

#### `agentbox.template_sync.main.Template.remove_tags`

```python
remove_tags(name, tags, **opts)
```

Remove tag(s) from a template.

:param name: Template name
:param tags: Tag or tags to remove

Example
```python
from agentbox import Template

# Remove a single tag
Template.remove_tags('my-template', 'production')

# Remove multiple tags
Template.remove_tags('my-template', ['production', 'stable'])
```

#### `agentbox.template_sync.main.Template.skip_cache`

```python
skip_cache()
```

Skip cache for all subsequent build instructions from this point.

:return: `TemplateBase` class

Example
```python
template.skip_cache().from_python_image('3.11')
```

#### `agentbox.template_sync.main.Template.to_dockerfile`

```python
to_dockerfile(template)
```

Convert a template to Dockerfile format.

Note: Templates based on other AgentBox templates cannot be converted to Dockerfile.

:param template: The template to convert (TemplateBuilder or TemplateFinal instance)

:return: Dockerfile string representation

:raises ValueError: If the template is based on another AgentBox template or has no base image

Example
```python
template = Template().from_python_image('3').copy('app.py', '/app/')
dockerfile = TemplateBase.to_dockerfile(template)
```

#### `agentbox.template_sync.main.Template.to_json`

```python
to_json(template)
```

Convert a template to JSON representation.

:param template: The template to convert (TemplateBuilder or TemplateFinal instance)

:return: JSON string representation of the template

Example
```python
template = Template().from_python_image('3').copy('app.py', '/app/')
json_str = TemplateBase.to_json(template)
```

