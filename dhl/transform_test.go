package dhl_test

import (
	// "regexp"

	"lightspeed-dhl/config"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/lightspeed"
	"strings"
	"testing"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestTranslation(t *testing.T) {
	order := lightspeed.Order{
		Id:             12345,
		Email:          "john.doe@example.com",
		Firstname:      "John",
		Lastname:       "Doe",
		Middlename:     "A.",
		CompanyName:    "Doe Inc.",
		Phone:          "+1234567890",
		ShipmentTitle:  "Standard Shipping",
		Number:         "ORD-98765",
		IsCompany:      true,
		Status:         "processing_awaiting_shipment",
		ShipmentStatus: "not_shipped",

		AddressBillingStreet:    "123 Main St",
		AddressBillingCity:      "New York",
		AddressBillingZipcode:   "10001",
		AddressBillingCountry:   lightspeed.CountryCode{Code: "US"},
		AddressBillingNumber:    "12",
		AddressBillingExtension: "B",

		AddressShippingStreet:    "456 Elm St",
		AddressShippingCity:      "Los Angeles",
		AddressShippingZipcode:   "5050 AJ",
		AddressShippingCountry:   lightspeed.CountryCode{Code: "US"},
		AddressShippingNumber:    "34",
		AddressShippingExtension: "C",
	}

	incoming := lightspeed.IncomingOrder{
		Order: order,
	}

	note := ""
	conf := &config.Secrets{
		CompanyInfo: config.CompanyInfo{PersonalNote: &note},
	}

	draft, err := dhl.WebhookToDraft(incoming, conf)
	check(err, t)

	if draft.Receiver.Address.Street != "456 Elm St" {
		t.Errorf("Receiver street should be set to AdressShippingStreet")
	}
	if draft.Receiver.Address.City != "Los Angeles" {
		t.Errorf("Receiver city should be set to AdressShippingCity")
	}
	if draft.OrderReference != "12345" {
		t.Errorf("Order reference should match reference: %v, expected 12345", draft.OrderReference)
	}

	if strings.Contains(draft.Receiver.Address.PostalCode, " ") {
		t.Errorf("Postal code shouldn't contain any whitespace")
	}
}
