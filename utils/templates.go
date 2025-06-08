package utils

import (
	"bytes"
	"math/rand"
	"mime/multipart"
	"strings"
	"text/template"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type FuncMapGenerator struct {
	bodyDataHeader string
	localFaker     *gofakeit.Faker
	funcMap        *template.FuncMap
}

func NewFuncMapGenerator(localRand *rand.Rand) *FuncMapGenerator {
	f := &FuncMapGenerator{
		localFaker: gofakeit.NewFaker(localRand, false),
	}
	f.funcMap = f.newFuncMap()

	return f
}

func (g *FuncMapGenerator) GetBodyDataHeader() string {
	tempHeader := g.bodyDataHeader
	g.bodyDataHeader = ""
	return tempHeader
}

func (g *FuncMapGenerator) GetFuncMap() *template.FuncMap {
	return g.funcMap
}

// NewFuncMap creates a template.FuncMap populated with string manipulation functions
// and data generation functions from gofakeit.
//
// It takes a random number generator that is used to initialize a localized faker
// instance, ensuring that random data generation is deterministic within a request context.
//
// All functions are prefixed to avoid naming conflicts:
//   - String functions: "strings_*"
//   - Dict functions: "dict_*"
//   - Body functions: "body_*"
//   - Data generation functions: "fakeit_*"
func (g *FuncMapGenerator) newFuncMap() *template.FuncMap {
	return &template.FuncMap{
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
		"strings_First": func(s string, n int) string {
			if n >= len(s) {
				return s
			}
			return s[:n]
		},
		"strings_Last": func(s string, n int) string {
			if n >= len(s) {
				return s
			}
			return s[len(s)-n:]
		},
		"strings_Truncate": func(s string, n int) string {
			if n >= len(s) {
				return s
			}
			return s[:n] + "..."
		},
		"strings_TrimPrefix": strings.TrimPrefix,
		"strings_TrimSuffix": strings.TrimSuffix,
		"strings_Join": func(sep string, values ...string) string {
			return strings.Join(values, sep)
		},

		// Dict
		"dict_Str": func(values ...string) map[string]string {
			dict := make(map[string]string)
			for i := 0; i < len(values); i += 2 {
				if i+1 < len(values) {
					key := values[i]
					value := values[i+1]
					dict[key] = value
				}
			}
			return dict
		},

		// Slice
		"slice_Str":  func(values ...string) []string { return values },
		"slice_Int":  func(values ...int) []int { return values },
		"slice_Uint": func(values ...uint) []uint { return values },

		// Body
		"body_FormData": func(kv map[string]string) string {
			var data bytes.Buffer
			writer := multipart.NewWriter(&data)

			for k, v := range kv {
				_ = writer.WriteField(k, v)
			}

			_ = writer.Close()
			g.bodyDataHeader = writer.FormDataContentType()

			return data.String()
		},

		// FakeIt / Product
		"fakeit_ProductName":        g.localFaker.ProductName,
		"fakeit_ProductDescription": g.localFaker.ProductDescription,
		"fakeit_ProductCategory":    g.localFaker.ProductCategory,
		"fakeit_ProductFeature":     g.localFaker.ProductFeature,
		"fakeit_ProductMaterial":    g.localFaker.ProductMaterial,
		"fakeit_ProductUPC":         g.localFaker.ProductUPC,
		"fakeit_ProductAudience":    g.localFaker.ProductAudience,
		"fakeit_ProductDimension":   g.localFaker.ProductDimension,
		"fakeit_ProductUseCase":     g.localFaker.ProductUseCase,
		"fakeit_ProductBenefit":     g.localFaker.ProductBenefit,
		"fakeit_ProductSuffix":      g.localFaker.ProductSuffix,

		// FakeIt / Person
		"fakeit_Name":           g.localFaker.Name,
		"fakeit_NamePrefix":     g.localFaker.NamePrefix,
		"fakeit_NameSuffix":     g.localFaker.NameSuffix,
		"fakeit_FirstName":      g.localFaker.FirstName,
		"fakeit_MiddleName":     g.localFaker.MiddleName,
		"fakeit_LastName":       g.localFaker.LastName,
		"fakeit_Gender":         g.localFaker.Gender,
		"fakeit_SSN":            g.localFaker.SSN,
		"fakeit_Hobby":          g.localFaker.Hobby,
		"fakeit_Email":          g.localFaker.Email,
		"fakeit_Phone":          g.localFaker.Phone,
		"fakeit_PhoneFormatted": g.localFaker.PhoneFormatted,

		// FakeIt / Auth
		"fakeit_Username": g.localFaker.Username,
		"fakeit_Password": g.localFaker.Password,

		// FakeIt / Address
		"fakeit_City":         g.localFaker.City,
		"fakeit_Country":      g.localFaker.Country,
		"fakeit_CountryAbr":   g.localFaker.CountryAbr,
		"fakeit_State":        g.localFaker.State,
		"fakeit_StateAbr":     g.localFaker.StateAbr,
		"fakeit_Street":       g.localFaker.Street,
		"fakeit_StreetName":   g.localFaker.StreetName,
		"fakeit_StreetNumber": g.localFaker.StreetNumber,
		"fakeit_StreetPrefix": g.localFaker.StreetPrefix,
		"fakeit_StreetSuffix": g.localFaker.StreetSuffix,
		"fakeit_Zip":          g.localFaker.Zip,
		"fakeit_Latitude":     g.localFaker.Latitude,
		"fakeit_LatitudeInRange": func(min, max float64) float64 {
			value, err := g.localFaker.LatitudeInRange(min, max)
			if err != nil {
				var zero float64
				return zero
			}
			return value
		},
		"fakeit_Longitude": g.localFaker.Longitude,
		"fakeit_LongitudeInRange": func(min, max float64) float64 {
			value, err := g.localFaker.LongitudeInRange(min, max)
			if err != nil {
				var zero float64
				return zero
			}
			return value
		},

		// FakeIt / Game
		"fakeit_Gamertag": g.localFaker.Gamertag,

		// FakeIt / Beer
		"fakeit_BeerAlcohol": g.localFaker.BeerAlcohol,
		"fakeit_BeerBlg":     g.localFaker.BeerBlg,
		"fakeit_BeerHop":     g.localFaker.BeerHop,
		"fakeit_BeerIbu":     g.localFaker.BeerIbu,
		"fakeit_BeerMalt":    g.localFaker.BeerMalt,
		"fakeit_BeerName":    g.localFaker.BeerName,
		"fakeit_BeerStyle":   g.localFaker.BeerStyle,
		"fakeit_BeerYeast":   g.localFaker.BeerYeast,

		// FakeIt / Car
		"fakeit_CarMaker":            g.localFaker.CarMaker,
		"fakeit_CarModel":            g.localFaker.CarModel,
		"fakeit_CarType":             g.localFaker.CarType,
		"fakeit_CarFuelType":         g.localFaker.CarFuelType,
		"fakeit_CarTransmissionType": g.localFaker.CarTransmissionType,

		// FakeIt / Words
		"fakeit_Noun":                      g.localFaker.Noun,
		"fakeit_NounCommon":                g.localFaker.NounCommon,
		"fakeit_NounConcrete":              g.localFaker.NounConcrete,
		"fakeit_NounAbstract":              g.localFaker.NounAbstract,
		"fakeit_NounCollectivePeople":      g.localFaker.NounCollectivePeople,
		"fakeit_NounCollectiveAnimal":      g.localFaker.NounCollectiveAnimal,
		"fakeit_NounCollectiveThing":       g.localFaker.NounCollectiveThing,
		"fakeit_NounCountable":             g.localFaker.NounCountable,
		"fakeit_NounUncountable":           g.localFaker.NounUncountable,
		"fakeit_Verb":                      g.localFaker.Verb,
		"fakeit_VerbAction":                g.localFaker.VerbAction,
		"fakeit_VerbLinking":               g.localFaker.VerbLinking,
		"fakeit_VerbHelping":               g.localFaker.VerbHelping,
		"fakeit_Adverb":                    g.localFaker.Adverb,
		"fakeit_AdverbManner":              g.localFaker.AdverbManner,
		"fakeit_AdverbDegree":              g.localFaker.AdverbDegree,
		"fakeit_AdverbPlace":               g.localFaker.AdverbPlace,
		"fakeit_AdverbTimeDefinite":        g.localFaker.AdverbTimeDefinite,
		"fakeit_AdverbTimeIndefinite":      g.localFaker.AdverbTimeIndefinite,
		"fakeit_AdverbFrequencyDefinite":   g.localFaker.AdverbFrequencyDefinite,
		"fakeit_AdverbFrequencyIndefinite": g.localFaker.AdverbFrequencyIndefinite,
		"fakeit_Preposition":               g.localFaker.Preposition,
		"fakeit_PrepositionSimple":         g.localFaker.PrepositionSimple,
		"fakeit_PrepositionDouble":         g.localFaker.PrepositionDouble,
		"fakeit_PrepositionCompound":       g.localFaker.PrepositionCompound,
		"fakeit_Adjective":                 g.localFaker.Adjective,
		"fakeit_AdjectiveDescriptive":      g.localFaker.AdjectiveDescriptive,
		"fakeit_AdjectiveQuantitative":     g.localFaker.AdjectiveQuantitative,
		"fakeit_AdjectiveProper":           g.localFaker.AdjectiveProper,
		"fakeit_AdjectiveDemonstrative":    g.localFaker.AdjectiveDemonstrative,
		"fakeit_AdjectivePossessive":       g.localFaker.AdjectivePossessive,
		"fakeit_AdjectiveInterrogative":    g.localFaker.AdjectiveInterrogative,
		"fakeit_AdjectiveIndefinite":       g.localFaker.AdjectiveIndefinite,
		"fakeit_Pronoun":                   g.localFaker.Pronoun,
		"fakeit_PronounPersonal":           g.localFaker.PronounPersonal,
		"fakeit_PronounObject":             g.localFaker.PronounObject,
		"fakeit_PronounPossessive":         g.localFaker.PronounPossessive,
		"fakeit_PronounReflective":         g.localFaker.PronounReflective,
		"fakeit_PronounDemonstrative":      g.localFaker.PronounDemonstrative,
		"fakeit_PronounInterrogative":      g.localFaker.PronounInterrogative,
		"fakeit_PronounRelative":           g.localFaker.PronounRelative,
		"fakeit_Connective":                g.localFaker.Connective,
		"fakeit_ConnectiveTime":            g.localFaker.ConnectiveTime,
		"fakeit_ConnectiveComparative":     g.localFaker.ConnectiveComparative,
		"fakeit_ConnectiveComplaint":       g.localFaker.ConnectiveComplaint,
		"fakeit_ConnectiveListing":         g.localFaker.ConnectiveListing,
		"fakeit_ConnectiveCasual":          g.localFaker.ConnectiveCasual,
		"fakeit_ConnectiveExamplify":       g.localFaker.ConnectiveExamplify,
		"fakeit_Word":                      g.localFaker.Word,
		"fakeit_Sentence":                  g.localFaker.Sentence,
		"fakeit_Paragraph":                 g.localFaker.Paragraph,
		"fakeit_LoremIpsumWord":            g.localFaker.LoremIpsumWord,
		"fakeit_LoremIpsumSentence":        g.localFaker.LoremIpsumSentence,
		"fakeit_LoremIpsumParagraph":       g.localFaker.LoremIpsumParagraph,
		"fakeit_Question":                  g.localFaker.Question,
		"fakeit_Quote":                     g.localFaker.Quote,
		"fakeit_Phrase":                    g.localFaker.Phrase,

		// FakeIt / Foods
		"fakeit_Fruit":     g.localFaker.Fruit,
		"fakeit_Vegetable": g.localFaker.Vegetable,
		"fakeit_Breakfast": g.localFaker.Breakfast,
		"fakeit_Lunch":     g.localFaker.Lunch,
		"fakeit_Dinner":    g.localFaker.Dinner,
		"fakeit_Snack":     g.localFaker.Snack,
		"fakeit_Dessert":   g.localFaker.Dessert,

		// FakeIt / Misc
		"fakeit_Bool":      g.localFaker.Bool,
		"fakeit_UUID":      g.localFaker.UUID,
		"fakeit_FlipACoin": g.localFaker.FlipACoin,

		// FakeIt / Colors
		"fakeit_Color":      g.localFaker.Color,
		"fakeit_HexColor":   g.localFaker.HexColor,
		"fakeit_RGBColor":   g.localFaker.RGBColor,
		"fakeit_SafeColor":  g.localFaker.SafeColor,
		"fakeit_NiceColors": g.localFaker.NiceColors,

		// FakeIt / Internet
		"fakeit_URL":                  g.localFaker.URL,
		"fakeit_DomainName":           g.localFaker.DomainName,
		"fakeit_DomainSuffix":         g.localFaker.DomainSuffix,
		"fakeit_IPv4Address":          g.localFaker.IPv4Address,
		"fakeit_IPv6Address":          g.localFaker.IPv6Address,
		"fakeit_MacAddress":           g.localFaker.MacAddress,
		"fakeit_HTTPStatusCode":       g.localFaker.HTTPStatusCode,
		"fakeit_HTTPStatusCodeSimple": g.localFaker.HTTPStatusCodeSimple,
		"fakeit_LogLevel":             g.localFaker.LogLevel,
		"fakeit_HTTPMethod":           g.localFaker.HTTPMethod,
		"fakeit_HTTPVersion":          g.localFaker.HTTPVersion,
		"fakeit_UserAgent":            g.localFaker.UserAgent,
		"fakeit_ChromeUserAgent":      g.localFaker.ChromeUserAgent,
		"fakeit_FirefoxUserAgent":     g.localFaker.FirefoxUserAgent,
		"fakeit_OperaUserAgent":       g.localFaker.OperaUserAgent,
		"fakeit_SafariUserAgent":      g.localFaker.SafariUserAgent,

		// FakeIt / HTML
		"fakeit_InputName": g.localFaker.InputName,

		// FakeIt / Date/Time
		"fakeit_Date":           g.localFaker.Date,
		"fakeit_PastDate":       g.localFaker.PastDate,
		"fakeit_FutureDate":     g.localFaker.FutureDate,
		"fakeit_DateRange":      g.localFaker.DateRange,
		"fakeit_NanoSecond":     g.localFaker.NanoSecond,
		"fakeit_Second":         g.localFaker.Second,
		"fakeit_Minute":         g.localFaker.Minute,
		"fakeit_Hour":           g.localFaker.Hour,
		"fakeit_Month":          g.localFaker.Month,
		"fakeit_MonthString":    g.localFaker.MonthString,
		"fakeit_Day":            g.localFaker.Day,
		"fakeit_WeekDay":        g.localFaker.WeekDay,
		"fakeit_Year":           g.localFaker.Year,
		"fakeit_TimeZone":       g.localFaker.TimeZone,
		"fakeit_TimeZoneAbv":    g.localFaker.TimeZoneAbv,
		"fakeit_TimeZoneFull":   g.localFaker.TimeZoneFull,
		"fakeit_TimeZoneOffset": g.localFaker.TimeZoneOffset,
		"fakeit_TimeZoneRegion": g.localFaker.TimeZoneRegion,

		// FakeIt / Payment
		"fakeit_Price":             g.localFaker.Price,
		"fakeit_CreditCardCvv":     g.localFaker.CreditCardCvv,
		"fakeit_CreditCardExp":     g.localFaker.CreditCardExp,
		"fakeit_CreditCardNumber":  g.localFaker.CreditCardNumber,
		"fakeit_CreditCardType":    g.localFaker.CreditCardType,
		"fakeit_CurrencyLong":      g.localFaker.CurrencyLong,
		"fakeit_CurrencyShort":     g.localFaker.CurrencyShort,
		"fakeit_AchRouting":        g.localFaker.AchRouting,
		"fakeit_AchAccount":        g.localFaker.AchAccount,
		"fakeit_BitcoinAddress":    g.localFaker.BitcoinAddress,
		"fakeit_BitcoinPrivateKey": g.localFaker.BitcoinPrivateKey,

		// FakeIt / Finance
		"fakeit_Cusip": g.localFaker.Cusip,
		"fakeit_Isin":  g.localFaker.Isin,

		// FakeIt / Company
		"fakeit_BS":            g.localFaker.BS,
		"fakeit_Blurb":         g.localFaker.Blurb,
		"fakeit_BuzzWord":      g.localFaker.BuzzWord,
		"fakeit_Company":       g.localFaker.Company,
		"fakeit_CompanySuffix": g.localFaker.CompanySuffix,
		"fakeit_JobDescriptor": g.localFaker.JobDescriptor,
		"fakeit_JobLevel":      g.localFaker.JobLevel,
		"fakeit_JobTitle":      g.localFaker.JobTitle,
		"fakeit_Slogan":        g.localFaker.Slogan,

		// FakeIt / Hacker
		"fakeit_HackerAbbreviation": g.localFaker.HackerAbbreviation,
		"fakeit_HackerAdjective":    g.localFaker.HackerAdjective,
		"fakeit_HackerNoun":         g.localFaker.HackerNoun,
		"fakeit_HackerPhrase":       g.localFaker.HackerPhrase,
		"fakeit_HackerVerb":         g.localFaker.HackerVerb,

		// FakeIt / Hipster
		"fakeit_HipsterWord":      g.localFaker.HipsterWord,
		"fakeit_HipsterSentence":  g.localFaker.HipsterSentence,
		"fakeit_HipsterParagraph": g.localFaker.HipsterParagraph,

		// FakeIt / App
		"fakeit_AppName":    g.localFaker.AppName,
		"fakeit_AppVersion": g.localFaker.AppVersion,
		"fakeit_AppAuthor":  g.localFaker.AppAuthor,

		// FakeIt / Animal
		"fakeit_PetName":    g.localFaker.PetName,
		"fakeit_Animal":     g.localFaker.Animal,
		"fakeit_AnimalType": g.localFaker.AnimalType,
		"fakeit_FarmAnimal": g.localFaker.FarmAnimal,
		"fakeit_Cat":        g.localFaker.Cat,
		"fakeit_Dog":        g.localFaker.Dog,
		"fakeit_Bird":       g.localFaker.Bird,

		// FakeIt / Emoji
		"fakeit_Emoji":            g.localFaker.Emoji,
		"fakeit_EmojiDescription": g.localFaker.EmojiDescription,
		"fakeit_EmojiCategory":    g.localFaker.EmojiCategory,
		"fakeit_EmojiAlias":       g.localFaker.EmojiAlias,
		"fakeit_EmojiTag":         g.localFaker.EmojiTag,

		// FakeIt / Language
		"fakeit_Language":             g.localFaker.Language,
		"fakeit_LanguageAbbreviation": g.localFaker.LanguageAbbreviation,
		"fakeit_ProgrammingLanguage":  g.localFaker.ProgrammingLanguage,

		// FakeIt / Number
		"fakeit_Number":       g.localFaker.Number,
		"fakeit_Int":          g.localFaker.Int,
		"fakeit_IntN":         g.localFaker.IntN,
		"fakeit_IntRange":     g.localFaker.IntRange,
		"fakeit_RandomInt":    g.localFaker.RandomInt,
		"fakeit_Int8":         g.localFaker.Int8,
		"fakeit_Int16":        g.localFaker.Int16,
		"fakeit_Int32":        g.localFaker.Int32,
		"fakeit_Int64":        g.localFaker.Int64,
		"fakeit_Uint":         g.localFaker.Uint,
		"fakeit_UintN":        g.localFaker.UintN,
		"fakeit_UintRange":    g.localFaker.UintRange,
		"fakeit_RandomUint":   g.localFaker.RandomUint,
		"fakeit_Uint8":        g.localFaker.Uint8,
		"fakeit_Uint16":       g.localFaker.Uint16,
		"fakeit_Uint32":       g.localFaker.Uint32,
		"fakeit_Uint64":       g.localFaker.Uint64,
		"fakeit_Float32":      g.localFaker.Float32,
		"fakeit_Float32Range": g.localFaker.Float32Range,
		"fakeit_Float64":      g.localFaker.Float64,
		"fakeit_Float64Range": g.localFaker.Float64Range,
		"fakeit_HexUint":      g.localFaker.HexUint,

		// FakeIt / String
		"fakeit_Digit":   g.localFaker.Digit,
		"fakeit_DigitN":  g.localFaker.DigitN,
		"fakeit_Letter":  g.localFaker.Letter,
		"fakeit_LetterN": g.localFaker.LetterN,
		"fakeit_LetterNN": func(min, max uint) string {
			return g.localFaker.LetterN(g.localFaker.UintRange(min, max))
		},
		"fakeit_Lexify":   g.localFaker.Lexify,
		"fakeit_Numerify": g.localFaker.Numerify,
		"fakeit_RandomString": func(values ...string) string {
			return g.localFaker.RandomString(values)
		},

		// FakeIt / Celebrity
		"fakeit_CelebrityActor":    g.localFaker.CelebrityActor,
		"fakeit_CelebrityBusiness": g.localFaker.CelebrityBusiness,
		"fakeit_CelebritySport":    g.localFaker.CelebritySport,

		// FakeIt / Minecraft
		"fakeit_MinecraftOre":             g.localFaker.MinecraftOre,
		"fakeit_MinecraftWood":            g.localFaker.MinecraftWood,
		"fakeit_MinecraftArmorTier":       g.localFaker.MinecraftArmorTier,
		"fakeit_MinecraftArmorPart":       g.localFaker.MinecraftArmorPart,
		"fakeit_MinecraftWeapon":          g.localFaker.MinecraftWeapon,
		"fakeit_MinecraftTool":            g.localFaker.MinecraftTool,
		"fakeit_MinecraftDye":             g.localFaker.MinecraftDye,
		"fakeit_MinecraftFood":            g.localFaker.MinecraftFood,
		"fakeit_MinecraftAnimal":          g.localFaker.MinecraftAnimal,
		"fakeit_MinecraftVillagerJob":     g.localFaker.MinecraftVillagerJob,
		"fakeit_MinecraftVillagerStation": g.localFaker.MinecraftVillagerStation,
		"fakeit_MinecraftVillagerLevel":   g.localFaker.MinecraftVillagerLevel,
		"fakeit_MinecraftMobPassive":      g.localFaker.MinecraftMobPassive,
		"fakeit_MinecraftMobNeutral":      g.localFaker.MinecraftMobNeutral,
		"fakeit_MinecraftMobHostile":      g.localFaker.MinecraftMobHostile,
		"fakeit_MinecraftMobBoss":         g.localFaker.MinecraftMobBoss,
		"fakeit_MinecraftBiome":           g.localFaker.MinecraftBiome,
		"fakeit_MinecraftWeather":         g.localFaker.MinecraftWeather,

		// FakeIt / Book
		"fakeit_BookTitle":  g.localFaker.BookTitle,
		"fakeit_BookAuthor": g.localFaker.BookAuthor,
		"fakeit_BookGenre":  g.localFaker.BookGenre,

		// FakeIt / Movie
		"fakeit_MovieName":  g.localFaker.MovieName,
		"fakeit_MovieGenre": g.localFaker.MovieGenre,

		// FakeIt / Error
		"fakeit_Error":           g.localFaker.Error,
		"fakeit_ErrorDatabase":   g.localFaker.ErrorDatabase,
		"fakeit_ErrorGRPC":       g.localFaker.ErrorGRPC,
		"fakeit_ErrorHTTP":       g.localFaker.ErrorHTTP,
		"fakeit_ErrorHTTPClient": g.localFaker.ErrorHTTPClient,
		"fakeit_ErrorHTTPServer": g.localFaker.ErrorHTTPServer,
		"fakeit_ErrorRuntime":    g.localFaker.ErrorRuntime,

		// FakeIt / School
		"fakeit_School": g.localFaker.School,

		// FakeIt / Song
		"fakeit_SongName":   g.localFaker.SongName,
		"fakeit_SongArtist": g.localFaker.SongArtist,
		"fakeit_SongGenre":  g.localFaker.SongGenre,
	}
}
