# Nano DB
A lightweight DB in golang.

# Features
- RESTful: `{GET|POST|DELETE} /{DB}/{TABLE}/{ID}`
- Cache: LRU
- Permission: ACL

# Usage
## Get cookie
`./nano-db -c {userName}`
then the cookie will print to stdout.

## Get data
eg: `GET https://example.com/novel/books/1`

## Add/Update
eg: `POST https://example.com/novel/books/1`
with JSON body

## Delete
eg: `DELETE https://example.com/novel/books/1`

## Status
`GET https://example.com/status`

# Attention
`{DB}/{TABLE}/{ID}` cant contains `/`,` `,`\\`,`..`, and its length cant be larger than 50.

