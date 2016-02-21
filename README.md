# Frala

## About

Frala is fragment and translation engine. Frala enables you to separate out files into "fragments" and import / re-use them elsewhere. Frala also enables your content to be multi-lingual through the use of "terms". With terms, you can specify the term name, languages and the language's value for the word.

**License:** `Apache 2.0`

TODO:

- [ ] Implement Term JSON to .po and .po to Term JSON converters.

### Syntax

Frala does all this by using a syntax similar to JSON and Mustache, in that it leverages `{` and `}`.

**Example Fragment:** In the example below, we are importing a Fragment. Note that we are using an HTML file in this example, but any file type works.

`{ type="fragment" src="file.html" }`

**Example Term:** In the example below, we are importing a Term. In our config, we have a term specified for "hello" in the Finnish `fi` language.

`{ type="term" lang="fi" src="hello" }`

### Config

Configuring Frala is simple.

1. We will automatically read `frala.json` from the directory (if it exists).
2. Using Fragments requires no configuration at all and you only need to specify what you use.
3. You can specify a default language and thus eliminate the need to pass `lang="the-language"` for a Term if you want it for the value of the term for the same language as your default.
4. Specify Terms.

Note: We will default to using `en` if not default language is specified in the config.

**Example Config:**

``` json
{
    "DefaultLanguage" : "en",
    "Terms" : {
        "good_morning" : {
            "en" : "good morning",
            "fi" : "hyvä huomenta"
        }
    }
}
```

## Contribute

This project leverages CodeUtils for development and adopts the CodeUtils Usage Spec. To learn how to contribute to this project and set up CodeUtils, read
[here](https://github.com/StroblIndustries/CodeUtils/blob/master/CodeUtils-Usage-Spec.md).

## Usage

### Import

You can use Frala in your Go software via: `import "github.com/JoshStrobl/frala"`

### Structs

#### Context

Context is a struct that has properties relating to the type and type's associated information. Used by the line parser.

``` go
type Context struct {
    Lang   string `json:"lang"` // Language of the term (if not a fragment)
    Source string `json:"src"`  // Source such as the word or link to fragment
    Type   string `json:"type"` // Type of the Context (fragment or term)
}
```

#### FralaConfig

This is the configuration for Frala.

``` go
type FralaConfig struct {
	DefaultLanguage string          // Default Language string, if not declared, default to en
	Languages       []string        // Languages is a list of languages (string)
	Terms           map[string]Term // Terms is a map of strings (term names) to individual Terms
}
```

#### ParseResponse

ParseResponse is a struct that contains both the content of a file and associated parsing error

``` go
type ParseResponse struct {
    Content string // Content of the parsed file
    Error   error  // Error that occurred during parsing
}
```
#### Term

Term is a map[string]string, as each Term has a map of language -> value (where language is a string and value is a string)

``` go
type Term map[string]string
```

### Variables

``` go
// Define CurrentParsingFile as the file we're currently parsing
var CurrentParsingFile string
// Define InitError as any potential error from initializing Frala
var InitError error
```

### Functions

#### Config

##### ReadConfig

This function will read any frala.json file and update the Config.

``` go
func ReadConfig() error
```

##### SaveConfig

This function will save the Config to frala.json.

``` go
func SaveConfig() error
```

#### Parsing

##### MultiParse

This function will parse multiple files provided and return a map of ParseResponses

``` go
func MultiParse(files []string) map[string]ParseResponse
```

##### Parse

This function will parse a file provided and return a ParseResponse.

``` go
func Parse(file string) ParseResponse
```

##### ParseLine

This function will parse an individual line.

``` go
func ParseLine(lineContent string) string
```

##### ParseSyntax

This function will parse a Frala syntax string and return the appropriate (if any) associated HTML content or term

``` go
func ParseSyntax(fralaSyntax string) string
```

#### Terms

##### DeleteTerm

This function will delete a Term from Terms.

``` go
func DeleteTerm(termName string)
```

##### DeleteValue

This function will delete a language / value from a Term

``` go
func DeleteValue(termName, language string)
```

##### GetValue

This function will get the value of a language from a Term, if it exists

``` go
func GetValue(termName, language string) string
```

##### SetTerm

This function will enable you to set a Term to Terms.

``` go
func SetTerm(termName string)
```

##### SetValue

This function will enable you to set the value of a Term language.

``` go
func SetValue(termName, language, value string)
```