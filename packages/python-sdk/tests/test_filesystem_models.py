import datetime
import io

from protobuf.wkt import Timestamp

from agentbox.envd.filesystem import filesystem_pb
from agentbox.sandbox.filesystem.filesystem import (
    multipart_body_is_streamed,
    FileType,
    WriteInfo,
    map_entry_info,
    map_file_type_str,
)


def test_write_info_from_dict_converts_type_to_enum():
    info = WriteInfo.from_dict(
        {"name": "a.txt", "type": "file", "path": "/home/user/a.txt"}
    )
    assert info.type is FileType.FILE

    info = WriteInfo.from_dict({"name": "dir", "type": "dir", "path": "/home/user/dir"})
    assert info.type is FileType.DIR


def test_write_info_from_dict_handles_missing_or_unknown_type():
    info = WriteInfo.from_dict({"name": "a.txt", "path": "/home/user/a.txt"})
    assert info.type is None

    info = WriteInfo.from_dict(
        {"name": "a.txt", "type": "symlink", "path": "/home/user/a.txt"}
    )
    assert info.type is None


def test_map_file_type_str():
    assert map_file_type_str("file") is FileType.FILE
    assert map_file_type_str("dir") is FileType.DIR
    assert map_file_type_str("unknown") is None
    assert map_file_type_str(None) is None


def test_map_entry_info_maps_every_known_file_type():
    def entry(file_type: filesystem_pb.FileType) -> filesystem_pb.EntryInfo:
        return filesystem_pb.EntryInfo(
            name="entry",
            type=file_type,
            path="/home/user/entry",
            size=0,
            mode=0o644,
            permissions="-rw-r--r--",
            owner="user",
            group="user",
        )

    assert map_entry_info(entry(filesystem_pb.FileType.FILE)).type is FileType.FILE
    assert map_entry_info(entry(filesystem_pb.FileType.DIRECTORY)).type is FileType.DIR
    assert (
        map_entry_info(entry(filesystem_pb.FileType.SYMLINK)).type is FileType.SYMLINK
    )


def test_map_entry_info_keeps_symlink_target_on_symlink_entries():
    entry = filesystem_pb.EntryInfo(
        name="link",
        type=filesystem_pb.FileType.SYMLINK,
        path="/home/user/link",
        size=0,
        mode=0o777,
        permissions="lrwxrwxrwx",
        owner="user",
        group="user",
    )
    entry.symlink_target = "/home/user/a.txt"

    info = map_entry_info(entry)

    assert info.type is FileType.SYMLINK
    assert info.symlink_target == "/home/user/a.txt"


def test_map_entry_info_modified_time_is_timezone_aware():
    entry = filesystem_pb.EntryInfo(
        name="a.txt",
        type=filesystem_pb.FileType.FILE,
        path="/home/user/a.txt",
        size=4,
        mode=0o644,
        permissions="-rw-r--r--",
        owner="user",
        group="user",
    )
    entry.modified_time = Timestamp.from_datetime(
        datetime.datetime(2026, 1, 2, 3, 4, 5, tzinfo=datetime.timezone.utc)
    )

    info = map_entry_info(entry)

    assert info.modified_time.tzinfo == datetime.timezone.utc
    assert info.modified_time == datetime.datetime(
        2026, 1, 2, 3, 4, 5, tzinfo=datetime.timezone.utc
    )


def test_multipart_body_is_streamed_only_for_binary_file_like_entries():
    # The multipart upload drops its request deadline only when httpx really
    # streams the body; `_to_httpx_file` reads text file-like data into memory,
    # so those uploads must stay bounded like str/bytes ones.
    assert not multipart_body_is_streamed([{"path": "a.txt", "data": "text"}])
    assert not multipart_body_is_streamed([{"path": "a.bin", "data": b"bytes"}])
    assert not multipart_body_is_streamed(
        [{"path": "a.txt", "data": io.StringIO("text")}]
    )
    assert multipart_body_is_streamed([{"path": "a.bin", "data": io.BytesIO(b"bytes")}])
    # Any streamed entry makes the whole body streamed.
    assert multipart_body_is_streamed(
        [
            {"path": "a.txt", "data": "text"},
            {"path": "b.bin", "data": io.BytesIO(b"bytes")},
        ]
    )
