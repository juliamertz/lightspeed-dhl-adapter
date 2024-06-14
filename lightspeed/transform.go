package lightspeed

import (
	"fmt"
	"jorismertz/lightspeed-dhl/config"
	"jorismertz/lightspeed-dhl/dhl"

	"github.com/google/uuid"
)

func WebhookToDraft(incoming IncomingOrder) dhl.Draft {
	conf, err := config.LoadSecrets("config.toml")
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
			{Key: "PERS_NOTE", Input: *conf.CompanyInfo.PersonalNote},
		},
		Pieces: []dhl.Piece{
			{ParcelType: "MEDIUM"},
		},

		Shipper:   dhl.ShipperFromConfig(conf.CompanyInfo),
		AccountId: conf.Dhl.AccountId,
	}
}
