package dhl

// Omit empty is required so the DHL api doesn't get confused by empty strings etc when serializing the response data
// This should never be used with boolean fields since 'false' values will also be ommited, wtf google? https://themirgleich.medium.com/how-golangs-omitempty-can-confuse-the-google-apis-c0e46d22d9ff

type Draft struct {
	Id                 string              `json:"id,omitempty"`
	ShipmentId         string              `json:"shipmentId,omitempty"`
	OrderReference     string              `json:"orderReference,omitempty"`
	Receiver           Contact             `json:"receiver,omitempty"`
	BccEmails          []string            `json:"bccEmails,omitempty"`
	Shipper            Shipper             `json:"shipper,omitempty"`
	AccountId          string              `json:"accountId,omitempty"`
	Options            []Option            `json:"options,omitempty"`
	OnBehalfOf         *Shipper            `json:"onBehalfOf,omitempty"`
	Product            string              `json:"product,omitempty"`
	CustomsDeclaration *CustomsDeclaration `json:"customsDeclaration,omitempty"`
	ReturnLabel        bool                `json:"returnLabel"`
	Pieces             []Piece             `json:"pieces,omitempty"`
	Metadata           *Metadata           `json:"metadata,omitempty"`
	// DeliveryArea       DeliveryArea       `json:"deliveryArea,omitempty"`
}

type Contact struct {
	Name        *Name    `json:"name,omitempty"`
	Address     *Address `json:"address,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phoneNumber,omitempty"`
	VatNumber   string   `json:"vatNumber,omitempty"`
	EoriNumber  string   `json:"eoriNumber,omitempty"`
}

type Shipper struct {
	Name        *Name    `json:"name,omitempty"`
	Address     *Address `json:"address,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phoneNumber,omitempty"`
	VatNumber   string   `json:"vatNumber,omitempty"`
	EoriNumber  string   `json:"eoriNumber,omitempty"`
	RexNumber   string   `json:"rexNumber,omitempty"`
}

type Name struct {
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	CompanyName    string `json:"companyName,omitempty"`
	AdditionalName string `json:"additionalName,omitempty"`
}

type Address struct {
	CountryCode           string `json:"countryCode,omitempty"`
	PostalCode            string `json:"postalCode,omitempty"`
	City                  string `json:"city,omitempty"`
	Street                string `json:"street,omitempty"`
	AdditionalAddressLine string `json:"additionalAddressLine,omitempty"`
	Number                string `json:"number,omitempty"`
	IsBusiness            bool   `json:"isBusiness"`
	Addition              string `json:"addition,omitempty"`
}

type Option struct {
	Key   string `json:"key,omitempty"`
	Input string `json:"input,omitempty"`
}

type Label struct {
	LabelId        string `json:"labelId,omitempty"`
	OrderReference string `json:"orderReference,omitempty"`
	ParcelType     string `json:"parcelType,omitempty"`
	LabelType      string `json:"labelType,omitempty"`
	PieceNumber    int    `json:"pieceNumber,omitempty"`
	TrackerCode    string `json:"trackerCode,omitempty"`
	RoutingCode    string `json:"routingCode,omitempty"`
	UserId         string `json:"userId,omitempty"`
	OrganisationId string `json:"organisationId,omitempty"`
	Application    string `json:"application,omitempty"`
	TimeCreated    string `json:"timeCreated,omitempty"`
	ShipmentId     string `json:"shipmentId,omitempty"`
	AccountNumber  string `json:"accountNumber,omitempty"`
}

type CustomsDeclaration struct {
	CertificateNumber          string        `json:"certificateNumber,omitempty"`
	Currency                   string        `json:"currency,omitempty"`
	InvoiceNumber              string        `json:"invoiceNumber,omitempty"`
	LicenceNumber              string        `json:"licenceNumber,omitempty"`
	Remarks                    string        `json:"remarks,omitempty"`
	InvoiceType                string        `json:"invoiceType,omitempty"`
	ExportType                 string        `json:"exportType,omitempty"`
	ExportReason               string        `json:"exportReason,omitempty"`
	CustomsGoods               []CustomsGood `json:"customsGoods,omitempty"`
	IncoTerms                  string        `json:"incoTerms,omitempty"`
	IncoTermsCity              string        `json:"incoTermsCity,omitempty"`
	SenderInboundVatNumber     string        `json:"senderInboundVatNumber,omitempty"`
	AttachmentIds              []string      `json:"attachmentIds,omitempty"`
	ShippingFee                ShippingFee   `json:"shippingFee,omitempty"`
	ImporterOfRecord           Importer      `json:"importerOfRecord,omitempty"`
	DefermentAccountVat        string        `json:"defermentAccountVat,omitempty"`
	DefermentAccountDuties     string        `json:"defermentAccountDuties,omitempty"`
	VatReverseCharge           bool          `json:"vatReverseCharge"`
	SenderHasInboundEoriNumber bool          `json:"senderHasInboundEoriNumber"`
}

type CustomsGood struct {
	Code        string  `json:"code,omitempty"`
	Description string  `json:"description,omitempty"`
	Origin      string  `json:"origin,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
}

type ShippingFee struct {
	Value float64 `json:"value,omitempty"`
}

type Importer struct {
	Name        Name     `json:"name,omitempty"`
	Address     *Address `json:"address,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phoneNumber,omitempty"`
	VatNumber   string   `json:"vatNumber,omitempty"`
	EoriNumber  string   `json:"eoriNumber,omitempty"`
}

type Piece struct {
	ParcelType string      `json:"parcelType,omitempty"`
	Quantity   int         `json:"quantity,omitempty"`
	Weight     int         `json:"weight,omitempty"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
}

type Dimensions struct {
	Length float64 `json:"length,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

type DeliveryArea struct {
	Type   string `json:"type,omitempty"`
	Remote bool   `json:"remote"`
}

type Metadata struct {
	ImportId       string `json:"importId,omitempty"`
	UserId         string `json:"userId,omitempty"`
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
	OrganizationId string `json:"organizationId,omitempty"`
	TimeCreated    string `json:"timeCreated,omitempty"`
}
