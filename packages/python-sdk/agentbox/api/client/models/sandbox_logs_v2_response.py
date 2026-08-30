from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.sandbox_log_entry import SandboxLogEntry


T = TypeVar("T", bound="SandboxLogsV2Response")


@_attrs_define
class SandboxLogsV2Response:
    """
    Attributes:
        logs (list['SandboxLogEntry']): Sandbox logs structured
        next_cursor (Union[Unset, str]): Opaque continuation cursor for the next page
    """

    logs: list["SandboxLogEntry"]
    next_cursor: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        logs = []
        for logs_item_data in self.logs:
            logs_item = logs_item_data.to_dict()
            logs.append(logs_item)

        next_cursor = self.next_cursor

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "logs": logs,
            }
        )
        if next_cursor is not UNSET:
            field_dict["nextCursor"] = next_cursor

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.sandbox_log_entry import SandboxLogEntry

        d = dict(src_dict)
        logs = []
        _logs = d.pop("logs")
        for logs_item_data in _logs:
            logs_item = SandboxLogEntry.from_dict(logs_item_data)

            logs.append(logs_item)

        next_cursor = d.pop("nextCursor", UNSET)

        sandbox_logs_v2_response = cls(
            logs=logs,
            next_cursor=next_cursor,
        )

        sandbox_logs_v2_response.additional_properties = d
        return sandbox_logs_v2_response

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
