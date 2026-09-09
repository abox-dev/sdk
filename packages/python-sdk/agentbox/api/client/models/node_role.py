from enum import Enum


class NodeRole(str, Enum):
    TEMPLATE_BUILDER = "template-builder"
    WORKER = "worker"

    def __str__(self) -> str:
        return str(self.value)
