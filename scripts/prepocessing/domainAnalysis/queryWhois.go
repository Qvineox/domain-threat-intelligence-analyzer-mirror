package domainAnalysis

import (
	"fmt"
	"github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
	"log/slog"
)

func QueryWhois(index int, name string) (IPData, error) {
	slog.Info(fmt.Sprintf("#%d \t | querying whois... \t | %s", index, name))

	raw, err := whois.Whois(name)
	if err != nil {
		slog.Warn("failed to query whois info: " + err.Error())
		return IPData{}, err
	}

	info, err := parser.Parse(raw)
	if err != nil {
		slog.Warn("failed to parse whois info: " + err.Error())
		return IPData{}, err
	}

	rec := NewIPDataFromWhoisData(name, info)
	return rec, nil
}
