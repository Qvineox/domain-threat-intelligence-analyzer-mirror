package record

import (
	"strconv"
	"strings"
)

type DomainRecord struct {
	FullName    string `json:"full_name"`
	TLD         string `json:"tld"`
	LevelsCount int    `json:"levels_count"`

	LevelsMAD float64 `json:"levels_mad"` // https://en.wikipedia.org/wiki/Median_absolute_deviation

	SymbolsCount    int `json:"symbols_count"`
	VowelsCount     int `json:"vowels_count"`
	ConsonantsCount int `json:"consonants_count"`
	NumbersCount    int `json:"numbers_count"`
	PointsCount     int `json:"points_count"`
	SpecialCount    int `json:"special_count"` // i.e. '-' and '_'

	VowelsRatio     float64 `json:"vowels_ratio"`     // vowels to total symbols
	ConsonantsRatio float64 `json:"consonants_ratio"` // consonants to total symbols
	NumbersRatio    float64 `json:"numbers_ratio"`    // numbers to total symbols
	PointsRatio     float64 `json:"points_ratio"`     // points to total symbols
	SpecialRatio    float64 `json:"special_ratio"`    // special symbols to total symbols

	UniqueRatio float64 `json:"unique_ratio"` // unique symbols to total symbols

	MaxRepeatedSymbols int `json:"max_repeated_symbols"`

	Lookups LookupsData `json:"lookups"`

	ICMPResponse bool `json:"icmp_response"`

	IsLegit bool `json:"is_legit"`
}

func NewDomainRecord(record string) (r DomainRecord, err error) {
	var fullName = record

	fullName = strings.ToLower(fullName)
	fullName = strings.Trim(fullName, ".")

	r.FullName = fullName

	r.SymbolsCount, r.VowelsCount, r.ConsonantsCount, r.NumbersCount, r.SpecialCount, r.PointsCount, r.VowelsRatio, r.ConsonantsRatio, r.NumbersRatio, r.SpecialRatio, r.PointsRatio, err = countSymbols(fullName)
	if err != nil {
		return DomainRecord{}, err
	}

	r.MaxRepeatedSymbols = maxRepeatedSymbolCount(fullName)
	r.UniqueRatio = uniqueSymbolsRation(fullName)

	parts := strings.Split(fullName, ".")

	r.TLD = parts[len(parts)-1]
	r.LevelsCount = len(parts)
	r.LevelsMAD, err = levelsMAD(parts)
	if err != nil {
		return DomainRecord{}, err
	}

	return r, nil
}

func CSVHeader() []string {
	return []string{
		"domain",
		"tld",
		"levels_count",
		"levels_mad",
		"symbols_count",
		"vowels_ratio",
		"consonants_ratio",
		"numbers_ratio",
		"points_ratio",
		"special_ratio",
		"unique_ratio",
		"max_repeated",
		"a_records",
		"mx_records",
		"cname_records",
		"txt_records",
		"ptr_records",
		"ptr_ratio",
		//"icmp_response",
		"is_legit",
	}
}

//func NewRecordFromCSV(line []string) DomainRecord {
//	return DomainRecord{
//		FullName:           line[0],
//		TLD:                line[1],
//		LevelsCount:        strconv.Atoi(line[2]),
//		LevelsMAD:          line[3],
//		SymbolsCount:       line[4],
//		VowelsRatio:        line[5],
//		ConsonantsRatio:    line[6],
//		NumbersRatio:       line[7],
//		PointsRatio:        line[8],
//		SpecialRatio:       line[9],
//		UniqueRatio:        line[10],
//		MaxRepeatedSymbols: line[11],
//		Lookups: LookupsData{
//			IPs:      line[12],
//			MXs:      line[13],
//			CNAMEs:   line[14],
//			TXTs:     line[15],
//			PTRs:     line[16],
//			PTRRatio: line[17],
//		},
//		ICMPResponse: false,
//		IsLegit:      line[18],
//	}
//}

func (r DomainRecord) ToCSV() []string {
	var isLegit = "0"
	if r.IsLegit {
		isLegit = "1"
	}

	//var icmp = "0"
	//if r.ICMPResponse {
	//	icmp = "1"
	//}

	return []string{
		r.FullName,
		r.TLD,
		strconv.Itoa(r.LevelsCount),
		strconv.FormatFloat(r.LevelsMAD, 'f', 4, 32),
		strconv.Itoa(r.SymbolsCount),
		strconv.FormatFloat(r.VowelsRatio, 'f', 4, 32),
		strconv.FormatFloat(r.ConsonantsRatio, 'f', 4, 32),
		strconv.FormatFloat(r.NumbersRatio, 'f', 4, 32),
		strconv.FormatFloat(r.PointsRatio, 'f', 4, 32),
		strconv.FormatFloat(r.SpecialRatio, 'f', 4, 32),
		strconv.FormatFloat(r.UniqueRatio, 'f', 4, 32),
		strconv.Itoa(r.MaxRepeatedSymbols),
		strconv.Itoa(len(r.Lookups.IPs)),
		strconv.Itoa(r.Lookups.MXs),
		strconv.Itoa(r.Lookups.CNAMEs),
		strconv.Itoa(r.Lookups.TXTs),
		strconv.Itoa(r.Lookups.PTRs),
		strconv.FormatFloat(r.Lookups.PTRRatio, 'f', 4, 32),
		//icmp,
		isLegit,
	}
}
