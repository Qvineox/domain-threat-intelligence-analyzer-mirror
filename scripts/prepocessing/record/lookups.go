package record

import (
	"context"
	"log/slog"
	"net"
	"strings"
	"time"
)

var resolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(5000),
		}
		return d.DialContext(ctx, network, "8.8.8.8:53")
	},
}

type LookupsData struct {
	IPs []net.IP

	// records
	MXs      int
	CNAMEs   int
	TXTs     int
	PTRs     int
	PTRRatio float64
}

func LookupRecords(name string) (l LookupsData, err error) {
	slog.Info("looking up... \t | " + name)

	l.IPs, _ = resolver.LookupIP(context.Background(), "ip", name)

	// query whois data
	//r.asnCount, r.regionCount, r.allocationMonths, r.registryCount, err = queryWhoIS(r.ips)
	//if err != nil {
	//	slog.Error(err.Error())
	//}

	//queryWhoISDomain(r.name)

	mxs, _ := resolver.LookupMX(context.Background(), name)
	l.MXs = len(mxs)

	cNames, _ := resolver.LookupCNAME(context.Background(), name)
	l.CNAMEs = len(cNames)

	txts, _ := resolver.LookupTXT(context.Background(), name)
	l.TXTs = len(txts)

	var matchedPtrRecords = 0

	for _, ip := range l.IPs {
		a, err := resolver.LookupAddr(context.Background(), ip.String())
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		l.PTRs += len(a)

		for _, ptr := range a {
			if strings.Contains(ptr[:len(ptr)-1], name) {
				matchedPtrRecords++
			}
		}
	}

	if l.PTRs > 0 {
		l.PTRRatio = float64(matchedPtrRecords) / float64(l.PTRs)
	} else {
		l.PTRRatio = -1
	}

	return l, nil
}
