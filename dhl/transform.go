package dhl

import (
	"fmt"
	"jorismertz/lightspeed-dhl/config"
	"jorismertz/lightspeed-dhl/lightspeed"

	"github.com/google/uuid"
)

func WebhookToDraft(incoming lightspeed.IncomingOrder) Draft {
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		panic(err)
	}

	orderId := fmt.Sprint(incoming.Order.Id)
	return Draft{
		Id:             uuid.New().String(),
		ShipmentId:     uuid.New().String(),
		OrderReference: &orderId,
		Receiver: &Contact{
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
				PostalCode:  incoming.Order.AddressShippingZipcode,
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
			{ParcelType: "MEDIUM"},
		},

		Shipper:   ShipperFromConfig(conf.CompanyInfo),
		AccountId: conf.Dhl.AccountId,
	}
}
