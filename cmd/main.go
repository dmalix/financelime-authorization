/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package main

import (
	"github.com/dmalix/financelime-authorization/server"
	"log"
	"math/rand"
	"time"
)

func init() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Llongfile)
	rand.Seed(time.Now().UTC().UnixNano())
}

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
