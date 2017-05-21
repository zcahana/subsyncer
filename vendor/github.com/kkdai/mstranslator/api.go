package mstranslator

type Client struct {
	translateProvider *TranslateProvider
	languageProvider  *LanguageProvider
	authenicator      *Authenicator
}

func NewClient(clientID, clientSecret string) *Client {
	auth := NewAuthenicator(clientID, clientSecret)
	if auth == nil {
		return nil
	}
	//Retreive token first to avoid double request on each provider.
	auth.GetToken()

	return &Client{
		translateProvider: NewTranslateProvider(auth),
		languageProvider:  NewLanguageProvider(auth),
		authenicator:      auth,
	}
}

//Translates a text string from one language to another.
func (c *Client) Translate(text, from, to string) (string, error) {
	return c.translateProvider.Translate(text, from, to)
}

//The TransformText method is a text normalization function for social media, which returns a normalized form of the input.
//The method can be used as a preprocessing step in Machine Translation or other applications, which expect clean input text than is typically found in social media or user-generated content. The function currently works only with English input.
func (c *Client) TransformText(lang, category, text string) (string, error) {
	return c.translateProvider.TransformText(lang, category, text)
}

// Returns a wave or mp3 stream of the passed-in text being spoken in the desired language.
func (c *Client) Speak(text, lang, outFormat string) ([]byte, error) {
	return c.translateProvider.Speak(text, lang, outFormat)
}

// Use the Detect Method to identify the language of a selected piece of text.
func (c *Client) Detect(text string) (string, error) {
	return c.languageProvider.Detect(text)
}

// Use the DetectArray Method to identify the language of an array of string at once. Performs independent detection of each individual array element and returns a result for each row of the array.
func (c *Client) DetectArray(textArray []string) ([]string, error) {
	return c.languageProvider.DetectArray(textArray)
}

//Retrieves an array of translations for a given language pair from the store and the MT engine. GetTranslations differs from Translate as it returns all available translations.
func (c *Client) GetTranslations(text, from, to string, maxTranslations int) ([]ResponseTranslationMatch, error) {
	return c.languageProvider.GetTranslations(text, from, to, maxTranslations)
}

//Retrieves friendly names for the languages passed in as the parameter languageCodes, and localized using the passed locale language.
func (c *Client) GetLanguageNames(codes []string) ([]string, error) {
	return c.languageProvider.GetLanguageNames(codes)
}

//Obtain a list of language codes representing languages that are supported by the Translation Service. Translate() and TranslateArray() can translate between any two of these languages.
func (c *Client) GetLanguagesForTranslate() ([]string, error) {
	return c.languageProvider.GetLanguagesForTranslate()
}

//Retrieves the languages available for speech synthesis.
func (c *Client) GetLanguagesForSpeak() ([]string, error) {
	return c.languageProvider.GetLanguagesForSpeak()
}
