package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"preprocessing/record"
	"strings"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		processAdGuardTXT()
		wg.Done()
	}()

	go func() {
		processCertPLCSV()
		wg.Done()
	}()

	go func() {
		processDNSWLCSV()
		wg.Done()
	}()

	go func() {
		processSefinekTXT()
		wg.Done()
	}()

	wg.Wait()
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

		//rec.ICMPResponse = record.PingDomain(rec.FullName)
		rec.IsLegit = false

		rec.Lookups, err = record.LookupRecords(rec.FullName)
		if err != nil {
			slog.Warn("failed to lookup: " + err.Error())
		}

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func processAdGuardTXT() {
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

		rec.Lookups, err = record.LookupRecords(rec.FullName)
		if err != nil {
			slog.Warn("failed to lookup: " + err.Error())
		}

		//rec.ICMPResponse = record.PingDomain(rec.FullName)
		rec.IsLegit = false

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func processDNSWLCSV() {
	inputFile, err := os.Open("../../data/dnswl/generic-dnswl.csv")
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

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_dnswl.csv", dir, now.Unix()))
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
	csvReader.Comma = ';'

	_, _ = csvReader.Read() // skip header

	var unique = make(map[string]bool)

	for {
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		name := line[3]
		_, ok := unique[name]
		if ok {
			continue
		} else {
			unique[name] = false
		}

		rec, err := record.NewDomainRecord(name)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		rec.Lookups, err = record.LookupRecords(rec.FullName)
		if err != nil {
			slog.Warn("failed to lookup: " + err.Error())
		}

		//rec.ICMPResponse = record.PingDomain(rec.FullName)
		rec.IsLegit = true

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func processSefinekTXT() {
	inputFile, err := os.Open("../../data/sefinek/hosts.fork.txt")
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

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_sefinek.csv", dir, now.Unix()))
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

		parts := strings.Split(line, " ")
		rec, err := record.NewDomainRecord(parts[1])
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		rec.Lookups, err = record.LookupRecords(rec.FullName)
		if err != nil {
			slog.Warn("failed to lookup: " + err.Error())
		}

		//rec.ICMPResponse = record.PingDomain(rec.FullName)
		rec.IsLegit = false

		err = csvWriter.Write(rec.ToCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
		time.Sleep(300 * time.Millisecond)
	}
}
