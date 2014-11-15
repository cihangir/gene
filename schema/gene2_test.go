package schema

const testJSON1 = `
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
