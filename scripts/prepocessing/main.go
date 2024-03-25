package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"preprocessing/domainAnalysis"
	"preprocessing/record"
	"strings"
	"time"
)

func main() {
	//wg := sync.WaitGroup{}

	//go func() {
	//	wg.Add(1)
	//	processAdGuardTXT()
	//	wg.Done()
	//}()
	//
	//go func() {
	//	wg.Add(1)
	//	processCertPLCSV()
	//	wg.Done()
	//}()
	//
	//go func() {
	//	wg.Add(1)
	//	processDNSWLCSV()
	//	wg.Done()
	//}()
	//
	//go func() {
	//	wg.Add(1)
	//	processSefinekTXT()
	//	wg.Done()
	//}()

	//go func() {
	//	processRGSTXT()
	//	wg.Done()
	//}()

	//wg.Wait()

	//mergeCSVFiles("output/03/24") // DOES NOT WORK, USE https://merge-csv.com/

	//filterMergedOnlyWithARecords()
	//filterMergedOnlyWithMXRecords()

	//populateWithWhoisData()

	processDomCopCSV()

	//data, err := domainAnalysis.QueryWhois("qvineox.ru")
	//if err != nil {
	//	slog.Warn(err.Error())
	//	data.ToCSV()
	//}
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

func processRGSTXT() {
	inputFile, err := os.Open("../../data/rgs/domains.txt")
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

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_rgs.csv", dir, now.Unix()))
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

		ip := net.ParseIP(line)
		if ip != nil {
			continue
		}

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
		time.Sleep(300 * time.Millisecond)
	}
}

func processDomCopCSV() {
	inputFile, err := os.Open("../../data/domcop/top10milliondomains.csv")
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

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_domcop.csv", dir, now.Unix()))
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

	_, _ = csvReader.Read() // skip header

	var maxIndex = 160000
	var index = 0

	for index < maxIndex {
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		index++
		name := line[1]

		rec, err := record.NewDomainRecord(name)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		rec.IsLegit = true

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

func mergeCSVFiles(dir string) {
	readDir, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("failed to read directory: " + err.Error())
		panic(err)
	}

	now := time.Now()

	outputDir := fmt.Sprintf("output/%s", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/merged/%d_merged.csv", outputDir, now.Unix()))
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

	var unique = make(map[string]bool)

	for _, entry := range readDir {
		if entry.IsDir() {
			continue
		}

		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			slog.Error("failed to open file: " + err.Error())
			panic(err)
		}

		slog.Info("reading file: " + entry.Name())

		csvReader := csv.NewReader(f)
		_, _ = csvReader.Read() // skip header

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

			err = csvWriter.Write(line)
			if err != nil {
				slog.Error(err.Error())
			}
		}

		csvWriter.Flush()
	}
}

func filterMergedOnlyWithARecords() {
	inputFile, err := os.Open("../../data/merged/merged_full_2024_03_24.csv")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var now = time.Now()

	dir := fmt.Sprintf("output/%s/filtered", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_with_A.csv", dir, now.Unix()))
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
	_, _ = csvReader.Read() // skip header

	for {
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		hasIPs := line[12]
		if hasIPs == "0" || hasIPs == "-1" {
			continue
		}

		err = csvWriter.Write(line)
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func filterMergedOnlyWithMXRecords() {
	inputFile, err := os.Open("../../data/merged/merged_full_2024_03_24.csv")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var now = time.Now()

	dir := fmt.Sprintf("output/%s/filtered", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_with_MX.csv", dir, now.Unix()))
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
	_, _ = csvReader.Read() // skip header

	for {
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		hasIPs := line[13]
		if hasIPs == "0" || hasIPs == "-1" {
			continue
		}

		err = csvWriter.Write(line)
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

func populateWithWhoisData() {
	inputFile, err := os.Open("../../data/merged/merged_full_2024_03_24.csv")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var now = time.Now()

	dir := fmt.Sprintf("output/%s/whois", now.Format("01/02"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_full_2024_03_24.csv", dir, now.Unix()))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	csvWriter := csv.NewWriter(outputFile)

	header := record.CSVHeader()
	header = append(header, domainAnalysis.CSVHeader()...)

	err = csvWriter.Write(header)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	csvReader := csv.NewReader(inputFile)
	_, _ = csvReader.Read() // skip header

	var index = 0

	for {
		index++
		line, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		data, err := domainAnalysis.QueryWhois(index, line[0])
		if err != nil {
			slog.Warn(err.Error())
			line = append(line, make([]string, 11)...)
		} else {
			line = append(line, data.ToCSV()...)
		}

		err = csvWriter.Write(line)
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
		time.Sleep(300 * time.Millisecond)
	}
}
