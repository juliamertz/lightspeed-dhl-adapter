package dhl

import (
	"lightspeed-dhl/config"
	"lightspeed-dhl/lightspeed"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func WebhookToDraft(incoming lightspeed.IncomingOrder, conf *config.Secrets) Draft {
	return Draft{
		Id:             uuid.New().String(),
		ShipmentId:     uuid.New().String(),
		OrderReference: strconv.Itoa(incoming.Order.Id),
		Receiver: Contact{
			Email:       incoming.Order.Email,
			PhoneNumber: incoming.Order.Phone,
			Name: &Name{
				FirstName:      incoming.Order.Firstname,
				LastName:       incoming.Order.Lastname,
				AdditionalName: incoming.Order.Middlename,
				CompanyName:    incoming.Order.CompanyName,
			},
			Address: &Address{
				IsBusiness:  incoming.Order.IsCompany,
				Street:      incoming.Order.AddressShippingStreet,
				City:        incoming.Order.AddressShippingCity,
				PostalCode:  strings.ReplaceAll(incoming.Order.AddressShippingZipcode, " ", ""),
				CountryCode: incoming.Order.AddressShippingCountry.Code,
				Number:      incoming.Order.AddressShippingNumber,
				Addition:    incoming.Order.AddressShippingExtension,
			},
		},

		Options: []Option{
			{Key: "REFERENCE", Input: incoming.Order.Number},
			{Key: "PERS_NOTE", Input: *conf.CompanyInfo.PersonalNote},
		},
		Pieces: []Piece{
			{
				ParcelType: "MEDIUM",
				// Lightspeed reports it's weight in kilo's while DHL expects grams
				Weight: incoming.Order.Weight / 1000,
			},
		},

		Shipper: Shipper{
			Name: &Name{
				CompanyName: conf.CompanyInfo.Name,
			},
			Address: &Address{
				IsBusiness:  true,
				Street:      conf.CompanyInfo.Street,
				Number:      conf.CompanyInfo.Number,
				Addition:    conf.CompanyInfo.Addition,
				PostalCode:  conf.CompanyInfo.PostalCode,
				City:        conf.CompanyInfo.City,
				CountryCode: conf.CompanyInfo.CountryCode,
			},
			Email:       conf.CompanyInfo.Email,
			PhoneNumber: conf.CompanyInfo.PhoneNumber,
		},
		AccountId: conf.Dhl.AccountId,
	}
}
