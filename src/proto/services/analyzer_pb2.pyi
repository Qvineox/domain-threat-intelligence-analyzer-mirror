from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ScoreFlag(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    DOMAIN_SCORE_BENIGN: _ClassVar[ScoreFlag]
    DOMAIN_SCORE_DUBIOUS: _ClassVar[ScoreFlag]
    DOMAIN_SCORE_SUSPICIOUS: _ClassVar[ScoreFlag]
    DOMAIN_SCORE_MALICIOUS: _ClassVar[ScoreFlag]
DOMAIN_SCORE_BENIGN: ScoreFlag
DOMAIN_SCORE_DUBIOUS: ScoreFlag
DOMAIN_SCORE_SUSPICIOUS: ScoreFlag
DOMAIN_SCORE_MALICIOUS: ScoreFlag

class Domain(_message.Message):
    __slots__ = ("name",)
    NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    def __init__(self, name: _Optional[str] = ...) -> None: ...

class FullDomainScoring(_message.Message):
    __slots__ = ("finalScore", "semanticScore", "resourceScore", "dgaScore", "tag")
    FINALSCORE_FIELD_NUMBER: _ClassVar[int]
    SEMANTICSCORE_FIELD_NUMBER: _ClassVar[int]
    RESOURCESCORE_FIELD_NUMBER: _ClassVar[int]
    DGASCORE_FIELD_NUMBER: _ClassVar[int]
    TAG_FIELD_NUMBER: _ClassVar[int]
    finalScore: float
    semanticScore: float
    resourceScore: float
    dgaScore: float
    tag: ScoreFlag
    def __init__(self, finalScore: _Optional[float] = ..., semanticScore: _Optional[float] = ..., resourceScore: _Optional[float] = ..., dgaScore: _Optional[float] = ..., tag: _Optional[_Union[ScoreFlag, str]] = ...) -> None: ...
