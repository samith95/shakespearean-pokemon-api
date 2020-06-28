package translation_domain

type TranslationRequest struct {
	Text string
}

//Used to parse and store json responses containing the translation in the form of:
//	{
//		"success": {
//			"total": 1
//		},
//		"contents": {
//			"translated": "\"CHARIZARD flies around the sky powerful opponents.\\nIt breathes fire of such great heat\\
//							fthat it melts anything. However, it\\nnever turns its fiery breath on any\\nopponent weaker
//							than itself.\"",
//			"text": "\"CHARIZARD flies around the sky powerful opponents.\\nIt breathes fire of such great heat\\fthat
//					it melts anything. However, it\\nnever turns its fiery breath on any\\nopponent weaker than itself.\"",
//			"translation": "shakespeare"
//		}
//	}
type TranslationResponse struct {
	Content ContentFields
}

type ContentFields struct {
	Translation string `json:"translated"`
}
