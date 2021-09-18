package skill

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var minimalSkillDef = Skill{
	Manifest: Manifest{
		Version: "1.0",
		Publishing: Publishing{
			Locales: map[string]LocaleDef{
				"de-DE": {
					Name:        "name",
					Description: "description",
					Summary:     "summary",
					Keywords:    []string{"Demo"},
					Examples:    []string{"tell me how much beer people drink in germany"},
				},
			},
			Category:  "mycategory",
			Countries: []string{"DE"},
		},
		Apis: &Apis{
			Custom: &Custom{
				Endpoint: &Endpoint{
					URI: "arn:...",
				},
			},
			Interfaces: []string{},
		},
		Permissions: []Permission{},
		Privacy: &Privacy{
			IsExportCompliant: true,
			Locales:           map[string]PrivacyLocaleDef{},
		},
	},
}

// https://developer.amazon.com/de/docs/smapi/skill-manifest.html#sample-manifests
var awsCustomSkillExample = []byte(`{
  "manifest": {
    "publishingInformation": {
      "locales": {
        "en-US": {
          "summary": "This is a sample Alexa custom skill.",
          "examplePhrases": [
            "Alexa, open sample custom skill.",
            "Alexa, play sample custom skill."
          ],
          "keywords": [
            "Descriptive_Phrase_1",
            "Descriptive_Phrase_2",
            "Descriptive_Phrase_3"
          ],
          "smallIconUri": "https://smallUri.com",
          "largeIconUri": "https://largeUri.com",
          "name": "Sample custom skill name.",
          "description": "This skill does interesting things."
        }
      },
      "isAvailableWorldwide": false,
      "testingInstructions": "1) Say 'Alexa, hello world'",
      "category": "HEALTH_AND_FITNESS",
      "distributionCountries": [
        "US",
        "GB",
        "DE"
      ]
    },
    "apis": {
      "custom": {
        "endpoint": {
          "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
        },
        "interfaces": [
          {
            "type":"ALEXA_PRESENTATION_APL"
          },		
          {
            "type":"AUDIO_PLAYER"
          },
          {
            "type":"CAN_FULFILL_INTENT_REQUEST"
          },
          {
            "type":"GADGET_CONTROLLER"
          },
          {
            "type":"GAME_ENGINE"
          },
          {
            "type":"RENDER_TEMPLATE"
          },
          {
            "type":"VIDEO_APP"
          }		  
        ],
        "regions": {
          "NA": {
            "endpoint": {
              "sslCertificateType": "Trusted",
              "uri": "https://customapi.sampleskill.com"
            }
          }
        }
      }
    },
    "manifestVersion": "1.0",
    "permissions": [
      {
        "name": "alexa::devices:all:address:full:read"
      },
      {
        "name": "alexa:devices:all:address:country_and_postal_code:read"
      },
      {
        "name": "alexa::household:lists:read"
      },
      {
        "name": "alexa::household:lists:write"
      },
      {
        "name": "alexa::alerts:reminders:skill:readwrite"
      }
    ],
    "privacyAndCompliance": {
      "allowsPurchases": false,
      "usesPersonalInfo": false,
      "isChildDirected": false,
      "isExportCompliant": true,
      "containsAds": false,
      "locales": {
        "en-US": {
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com",
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com"
        }
      }
    },
    "events": {
      "endpoint": {
        "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
      },
      "subscriptions": [
        {
          "eventName": "SKILL_ENABLED"
        },
        {
          "eventName": "SKILL_DISABLED"
        },
        {
          "eventName": "SKILL_PERMISSION_ACCEPTED"
        },
        {
          "eventName": "SKILL_PERMISSION_CHANGED"
        },
        {
          "eventName": "SKILL_ACCOUNT_LINKED"
        }
      ],
      "regions": {
        "NA": {
          "endpoint": {
            "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
          }
        }
      }
    }
  }
}`)

var awsListSkillExample = []byte(`
{
  "manifest": {
    "publishingInformation": {
      "locales": {
        "en-US": {
          "summary": "This is a sample Alexa skill.",
          "examplePhrases": [
            "Alexa, open sample skill.",
            "Alexa, play sample skill."
          ],
          "keywords": [
            "Descriptive_Phrase_1",
            "Descriptive_Phrase_2",
            "Descriptive_Phrase_3"
          ],
          "smallIconUri": "https://smallUri.com",
          "largeIconUri": "https://largeUri.com",
          "name": "Sample skill name.",
          "description": "This skill does interesting things."
        }
      },
      "isAvailableWorldwide": false,
      "testingInstructions": "1) Say 'Alexa, hello world'",
      "category": "HEALTH_AND_FITNESS",
      "distributionCountries": [
        "US",
        "GB",
        "DE"
      ]
    },
    "apis": {
      "householdList": {}
    },
    "manifestVersion": "1.0",
    "permissions": [
      {
        "name": "alexa::devices:all:address:full:read"
      },
      {
        "name": "alexa:devices:all:address:country_and_postal_code:read"
      },
      {
        "name": "alexa::household:lists:read"
      },
      {
        "name": "alexa::household:lists:write"
      }
    ],
    "privacyAndCompliance": {
      "allowsPurchases": false,
      "locales": {
        "en-US": {
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com",
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com"
        }
      },
      "isExportCompliant": true,
      "containsAds": false,
      "isChildDirected": false,
      "usesPersonalInfo": false
    },
    "events": {
      "endpoint": {
        "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
      },
      "subscriptions": [
        {
          "eventName": "SKILL_ENABLED"
        },
        {
          "eventName": "SKILL_DISABLED"
        },
        {
          "eventName": "SKILL_PERMISSION_ACCEPTED"
        },
        {
          "eventName": "SKILL_PERMISSION_CHANGED"
        },
        {
          "eventName": "SKILL_ACCOUNT_LINKED"
        },
        {
          "eventName": "ITEMS_CREATED"
        },
        {
          "eventName": "ITEMS_UPDATED"
        },
        {
          "eventName": "ITEMS_DELETED"
        }
      ],
      "regions": {
        "NA": {
          "endpoint": {
            "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
          }
        }
      }
    }
  }
}`)

var awsMeetingSkillExample = []byte(`{
  "manifest": {
    "apis": {
      "alexaForBusiness": {
        "regions": {
          "NA": {
            "endpoint": {
              "uri": "arn:aws:lambda:us-east-1:123456789:function:myFunctionName1"
            }
          }
        },
        "endpoint": {
          "uri": "arn:aws:lambda:us-east-1:123456789:function:myFunctionName1"
        },
        "interfaces": [
          {
            "namespace": "Alexa.Business.Reservation.Room",
            "version": "1.0",
            "requests": [
              {
                "name": "Search"
              },
              {
                "name": "Create"
              },
              {
                "name": "Update"
              }
            ]
          }
        ]
      }
    },
    "manifestVersion": "1.0",
    "privacyAndCompliance": {
      "locales": {
        "en-US": {
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com",
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com"
        }
      },
      "allowsPurchases": false,
      "usesPersonalInfo": false,
      "isChildDirected": false,
      "isExportCompliant": true,
      "containsAds": false
    },
    "publishingInformation": {
      "locales": {
        "en-US": {
          "name": "Room Booking Skill",
          "smallIconUri": "https://smallUri.example.com/small1.png",
          "largeIconUri": "https://largeUri.example.com/large1.png",
          "summary": "This is a sample Alexa skill.",
          "description": "This skill has Alexa for Business reservations features.",
          "examplePhrases": [
            "Alexa, book this room.",
            "Alexa, find a room at 3pm tomorrow."
          ],
          "keywords": [
            "Meetings",
            "Booking",
            "Alexa For Business"
          ],
          "updatesDescription": "This skill has updates that fix feature bugs."
        }
      },
      "isAvailableWorldwide": false,
      "testingInstructions": "1) Say 'Alexa, Book this room'",
      "category": "CALENDARS_AND_REMINDERS",
      "distributionCountries": [
        "US"
      ]
    }
  }
}`)

var awsSmarthomeExample = []byte(`{
  "manifest": {
    "manifestVersion": "1.0",
    "publishingInformation": {
      "locales": {
        "en-US": {
          "name": "Sample skill name.",
          "summary": "This is a sample Alexa skill.",
          "description": "This skill has basic and advanced smart devices control features.",
          "smallIconUri": "https://smallUri.com",
          "largeIconUri": "https://largeUri.com",
          "examplePhrases": [
            "Alexa, open sample skill.",
            "Alexa, blink kitchen lights."
          ],
          "keywords": [
            "Smart Home",
            "Lights",
            "Smart Devices"
          ]
        }
      },
      "distributionCountries": [
        "US",
        "GB",
        "DE"
      ],
      "isAvailableWorldwide": false,
      "testingInstructions": "1) Say 'Alexa, turn on sample lights'",
      "category": "SMART_HOME"
    },
    "privacyAndCompliance": {
      "allowsPurchases": false,
      "usesPersonalInfo": false,
      "isChildDirected": false,
      "isExportCompliant": true,
      "containsAds": false,
      "locales": {
        "en-US": {
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com",
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com"
        }
      }
    },
    "apis": {
      "smartHome": {
        "endpoint": {
          "uri": "arn:aws:lambda:us-east-1:040623927470:function:sampleSkill"
        },
        "regions": {
          "NA": {
            "endpoint": {
              "uri": "arn:aws:lambda:us-west-2:010623927470:function:sampleSkillWest"
            }
          }
        }
      }
    }
  }
}`)

var awsVideoSkillExample = []byte(`{
  "manifest": {
    "publishingInformation": {
      "locales": {
        "en-US": {
          "summary": "This is a sample Alexa skill.",
          "examplePhrases": [
            "Alexa, tune to channel 206",
            "Alexa, search for comedy movies",
            "Alexa, pause."
          ],
          "keywords": [
            "Video",
            "TV"
          ],
          "name": "VideoSampleSkill",
          "smallIconUri": "https://smallUri.com",
          "largeIconUri": "https://smallUri.com",
          "description": "This skill has video control features."
        }
      },
      "isAvailableWorldwide": false,
      "testingInstructions": "",
      "category": "SMART_HOME",
      "distributionCountries": [
        "US",
        "GB",
        "DE"
      ]
    },
    "apis": {
      "video": {
        "locales": {
          "en-US": {
            "videoProviderTargetingNames": [
              "TV provider"
            ],
            "catalogInformation": [
              {
                "sourceId": "1234",
                "type": "FIRE_TV"
              }
            ]
          }
        },
        "endpoint": {
          "uri": "arn:aws:lambda:us-east-1:452493640596:function:sampleSkill"
        },
        "regions": {
          "NA": {
            "endpoint": {
              "uri": "arn:aws:lambda:us-east-1:452493640596:function:sampleSkill"
            },
            "upchannel": [
              {
                "uri": "arn:aws:sns:us-east-1:291420629295:sampleSkill",
                "type": "SNS"
              }
            ]
          }
        }
      }
    },
    "manifestVersion": "1.0",
    "privacyAndCompliance": {
      "allowsPurchases": false,
      "locales": {
        "en-US": {
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com",
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com"
        }
      },
      "isExportCompliant": true,
      "isChildDirected": false,
      "usesPersonalInfo": false,
      "containsAds": false
    }
  }
}`)

var awsBabyActivityExample = []byte(`{
   "skillManifest": {
            "publishingInformation": {
                "locales": {
                    "en-US": {
                        "summary": "Baby Activity Skill 1",
                        "examplePhrases": [
                        "\"Alexa, log a diaper change for Jane\"",
                        "\"Alexa, what is Jane's weight\""
                        ],
                        "description": "A skill that logs and tracks baby activities",
                        "keywords": [
                            "Family",
                            "Infant Tracking"
                        ],
                        "name": "Baby Activity Test Skill",
                        "smallIconUri": "iconUri",
                        "largeIconUri": "iconUri"
                    }
                },
                "isAvailableWorldwide": true,
                "category": "HEALTH_AND_FITNESS",
                "testingInstructions": "Alexa, log a diaper change",
                "distributionCountries": []
            },
            "permissions": [
                {
                    "name": "alexa::health:profile:write"
                }
            ],
            "apis": {
                "health": {
                    "endpoint": {
                          "uri": "lambda-endpoint"
                    },
                    "regions": {
                            "NA": {
                               "endpoint": {
                                  "uri": "lambda-endpoint"
                                }
                            }
                    }
                }
            },
            "privacyAndCompliance": {
                "locales": {
                    "en-US": {
                        "privacyPolicyUrl": "https://example.com/privacy"
                    }
                },
                "allowsPurchases": false,
                "isExportCompliant": true,
                "containsAds": false,
                "isChildDirected": false,
                "usesPersonalInfo": true
            },
            "manifestVersion": "1.0"
        },
        "vendorId" : "your-vendor-id"
}`)

var awsFlashBriefingExample = []byte(`{
  "manifest": {
    "manifestVersion": "1.0",
    "publishingInformation": {
      "locales": {
        "en-US": {
          "name": "Sample skill name.",
          "summary": "This is a sample Alexa skill.",
          "description": "This skill has basic and advanced features.",
          "smallIconUri": "https://smallUri.com",
          "largeIconUri": "https://largeUri.com",
          "examplePhrases": [],
          "keywords": [
            "Flash Briefing",
            "News",
            "Happenings"
          ]
        }
      },
      "distributionCountries": [
        "US",
        "GB",
        "DE"
      ],
      "isAvailableWorldwide": false,
      "testingInstructions": "1) Say 'Alexa, hello world'",
      "category": "HEALTH_AND_FITNESS"
    },
    "privacyAndCompliance": {
      "allowsPurchases": false,
      "usesPersonalInfo": false,
      "isChildDirected": false,
      "isExportCompliant": true,
	  "containsAds": false,
      "locales": {
        "en-US": {
          "privacyPolicyUrl": "http://www.myprivacypolicy.sampleskill.com",
          "termsOfUseUrl": "http://www.termsofuse.sampleskill.com"
        }
      }
    },
    "apis": {
      "flashBriefing": {
        "locales": {
          "en-US": {
            "customErrorMessage": "Error message",
            "feeds": [
              {
                "name": "feed name",
                "isDefault": true,
                "vuiPreamble": "In this skill",
                "updateFrequency": "HOURLY",
                "genre": "POLITICS",
                "imageUri": "https://fburi.com",
                "contentType": "TEXT",
                "url": "https://feeds.sampleskill.com/feedX"
              }
            ]
          }
        }
      }
    }
  }
}`)

func TestMinimalSkillDefinition(t *testing.T) {
	res, _ := json.Marshal(minimalSkillDef)
	assert.NotEmpty(t, string(res), "Generated JSON must not be empty")

}

func TestSampleManifest(t *testing.T) {
	testHelper(t, "Custom Skill Example", awsCustomSkillExample)
}

func TestSampleListSkill(t *testing.T) {
	testHelper(t, "List Skill Example", awsListSkillExample)
}

func TestSmartHome(t *testing.T) {
	testHelper(t, "SmartHome Example", awsSmarthomeExample)
}

func TestBabyActivity(t *testing.T) {
	testHelper(t, "Baby Activity Example", awsBabyActivityExample)
}

func TestMeetingSkill(t *testing.T) {
	testHelper(t, "Meeting Skill Example", awsMeetingSkillExample)
}

func TestFlashBriefing(t *testing.T) {
	res := testHelper(t, "Flash Briefing Example", awsFlashBriefingExample)
	assert.Contains(t, string(res), "flashBriefing")
	assert.Contains(t, string(res), "feeds")
}

func TestVideoSkill(t *testing.T) {
	testHelper(t, "Video Skill Example", awsVideoSkillExample)
}

func testHelper(t *testing.T, name string, manifest []byte) []byte {
	var skill Skill
	err := json.Unmarshal(manifest, &skill)
	assert.Nil(t, err, "Unmarshal of %s returned error: %s", name, err)
	res, err := json.Marshal(skill)
	assert.NoError(t, err)
	assert.NotEmpty(t, string(res), "Marshal of %s must return JSON", name)
	return res
}
