package me

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/noteshare/config"
	"github.com/noteshare/log"
	modelaccount "github.com/noteshare/model/account"
	"github.com/noteshare/session"
)

//
// Get returns information about the account the user belongs to.
//
var Get = session.Authenticate(
	func(
		w http.ResponseWriter,
		r *http.Request,
		s session.Session,
	) {

		if config.BuildDebug == true {
			fmt.Println(`==> GET: ` + r.URL.Path)
		}

		response, err := modelaccount.GetAccount(s.AccountID)
		if err == modelaccount.ErrAccountNotFound {
			log.NotifyError(err, http.StatusNotFound)
			log.RespondJSON(w, `{}`, http.StatusNotFound)
			return
		} else if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(response)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		log.RespondJSON(w, string(jsonBytes), http.StatusOK)

	},
)
