/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"io/ioutil"
	"os"
)

func initLanguageContent() (models.LanguageContent, error) {

	var (
		file                string
		fileLanguageContent models.FileLanguageContent
		languageContent     models.LanguageContent
		body                []byte
		err                 error
	)

	// Load the content of the file and convert it to structure

	file = os.Getenv("LANGUAGE_CONTENT_FILE")

	body, err = ioutil.ReadFile(file)
	if err != nil {
		return languageContent,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"bHY3kazX",
				"Failed to read the content.json file",
				err))
	}

	err = json.Unmarshal(body, &fileLanguageContent)
	if err != nil {
		return languageContent,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"zY3XbHka",
				"Failed to convert the content.json file",
				err))
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
