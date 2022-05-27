## Nano DB
A lightweight DB in golang.

## Features
- RESTful: `{GET|POST|DELETE} /{DB}/{TABLE}/{ID}`
- Cache: LRU
- Permission: ACL
- Lightweight: can run on a Raspberry Pi Zero

## Usage
### Get cookie
`./nano-db -c {userName}`
then the cookie will print to stdout.

### Control DB
The first user who access the {DB} will be the {DB}'s admin.
#### Get data
`GET https://example.com/novel/books/1`

#### Add/Update
`POST https://example.com/novel/books/1`
with JSON body

#### Delete
`DELETE https://example.com/novel/books/1`

#### Status
`GET https://example.com/status`

## Attention
`{DB}/{TABLE}/{ID}` cant contains `/`,` `,`\\`,`..`, and its length cant be larger than 50.

