package dhl

type Draft struct {
	Id                 string             `json:"id"`
	ShipmentId         string             `json:"shipmentId"`
	OrderReference     string             `json:"orderReference"`
	Receiver           Contact            `json:"receiver"`
	BccEmails          []string           `json:"bccEmails"`
	Shipper            Shipper            `json:"shipper"`
	AccountId          string             `json:"accountId"`
	Options            []Option           `json:"options"`
	OnBehalfOf         Shipper            `json:"onBehalfOf"`
	Product            string             `json:"product"`
	CustomsDeclaration CustomsDeclaration `json:"customsDeclaration"`
	ReturnLabel        bool               `json:"returnLabel"`
	Pieces             []Piece            `json:"pieces"`
	DeliveryArea       DeliveryArea       `json:"deliveryArea"`
	Metadata           Metadata           `json:"metadata"`
}

type Contact struct {
	Name        Name    `json:"name"`
	Address     Address `json:"address"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	VatNumber   string  `json:"vatNumber"`
	EoriNumber  string  `json:"eoriNumber"`
}

type Shipper struct {
	Name        Name    `json:"name"`
	Address     Address `json:"address"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	VatNumber   string  `json:"vatNumber"`
	EoriNumber  string  `json:"eoriNumber"`
	RexNumber   string  `json:"rexNumber"`
}

type Name struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	CompanyName    string `json:"companyName"`
	AdditionalName string `json:"additionalName"`
}

type Address struct {
	CountryCode           string `json:"countryCode"`
	PostalCode            string `json:"postalCode"`
	City                  string `json:"city"`
	Street                string `json:"street"`
	AdditionalAddressLine string `json:"additionalAddressLine"`
	Number                string `json:"number"`
	IsBusiness            bool   `json:"isBusiness"`
	Addition              string `json:"addition"`
}

type Option struct {
	Key   string `json:"key"`
	Input string `json:"input"`
}

type CustomsDeclaration struct {
	CertificateNumber          string        `json:"certificateNumber"`
	Currency                   string        `json:"currency"`
	InvoiceNumber              string        `json:"invoiceNumber"`
	LicenceNumber              string        `json:"licenceNumber"`
	Remarks                    string        `json:"remarks"`
	InvoiceType                string        `json:"invoiceType"`
	ExportType                 string        `json:"exportType"`
	ExportReason               string        `json:"exportReason"`
	CustomsGoods               []CustomsGood `json:"customsGoods"`
	IncoTerms                  string        `json:"incoTerms"`
	IncoTermsCity              string        `json:"incoTermsCity"`
	SenderInboundVatNumber     string        `json:"senderInboundVatNumber"`
	AttachmentIds              []string      `json:"attachmentIds"`
	ShippingFee                ShippingFee   `json:"shippingFee"`
	ImporterOfRecord           Importer      `json:"importerOfRecord"`
	DefermentAccountVat        string        `json:"defermentAccountVat"`
	DefermentAccountDuties     string        `json:"defermentAccountDuties"`
	VatReverseCharge           bool          `json:"vatReverseCharge"`
	SenderHasInboundEoriNumber bool          `json:"senderHasInboundEoriNumber"`
}

type CustomsGood struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Origin      string  `json:"origin"`
	Quantity    int     `json:"quantity"`
	Value       float64 `json:"value"`
	Weight      float64 `json:"weight"`
}

type ShippingFee struct {
	Value float64 `json:"value"`
}

type Importer struct {
	Name        Name    `json:"name"`
	Address     Address `json:"address"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	VatNumber   string  `json:"vatNumber"`
	EoriNumber  string  `json:"eoriNumber"`
}

type Piece struct {
	ParcelType string      `json:"parcelType"`
	Quantity   int32       `json:"quantity"`
	Weight     float64     `json:"weight"`
	Dimensions *Dimensions `json:"dimensions"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type DeliveryArea struct {
	Type   string `json:"type"`
	Remote bool   `json:"remote"`
}

type Metadata struct {
	ImportId       string `json:"importId"`
	UserId         string `json:"userId"`
	LastModifiedBy string `json:"lastModifiedBy"`
	OrganizationId string `json:"organizationId"`
	TimeCreated    string `json:"timeCreated"`
}
