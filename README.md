# go-alexa-lambda
Alexa golang library to generate skill + interaction model as well as serve requests with lambda or as a server.
The packages can also be used standalone as a request/response abstraction for Alexa requests.

# Purpose



# Usage
## Build skill and interaction model
Run it on go playground: https://play.golang.org/p/VfHW4RcUVwn
```go
package main

import (
	"log"

	alexa "github.com/drpsychick/go-alexa-lambda"
	"github.com/drpsychick/go-alexa-lambda/l10n"
	"github.com/drpsychick/go-alexa-lambda/skill"
	jsoniter "github.com/json-iterator/go"
)

// Configure locale registry.
var enUS = &l10n.Locale{
	Name: "en-US",
	TextSnippets: map[string][]string{
		l10n.KeySkillName:   {"This is my awesome skill."},
		l10n.KeySkillDescription: {"Description"},
		l10n.KeySkillSummary: {"Skill summary"},
		l10n.KeySkillSmallIconURI: {"https://my-url.com/small.png"},
		l10n.KeySkillLargeIconURI: {"https://my-url.com/large.png"},
		l10n.KeySkillTestingInstructions: {"Testing instructions"},
		l10n.KeySkillInvocation: {"awesome skill"},
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
		WithPrivacyFlag(skill.FlagIsExportCompliant, true).
		WithModel()

	// Configure model.
	s.Model().
		WithDelegationStrategy(skill.DelegationSkillResponse).
		WithIntent(alexa.StopIntent)

	// Build `skill.json`
	sj, err := s.Build()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := jsoniter.MarshalIndent(sj, "", "  ")
	log.Printf("Skill:\n%s", res)

	// Build interaction model: `en-US.json`
	ms, err := s.BuildModels()
	if err != nil {
		log.Fatal(err)
	}
	for l, m := range ms {
		res, _ := jsoniter.MarshalIndent(m, "", "  ")
		log.Printf("Locale %s:\n%s", l, res)
	}

}
```

## Respond to alexa requests with lambda
```go

```

## Projects using `go-alexa-lambda`
* [alexa-go-cloudformation-demo](https://github.com/DrPsychick/alexa-go-cloudformation-demo) : the demo project that lead to developing this library. A fully automated build and deploy of an Alexa skill including lambda function via Cloudformation.

## Project template
To give you a head start, check out the template project:
* `go-alexa-lambda-template` : generate skill and deploy it including lambda with cloudformation.

### Create your own project based on the template
```shell
git clone https://github.com/drpsychick/go-alexa-lambda-template
mv go-alexa-lambda-template alexa-project
cd alexa-project
rm -rf .git
```

### Required variables for your pipeline
```shell
ASKClientId= # amzn1.application-oa2-client...
ASKClientSecret= # 6ef4...
ASKAccessToken= # Atza|IwE...
ASKRefreshToken= # 	Atzr|IwE...
ASKVendorId= # M3D....
ASKS3Bucket= # my-bucket
ASKS3Key= # my-skill.zip
ASKSkillId= # amzn1.ask.skill.fe19...
ASK_CONFIG= # cat ~/.ask/cli_config | jq -c | sed -e 's#\(["{}|]\)#\\\1#g'
AWS_DEFAULT_REGION= # eu-central-1
AWS_ACCESS_KEY_ID= # AKD...
AWS_SECRET_ACCESS_KEY= # 3D+...
CF_STACK_NAME= # skill-stack
KEEP_STACK= # if empty, the CF stack will be deleted after deploy (for tests)
```

### Links
* https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html
* multiple intents in one dialog: https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html#pass-a-new-intent
* https://developer.amazon.com/blogs/alexa/post/cfbd2f5e-c72f-4b03-8040-8628bbca204c/alexa-skill-teardown-understanding-entity-resolution-with-pet-match

### Credits
initially inspired by: https://github.com/soloworks/go-alexa-models