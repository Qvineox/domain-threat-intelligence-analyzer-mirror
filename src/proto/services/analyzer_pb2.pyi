from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

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

class FullDomainScoringParams(_message.Message):
    __slots__ = ("semantics", "resources", "names")
    SEMANTICS_FIELD_NUMBER: _ClassVar[int]
    RESOURCES_FIELD_NUMBER: _ClassVar[int]
    NAMES_FIELD_NUMBER: _ClassVar[int]
    semantics: _containers.RepeatedCompositeFieldContainer[SemanticData]
    resources: _containers.RepeatedCompositeFieldContainer[ResourceRecordsData]
    names: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, semantics: _Optional[_Iterable[_Union[SemanticData, _Mapping]]] = ..., resources: _Optional[_Iterable[_Union[ResourceRecordsData, _Mapping]]] = ..., names: _Optional[_Iterable[str]] = ...) -> None: ...

class DomainScoring(_message.Message):
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

class SemanticData(_message.Message):
    __slots__ = ("levelsCount", "levelsMAD", "symbolsCount", "vowelsRatio", "consonantsRatio", "numbersRatio", "pointsRatio", "specialRatio", "uniqueRatio", "maxRepeated")
    LEVELSCOUNT_FIELD_NUMBER: _ClassVar[int]
    LEVELSMAD_FIELD_NUMBER: _ClassVar[int]
    SYMBOLSCOUNT_FIELD_NUMBER: _ClassVar[int]
    VOWELSRATIO_FIELD_NUMBER: _ClassVar[int]
    CONSONANTSRATIO_FIELD_NUMBER: _ClassVar[int]
    NUMBERSRATIO_FIELD_NUMBER: _ClassVar[int]
    POINTSRATIO_FIELD_NUMBER: _ClassVar[int]
    SPECIALRATIO_FIELD_NUMBER: _ClassVar[int]
    UNIQUERATIO_FIELD_NUMBER: _ClassVar[int]
    MAXREPEATED_FIELD_NUMBER: _ClassVar[int]
    levelsCount: float
    levelsMAD: float
    symbolsCount: float
    vowelsRatio: float
    consonantsRatio: float
    numbersRatio: float
    pointsRatio: float
    specialRatio: float
    uniqueRatio: float
    maxRepeated: float
    def __init__(self, levelsCount: _Optional[float] = ..., levelsMAD: _Optional[float] = ..., symbolsCount: _Optional[float] = ..., vowelsRatio: _Optional[float] = ..., consonantsRatio: _Optional[float] = ..., numbersRatio: _Optional[float] = ..., pointsRatio: _Optional[float] = ..., specialRatio: _Optional[float] = ..., uniqueRatio: _Optional[float] = ..., maxRepeated: _Optional[float] = ...) -> None: ...

class ResourceRecordsData(_message.Message):
    __slots__ = ("aRecords", "mxRecords", "cnameRecords", "txtRecords", "ptrRecords", "ptrRatio")
    ARECORDS_FIELD_NUMBER: _ClassVar[int]
    MXRECORDS_FIELD_NUMBER: _ClassVar[int]
    CNAMERECORDS_FIELD_NUMBER: _ClassVar[int]
    TXTRECORDS_FIELD_NUMBER: _ClassVar[int]
    PTRRECORDS_FIELD_NUMBER: _ClassVar[int]
    PTRRATIO_FIELD_NUMBER: _ClassVar[int]
    aRecords: float
    mxRecords: float
    cnameRecords: float
    txtRecords: float
    ptrRecords: float
    ptrRatio: float
    def __init__(self, aRecords: _Optional[float] = ..., mxRecords: _Optional[float] = ..., cnameRecords: _Optional[float] = ..., txtRecords: _Optional[float] = ..., ptrRecords: _Optional[float] = ..., ptrRatio: _Optional[float] = ...) -> None: ...
