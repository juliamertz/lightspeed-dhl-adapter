use uuid::Uuid;

use crate::{config, dhl, lightspeed};

impl From<lightspeed::IncomingOrder> for dhl::Draft {
    fn from(incoming: lightspeed::IncomingOrder) -> Self {
        let conf = config::get();
        dhl::Draft {
            id: Uuid::new_v4().to_string(),
            shipment_id: Uuid::new_v4().to_string(),
            order_reference: incoming.order.id.to_string(),

            shipper: Some(dhl::Shipper {
                name: Some(dhl::Name {
                    company_name: Some(conf.company_info.name.clone()),
                    first_name: None,
                    last_name: None,
                    additional_name: None,
                }),
                address: Some(dhl::Address {
                    is_business: true,
                    street: Some(conf.company_info.street.clone()),
                    number: Some(conf.company_info.number.clone()),
                    addition: Some(conf.company_info.addition.clone()),
                    postal_code: Some(conf.company_info.postal_code.clone()),
                    city: Some(conf.company_info.city.clone()),
                    country_code: Some(conf.company_info.country_code.clone()),
                    additional_address_line: None,
                }),
                email: Some(conf.company_info.email.clone()),
                phone_number: Some(conf.company_info.phone_number.clone()),
                vat_number: None,
                rex_number: None,
                eori_number: None,
            }),

            receiver: Some(dhl::Contact {
                email: Some(incoming.order.email),
                phone_number: Some(incoming.order.phone),
                name: Some(dhl::Name {
                    first_name: Some(incoming.order.firstname),
                    last_name: Some(incoming.order.lastname),
                    additional_name: Some(incoming.order.middlename),
                    company_name: Some(incoming.order.company_name),
                }),
                address: Some(dhl::Address {
                    is_business: incoming.order.is_company,
                    street: Some(incoming.order.address_shipping_street),
                    city: Some(incoming.order.address_shipping_city),
                    postal_code: Some(incoming.order.address_shipping_zipcode.replace(" ", "")),
                    country_code: Some(incoming.order.address_shipping_country.code),
                    number: Some(incoming.order.address_shipping_number),
                    addition: Some(incoming.order.address_shipping_extension),
                    additional_address_line: None,
                }),
                vat_number: None,
                eori_number: None,
            }),

            options: vec![
                dhl::OptionField {
                    key: "REFERENCE".to_string(),
                    input: incoming.order.number,
                },
                dhl::OptionField {
                    key: "PERS_NOTE".to_string(),
                    input: conf.company_info.personal_note.clone(),
                },
            ],

            pieces: vec![dhl::Piece {
                parcel_type: Some(dhl::ParcelType::Medium),
                weight: Some(incoming.order.weight / 1000),
                quantity: None,
                dimensions: None,
            }],

            account_id: Some(conf.dhl.account_id.clone()),

            return_label: false,
            bcc_emails: None,
            product: None,
            metadata: None,
            on_behalf_of: None,
            customs_declaration: None,
        }
    }
}
