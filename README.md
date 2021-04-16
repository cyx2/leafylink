# Leafylink

## Introduction
Leafylink is a URL shortener written in Golang, utilizing MongoDB Atlas.  It provides a basic HTML template-based UI, and a RESTful POST method to insert new mappings.

A mapping is defined as a relationship between a `longUrl`, such as https://www.google.com/, and a Leafylink, such as https://leafylink.herokuapp.com/a6b7f8.

This project was started as part of Skunkworks 2021.

## Configuration
Leafylink requires either a configuration file, `.env`, or serverless environment configuration.  Required configuration variables are as follows.
* `DB_PASSWORD`: (string)
* `DB_URL`: (string)
* `DB_USER`: (string)
* `ENV`: (string) [LOCAL / DEV / PROD]
* `APP_URL`: (string)
* `PORT`: (int)

The application uses a single Atlas cluster, and the `ENV` variable specifies the Atlas `db` used by the application.  Within each `db`, the `collection` is a fixed value, `mappings`.

## Usage
Leafylink is deployed via Heroku.  Locally, it can be run without arguments (as long as a `.env` file exists in the project root directory), via `go run .`

There are two methods of usage:
* Application root URL
* HTTP POST API

The POST API can be reached via `/api/create`, and called with the following syntax
```json
{
    "LongUrl": "https://leafylink.herokuapp.com/"
}
```

## Database
Leafylink stores mappings in MongoDB Atlas documents in the following format.  In addition to the `key` to `longUrl` relationship (`redirect` in the database), a few more pieces of metadata are stored, such as the `createdate` and `usecount`.  Each time a Leafylink is used, the `usecount` is incremented.

```json
{
    "_id":{"$oid":"607860c8eed834a3b952437e"},
    "createdate":{"$date":{"$numberLong":"1618501832614"}},
    "key":"a6b7f8",
    "redirect":"https://leafylink.herokuapp.com/",
    "leafyurl":"https://leafylink.herokuapp.com/a6b7f8",
    "usecount":{"$numberInt":"1"}
}
```