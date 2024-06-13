package lightspeed

import (
	"fmt"
	"jorismertz/lightspeed-dhl/dhl"
	"jorismertz/lightspeed-dhl/secrets"

	"github.com/google/uuid"
)

func WebhookToDraft(incoming IncomingOrder) dhl.Draft {
	config, err := secrets.LoadSecrets("secrets.toml")
	if err != nil {
		panic(err)
	}

	orderId := fmt.Sprint(incoming.Order.Id)
	return dhl.Draft{
		Id:             uuid.New().String(),
		ShipmentId:     uuid.New().String(),
		OrderReference: &orderId,
		Receiver: &dhl.Contact{
			Email:       incoming.Order.Email,
			PhoneNumber: incoming.Order.Phone,
			Name: &dhl.Name{
				FirstName:      incoming.Order.Firstname,
				LastName:       incoming.Order.Lastname,
				AdditionalName: incoming.Order.Middlename,
				CompanyName:    incoming.Order.CompanyName,
			},
			Address: &dhl.Address{
				IsBusiness:  incoming.Order.IsCompany,
				Street:      incoming.Order.AddressShippingStreet,
				City:        incoming.Order.AddressShippingCity,
				PostalCode:  incoming.Order.AddressShippingZipcode,
				CountryCode: incoming.Order.AddressShippingCountry.Code,
				Number:      incoming.Order.AddressShippingNumber,
				Addition:    incoming.Order.AddressShippingExtension,
			},
		},

		Options: []dhl.Option{
			{Key: "REFERENCE", Input: incoming.Order.Number},
			{Key: "PERS_NOTE", Input: "Uw bestelling bij nettenshop.nl is met DHL onderweg! Via de bijgevoegde link kunt u uw pakket volgen. Mocht u vragen hebben, neem dan contact met ons op via de klantenservice. Met vriendelijke groet, Team Nettenshop.nl"},
		},
		Pieces: []dhl.Piece{
			{ParcelType: "MEDIUM"},
		},

		Shipper:   dhl.ShipperFromConfig(config.CompanyInfo),
		AccountId: config.Dhl.AccountId,
	}
}
