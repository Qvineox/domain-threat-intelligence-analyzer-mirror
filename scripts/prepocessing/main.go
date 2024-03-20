package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"preprocessing/record"
	"time"
)

func main() {
	processAdGuardCSV()
	processCertPLCSV()
}

func processCertPLCSV() {
	inputFile, err := os.Open("../../data/certpl/domains.csv")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var now = time.Now()

	dir := fmt.Sprintf("output/%s", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_certpl.csv", dir, now.Unix()))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	csvWriter := csv.NewWriter(outputFile)
	err = csvWriter.Write(record.CSVHeader())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	csvReader := csv.NewReader(inputFile)
	csvReader.Comma = '\t'

	_, _ = csvReader.Read() // skip header

	for {
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		name := line[1]

		rec, err := record.NewDomainRecord(name)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		rec.IsLegit = false

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func processAdGuardCSV() {
	inputFile, err := os.Open("../../data/adguard/adguarddns-justdomains.txt")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var now = time.Now()

	dir := fmt.Sprintf("output/%s", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_adguard.csv", dir, now.Unix()))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	csvWriter := csv.NewWriter(outputFile)
	err = csvWriter.Write(record.CSVHeader())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		rec, err := record.NewDomainRecord(line)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		rec.IsLegit = false

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}
