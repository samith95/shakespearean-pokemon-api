package pokemon_domain

type PokemonInfoRequest struct {
	Name string
}

//Used to parse and store json responses containing the desired pokemon info in the form of:
//	{
//		flavor_text_entries": [
//			{
//				"flavor_text": "Spits fire that\nis hot enough to\nmelt boulders.\fKnown to cause\nforest fires\nunintentionally.",
//				"language": {
//					"name": "en",
//					"url": "https://pokeapi.co/api/v2/language/9/"
//				},
//			},
//			{
//				"flavor_text": "Spits fire that\nis hot enough to\nmelt boulders.\fKnown to cause\nforest fires\nunintentionally.",
//				"language": {
//					"name": "en",
//					"url": "https://pokeapi.co/api/v2/language/9/"
//				},
//			}
//		],
//		"name": "charizard",
//	}
type PokemonInfoResponse struct {
	Name        string          `json:"name"`
	Description FlavourTextList `json:"flavor_text_entries"`
}

type FlavourTextList []FlavourText

type FlavourText struct {
	Text     string `json:"flavor_text"`
	Language LanguageFields
}

type LanguageFields struct {
	Name string `json:"name"`
}
