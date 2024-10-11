# Privy Raffle Joiner

## About This Actor

This Actor is a powerful, user-fiendly tool made to sign-up emails on a (Privy.com)[https://privy.com] raffle. This tool will save you time on joining large lists of emails on raffles.

Made with Golang 1.22.1

## Tutorial

Basic Usage

```json
{
    "buisnessID": "123123",
    "campaignID": "123123",
    "formName": "Test Form",
    "displayID": "1231212",

    "useCatchAll": true,
    "catchAllEmail": "@awesome.com",
    "catchAllLimit": 10,

    "emails": ["test@awesome.com"]
}
```

| parameter | type | argument | description |
| --------- | ----- | ------------------------- | ---------------------------- |
| buisnessID | string | _123123_ | The Privy Raffle Buisness ID |
| campaignID | string | _123123_ | The Privy Raffle Capaign ID |
| formName | string | _123123_ | The Privy Raffle Form Name |
| displayID | string | _123123_ | The Privy Raffle Display ID |
| useCatchAll | bool | _default=true_ | Wether or not to use catch alls |
| catchAllEmail | bool | _default=true_ | The catch-all email to use |
| catchAllLimit | bool | _default=true_ | Limit the amount of signups on catch-all |
| emails | array | _["test@awesome.com", ...]_ | An array of Emails (leave blank if using catch-all) |

### Output Sample

```json
[
  {
    "email": "test@awesome.com"
  }
]
```
