// Code generated by "genSchema --private Conditions"; DO NOT EDIT.

package condition

// conditionsJSONSchema is the JSON schema for the Conditions struct.
var conditionsJSONSchema = `{
  "type": "object",
  "properties": {
    "ifAncestorFiles": {"type": "array", "description": "IfAncestorFiles is a list of files to search for in the project folder, or another folder higher up in the directory structure.", "items": {"type": "string", "description": ""}},
    "ifFiles": {"type": "array", "description": "IfFiles is a list of files to search for in the project folder.", "items": {"type": "string", "description": ""}},
    "ifExtensions": {"type": "array", "description": "IfExtensions is a list of extensions to search for in the project folder.", "items": {"type": "string", "description": ""}},
    "ifOS": {"type": "array", "description": "IfOS is a list of operating systems.  If the current GOOS is not in the list, then this project type is not matched.", "items": {"type": "string", "description": ""}},
    "ifNotOS": {"type": "array", "description": "IfNotOS is a list of operating systems.  If the current GOOS is in the list, then this project type is not matched.", "items": {"type": "string", "description": ""}}
  },
  "additionalProperties": false
}`
