## Adapter for Lightspeed Retail and DHL

This service listens for order creation webhooks from Lightspeed Retail and automatically creates corresponding shipment drafts in DHL.

It also periodically polls DHL to check if any drafts have been promoted to “label” status. When this happens, the associated order status in Lightspeed is updated accordingly.

> [!NOTE]
> If you don't want your metrics to be publicly accesible, make sure to block the `/metrics` route in your reverse-proxy

## Example configuration
```toml
[Options]
DryRun      = false # Don't mutate any data.
Port        = 3000 # Where to start web server
Environment = "production" # development | production
Debug       = false # Turn on debug logging
PollingInterval = 5 # Interval to check for shipped orders, in minutes

[DHL]
UserId    = "uuid"
ApiKey    = "uuid"
AccountId = "0700000"

[Lightspeed]
Cluster = "https://api.webshopapp.com/en/"
Key     = ""
Secret  = ""
ShopId = ""
ClusterId = ""

[CompanyInfo]
Name         = ""
Street       = ""
City         = ""
PostalCode   = ""
CountryCode  = ""
Number       = ""
Addition     = ""
Email        = ""
PhoneNumber  = ""
PersonalNote = "Your order is on it's way!"
```
