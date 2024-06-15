## Idee

Nieuwe lightspeed order  -> Draft aanmaken in DHL
DHL label geprint -> Order in lightspeed op status "verzonden" zetten

## Todo

- [ ] Test potential rate limiting by DHL when fetching labels every poll
- [ ] Test polling logic with dummy order
- [ ] Use refresh token instead of authenticating every request
- [x] Add error handling/logging
- [x] Add DHL polling logic
  - [x] When a webhook is received, add order to database.
  - [x] Poll DHL for all drafts, if any entry in the database is missing here we can assume it's been shipped.
  - [x] Update the order status in lightspeed
