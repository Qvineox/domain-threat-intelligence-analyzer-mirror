package main

import (
	"encoding/csv"
	"ensembled/models"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	inputFile, err := os.Open("../../data/merged/merged_full_2024_03_27.csv")
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

	outputFile, err := os.Create(fmt.Sprintf("%s/%d_ensembled.csv", dir, now.Unix()))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	csvWriter := csv.NewWriter(outputFile)
	err = csvWriter.Write([]string{"domain", "dga_score", "semantics_score", "records_score", "is_legit"})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fmt.Print(models.ScoreResourceRecords([]float64{2, 0, 21, 1, 3, 0.0, 1}))

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

		semanticValues := make([]float64, 0)
		for _, v := range line[2:12] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}

			semanticValues = append(semanticValues, f)
		}

		semanticScore := models.ScoreSemantics(semanticValues)

		recordValues := make([]float64, 0)
		for _, v := range line[11:18] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}

			recordValues = append(recordValues, f)
		}

		recordScore := models.ScoreResourceRecords(recordValues)

		err = csvWriter.Write(
			[]string{line[0],
				"-1",
				strconv.FormatFloat(semanticScore, 'f', 4, 32),
				strconv.FormatFloat(recordScore, 'f', 4, 32),
				line[18],
			})

		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
	}
}

//func main() {
//	params := []float64{2, 1, 13, 5, 2, 1.0, 1} // --> 0.982995
//	// moves [5.000000 13.000000 1.000000 1.000000 1.000000 2.000000 2.000000]
//
//	fmt.Print(models.ScoreResourceRecords([]float64{params[3], params[2], params[1], params[5], params[6], params[0], params[4]}))
//
//	//params := []float64{2.0, 0.0, 21.0, 1.0, 3.0, 0.0, 1.0} // --> 0.458081
//	// moves [0.000000 1.000000 3.000000 0.000000 1.000000 21.000000 2.000000]
//
//	cs := combin.Permutations(len(params), 7)
//	for _, c := range cs {
//		//fmt.Printf("%f-%f-%f-%f-%f-%f-%f\n", params[c[0]], params[c[1]], params[c[2]], params[c[3]], params[c[4]], params[c[5]], params[c[6]])
//
//		params_ := []float64{params[c[0]], params[c[1]], params[c[2]], params[c[3]], params[c[4]], params[c[5]], params[c[6]]}
//
//		prediction := models.ScoreResourceRecords(params_)
//
//		if 0.982 < prediction && prediction < 0.983 {
//			fmt.Printf("%f : %f\n", params_, prediction)
//		}
//	}
//
//	//prediction := models.ScoreResourceRecords(params)
//	//
//	//if 0.458 < prediction && prediction < 0.459 {
//	//	fmt.Println(params)
//	//	panic("found!")
//	//}
//	//
//	//fmt.Println(prediction)
//
//	// [2 2 5 13 1 2]
//
//	fmt.Print(models.ScoreResourceRecords(params))
//}
