# go-alexa-lambda
Alexa golang library to generate skill + interaction model as well as serve requests with lambda or as a server.
The packages can also be used standalone as a request/response abstraction for Alexa requests.

# Purpose

## Projects using `go-alexa-lambda`
* [alexa-go-cloudformation-demo](https://github.com/DrPsychick/alexa-go-cloudformation-demo) : the demo project that lead to developing this library. A fully automated build and deploy of an Alexa skill including lambda function via Cloudformation. 

# Usage
## Build skill and interaction model
```go
package main

import(
	"log"
	"github.com/drpsychick/alexa-go-lambda/l10n"
	"github.com/drpsychick/alexa-go-lambda/skill"
)

// Configure locale registry.
var enUS = &l10n.Locale{
    Name: "en-US",
    TextSnippets: map[string][]string{
        l10n.KeySkillName:   {"This is my awesome skill."},
        l10n.KeyLaunchTitle: {"Welcome!"},
        // ... for each locale, define the required keys
    },
}

func main() {
    reg := l10n.NewRegistry()
    reg.Register(enUS)
    
    // Create and configure skill builder.
    s := skill.NewSkillBuilder().
        WithLocaleRegistry(reg).
        WithCategory(skill.CategoryGames).
        WithPrivacyFlag(skill.FlagIsExportCompliant, true)
    
    // Configure model.
    m := s.Model().
        WithDelegationStrategy(skill.DelegationSkillResponse)
    
    // Build `skill.json`
    sj, err := s.Build()
    if err != nil {
        log.Fatal(err)
    }
    res, _ := json.MarshalIndent(sj, "", "  ")
        log.Printf("Skill:\n%s", res)
    
    // Build interaction models
    ms := s.BuildModels()
    for l, m := range ms {
        log.Printf("Locale %s:\n%s", l, m)
    }
}
```

## Respond to alexa requests with lambda
```go

```

### Links
* https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html
* multiple intents in one dialog: https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html#pass-a-new-intent
* https://developer.amazon.com/blogs/alexa/post/cfbd2f5e-c72f-4b03-8040-8628bbca204c/alexa-skill-teardown-understanding-entity-resolution-with-pet-match

### Credits
initially inspired by: https://github.com/soloworks/go-alexa-models