import logging
from concurrent import futures

import servers
from models.loaders import *

import grpc

from src.proto.services import analyzer_pb2, analyzer_pb2_grpc

logger = logging.getLogger(__name__)


# ref: https://grpc.io/docs/languages/python/quickstart/
# ref: https://www.freecodecamp.org/news/googles-protocol-buffers-in-python/


def serve():
    logging.basicConfig(
        filename='analyzer_logs.log',
        level=logging.DEBUG,
        format='%(asctime)s.%(msecs)03d %(levelname)s %(module)s - %(funcName)s: %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
    )

    logger.info("starting server...")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # load models

    scoring_server = servers.DomainScoringServerImpl(
        load_semantic_regression_model(),
        load_resource_records_regression_model(),
        load_keras_lstm_model()
    )

    # print(scoring_server.semantic_score_calculation('test'))
    # print(scoring_server.resource_records_score_calculation('test'))
    # print(scoring_server.dga_score_calculation('cvyh1po636avyrsxebwbkn7.ddns.net'))

    analyzer_pb2_grpc.add_DomainAnalysisServiceServicer_to_server(scoring_server, server)

    server.add_insecure_port("[::]:50051")
    server.start()

    logger.info("server started.")
    server.wait_for_termination()


def main():
    serve()


if __name__ == "__main__":
    main()
