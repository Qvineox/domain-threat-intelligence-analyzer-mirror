from src.proto.services import analyzer_pb2, analyzer_pb2_grpc

import pandas as pd


class DomainScoringServerImpl(analyzer_pb2_grpc.DomainAnalysisServiceServicer):
    def __init__(self, semantic_model, resource_model, dga_model):
        self.semantic_model = semantic_model
        self.resource_model = resource_model
        self.dga_model = dga_model

    def GetFullScoring(self, request, context):
        print(request.semantics)

        return analyzer_pb2.DomainScoring(finalScore=0,
                                          semanticScore=self.semantic_score_calculation(request.semantics[0]),
                                          resourceScore=self.resource_records_score_calculation(request.resources[0]),
                                          dgaScore=0,
                                          tag=analyzer_pb2.DOMAIN_SCORE_BENIGN)

    def semantic_score_calculation(self, params):
        df = pd.DataFrame({"levels_count": [params.levelsCount],
                           "levels_mad": [params.levelsMAD],
                           "symbols_count": [params.symbolsCount],
                           "vowels_ratio": [params.vowelsRatio],
                           "consonants_ratio": [params.consonantsRatio],
                           "numbers_ratio": [params.numbersRatio],
                           "points_ratio": [params.pointsRatio],
                           "special_ratio": [params.specialRatio],
                           "unique_ratio": [params.uniqueRatio],
                           "max_repeated": [params.maxRepeated]})

        return self.semantic_model.predict(df)[0]

    def resource_records_score_calculation(self, params):
        df = pd.DataFrame({"a_records": [params.aRecords],
                           "mx_records": [params.mxRecords],
                           "cname_records": [params.cnameRecords],
                           "txt_records": [params.txtRecords],
                           "ptr_records": [params.ptrRecords],
                           "ptr_ratio": [params.ptrRatio]})

        return self.resource_model.predict(df)[0]
