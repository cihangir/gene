// Package testdata provides test data to the other packages
package testdata

//  JSON1 holds a primitive json
const JSON1 = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Message",
  "description": "Message represents a simple post",
  "type": "object",
  "properties": {
    "Id": {
      "description": "The unique identifier for a message",
      "type": "number",
      "format":"int64"
    },
    "Token": {
      "description": "The token for a message security",
      "type": "string"
    },
    "Body": {
      "description": "The body for a message",
      "type": "string",
      "pattern": "^(/[^/]+)+$",
      "minLength": 2,
      "maxLength": 3
    },
    "Age": {
      "type": "integer",
      "minimum": 0,
      "maximum": 100,
      "exclusiveMaximum": true
    },
    "Enabled": {
      "type": "boolean"
    },
    "StatusConstant": {
      "type": "string",
      "enum": [
        "active",
        "deleted"
      ]
    },
    "CreatedAt": {
      "type": "string",
      "format":"date-time"
    }
  },
  "required": [
    "id",
    "body"
  ]
}
`

// JSONWithModule holds a json with module support
const JSONWithModule = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "id": "http://savas.io/account",
  "type": "object",
  "additionalProperties": true,
  "title": "Account",
  "description": "Account module handles all the operations regarding account management.",
  "properties": {},
  "definitions": {
    "Address": {
      "id": "http: //savas.io/account/address",
      "type": "object",
      "additionalProperties": true,
      "title": "Address",
      "description": "Address holds the address of an account.",
      "properties": {
        "Street": {
          "id": "http: //savas.io/account/address/street",
          "type": "string",
          "minLength": 0,
          "title": "Street",
          "description": "Street name holds the name of the street for an address.",
          "default": "2nd Street"
        },
        "City": {
          "id": "http: //savas.io/account/address/city",
          "type": "string",
          "minLength": 0,
          "title": "City",
          "description": "City holds the name of the city for the address.",
          "default": "Manisa"
        }
      }
    },
    "PhoneNumber": {
      "id": "http: //savas.io/account/address/phoneNumber",
      "type": "object",
      "minItems": 1,
      "uniqueItems": false,
      "title": "PhoneNumber",
      "description": "Phone number holds a general data for a phone number.",
      "properties": {
        "Location": {
          "id": "http: //savas.io/account/address/phoneNumber/location",
          "type": "string",
          "minLength": 0,
          "title": "Location",
          "description": "Location holds the location data for the phone number.",
          "name": "location",
          "default": "home"
        },
        "Code": {
          "id": "http: //savas.io/account/address/phoneNumber/code",
          "type": "integer",
          "multipleOf": 1,
          "maximum": 100,
          "minimum": 1,
          "exclusiveMaximum": false,
          "exclusiveMinimum": false,
          "title": "Code",
          "description": "Code holds the area code for the phone number",
          "default": 44
        }
      }
    }
  }
}`
