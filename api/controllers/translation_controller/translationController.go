package translation_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_domain"
	"shakespearing-pokemon/api/services"
)

func HandleShakespeareanPokemonTranslationRequest(c *gin.Context) {
	request := shksprean_pokemon_domain.ShakespeareanPokemonRequest{
		Name: c.Param("pokemonName"),
	}

	response, apiError := services.TranslationService.GetShakespeareanPokemonTranslation(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}

	c.JSON(http.StatusOK, response)
}
