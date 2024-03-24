package domainAnalysis

import (
	whoisparser "github.com/likexian/whois-parser"
	"strconv"
	"time"
)

type IPData struct {
	Name string
	IP   string

	// domain data
	DNSSec      bool
	Domain      string
	Extension   string
	WhoisServer string
	ID          string

	CreatedDate, ExpirationDate *time.Time

	// registrar data
	RegistrarCountry      string
	RegistrarOrganization string
	RegistrarPostalCode   string
	RegistrarName         string

	// registrant data
	RegistrantCountry      string
	RegistrantOrganization string
	RegistrantPostalCode   string

	// admin data
	AdminName         string
	AdminCountry      string
	AdminOrganization string
	AdminPostalCode   string

	// billing data
	BillingCountry      string
	BillingOrganization string
	BillingPostalCode   string

	//
}

func NewIPDataFromWhoisData(domain string, whoisData whoisparser.WhoisInfo) (data IPData) {
	if whoisData.Administrative != nil {
		data.AdminName = whoisData.Administrative.Name
		data.AdminCountry = whoisData.Administrative.Country
		data.AdminOrganization = whoisData.Administrative.Organization
		data.AdminPostalCode = whoisData.Administrative.PostalCode
	}

	if whoisData.Registrant != nil {
		data.RegistrantCountry = whoisData.Registrant.Country
		data.RegistrantOrganization = whoisData.Registrant.Organization
		data.RegistrantPostalCode = whoisData.Registrant.PostalCode
	}

	if whoisData.Registrar != nil {
		data.RegistrarName = whoisData.Registrar.Name
		data.RegistrarCountry = whoisData.Registrar.Country
		data.RegistrarOrganization = whoisData.Registrar.Organization
		data.RegistrarPostalCode = whoisData.Registrar.PostalCode
	}

	if whoisData.Domain != nil {
		data.ID = whoisData.Domain.ID
		data.Domain = whoisData.Domain.Domain
		data.Extension = whoisData.Domain.Extension
		data.WhoisServer = whoisData.Domain.WhoisServer

		data.CreatedDate = whoisData.Domain.CreatedDateInTime
		data.ExpirationDate = whoisData.Domain.ExpirationDateInTime

		data.DNSSec = whoisData.Domain.DNSSec
	}

	if whoisData.Billing != nil {
		data.BillingCountry = whoisData.Billing.Country
		data.BillingOrganization = whoisData.Billing.Organization
		data.BillingPostalCode = whoisData.Billing.PostalCode
	}

	return data
}

func CSVHeader() []string {
	return []string{
		"dnssec_active",
		"whois_id",
		"created_at",
		"expires_at",
		"admin_name",
		"admin_country",
		"admin_org",
		"registrar_name",
		"registrar_country",
		"registrar_org",
		"registrant_country",
		"registrant_org",
	}
}

func (data IPData) ToCSV() []string {
	var dnsSecActive = "0"
	if data.DNSSec {
		dnsSecActive = "1"
	}

	if data.CreatedDate == nil || data.ExpirationDate == nil {
		return []string{
			dnsSecActive,
			data.ID,
			"",
			"",
			data.AdminName,
			data.AdminCountry,
			data.AdminOrganization,
			data.RegistrarName,
			data.RegistrarCountry,
			data.RegistrarOrganization,
			data.RegistrantCountry,
			data.RegistrantOrganization,
		}
	}

	return []string{
		dnsSecActive,
		data.ID,
		strconv.FormatInt(data.CreatedDate.Unix(), 10),
		strconv.FormatInt(data.ExpirationDate.Unix(), 10),
		data.AdminName,
		data.AdminCountry,
		data.AdminOrganization,
		data.RegistrarName,
		data.RegistrarCountry,
		data.RegistrarOrganization,
		data.RegistrantCountry,
		data.RegistrantOrganization,
	}
}
