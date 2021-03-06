package id

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/noteshare/config"
	"github.com/noteshare/log"
	modelfolder "github.com/noteshare/model/folder"
	"github.com/noteshare/session"
)

//
// Delete handles deletion of single folder.
//
var Delete = session.Authenticate(
	func(
		w http.ResponseWriter,
		r *http.Request,
		s session.Session,
	) {

		if config.BuildDebug == true {
			fmt.Println(`==> DELETE: ` + r.URL.Path)
		}

		folderID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			log.NotifyError(err, http.StatusBadRequest)
			log.RespondJSON(w, `{}`, http.StatusBadRequest)
			return
		}

		err = modelfolder.DeleteFolder(
			folderID,
			s.UserID,
			s.AccountID,
		)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		log.RespondJSON(w, `{}`, http.StatusOK)

	},
)
