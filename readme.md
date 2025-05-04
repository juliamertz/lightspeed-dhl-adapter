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
