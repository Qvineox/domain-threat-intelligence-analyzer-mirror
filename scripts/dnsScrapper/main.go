package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ammario/ipisp"
	"github.com/miekg/dns"
	"github.com/shlin168/go-whois/whois"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"slices"
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
	//prepareWhoisClient()

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

	var recordsBuffer []record

	for {
		index++

		l, err := csvReader.Read()
		if err != nil {
			slog.Error(err.Error())
			break
		}

		r := record{
			cidr:             l[0],
			name:             l[3],
			asnCount:         -1,
			regionCount:      -1,
			allocationMonths: -1,
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
		//r.asnCount, r.regionCount, r.allocationMonths, r.registryCount, err = queryWhoIS(r.ips)
		//if err != nil {
		//	slog.Error(err.Error())
		//}

		//queryWhoISDomain(r.name)

		mxs, _ := resolver.LookupMX(context.Background(), r.name)
		r.mxCount = len(mxs)

		cNames, _ := resolver.LookupCNAME(context.Background(), r.name)
		r.cnameCount = len(cNames)

		txts, _ := resolver.LookupTXT(context.Background(), r.name)
		r.txtCount = len(txts)

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

		recordsBuffer = append(recordsBuffer, r)
		if len(recordsBuffer) >= 1000 {
			slog.Info("starting whois query for next 100 records...")

			var recordsToWrite []record

			recordsToWrite, err = queryRecords(recordsBuffer)
			if err != nil {
				slog.Error(err.Error())
				recordsToWrite = recordsBuffer
			}

			for _, rec := range recordsToWrite {
				err = csvWriter.Write(rec.toCSV())
				if err != nil {
					slog.Warn(err.Error())
				}
			}

			csvWriter.Flush()
			recordsBuffer = []record{}
		}

		time.Sleep((100 + time.Duration(rand.Int63n(500))) * time.Millisecond)
	}

	if len(recordsBuffer) > 0 {
		var recordsToWrite []record

		recordsToWrite, err = queryRecords(recordsBuffer)
		if err != nil {
			slog.Error(err.Error())
			recordsToWrite = recordsBuffer
		}

		for _, rec := range recordsToWrite {
			err = csvWriter.Write(rec.toCSV())
			if err != nil {
				slog.Warn(err.Error())
			}
		}
	}

	csvWriter.Flush()

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
	txtCount   int

	ptrCount int
	ptrRatio int

	asnCount         int
	regionCount      int
	allocationMonths int
}

func csvHeader() []string {
	return []string{"domain", "tld", "a_count", "mx_count", "cname_count", "txt_count", "ptr_count", "ptr_ratio", "asn_count", "region_count", "allocation_duration_months", "is_legit"}
}

func (r record) toCSV() []string {
	return []string{
		r.name,
		r.tld,
		strconv.Itoa(len(r.ips)),
		strconv.Itoa(r.mxCount),
		strconv.Itoa(r.cnameCount),
		strconv.Itoa(r.txtCount),
		strconv.Itoa(r.ptrCount),
		strconv.Itoa(r.ptrRatio),
		strconv.Itoa(r.asnCount),
		strconv.Itoa(r.regionCount),
		strconv.Itoa(r.allocationMonths),
		"1",
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

func queryRecords(records []record) ([]record, error) {
	var ipMap = make(map[string][3]string) // map[ip][[ASNs], [countries], [allocationDuration]]
	var domainMap = make(map[string]*domainParcel)

	var ips = make([]net.IP, 0)

	for _, r := range records {
		for _, ip := range r.ips {
			ips = append(ips, ip)
		}

		domainMap[r.name] = &domainParcel{
			record:      r,
			asns:        []string{},
			countries:   []string{},
			allocatedAt: time.Now(),
			ips:         r.ips,
		}
	}

	retiesLeft := 3
	var response []ipisp.Response
	var err error

	for retiesLeft > 0 {
		response, err = whoIsClient.LookupIPs(ips)
		if err == nil {
			break
		}

		retiesLeft--
	}

	if err != nil {
		return nil, err
	}

	for _, res := range response {
		//ip_ := res.IP.String()

		ipMap[res.IP.String()] = [3]string{
			fmt.Sprintf("%d", res.ASN),
			res.Country,
			fmt.Sprintf("%d", res.AllocatedAt.Unix()),
		}

		//v, ok := ipMap[ip_]
		//if ok {
		//	v[0] = append(v[0], fmt.Sprintf("%d", res.ASN))
		//	v[1] = append(v[1], res.Country)
		//	v[2] = append(v[2], fmt.Sprintf("%d", res.AllocatedAt.Unix()))
		//
		//	ipMap[ip_] = v
		//} else {
		//	ipMap[ip_] = [3][]string{
		//		{fmt.Sprintf("%d", res.ASN)},
		//		{res.Country},
		//		{fmt.Sprintf("%d", res.AllocatedAt.Unix())},
		//	}
		//}
	}

	for _, r := range domainMap {
		for _, ip := range r.ips {

			v, ok := ipMap[ip.String()]
			if ok {
				r.asns = append(r.asns, v[0])
				r.countries = append(r.countries, v[1])

				timestamp, _ := strconv.Atoi(v[2])

				t := time.Unix(int64(timestamp), 0)

				if r.allocatedAt.After(t) {
					r.allocatedAt = t
				}
			}

		}
	}

	var result = make([]record, 0)

	for _, r := range records {
		v, ok := domainMap[r.name]
		if ok {
			slices.Sort(v.countries)
			v.countries = slices.Compact[[]string, string](v.countries)

			slices.Sort(v.asns)
			v.asns = slices.Compact[[]string, string](v.asns)

			r.asnCount = len(v.asns)
			r.regionCount = len(v.countries)

			result = append(result, r)
		}
	}

	return result, err
}

type domainParcel struct {
	record      record
	asns        []string
	countries   []string
	allocatedAt time.Time
	ips         []net.IP
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

var whoisClient *whois.Client

func prepareWhoisClient() {
	var err error
	whoisClient, err = whois.NewClient()
	if err != nil {
		panic(err)
	}
}

func queryWhoISDomain(domain string) {
	ctx := context.Background()

	whoisDomain, err := whoisClient.Query(ctx, domain)
	if err == nil {
		slog.Info("queried whois from: " + whoisDomain.WhoisServer)

		fmt.Println("rawtext:", whoisDomain.RawText)
		fmt.Printf("parsed whois: %+v\n", whoisDomain.ParsedWhois)

		if whoisDomain.IsAvailable != nil {
			fmt.Println("available:", *whoisDomain.IsAvailable)
		}
	}
}

func queryWhoISIP(ip string) {
	ctx := context.Background()

	whoisIP, err := whoisClient.QueryIP(ctx, ip)
	if err == nil {
		slog.Info("queried whois from: " + whoisIP.WhoisServer)

		fmt.Printf("parsed whois: %+v\n", whoisIP.ParsedWhois)
	}
}
