package lightspeed

type IncomingOrder struct {
	Order Order `json:"order"`
}

// There is a lot more data coming in, but for now this is all we need to parse.
// https://developers.lightspeedhq.com/ecom/endpoints/order
type Order struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Middlename    string `json:"middlename"`
	CompanyName   string `json:"companyName"`
	Phone         string `json:"phone"`
	ShipmentTitle string `json:"shipmentTitle"`
	Number        string `json:"number"`
	IsCompany     bool   `json:"isCompany"`

	// "cancelled" / "completed_shipped" / "processing_awaiting_shipment"
	Status string `json:"status"`
	// "shipped" / "not_shipped"
	ShipmentStatus string `json:"shipmentStatus"`

	// PaymentStatus  string `json:"paymentStatus"`  // "paid" / "not_paid"

	// Combined weight of all products included in order
	Weight int `json:"weight"`

	AddressBillingStreet    string      `json:"addressBillingStreet"`
	AddressBillingCity      string      `json:"addressBillingCity"`
	AddressBillingZipcode   string      `json:"addressBillingZipcode"`
	AddressBillingCountry   CountryCode `json:"addressBillingCountry"`
	AddressBillingNumber    string      `json:"addressBillingNumber"`
	AddressBillingExtension string      `json:"addressBillingExtension"`

	AddressShippingStreet    string      `json:"addressShippingStreet"`
	AddressShippingCity      string      `json:"addressShippingCity"`
	AddressShippingZipcode   string      `json:"addressShippingZipcode"`
	AddressShippingCountry   CountryCode `json:"addressShippingCountry"`
	AddressShippingNumber    string      `json:"addressShippingNumber"`
	AddressShippingExtension string      `json:"addressShippingExtension"`
}

type CountryCode struct {
	Code string `json:"code"`
}

type Variant struct {
	Id         int `json:"id"`
	StockLevel int `json:"stockLevel"`
	StockAlert int `json:"stockAlert"`
}

type Product struct {
	Id       int                `json:"id"`
	Title    string             `json:"title"`
	Variants map[string]Variant `json:"variants"`
}

type CatalogResponse struct {
	Products []Product `json:"products"`
}
