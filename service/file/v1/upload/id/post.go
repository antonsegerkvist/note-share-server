package id

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/noteshare/config"
	"github.com/noteshare/log"
	modelfile "github.com/noteshare/model/file"
	"github.com/noteshare/session"
)

//
// ResponseData contains the fields of a response.
//
type ResponseData struct {
	Filesize uint64 `json:"filesize"`
	Checksum uint32 `json:"checksum"`
}

//
// Post handles file uploads.
//
var Post = session.Authenticate(
	func(
		w http.ResponseWriter,
		r *http.Request,
		s session.Session,
	) {

		if config.BuildDebug == true {
			fmt.Println(`==> POST: ` + r.URL.Path)
		}

		fileID, err := strconv.ParseUint(mux.Vars(r)["fid"], 10, 64)
		if err != nil {
			log.NotifyError(err, http.StatusBadRequest)
			log.RespondJSON(w, `{}`, http.StatusBadRequest)
			return
		}

		uploadFile, err := modelfile.LookupUploadFile(fileID, s.UserID, s.AccountID)
		if err == modelfile.ErrFileNotFound {
			log.NotifyError(err, http.StatusNotFound)
			log.RespondJSON(w, `{}`, http.StatusNotFound)
			return
		} else if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		fds, err := os.Create(
			path.Join(
				config.FileRootDir,
				config.FileOriginalDir,
				strconv.FormatUint(fileID, 10),
			),
		)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}
		defer fds.Close()

		filesizeOK, err := io.Copy(fds, r.Body)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		_, err = fds.Seek(0, 0)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		adler32Hash := adler32.New()
		_, err = io.Copy(adler32Hash, fds)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}
		checksumOK := adler32Hash.Sum32()

		if uploadFile.Filesize != 0 && uploadFile.Filesize != uint64(filesizeOK) {
			log.NotifyError(
				errors.New(`Filesize missmatch`),
				http.StatusBadRequest,
			)
			log.RespondJSON(w, `{}`, http.StatusBadRequest)
			return
		}

		if uploadFile.Checksum != 0 && uploadFile.Checksum != uint32(checksumOK) {
			log.NotifyError(
				errors.New(`Checksum missmatch`),
				http.StatusBadRequest,
			)
			log.RespondJSON(w, `{}`, http.StatusBadRequest)
			return
		}

		err = modelfile.MarkFileAsUploaded(
			fileID,
			uploadFile.FolderID,
			uploadFile.Name,
			uploadFile.Filename,
			uint64(filesizeOK),
			uint32(checksumOK),
			s.UserID,
			s.AccountID,
		)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		responseData := ResponseData{
			Filesize: uint64(filesizeOK),
			Checksum: checksumOK,
		}
		jsonBytes, err := json.Marshal(responseData)
		if err != nil {
			log.NotifyError(err, http.StatusInternalServerError)
			log.RespondJSON(w, `{}`, http.StatusInternalServerError)
			return
		}

		log.RespondJSON(w, string(jsonBytes), http.StatusOK)

	},
)
