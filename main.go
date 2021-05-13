/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package main

import (
	"context"
	authorizationapp "github.com/dmalix/financelime-authorization/app"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"log"
	"math/rand"
	"time"
)

func init() {
	//log.SetFlags(log.Llongfile)
	//log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Lmsgprefix)
	rand.Seed(time.Now().UTC().UnixNano())
}

// @title Financelime Authorization
// @version v0.2.0-alpha
// @description Financelime Authorization RESTful API service
// @contact.name API Support
// @contact.email dmalix@financelime.com
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @host api.auth.dev.financelime.com
// @securityDefinitions.apikey authorization
// @in header
// @name authorization
// @schemes https
// @BasePath /

func main() {

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := authorizationapp.NewApp()
	if err != nil {
		log.Fatalf("%s: %s %s [%s]", "FATAL", trace.GetCurrentPoint(), "Failed to get a new App", err)
	}
	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("%s: %s %s [%s]", "FATAL", trace.GetCurrentPoint(), "Failed to run the App", err)
	}
}
