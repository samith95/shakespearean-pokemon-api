package app

import "shakespearing-pokemon/api/controllers/translation_controller"

func routes() {
	router.GET("/pokemon/:pokemonName", translation_controller.HandleShakespeareanPokemonTranslationRequest)
}
