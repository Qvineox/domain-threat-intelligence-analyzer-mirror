package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ammario/ipisp"
	"github.com/miekg/dns"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// ref: https://dev.to/ductnn/simple-tool-check-domain-information-with-golang-54j8
// ref: https://www.golangprograms.com/find-dns-records-programmatically.html

var whoIsClient ipisp.Client
var dnsClient ipisp.Client

func main() {
	resolver := prepareResolver("8.8.8.8:53")

	inputFile, err := os.Open("input/input_dns.csv")
	if err != nil {
		slog.Error(err.Error())
	}

	timestamp := time.Now().Unix()

	outputFileCSV, err := os.Create(fmt.Sprintf("output/%d_output_dns.csv", timestamp))
	if err != nil {
		slog.Error(err.Error())
	}

	//outputFileTXT, err := os.Create(fmt.Sprintf("output/%d_output_dns_raw.txt", timestamp))
	//if err != nil {
	//	slog.Error(err.Error())
	//}

	whoIsClient, err = ipisp.NewWhoisClient()
	dnsClient, err = ipisp.NewDNSClient()

	var resolved = make(map[string]bool)

	csvWriter := csv.NewWriter(outputFileCSV)
	err = csvWriter.Write(csvHeader())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	csvReader := csv.NewReader(inputFile)
	_, _ = csvReader.Read() // skip header
	var index = 0

	for {
		index++

		l, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		r := record{
			cidr: l[0],
			name: l[3],
		}

		_, ok := resolved[r.name]
		if ok {
			continue
		}
		resolved[r.name] = false

		slog.Info(fmt.Sprintf("#%d \t | scanning %s", index, r.name))

		// get top level domain
		parts := strings.Split(r.name, ".")
		r.tld = parts[len(parts)-1]

		//dnsData, err := queryDNS(r.name)
		//if err == nil && len(dnsData) > 0 {
		//	_, _ = outputFileTXT.Write([]byte(dnsData))
		//} else {
		//	_, _ = outputFileTXT.Write([]byte("-"))
		//}

		result := make([]string, 8)
		result[0] = r.name

		// getting A records
		r.ips, _ = resolver.LookupIP(context.Background(), "ip", r.name)

		// query whois data
		r.asnCount, r.regionCount, r.allocationMonths, r.registryCount, err = queryWhoIS(r.ips)
		if err != nil {
			slog.Error(err.Error())
		}

		mxs, _ := resolver.LookupMX(context.Background(), r.name)
		r.mxCount = len(mxs)

		cNames, _ := resolver.LookupCNAME(context.Background(), r.name)
		r.cnameCount = len(cNames)

		var matchedPtrRecords = 0

		for _, ip := range r.ips {
			a, err := resolver.LookupAddr(context.Background(), ip.String())
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			r.ptrCount += len(a)

			for _, ptr := range a {
				if ptr[:len(ptr)-1] == r.name {
					matchedPtrRecords++
				}
			}
		}

		if r.ptrCount > 0 {
			r.ptrRatio = matchedPtrRecords / r.ptrCount
		} else {
			r.ptrRatio = -1
		}

		err = csvWriter.Write(r.toCSV())
		if err != nil {
			slog.Error(err.Error())
		}

		csvWriter.Flush()
		time.Sleep(300 * time.Millisecond)
	}

	_ = inputFile.Close()
	_ = outputFileCSV.Close()
	//_ = outputFileTXT.Close()
}

type record struct {
	cidr     string
	category uint8
	score    uint8
	name     string
	DNSwlID  uint64

	tld string

	ips []net.IP

	cnameCount int
	mxCount    int

	ptrCount int
	ptrRatio int

	registryCount int

	asnCount         int
	regionCount      int
	allocationMonths int
}

func csvHeader() []string {
	return []string{"domain", "tld", "a_count", "mx_count", "cname_count", "ptr_count", "ptr_ratio", "registry_count", "asn_count", "region_count", "allocation_duration_months"}
}

func (r record) toCSV() []string {
	return []string{
		r.name,
		r.tld,
		strconv.Itoa(len(r.ips)),
		strconv.Itoa(r.mxCount),
		strconv.Itoa(r.cnameCount),
		strconv.Itoa(r.ptrCount),
		strconv.Itoa(r.ptrRatio),
		strconv.Itoa(r.registryCount),
		strconv.Itoa(r.asnCount),
		strconv.Itoa(r.regionCount),
		strconv.Itoa(r.allocationMonths),
	}
}

//func queryDomain(domain string) (ASNs, countries string, err error) {
//
//
//	c := make(map[string]bool)
//	a := make(map[int]bool)
//
//	for _, v := range data {
//		a[v.ASN] = false
//		c[v.Country] = false
//	}
//
//	return strconv.Itoa(len(a)), strconv.Itoa(len(c)), err
//}

func prepareResolver(dnsServer string) *net.Resolver {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, dnsServer)
		},
	}

	return r
}

func queryWhoIS(ips []net.IP) (asns, countries, allocMonths, registries int, err error) {
	if len(ips) == 0 {
		return 0, 0, 0, 0, errors.New("no ips")
	}

	response, err := whoIsClient.LookupIPs(ips)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var oldestAllocation time.Time

	c := make(map[string]bool)
	a := make(map[string]bool)
	r := make(map[string]bool)

	for _, v := range response {
		a[v.ASN.String()] = false
		c[v.Country] = false
		r[v.Registry] = false

		if oldestAllocation.IsZero() || v.AllocatedAt.Before(oldestAllocation) {
			oldestAllocation = v.AllocatedAt
		}
	}

	difference := oldestAllocation.Sub(time.Now())

	return len(a), len(c), -int(difference.Hours() / 24 / 30), len(r), nil
}

func queryDNS(domain string) (string, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{Name: domain + ".", Qtype: dns.TypeA}

	c := new(dns.Client)
	msg, _, err := c.Exchange(m1, "8.8.8.8:53")
	if err != nil {
		return "", err
	}

	return msg.String(), nil
}
