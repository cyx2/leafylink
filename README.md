# Leafylink

## Introduction
Leafylink is a URL shortener written in Golang, utilizing MongoDB Atlas.  It provides a basic HTML template-based UI, and a RESTful POST method to insert new mappings.

A mapping is defined as a relationship between a longUrl, such as https://www.google.com/, and a Leafylink, such as https://leafylink.herokuapp.com/a6b7f8.

This project was started as part of Skunkworks 2021.

## Configuration
Leafylink requires either a configuration file, `.env`, or serverless environment configuration.  Required configuration variables are as follows.
* DB_PASSWORD: (string)
* DB_URL: (string)
* DB_USER: (string)
* ENV: (string) [LOCAL / DEV / PROD]
* APP_URL: (string)
* PORT

## Usage
There are two methods of usage:
* Application root URL
* HTTP POST API

The POST API can be reached via `/api/create`, and called with the following syntax
```json
{
    "LongUrl": "https://leafylink.herokuapp.com/"
}
```