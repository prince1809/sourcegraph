// Code generated by stringdata. DO NOT EDIT.

package schema

// PhabricatorSchemaJSON is the content of the file "phabricator.schema.json".
const PhabricatorSchemaJSON = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "phabricator.schema.json#",
  "title": "PhabricatorConnection",
  "description": "Configuration for a connection to Phabricator.",
  "type": "object",
  "additionalProperties": false,
  "anyOf": [{ "required": ["token"] }, { "required": ["repos"] }]
}
`