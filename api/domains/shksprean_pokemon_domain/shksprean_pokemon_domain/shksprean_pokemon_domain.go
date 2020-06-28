package shksprean_pokemon_domain

type ShakespeareanPokemonRequest struct {
	Name string
}

//Used to store and generate Shakespearean translation of the pokemon's description in the form of:
//		{
//			"name": "charizard",
//			"description": "translated version of charizard's description",
//		}
type ShakespeareanPokemonResponse struct {
	Name        string `json:"name"`
	Translation string `json:"description"`
}
