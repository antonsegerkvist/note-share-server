package folder

import (
	"database/sql"

	"github.com/noteshare/mysql"
)

//
// GetFoldersFromFolderID returns a list of folder with the specified folder ID
// as a parent.
//
func GetFoldersFromFolderID(folderID, userID uint64) (*[]ModelFolder, error) {

	const query = `
		select c_id, c_parent, c_name, c_created_by_user_id, c_modified_by_user_id, c_created_date, c_modified_date
		from t_folder
		inner join t_user on t_folder.c_account_id = t_user.c_account_id
		where t_folder.c_parent = ? and t_user.c_id = ?
	`

	db := mysql.Open()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(folderID, userID)
	if err == sql.ErrNoRows {
		return &[]ModelFolder{}, nil
	} else if err != nil {
		return nil, err
	}

	ret := []ModelFolder{}
	for rows.Next() {
		buffer := ModelFolder{}
		err = rows.Scan(
			&buffer.ID,
			&buffer.Parent,
			&buffer.Name,
			&buffer.CreatedByUserID,
			&buffer.ModifiedByUserID,
			&buffer.CreatedDate,
			&buffer.ModifiedDate,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, buffer)
	}

	return &ret, nil
}
