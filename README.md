# Frala

## About

Frala is fragment and translation engine.

**Features:**

- Frala enables you to separate out files into "fragments" and import / re-use them elsewhere.
- Frala enables your content to be multi-lingual through the use of "Terms". With Terms, you can specify the term name, languages and the language's value for the word.
- Frala interfaces with the standarized gettext PO files for translations, enabling you to convert to and from PO files / Frala Terms.

**Current Version:** `0.3`

**License:** `Apache 2.0`

**TODO:**

- [ ] Implement CLI tool that leverages the core Frala package to enable usage of functionality without needing to implement the package
 - *though really, you should -- totally no bias /s*.

### Syntax

Frala does all this by using a syntax similar to JSON and Mustache, in that it leverages `{{` and `}}`.

**Example Fragment:** In the example below, we are importing a Fragment. Note that we are using an HTML file in this example, but any file type works.

`{{ type="fragment" src="file.html" }}`

**Example Term:** In the example below, we are importing a Term. In our config, we have a term specified for "hello" in the Finnish `fi` language.

`{{ type="term" lang="fi" src="hello" }}`

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
    "Languages" : ["en", "fi"],
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

## Usage: HTML

### Fragments

You can declare the importing of a Fragment by using the Frala Fragment syntax anywhere in an HTML file, including other Fragments. If you specify only the file, without a path, it will look in the directory of the same file that the Fragment it is.

#### Examples

In the example below, we are importing `info.html` from the current directory and `innerfragment.html` from `innerdir`.

``` html
<div>
    {{type="fragment" src="info.html"}}
    {{type="fragment" src="innerdir/innerfragment.html"}}
</div>
```

In `innerfragment.html`, the content imports `innerdirref.html`. Because of our use of relative URLs, we are actually importing the file from the same directory as `innerfragment.html`

``` html
{{ type="fragment" src="innerdirref.html" }}
```

### Terms

You can declare the use of a Term by using the Frala Term Syntax anywhere in an HTML file, including Fragments. Terms can be specified without a language, meaning it will use the `DefaultLanguage` from the Config, or if specified with a language, it will attempt to use the value from that Term's language key / val.

#### Examples

In the example below, we are importing the "hello" Term with English (default language) and Finnish.

``` html
<div>
    This is our message in English: {{type="term" src="hello" }}, Josh!
    This is our message in Finnish: {{type="term" src="hello" lang="fi"}}, Josh!
</div>
```

## Usage: Go Package

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
type ConfigOptions struct {
	DefaultLanguage string          // Default Language string, if not declared, default to en
	Languages       []string        // Languages is a list of languages (string) -- Optional
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

Term is a `map[string]string`, as each Term has a map of language -> value (where language is a string and value is a string)

``` go
type Term map[string]string
```

### Variables

``` go
// Config is the configuration of FralaConfig
var Config ConfigOptions

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

#### Po Conversion

##### ConvertFromPo

ConvertFromPo reads a .po file and convert its content to Frala Terms, automatically adding them to the config.

``` go
func ConvertFromPo(fileName string) error
```

##### ConvertToPo

ConvertToPo converts Frala Terms into msgid / msgstr context for usage
in a .po file. It returns the po file content.

``` go
func ConvertToPo(language string) string
```

#### Terms

Really these are all helper functions. They aren't necessary to use, but handy if desired.

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