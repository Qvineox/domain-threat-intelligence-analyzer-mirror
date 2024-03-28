import logging
from concurrent import futures

import grpc

from cmd.proto.services import analyzer_pb2, analyzer_pb2_grpc

logger = logging.getLogger(__name__)


# ref: https://grpc.io/docs/languages/python/quickstart/
# ref: https://www.freecodecamp.org/news/googles-protocol-buffers-in-python/

class DomainScoringServerImpl(analyzer_pb2_grpc.DomainAnalysisServiceServicer):
    def GetFullScoring(self, request, context):
        return analyzer_pb2.FullDomainScoring(finalScore=0,
                                              semanticScore=0,
                                              resourceScore=0,
                                              dgaScore=0,
                                              tag=analyzer_pb2.DOMAIN_SCORE_BENIGN)


def serve():
    logging.basicConfig(
        filename='analyzer_logs.log',
        level=logging.DEBUG,
        format='%(asctime)s.%(msecs)03d %(levelname)s %(module)s - %(funcName)s: %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
    )

    logger.info("starting server...")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    analyzer_pb2_grpc.add_DomainAnalysisServiceServicer_to_server(DomainScoringServerImpl(), server)

    server.add_insecure_port("[::]:50051")
    server.start()

    logger.info("server started.")
    server.wait_for_termination()


serve()
