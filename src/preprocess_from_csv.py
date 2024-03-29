import csv

from src import servers
from src.models.loaders import *
from src.proto.services.analyzer_pb2 import SemanticData, ResourceRecordsData


def preprocess_from_csv():
    scoring_server = servers.DomainScoringServerImpl(
        load_semantic_regression_model(),
        load_resource_records_regression_model(),
        load_keras_lstm_model()
    )

    with open('merged_full_2024_03_27.csv', 'r', newline='') as csvfile_read:
        with open('output.csv', 'w', newline='') as csvfile_write:
            reader = csv.reader(csvfile_read, delimiter=',', quotechar='"')
            headers = next(reader)
            print(headers)

            writer = csv.writer(csvfile_write, delimiter=',')
            writer.writerow(['domain', 'semantic_score', 'resource_records_score', 'dga_score', 'is_legit'])

            for row in reader:
                domain = row[0]
                is_legit = row[-1]

                dga = scoring_server.dga_score_calculation(domain)

                s = SemanticData
                s.levelsCount = row[2]
                s.levelsMAD = row[3]
                s.symbolsCount = row[4]
                s.vowelsRatio = row[5]
                s.consonantsRatio = row[6]
                s.numbersRatio = row[7]
                s.pointsRatio = row[8]
                s.specialRatio = row[9]
                s.uniqueRatio = row[10]
                s.maxRepeated = row[11]

                r = ResourceRecordsData
                r.aRecords = row[12]
                r.mxRecords = row[13]
                r.cnameRecords = row[14]
                r.txtRecords = row[15]
                r.ptrRecords = row[16]
                r.ptrRatio = row[17]

                writer.writerow([domain,
                                 scoring_server.semantic_score_calculation(s),
                                 scoring_server.resource_records_score_calculation(r),
                                 dga,
                                 is_legit])


if __name__ == "__main__":
    preprocess_from_csv()
