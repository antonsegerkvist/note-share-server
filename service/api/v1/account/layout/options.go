package layout

import (
	"fmt"
	"net/http"

	"github.com/noteshare/config"
	"github.com/noteshare/log"
)

//
// Options simply tells the client which methods, headers and origins
// are allowed.
//
func Options(w http.ResponseWriter, r *http.Request) {

	if config.BuildDebug == true {
		fmt.Println(`==> OPTIONS: /service/api/v1/account/layout`)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Origin")
	log.RespondJSON(w, ``, http.StatusOK)
}
