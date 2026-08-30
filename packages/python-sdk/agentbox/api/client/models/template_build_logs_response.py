from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.logs_source import LogsSource
from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.build_log_entry import BuildLogEntry


T = TypeVar("T", bound="TemplateBuildLogsResponse")


@_attrs_define
class TemplateBuildLogsResponse:
    """
    Attributes:
        logs (list['BuildLogEntry']): Build logs structured
        next_cursor (Union[Unset, str]): Opaque continuation cursor for the next page
        source (Union[Unset, LogsSource]): Source of the logs that should be returned
    """

    logs: list["BuildLogEntry"]
    next_cursor: Union[Unset, str] = UNSET
    source: Union[Unset, LogsSource] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        logs = []
        for logs_item_data in self.logs:
            logs_item = logs_item_data.to_dict()
            logs.append(logs_item)

        next_cursor = self.next_cursor

        source: Union[Unset, str] = UNSET
        if not isinstance(self.source, Unset):
            source = self.source.value

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "logs": logs,
            }
        )
        if next_cursor is not UNSET:
            field_dict["nextCursor"] = next_cursor
        if source is not UNSET:
            field_dict["source"] = source

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.build_log_entry import BuildLogEntry

        d = dict(src_dict)
        logs = []
        _logs = d.pop("logs")
        for logs_item_data in _logs:
            logs_item = BuildLogEntry.from_dict(logs_item_data)

            logs.append(logs_item)

        next_cursor = d.pop("nextCursor", UNSET)

        _source = d.pop("source", UNSET)
        source: Union[Unset, LogsSource]
        if isinstance(_source, Unset):
            source = UNSET
        else:
            source = LogsSource(_source)

        template_build_logs_response = cls(
            logs=logs,
            next_cursor=next_cursor,
            source=source,
        )

        template_build_logs_response.additional_properties = d
        return template_build_logs_response

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
