package skill

// Skill is the Alexa `skill.json` top element.
type Skill struct {
	Manifest Manifest `json:"manifest"`
}

// Manifest is the parent for all other elements.
type Manifest struct {
	Version     string       `json:"manifestVersion"`
	Publishing  Publishing   `json:"publishingInformation"`
	Apis        *Apis        `json:"apis,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	Privacy     *Privacy     `json:"privacyAndCompliance"`
}

// Publishing information.
type Publishing struct {
	Locales             map[string]LocaleDef `json:"locales"`
	Worldwide           bool                 `json:"isAvailableWorldwide"`
	Category            Category             `json:"category"`
	Countries           []string             `json:"distributionCountries,omitempty"`
	TestingInstructions string               `json:"testingInstructions"`
}

// LocaleDef description of each locale.
type LocaleDef struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Summary      string   `json:"summary"`
	Examples     []string `json:"examplePhrases"`
	Keywords     []string `json:"keywords"`
	SmallIconURI string   `json:"smallIconUri"`
	LargeIconURI string   `json:"largeIconUri"`
}

const (
	// CountryAustralia is AU.
	CountryAustralia string = "AU"
	// CountryCanada is CA.
	CountryCanada string = "CA"
	// CountryGermany is DE.
	CountryGermany string = "DE"
	// CountryFrance is FR.
	CountryFrance string = "FR"
	// CountryGreatBritain is GB.
	CountryGreatBritain string = "GB"
	// CountryIndia is IN.
	CountryIndia string = "IN"
	// CountryItaly is IT.
	CountryItaly string = "IT"
	// CountryJapan is JP.
	CountryJapan string = "JP"
	// CountryUnitedStates is US.
	CountryUnitedStates string = "US"
)

// Category of the Skill that is used for filtering in the Alexa App
//
// see https://developer.amazon.com/de/docs/smapi/skill-manifest.html#category-enum
type Category string

const (
	// CategoryAlarmsAndClocks is ALARMS_AND_CLOCKS.
	CategoryAlarmsAndClocks Category = "ALARMS_AND_CLOCKS"
	// CategoryAstrology is ASTROLOGY.
	CategoryAstrology Category = "ASTROLOGY"
	// CategoryBusinessAndFinance is BUSINESS_AND_FINANCE.
	CategoryBusinessAndFinance Category = "BUSINESS_AND_FINANCE"
	// CategoryCalculators is CALCULATORS.
	CategoryCalculators Category = "CALCULATORS"
	// CategoryCalendarsAndReminders is CALENDARS_AND_REMINDERS.
	CategoryCalendarsAndReminders Category = "CALENDARS_AND_REMINDERS"
	// CategoryChildrensEducationAndReference is CHILDRENS_EDUCATION_AND_REFERENCE.
	CategoryChildrensEducationAndReference Category = "CHILDRENS_EDUCATION_AND_REFERENCE"
	// CategoryChildrensGames is CHILDRENS_GAMES.
	CategoryChildrensGames Category = "CHILDRENS_GAMES"
	// CategoryChildrensMusicAndAudio is CHILDRENS_MUSIC_AND_AUDIO.
	CategoryChildrensMusicAndAudio Category = "CHILDRENS_MUSIC_AND_AUDIO"
	// CategoryChildrensNoveltyAndHumor is CHILDRENS_NOVELTY_AND_HUMOR.
	CategoryChildrensNoveltyAndHumor Category = "CHILDRENS_NOVELTY_AND_HUMOR"
	// CategoryCommunication is COMMUNICATION.
	CategoryCommunication Category = "COMMUNICATION"
	// CategoryConnectedCar is CONNECTED_CAR.
	CategoryConnectedCar Category = "CONNECTED_CAR"
	// CategoryCookingAndRecipe is COOKING_AND_RECIPE.
	CategoryCookingAndRecipe Category = "COOKING_AND_RECIPE"
	// CategoryCurrencyGuidesAndConverters is CURRENCY_GUIDES_AND_CONVERTERS.
	CategoryCurrencyGuidesAndConverters Category = "CURRENCY_GUIDES_AND_CONVERTERS"
	// CategoryDating is DATING.
	CategoryDating Category = "DATING"
	// CategoryDeliveryAndTakeout is DELIVERY_AND_TAKEOUT.
	CategoryDeliveryAndTakeout Category = "DELIVERY_AND_TAKEOUT"
	// CategoryDeviceTracking is DEVICE_TRACKING.
	CategoryDeviceTracking Category = "DEVICE_TRACKING"
	// CategoryEducationAndReference is EDUCATION_AND_REFERENCE.
	CategoryEducationAndReference Category = "EDUCATION_AND_REFERENCE"
	// CategoryEventFinders is EVENT_FINDERS.
	CategoryEventFinders Category = "EVENT_FINDERS"
	// CategoryExerciseAndWorkout is EXERCISE_AND_WORKOUT.
	CategoryExerciseAndWorkout Category = "EXERCISE_AND_WORKOUT"
	// CategoryFashionAndStyle is FASHION_AND_STYLE.
	CategoryFashionAndStyle Category = "FASHION_AND_STYLE"
	// CategoryFlightFinders is FLIGHT_FINDERS.
	CategoryFlightFinders Category = "FLIGHT_FINDERS"
	// CategoryFriendsAndFamily is FRIENDS_AND_FAMILY.
	CategoryFriendsAndFamily Category = "FRIENDS_AND_FAMILY"
	// CategoryGameInfoAndAccessory is GAME_INFO_AND_ACCESSORY.
	CategoryGameInfoAndAccessory Category = "GAME_INFO_AND_ACCESSORY"
	// CategoryGames is GAMES.
	CategoryGames Category = "GAMES"
	// CategoryHealthAndFitness is HEALTH_AND_FITNESS.
	CategoryHealthAndFitness Category = "HEALTH_AND_FITNESS"
	// CategoryHotelFinders is HOTEL_FINDERS.
	CategoryHotelFinders Category = "HOTEL_FINDERS"
	// CategoryKnowledgeAndTrivia is KNOWLEDGE_AND_TRIVIA.
	CategoryKnowledgeAndTrivia Category = "KNOWLEDGE_AND_TRIVIA"
	// CategoryMovieAndTvKnowledgeAndTrivia is MOVIE_AND_TV_KNOWLEDGE_AND_TRIVIA.
	CategoryMovieAndTvKnowledgeAndTrivia Category = "MOVIE_AND_TV_KNOWLEDGE_AND_TRIVIA"
	// CategoryMovieInfoAndReviews is MOVIE_INFO_AND_REVIEWS.
	CategoryMovieInfoAndReviews Category = "MOVIE_INFO_AND_REVIEWS"
	// CategoryMovieShowtimes is MOVIE_SHOWTIMES.
	CategoryMovieShowtimes Category = "MOVIE_SHOWTIMES"
	// CategoryMusicAndAudioAccessories is MUSIC_AND_AUDIO_ACCESSORIES.
	CategoryMusicAndAudioAccessories Category = "MUSIC_AND_AUDIO_ACCESSORIES"
	// CategoryMusicAndAudioKnowledgeAndTrivia is MUSIC_AND_AUDIO_KNOWLEDGE_AND_TRIVIA.
	CategoryMusicAndAudioKnowledgeAndTrivia Category = "MUSIC_AND_AUDIO_KNOWLEDGE_AND_TRIVIA"
	// CategoryMusicInfoReviewsAndRecognitionService is MUSIC_INFO_REVIEWS_AND_RECOGNITION_SERVICE.
	CategoryMusicInfoReviewsAndRecognitionService Category = "MUSIC_INFO_REVIEWS_AND_RECOGNITION_SERVICE"
	// CategoryNavigationAndTripPlanner is NAVIGATION_AND_TRIP_PLANNER.
	CategoryNavigationAndTripPlanner Category = "NAVIGATION_AND_TRIP_PLANNER"
	// CategoryNews is NEWS.
	CategoryNews Category = "NEWS"
	// CategoryNovelty is NOVELTY.
	CategoryNovelty Category = "NOVELTY"
	// CategoryOrganizersAndAssistants is ORGANIZERS_AND_ASSISTANTS.
	CategoryOrganizersAndAssistants Category = "ORGANIZERS_AND_ASSISTANTS"
	// CategoryPetsAndAnimal is PETS_AND_ANIMAL.
	CategoryPetsAndAnimal Category = "PETS_AND_ANIMAL"
	// CategoryPodcast is PODCAST.
	CategoryPodcast Category = "PODCAST"
	// CategoryPublicTransportation is PUBLIC_TRANSPORTATION.
	CategoryPublicTransportation Category = "PUBLIC_TRANSPORTATION"
	// CategoryReligionAndSpirituality is RELIGION_AND_SPIRITUALITY.
	CategoryReligionAndSpirituality Category = "RELIGION_AND_SPIRITUALITY"
	// CategoryRestaurantBookingInfoAndReview is RESTAURANT_BOOKING_INFO_AND_REVIEW.
	CategoryRestaurantBookingInfoAndReview Category = "RESTAURANT_BOOKING_INFO_AND_REVIEW"
	// CategorySchools is SCHOOLS.
	CategorySchools Category = "SCHOOLS"
	// CategoryScoreKeeping is SCORE_KEEPING.
	CategoryScoreKeeping Category = "SCORE_KEEPING"
	// CategorySelfImprovement is SELF_IMPROVEMENT.
	CategorySelfImprovement Category = "SELF_IMPROVEMENT"
	// CategoryShopping is SHOPPING.
	CategoryShopping Category = "SHOPPING"
	// CategorySmartHome is SMART_HOME.
	CategorySmartHome Category = "SMART_HOME"
	// CategorySocialNetworking is SOCIAL_NETWORKING.
	CategorySocialNetworking Category = "SOCIAL_NETWORKING"
	// CategorySportsGames is SPORTS_GAMES.
	CategorySportsGames Category = "SPORTS_GAMES"
	// CategorySportsNews is SPORTS_NEWS.
	CategorySportsNews Category = "SPORTS_NEWS"
	// CategoryStreamingService is STREAMING_SERVICE.
	CategoryStreamingService Category = "STREAMING_SERVICE"
	// CategoryTaxiAndRidesharing is TAXI_AND_RIDESHARING.
	CategoryTaxiAndRidesharing Category = "TAXI_AND_RIDESHARING"
	// CategoryToDoListsAndNotes is TO_DO_LISTS_AND_NOTES.
	CategoryToDoListsAndNotes Category = "TO_DO_LISTS_AND_NOTES"
	// CategoryTranslators is TRANSLATORS.
	CategoryTranslators Category = "TRANSLATORS"
	// CategoryTvGuides is TV_GUIDES.
	CategoryTvGuides Category = "TV_GUIDES"
	// CategoryUnitConverters is UNIT_CONVERTERS.
	CategoryUnitConverters Category = "UNIT_CONVERTERS"
	// CategoryWeather is WEATHER.
	CategoryWeather Category = "WEATHER"
	// CategoryWineAndBeverage is WINE_AND_BEVERAGE.
	CategoryWineAndBeverage Category = "WINE_AND_BEVERAGE"
	// CategoryZipCodeLookup is ZIP_CODE_LOOKUP.
	CategoryZipCodeLookup Category = "ZIP_CODE_LOOKUP"
)

// Apis of the Alexa Skill https://developer.amazon.com/de/docs/smapi/skill-manifest.html#apis
type Apis struct {
	ForBusiness *ForBusiness `json:"alexaForBusiness,omitempty"`
	Custom      *Custom      `json:"custom,omitempty"`
	// SmartHome *SmartHome `json:"smartHome"`
	FlashBriefing *FlashBriefing `json:"flashBriefing"`
	// Health     *Health	`json:"health"`
	// HouseholdList *HouseholdList `json:"householdList"`
	// Video *Video `json:"video"`
	Interfaces []string `json:"interfaces,omitempty"`
}

// ForBusiness API are available in English only.
type ForBusiness struct {
	Endpoint   *Endpoint             `json:"endpoint"`
	Regions    *map[Region]RegionDef `json:"regions"`
	Interfaces []Interface           `json:"interfaces,omitempty"`
}

// Custom API endpoint.
type Custom struct {
	Endpoint   *Endpoint             `json:"endpoint"`
	Regions    *map[Region]RegionDef `json:"regions"`
	Interfaces []Interface           `json:"interfaces,omitempty"`
}

// FlashBriefing API endpoint.
type FlashBriefing struct {
	Locales map[string]FlashBriefingLocaleDef `json:"locales"`
}

// FlashBriefingLocaleDef is locale definition for flashBriefing API.
type FlashBriefingLocaleDef struct {
	CustomErrorMessage string                    `json:"customErrorMessage"`
	Feeds              []FlashBriefingLocaleFeed `json:"feeds"`
}

// FlashBriefingLocaleFeed is a feed definition for flashBriefing API.
type FlashBriefingLocaleFeed struct {
	Name            string          `json:"name"`
	IsDefault       bool            `json:"isDefault"`
	VuiPreamble     string          `json:"vuiPreamble"`
	UpdateFrequency UpdateFrequency `json:"updateFrequency"`
	Genre           ContentGenre    `json:"genre"`
	ImageURI        string          `json:"imageUri"`
	ContentType     ContentType     `json:"contentType"`
	URL             string          `json:"url"`
}

// UpdateFrequency is an enum for flashBriefing feed update frequency.
type UpdateFrequency string

const (
	// UpdateFrequencyDaily is DAILY.
	UpdateFrequencyDaily UpdateFrequency = "DAILY"
	// UpdateFrequencyHourly is HOURLY.
	UpdateFrequencyHourly UpdateFrequency = "HOURLY"
	// UpdateFrequencyWeekly is WEEKLY.
	UpdateFrequencyWeekly UpdateFrequency = "WEEKLY"
)

// ContentGenre is an enum for flashBriefing feed content genre.
type ContentGenre string

const (
	// ContentGenreHeadlineNews is HEADLINE_NEWS.
	ContentGenreHeadlineNews ContentGenre = "HEADLINE_NEWS"
	// ContentGenreBusiness is BUSINESS.
	ContentGenreBusiness ContentGenre = "BUSINESS"
	// ContentGenrePolitics is POLITICS.
	ContentGenrePolitics ContentGenre = "POLITICS"
	// ContentGenreEntertainment is ENTERTAINMENT.
	ContentGenreEntertainment ContentGenre = "ENTERTAINMENT"
	// ContentGenreTechnology is TECHNOLOGY.
	ContentGenreTechnology ContentGenre = "TECHNOLOGY"
	// ContentGenreHumor is HUMOR.
	ContentGenreHumor ContentGenre = "HUMOR"
	// ContentGenreLifestyle is LIFESTYLE.
	ContentGenreLifestyle ContentGenre = "LIFESTYLE"
	// ContentGenreSports is SPORTS.
	ContentGenreSports ContentGenre = "SPORTS"
	// ContentGenreScience is SCIENCE.
	ContentGenreScience ContentGenre = "SCIENCE"
	// ContentGenreHealthAndFitness is HEALTH_AND_FITNESS.
	ContentGenreHealthAndFitness ContentGenre = "HEALTH_AND_FITNESS"
	// ContentGenreArtsAndCulture is ARTS_AND_CULTURE.
	ContentGenreArtsAndCulture ContentGenre = "ARTS_AND_CULTURE"
	// ContentGenreProductivityAndUtilities is PRODUCTIVITY_AND_UTILITIES.
	ContentGenreProductivityAndUtilities ContentGenre = "PRODUCTIVITY_AND_UTILITIES"
	// ContentGenreOther is OTHER.
	ContentGenreOther ContentGenre = "OTHER"
)

// ContentType is an enum for flashBriefing feed content type.
type ContentType string

const (
	// ContentTypeText is TEXT.
	ContentTypeText ContentType = "TEXT"
	// ContentTypeAudio is AUDIO.
	ContentTypeAudio ContentType = "AUDIO"
)

// Endpoint definition.
type Endpoint struct {
	URI                string `json:"uri"`
	SslCertificateType string `json:"sslCertificateType,omitempty"`
}

// Region for Alexa.
type Region string

const (
	// RegionNorthAmerica is NA.
	RegionNorthAmerica Region = "NA"
	// RegionEurope is EU.
	RegionEurope Region = "EU"
	// RegionFarEast is FE.
	RegionFarEast Region = "FE"
)

// RegionDef for regional endpoints.
type RegionDef struct {
	Endpoint *Endpoint `json:"endpoint"`
}

// Interface definition for API.
type Interface struct {
	Type InterfaceType `json:"type"`
}

// InterfaceType string reference.
type InterfaceType string

const (
	// InterfaceTypeAlexaPresentationAPL is ALEXA_PRESENTATION_APL.
	InterfaceTypeAlexaPresentationAPL InterfaceType = "ALEXA_PRESENTATION_APL"
	// InterfaceTypeAudioPlayer is AUDIO_PLAYER.
	InterfaceTypeAudioPlayer InterfaceType = "AUDIO_PLAYER"
	// InterfaceTypeCanFulfillIntentRequest is CAN_FULFILL_INTENT_REQUEST.
	InterfaceTypeCanFulfillIntentRequest InterfaceType = "CAN_FULFILL_INTENT_REQUEST"
	// InterfaceTypeGadgetController is GADGET_CONTROLLER.
	InterfaceTypeGadgetController InterfaceType = "GADGET_CONTROLLER"
	// InterfaceTypeGameEngine is GAME_ENGINE.
	InterfaceTypeGameEngine InterfaceType = "GAME_ENGINE"
	// InterfaceTypeRenderTemplate is RENDER_TEMPLATE.
	InterfaceTypeRenderTemplate InterfaceType = "RENDER_TEMPLATE"
	// InterfaceTypeVideoApp is VIDEO_APP.
	InterfaceTypeVideoApp InterfaceType = "VIDEO_APP"
)

// Permission string.
type Permission struct {
	Name string `json:"name"`
}

// Privacy definition.
type Privacy struct {
	IsExportCompliant bool                        `json:"isExportCompliant"`
	ContainsAds       bool                        `json:"containsAds"`
	AllowsPurchases   bool                        `json:"allowsPurchases"`
	UsesPersonalInfo  bool                        `json:"usesPersonalInfo"`
	IsChildDirected   bool                        `json:"isChildDirected"`
	Locales           map[string]PrivacyLocaleDef `json:"locales,omitempty"`
}

// PrivacyLocaleDef defines.
type PrivacyLocaleDef struct {
	PrivacyPolicyURL string `json:"privacyPolicyUrl,omitempty"`
	TermsOfUse       string `json:"termsOfUse,omitempty"`
}
