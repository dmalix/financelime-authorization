package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func InitLanguageContent(file string) (LanguageContent, error) {

	var (
		fileLanguageContent FileLanguageContent
		languageContent     LanguageContent
		body                []byte
		err                 error
	)

	// Load the content of the file and convert it to structure

	body, err = ioutil.ReadFile(file)
	if err != nil {
		return LanguageContent{},
			fmt.Errorf("failed to read the language content file: %s", err)
	}

	err = json.Unmarshal(body, &fileLanguageContent)
	if err != nil {
		return LanguageContent{},
			fmt.Errorf("failed to umarshal the language content body: %s", err)
	}

	// Init content in different languages
	// An example of a call after initialization:
	// content.Data.User.Login.Email.Subject[content.Language["en"]]

	languageContent.Data = fileLanguageContent.Data
	languageContent.Language = make(map[string]int)
	for ID := 0; ID < len(fileLanguageContent.Language); ID++ {
		languageContent.Language[fileLanguageContent.Language[ID]] = ID
	}

	return languageContent, nil
}
