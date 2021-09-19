[![Go Report Card](https://goreportcard.com/badge/github.com/drpsychick/go-alexa-lambda)](https://goreportcard.com/report/github.com/drpsychick/go-alexa-lambda)
[![Build Status](https://app.travis-ci.com/DrPsychick/go-alexa-lambda.svg?branch=master)](https://app.travis-ci.com/DrPsychick/go-alexa-lambda)
[![Coverage Status](https://coveralls.io/repos/github/DrPsychick/go-alexa-lambda/badge.svg?branch=master)](https://coveralls.io/github/DrPsychick/go-alexa-lambda?branch=master)
[![license](https://img.shields.io/github/license/drpsychick/go-alexa-lambda.svg)](https://github.com/drpsychick/go-alexa-lambda/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/drpsychick/go-alexa-lambda.svg)](https://github.com/drpsychick/go-alexa-lambda)
[![Paypal](https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=FTXDN7LCDWUEA&source=url)
[![GitHub Sponsor](https://img.shields.io/badge/github-sponsor-blue?logo=github)](https://github.com/sponsors/DrPsychick)

# go-alexa-lambda
Alexa golang library to generate skill + interaction model as well as serve requests with lambda or as a server.
The packages can also be used standalone as a request/response abstraction for Alexa requests.

# Purpose
The Alexa skill and interaction model is tightly coupled with the actual intents a lambda function will process.
As developing a skill with a golang backend impacts the skill definition in many cases and they share localization scope,
this package allows defining and generating the skill and model for deployment and using the same source (with intents, localization) to
build the lambda function that responds to requests from Alexa.

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
Run it on go playground: https://play.golang.org/p/fRrzn_kmaBi
```go
package main

import (
	"os"
	"context"
	alexa "github.com/drpsychick/go-alexa-lambda"
	"github.com/drpsychick/go-alexa-lambda/skill"
	log "github.com/hamba/logger/v2"
)

var request = `{
  "version": "1.0",
  "session": {},
  "context": {},
  "request": {
    "type": "IntentRequest",
    "requestId": "amzn1.echo-api.request.1234",
    "timestamp": "2016-10-27T21:06:28Z",
    "locale": "en-US",
    "intent": {
      "name": "AMAZON.HelpIntent"
    }
  }
}
`

func handleHelp(sb *skill.SkillBuilder) alexa.HandlerFunc {
	sb.Model().WithIntent(alexa.HelpIntent)

	return alexa.HandlerFunc(func(b *alexa.ResponseBuilder, r *alexa.RequestEnvelope) {
		b.WithSimpleCard("Help Title", "Text explaining how it works.")
	})
}

func main() {
	sb := skill.NewSkillBuilder()
	mux := alexa.NewServerMux(log.New(os.Stdout, log.ConsoleFormat(), log.Info))
	sb.WithModel()

	mux.HandleIntent(alexa.HelpIntent, handleHelp(sb))
	
	// actually, one would call `alexa.Serve(mux)`
	// but we want to pass a request and get a response
	s := &alexa.Server{Handler: mux}
	ctx := context.Background()
	response, err := s.Invoke(ctx, []byte(request))
	if err != nil {
		mux.Logger().Error(err.Error())
	}
	mux.Logger().Info(string(response))
}
```

# Projects using `go-alexa-lambda`
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

# References
### Links
* https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html
* multiple intents in one dialog: https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html#pass-a-new-intent
* https://developer.amazon.com/blogs/alexa/post/cfbd2f5e-c72f-4b03-8040-8628bbca204c/alexa-skill-teardown-understanding-entity-resolution-with-pet-match

### Credits
initially inspired by: https://github.com/soloworks/go-alexa-models