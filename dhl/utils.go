package dhl

import (
	"lightspeed-dhl/config"
)

func ShipperFromConfig(d config.CompanyInfo) Shipper {
	return Shipper{
		Name: Name{
			CompanyName: d.Name,
		},
		Address: Address{
			IsBusiness:  true,
			Street:      d.Street,
			Number:      d.Number,
			Addition:    d.Addition,
			PostalCode:  d.PostalCode,
			City:        d.City,
			CountryCode: d.CountryCode,
		},
		Email:       d.Email,
		PhoneNumber: d.PhoneNumber,
	}
}
