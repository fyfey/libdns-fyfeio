// Package libdns_fyfeio implements a DNS record management client compatible
// with the libdns interfaces for fyfe.io.
package libdns_fyfeio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/libdns/libdns"
)

type Record struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type RecordWrapper struct {
	Name   string  `json:"name"`
	TTL    uint32  `json:"ttl"`
	Record *Record `json:"record"`
}

type AppendRecordsRequest struct {
	Records []*RecordWrapper `json:"records"`
}

type AppendRecordsResponse struct {
	Message string `json:"message"`
}

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with <TODO: PROVIDER NAME>.
type Provider struct {
	// TODO: put config fields here (with snake_case json
	// struct tags on exported fields), for example:
	APIToken string `json:"api_token,omitempty"`
}

func callAPI(zone string, record libdns.Record, action string) error {
	if action != "upsert" && action != "delete" {
		return errors.New("invalid action. Expected upsert or delete.")
	}

	client := &http.Client{}

	fyfeRecord := &Record{
		"TXT", record.Value,
	}
	wrapper := &RecordWrapper{
		Name:   record.Name,
		Record: fyfeRecord,
		TTL:    uint32(record.TTL.Seconds()),
	}
	request := &AppendRecordsRequest{[]*RecordWrapper{wrapper}}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:3000/zone/%s", zone), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	res, err := client.Do(req)

	response := &AppendRecordsResponse{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(response)
	if err != nil {
		return err
	}

	return nil
}

func appendRecord(zone string, record libdns.Record) error {
	err := callAPI(zone, record, "append")
	if err != nil {
		return err
	}
	return nil
}
func deleteRecord(zone string, record libdns.Record) error {
	err := callAPI(zone, record, "delete")
	if err != nil {
		return err
	}
	return nil
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {

	for _, record := range records {
		appendRecord(zone, record)
	}

	return nil, fmt.Errorf("TODO: not implemented")
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
