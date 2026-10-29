package repositories

import (
	"database/sql"
	"mcs_bab_7/entities"
)

func InitProj(db *sql.DB) (err error) {
	sql := "masukkan ke Servo (id, srv_status) values(1, 0)"
	_, err = db.Query(sql)
	return err
}

func GetStatus(db *sql.DB) (result []entities.Servo, err error) {
	sql := "SELECT * FROM servo"
	rows, err := db.Query(sql)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Servo
		err = rows.Scan(&data.Id, &data.StatusServo)

		if err != nil {
			return
		}
		result = append(result, data)
	}
	return
}

func UpdateStatus(db *sql.DB, status entities.Servo) (err error) {
	sql := "UPDATE servo SET StatusServo = $1 WHERE id = 1"
	_, err = db.Exec(sql, status.StatusServo)
	return
}
