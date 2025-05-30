package utils

import (
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// NewFuncMap creates a template.FuncMap populated with string manipulation functions
// and data generation functions from gofakeit.
//
// It takes a random number generator that is used to initialize a localized faker
// instance, ensuring that random data generation is deterministic within a request context.
//
// All functions are prefixed to avoid naming conflicts:
//   - String functions: "strings_*"
//   - Data generation functions: "fakeit_*"
func NewFuncMap(localRand *rand.Rand) template.FuncMap {
	localFaker := gofakeit.NewFaker(localRand, false)

	return template.FuncMap{
		// Strings
		"strings_ToUpper":      strings.ToUpper,
		"strings_ToLower":      strings.ToLower,
		"strings_RemoveSpaces": func(s string) string { return strings.ReplaceAll(s, " ", "") },
		"strings_Replace":      strings.Replace,
		"strings_ToDate": func(dateString string) time.Time {
			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				return time.Now()
			}
			return date
		},

		// FakeIt / Product
		"fakeit_ProductName":        localFaker.ProductName,
		"fakeit_ProductDescription": localFaker.ProductDescription,
		"fakeit_ProductCategory":    localFaker.ProductCategory,
		"fakeit_ProductFeature":     localFaker.ProductFeature,
		"fakeit_ProductMaterial":    localFaker.ProductMaterial,
		"fakeit_ProductUPC":         localFaker.ProductUPC,
		"fakeit_ProductAudience":    localFaker.ProductAudience,
		"fakeit_ProductDimension":   localFaker.ProductDimension,
		"fakeit_ProductUseCase":     localFaker.ProductUseCase,
		"fakeit_ProductBenefit":     localFaker.ProductBenefit,
		"fakeit_ProductSuffix":      localFaker.ProductSuffix,

		// FakeIt / Person
		"fakeit_Name":           localFaker.Name,
		"fakeit_NamePrefix":     localFaker.NamePrefix,
		"fakeit_NameSuffix":     localFaker.NameSuffix,
		"fakeit_FirstName":      localFaker.FirstName,
		"fakeit_MiddleName":     localFaker.MiddleName,
		"fakeit_LastName":       localFaker.LastName,
		"fakeit_Gender":         localFaker.Gender,
		"fakeit_SSN":            localFaker.SSN,
		"fakeit_Hobby":          localFaker.Hobby,
		"fakeit_Email":          localFaker.Email,
		"fakeit_Phone":          localFaker.Phone,
		"fakeit_PhoneFormatted": localFaker.PhoneFormatted,

		// FakeIt / Auth
		"fakeit_Username": localFaker.Username,
		"fakeit_Password": localFaker.Password,

		// FakeIt / Address
		"fakeit_City":         localFaker.City,
		"fakeit_Country":      localFaker.Country,
		"fakeit_CountryAbr":   localFaker.CountryAbr,
		"fakeit_State":        localFaker.State,
		"fakeit_StateAbr":     localFaker.StateAbr,
		"fakeit_Street":       localFaker.Street,
		"fakeit_StreetName":   localFaker.StreetName,
		"fakeit_StreetNumber": localFaker.StreetNumber,
		"fakeit_StreetPrefix": localFaker.StreetPrefix,
		"fakeit_StreetSuffix": localFaker.StreetSuffix,
		"fakeit_Zip":          localFaker.Zip,
		"fakeit_Latitude":     localFaker.Latitude,
		"fakeit_LatitudeInRange": func(min, max float64) float64 {
			value, err := localFaker.LatitudeInRange(min, max)
			if err != nil {
				var zero float64
				return zero
			}
			return value
		},
		"fakeit_Longitude": localFaker.Longitude,
		"fakeit_LongitudeInRange": func(min, max float64) float64 {
			value, err := localFaker.LongitudeInRange(min, max)
			if err != nil {
				var zero float64
				return zero
			}
			return value
		},

		// FakeIt / Game
		"fakeit_Gamertag": localFaker.Gamertag,

		// FakeIt / Beer
		"fakeit_BeerAlcohol": localFaker.BeerAlcohol,
		"fakeit_BeerBlg":     localFaker.BeerBlg,
		"fakeit_BeerHop":     localFaker.BeerHop,
		"fakeit_BeerIbu":     localFaker.BeerIbu,
		"fakeit_BeerMalt":    localFaker.BeerMalt,
		"fakeit_BeerName":    localFaker.BeerName,
		"fakeit_BeerStyle":   localFaker.BeerStyle,
		"fakeit_BeerYeast":   localFaker.BeerYeast,

		// FakeIt / Car
		"fakeit_CarMaker":            localFaker.CarMaker,
		"fakeit_CarModel":            localFaker.CarModel,
		"fakeit_CarType":             localFaker.CarType,
		"fakeit_CarFuelType":         localFaker.CarFuelType,
		"fakeit_CarTransmissionType": localFaker.CarTransmissionType,

		// FakeIt / Words
		"fakeit_Noun":                      localFaker.Noun,
		"fakeit_NounCommon":                localFaker.NounCommon,
		"fakeit_NounConcrete":              localFaker.NounConcrete,
		"fakeit_NounAbstract":              localFaker.NounAbstract,
		"fakeit_NounCollectivePeople":      localFaker.NounCollectivePeople,
		"fakeit_NounCollectiveAnimal":      localFaker.NounCollectiveAnimal,
		"fakeit_NounCollectiveThing":       localFaker.NounCollectiveThing,
		"fakeit_NounCountable":             localFaker.NounCountable,
		"fakeit_NounUncountable":           localFaker.NounUncountable,
		"fakeit_Verb":                      localFaker.Verb,
		"fakeit_VerbAction":                localFaker.VerbAction,
		"fakeit_VerbLinking":               localFaker.VerbLinking,
		"fakeit_VerbHelping":               localFaker.VerbHelping,
		"fakeit_Adverb":                    localFaker.Adverb,
		"fakeit_AdverbManner":              localFaker.AdverbManner,
		"fakeit_AdverbDegree":              localFaker.AdverbDegree,
		"fakeit_AdverbPlace":               localFaker.AdverbPlace,
		"fakeit_AdverbTimeDefinite":        localFaker.AdverbTimeDefinite,
		"fakeit_AdverbTimeIndefinite":      localFaker.AdverbTimeIndefinite,
		"fakeit_AdverbFrequencyDefinite":   localFaker.AdverbFrequencyDefinite,
		"fakeit_AdverbFrequencyIndefinite": localFaker.AdverbFrequencyIndefinite,
		"fakeit_Preposition":               localFaker.Preposition,
		"fakeit_PrepositionSimple":         localFaker.PrepositionSimple,
		"fakeit_PrepositionDouble":         localFaker.PrepositionDouble,
		"fakeit_PrepositionCompound":       localFaker.PrepositionCompound,
		"fakeit_Adjective":                 localFaker.Adjective,
		"fakeit_AdjectiveDescriptive":      localFaker.AdjectiveDescriptive,
		"fakeit_AdjectiveQuantitative":     localFaker.AdjectiveQuantitative,
		"fakeit_AdjectiveProper":           localFaker.AdjectiveProper,
		"fakeit_AdjectiveDemonstrative":    localFaker.AdjectiveDemonstrative,
		"fakeit_AdjectivePossessive":       localFaker.AdjectivePossessive,
		"fakeit_AdjectiveInterrogative":    localFaker.AdjectiveInterrogative,
		"fakeit_AdjectiveIndefinite":       localFaker.AdjectiveIndefinite,
		"fakeit_Pronoun":                   localFaker.Pronoun,
		"fakeit_PronounPersonal":           localFaker.PronounPersonal,
		"fakeit_PronounObject":             localFaker.PronounObject,
		"fakeit_PronounPossessive":         localFaker.PronounPossessive,
		"fakeit_PronounReflective":         localFaker.PronounReflective,
		"fakeit_PronounDemonstrative":      localFaker.PronounDemonstrative,
		"fakeit_PronounInterrogative":      localFaker.PronounInterrogative,
		"fakeit_PronounRelative":           localFaker.PronounRelative,
		"fakeit_Connective":                localFaker.Connective,
		"fakeit_ConnectiveTime":            localFaker.ConnectiveTime,
		"fakeit_ConnectiveComparative":     localFaker.ConnectiveComparative,
		"fakeit_ConnectiveComplaint":       localFaker.ConnectiveComplaint,
		"fakeit_ConnectiveListing":         localFaker.ConnectiveListing,
		"fakeit_ConnectiveCasual":          localFaker.ConnectiveCasual,
		"fakeit_ConnectiveExamplify":       localFaker.ConnectiveExamplify,
		"fakeit_Word":                      localFaker.Word,
		"fakeit_Sentence":                  localFaker.Sentence,
		"fakeit_Paragraph":                 localFaker.Paragraph,
		"fakeit_LoremIpsumWord":            localFaker.LoremIpsumWord,
		"fakeit_LoremIpsumSentence":        localFaker.LoremIpsumSentence,
		"fakeit_LoremIpsumParagraph":       localFaker.LoremIpsumParagraph,
		"fakeit_Question":                  localFaker.Question,
		"fakeit_Quote":                     localFaker.Quote,
		"fakeit_Phrase":                    localFaker.Phrase,

		// FakeIt / Foods
		"fakeit_Fruit":     localFaker.Fruit,
		"fakeit_Vegetable": localFaker.Vegetable,
		"fakeit_Breakfast": localFaker.Breakfast,
		"fakeit_Lunch":     localFaker.Lunch,
		"fakeit_Dinner":    localFaker.Dinner,
		"fakeit_Snack":     localFaker.Snack,
		"fakeit_Dessert":   localFaker.Dessert,

		// FakeIt / Misc
		"fakeit_Bool":      localFaker.Bool,
		"fakeit_UUID":      localFaker.UUID,
		"fakeit_FlipACoin": localFaker.FlipACoin,

		// FakeIt / Colors
		"fakeit_Color":      localFaker.Color,
		"fakeit_HexColor":   localFaker.HexColor,
		"fakeit_RGBColor":   localFaker.RGBColor,
		"fakeit_SafeColor":  localFaker.SafeColor,
		"fakeit_NiceColors": localFaker.NiceColors,

		// FakeIt / Internet
		"fakeit_URL":                  localFaker.URL,
		"fakeit_DomainName":           localFaker.DomainName,
		"fakeit_DomainSuffix":         localFaker.DomainSuffix,
		"fakeit_IPv4Address":          localFaker.IPv4Address,
		"fakeit_IPv6Address":          localFaker.IPv6Address,
		"fakeit_MacAddress":           localFaker.MacAddress,
		"fakeit_HTTPStatusCode":       localFaker.HTTPStatusCode,
		"fakeit_HTTPStatusCodeSimple": localFaker.HTTPStatusCodeSimple,
		"fakeit_LogLevel":             localFaker.LogLevel,
		"fakeit_HTTPMethod":           localFaker.HTTPMethod,
		"fakeit_HTTPVersion":          localFaker.HTTPVersion,
		"fakeit_UserAgent":            localFaker.UserAgent,
		"fakeit_ChromeUserAgent":      localFaker.ChromeUserAgent,
		"fakeit_FirefoxUserAgent":     localFaker.FirefoxUserAgent,
		"fakeit_OperaUserAgent":       localFaker.OperaUserAgent,
		"fakeit_SafariUserAgent":      localFaker.SafariUserAgent,

		// FakeIt / HTML
		"fakeit_InputName": localFaker.InputName,

		// FakeIt / Date/Time
		"fakeit_Date":           localFaker.Date,
		"fakeit_PastDate":       localFaker.PastDate,
		"fakeit_FutureDate":     localFaker.FutureDate,
		"fakeit_DateRange":      localFaker.DateRange,
		"fakeit_NanoSecond":     localFaker.NanoSecond,
		"fakeit_Second":         localFaker.Second,
		"fakeit_Minute":         localFaker.Minute,
		"fakeit_Hour":           localFaker.Hour,
		"fakeit_Month":          localFaker.Month,
		"fakeit_MonthString":    localFaker.MonthString,
		"fakeit_Day":            localFaker.Day,
		"fakeit_WeekDay":        localFaker.WeekDay,
		"fakeit_Year":           localFaker.Year,
		"fakeit_TimeZone":       localFaker.TimeZone,
		"fakeit_TimeZoneAbv":    localFaker.TimeZoneAbv,
		"fakeit_TimeZoneFull":   localFaker.TimeZoneFull,
		"fakeit_TimeZoneOffset": localFaker.TimeZoneOffset,
		"fakeit_TimeZoneRegion": localFaker.TimeZoneRegion,

		// FakeIt / Payment
		"fakeit_Price":             localFaker.Price,
		"fakeit_CreditCardCvv":     localFaker.CreditCardCvv,
		"fakeit_CreditCardExp":     localFaker.CreditCardExp,
		"fakeit_CreditCardNumber":  localFaker.CreditCardNumber,
		"fakeit_CreditCardType":    localFaker.CreditCardType,
		"fakeit_CurrencyLong":      localFaker.CurrencyLong,
		"fakeit_CurrencyShort":     localFaker.CurrencyShort,
		"fakeit_AchRouting":        localFaker.AchRouting,
		"fakeit_AchAccount":        localFaker.AchAccount,
		"fakeit_BitcoinAddress":    localFaker.BitcoinAddress,
		"fakeit_BitcoinPrivateKey": localFaker.BitcoinPrivateKey,

		// FakeIt / Finance
		"fakeit_Cusip": localFaker.Cusip,
		"fakeit_Isin":  localFaker.Isin,

		// FakeIt / Company
		"fakeit_BS":            localFaker.BS,
		"fakeit_Blurb":         localFaker.Blurb,
		"fakeit_BuzzWord":      localFaker.BuzzWord,
		"fakeit_Company":       localFaker.Company,
		"fakeit_CompanySuffix": localFaker.CompanySuffix,
		"fakeit_JobDescriptor": localFaker.JobDescriptor,
		"fakeit_JobLevel":      localFaker.JobLevel,
		"fakeit_JobTitle":      localFaker.JobTitle,
		"fakeit_Slogan":        localFaker.Slogan,

		// FakeIt / Hacker
		"fakeit_HackerAbbreviation": localFaker.HackerAbbreviation,
		"fakeit_HackerAdjective":    localFaker.HackerAdjective,
		"fakeit_HackerNoun":         localFaker.HackerNoun,
		"fakeit_HackerPhrase":       localFaker.HackerPhrase,
		"fakeit_HackerVerb":         localFaker.HackerVerb,

		// FakeIt / Hipster
		"fakeit_HipsterWord":      localFaker.HipsterWord,
		"fakeit_HipsterSentence":  localFaker.HipsterSentence,
		"fakeit_HipsterParagraph": localFaker.HipsterParagraph,

		// FakeIt / App
		"fakeit_AppName":    localFaker.AppName,
		"fakeit_AppVersion": localFaker.AppVersion,
		"fakeit_AppAuthor":  localFaker.AppAuthor,

		// FakeIt / Animal
		"fakeit_PetName":    localFaker.PetName,
		"fakeit_Animal":     localFaker.Animal,
		"fakeit_AnimalType": localFaker.AnimalType,
		"fakeit_FarmAnimal": localFaker.FarmAnimal,
		"fakeit_Cat":        localFaker.Cat,
		"fakeit_Dog":        localFaker.Dog,
		"fakeit_Bird":       localFaker.Bird,

		// FakeIt / Emoji
		"fakeit_Emoji":            localFaker.Emoji,
		"fakeit_EmojiDescription": localFaker.EmojiDescription,
		"fakeit_EmojiCategory":    localFaker.EmojiCategory,
		"fakeit_EmojiAlias":       localFaker.EmojiAlias,
		"fakeit_EmojiTag":         localFaker.EmojiTag,

		// FakeIt / Language
		"fakeit_Language":             localFaker.Language,
		"fakeit_LanguageAbbreviation": localFaker.LanguageAbbreviation,
		"fakeit_ProgrammingLanguage":  localFaker.ProgrammingLanguage,

		// FakeIt / Number
		"fakeit_Number":       localFaker.Number,
		"fakeit_Int":          localFaker.Int,
		"fakeit_IntN":         localFaker.IntN,
		"fakeit_Int8":         localFaker.Int8,
		"fakeit_Int16":        localFaker.Int16,
		"fakeit_Int32":        localFaker.Int32,
		"fakeit_Int64":        localFaker.Int64,
		"fakeit_Uint":         localFaker.Uint,
		"fakeit_UintN":        localFaker.UintN,
		"fakeit_Uint8":        localFaker.Uint8,
		"fakeit_Uint16":       localFaker.Uint16,
		"fakeit_Uint32":       localFaker.Uint32,
		"fakeit_Uint64":       localFaker.Uint64,
		"fakeit_Float32":      localFaker.Float32,
		"fakeit_Float32Range": localFaker.Float32Range,
		"fakeit_Float64":      localFaker.Float64,
		"fakeit_Float64Range": localFaker.Float64Range,
		"fakeit_HexUint":      localFaker.HexUint,

		// FakeIt / String
		"fakeit_Digit":    localFaker.Digit,
		"fakeit_DigitN":   localFaker.DigitN,
		"fakeit_Letter":   localFaker.Letter,
		"fakeit_LetterN":  localFaker.LetterN,
		"fakeit_Lexify":   localFaker.Lexify,
		"fakeit_Numerify": localFaker.Numerify,

		// FakeIt / Celebrity
		"fakeit_CelebrityActor":    localFaker.CelebrityActor,
		"fakeit_CelebrityBusiness": localFaker.CelebrityBusiness,
		"fakeit_CelebritySport":    localFaker.CelebritySport,

		// FakeIt / Minecraft
		"fakeit_MinecraftOre":             localFaker.MinecraftOre,
		"fakeit_MinecraftWood":            localFaker.MinecraftWood,
		"fakeit_MinecraftArmorTier":       localFaker.MinecraftArmorTier,
		"fakeit_MinecraftArmorPart":       localFaker.MinecraftArmorPart,
		"fakeit_MinecraftWeapon":          localFaker.MinecraftWeapon,
		"fakeit_MinecraftTool":            localFaker.MinecraftTool,
		"fakeit_MinecraftDye":             localFaker.MinecraftDye,
		"fakeit_MinecraftFood":            localFaker.MinecraftFood,
		"fakeit_MinecraftAnimal":          localFaker.MinecraftAnimal,
		"fakeit_MinecraftVillagerJob":     localFaker.MinecraftVillagerJob,
		"fakeit_MinecraftVillagerStation": localFaker.MinecraftVillagerStation,
		"fakeit_MinecraftVillagerLevel":   localFaker.MinecraftVillagerLevel,
		"fakeit_MinecraftMobPassive":      localFaker.MinecraftMobPassive,
		"fakeit_MinecraftMobNeutral":      localFaker.MinecraftMobNeutral,
		"fakeit_MinecraftMobHostile":      localFaker.MinecraftMobHostile,
		"fakeit_MinecraftMobBoss":         localFaker.MinecraftMobBoss,
		"fakeit_MinecraftBiome":           localFaker.MinecraftBiome,
		"fakeit_MinecraftWeather":         localFaker.MinecraftWeather,

		// FakeIt / Book
		"fakeit_BookTitle":  localFaker.BookTitle,
		"fakeit_BookAuthor": localFaker.BookAuthor,
		"fakeit_BookGenre":  localFaker.BookGenre,

		// FakeIt / Movie
		"fakeit_MovieName":  localFaker.MovieName,
		"fakeit_MovieGenre": localFaker.MovieGenre,

		// FakeIt / Error
		"fakeit_Error":           localFaker.Error,
		"fakeit_ErrorDatabase":   localFaker.ErrorDatabase,
		"fakeit_ErrorGRPC":       localFaker.ErrorGRPC,
		"fakeit_ErrorHTTP":       localFaker.ErrorHTTP,
		"fakeit_ErrorHTTPClient": localFaker.ErrorHTTPClient,
		"fakeit_ErrorHTTPServer": localFaker.ErrorHTTPServer,
		"fakeit_ErrorRuntime":    localFaker.ErrorRuntime,

		// FakeIt / School
		"fakeit_School": localFaker.School,

		// FakeIt / Song
		"fakeit_SongName":   localFaker.SongName,
		"fakeit_SongArtist": localFaker.SongArtist,
		"fakeit_SongGenre":  localFaker.SongGenre,
	}
}
