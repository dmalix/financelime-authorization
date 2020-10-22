/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package main

import (
	"github.com/dmalix/financelime-rest-api/server"
	"log"
)

func main() {

	var err error

	app, err := server.NewApp()
	if err != nil {
		log.Fatalf("FATAL [%s: Failed to get a new App [%s]]", "hiaX3GMR", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("FATAL [%s: Failed to run the App [%s]]", "kxXG0YqR", err)
	}
}
