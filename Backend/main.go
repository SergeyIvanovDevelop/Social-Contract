package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// Взаимодействие с БД
// ===========================================================================================
// ===========================================================================================

// Регистрация (БД)
// ===========================================================================================

func SQL_Work_for_ProcessPostsRegistration(registration_struct *Registration_struct) string {

	// Запрос всех логинов
	query := `SELECT "user"."login" FROM "view_auth_data" AS "user"`
	fmt.Println("# Достать все логины пользователей приложения")
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()
	var logins []string
	for result.Next() {
		var login string
		err2 := result.Scan(&login)
		if err2 != nil {
			panic(err2)
		}
		logins = append(logins, login)
	}
	fmt.Println("logins = ", logins)
	var login_busy bool = false
	for _, login := range logins {
		if (*registration_struct).Login_parent == login {
			login_busy = true
			break
		}
	}
	if login_busy == true {
		fmt.Println("LOGIN_BUSY")
		return "LOGIN_BUSY"
	} else if login_busy == false {
		// Сохраняем в файл string-строку
		var buf bytes.Buffer
		b, _ := base64.StdEncoding.DecodeString((*registration_struct).Photo_parent)
		io.Copy(&buf, bytes.NewBuffer(b))
		var path_to_buf_file string = "./Photo_parent.jpg"
		f1, err := os.Create(path_to_buf_file)
		if err != nil {
			panic(err.Error())
		}
		_, err = io.Copy(f1, &buf)
		if err != nil {
			panic(err.Error())
		}
		f1.Close()

		// Считаем хэш от пароля
		sha256 := sha256.Sum256([]byte((*registration_struct).Password_parent))
		hash_str := base64.StdEncoding.EncodeToString(sha256[:])
		// Отправляем в БД все данные для регистрации
		query = `INSERT INTO "users" (
			"login",
			"hash",
			"surname",
			"firstname",
			"photo"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		);`
		fmt.Println("# Отправляем в БД все данные для регистрации")
		result_, err := db.Exec(query,
			(*registration_struct).Login_parent,
			hash_str,
			(*registration_struct).Secondname_parent,
			(*registration_struct).Firstname_parent,
			(*registration_struct).Photo_parent) //"bytea("+path_to_buf_file+")")
		if err != nil {
			panic(err)
		}
		//defer result_.Close()
		//defer os.Remove(path_to_buf_file)
		fmt.Println("# result_ = ", result_)
		return "OK"
	}
	return "Impossible situation"
}

// ===========================================================================================

func Auth(enter_struct *Enter_struct) string {
	// Запрашиваю из БД хэш пароля для данного пользователя
	query := `SELECT "user"."hash" FROM "view_auth_data" AS "user" WHERE "user"."login" = $1;`
	fmt.Println("# Достать все логины пользователей приложения")
	rows, err := db.Query(query, (*enter_struct).Login_parent)
	//result, err := db.Exec(query, (*enter_struct).Login_parent)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var hash_sum string
		err2 := rows.Scan(&hash_sum)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, hash_sum)
	}

	var hash string
	if len(result) >= 1 {
		hash = result[0]
	}

	fmt.Println("hash = ", hash)
	// Считаем хэш от пароля
	sha256 := sha256.Sum256([]byte((*enter_struct).Password_parent))
	hash_str := base64.StdEncoding.EncodeToString(sha256[:])

	if hash == hash_str {
		return "OK"
	} else {
		return "AUTH_FALL"
	}
}

// ===========================================================================================
func SQL_Work_for_ProcessPostsEnter(enter_struct *Enter_struct) (EnterAnswer_struct, string) {

	var auth_result string = Auth(enter_struct)
	if auth_result == "AUTH_FALL" {
		return EnterAnswer_struct{}, "AUTH_FALL"
	} else {

		var enterAnswer_struct EnterAnswer_struct = EnterAnswer_struct{}

		// Запрашиваем из БД все данные пользователя
		query := `SELECT "users"."photo" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД фото пользователя")
		rows, err := db.Query(query, (*enter_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var result []string
		for rows.Next() {
			var photo string
			err2 := rows.Scan(&photo)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, photo)
		}
		var photo_parent string
		if len(result) != 0 {
			photo_parent = result[0]
		}
		enterAnswer_struct.Photo_parent = photo_parent
		fmt.Println("photo_parent[:5] = ", photo_parent[:5])
		query = `SELECT "users"."firstname", "users"."surname" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД имя и фамилию пользователя")
		rows_, err := db.Query(query, (*enter_struct).Login_parent)
		//result, err := db.Exec(query, (*enter_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		defer rows_.Close()
		var firstname string
		var surname string
		for rows_.Next() {

			err2 := rows_.Scan(&firstname, &surname)
			if err2 != nil {
				panic(err2)
			}
		}

		enterAnswer_struct.Firstname_parent = firstname
		enterAnswer_struct.Secondname_parent = surname

		fmt.Println("firstname = ", firstname)
		fmt.Println("surname = ", surname)

		query = `SELECT
		"rec"."payer_surname",
		"rec"."payer_firstname",
		"rec"."payer_photo",
		"rec"."recipient_surname",
		"rec"."recipient_firstname",
		"rec"."recipient_photo",
		"rec"."Name",
		"rec"."PersonalAcc",
		"rec"."BankName",
		"rec"."BIC",
		"rec"."CorrespAcc",
		"rec"."KPP",
		"rec"."PayeeINN"
	  FROM "view_recipients_for_payers" AS "rec"
	  WHERE "rec"."payer_login" = $1;`

		fmt.Println("# Запрашиваем из БД всех детей для конкретного родителя")
		rows_2, err := db.Query(query, (*enter_struct).Login_parent)
		if err != nil {
			panic(err)
		}

		var payer_surname string
		var payer_firstname string
		var payer_photo string
		var recipient_surname string
		var recipient_firstname string
		var recipient_photo string
		var Name_ string
		var PersonalAcc_ string
		var BankName_ string
		var BIC_ string
		var CorrespAcc_ string
		var KPP_ string
		var PayeeINN_ string
		for rows_2.Next() {

			err2 := rows_2.Scan(&payer_surname, &payer_firstname, &payer_photo,
				&recipient_surname, &recipient_firstname, &recipient_photo, &Name_,
				&PersonalAcc_, &BankName_, &BIC_, &CorrespAcc_, &KPP_, &PayeeINN_)
			if err2 != nil {
				panic(err2)
			}

			cci := Children_Card_info{Name: Name_, PersonalAcc: PersonalAcc_,
				BankName: BankName_, BIC: BIC_, CorrespAcc: CorrespAcc_, KPP: KPP_, PayeeINN: PayeeINN_}

			// For testing
			fmt.Println("BIC = ", BIC_)
			fmt.Println("payer_surname = ", payer_surname)

			enterAnswer_struct.Childrens = append(enterAnswer_struct.Childrens,
				Children{Firstname_children: recipient_firstname, Secondname_children: recipient_surname,
					Photo_children: recipient_photo, Children_CARD_info: cci})

		}
		defer rows_2.Close()

		query = `SELECT
		"con"."contract_name",
		"con"."recipient_surname",
		"con"."recipient_firstname",
		"con"."recipient_photo",
		"con"."g1",
		"con"."g2",
		"con"."g3",
		"con"."g4",
		"con"."g5"
	  FROM "view_contracts_with_price" AS "con"
	  WHERE "con"."payer_login" = $1;`

		fmt.Println("# Запрашиваем из БД все контракты конкретного родителя")
		rows_3, err := db.Query(query, (*enter_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("rows_3: %v\n", rows_3)
		defer rows_3.Close()
		var contract_name string
		var recipient_surname_ string
		var recipient_firstname_ string
		var recipient_photo_ string
		var g1 int
		var g2 int
		var g3 int
		var g4 int
		var g5 int

		for rows_3.Next() {
			err2 := rows_3.Scan(&contract_name, &recipient_surname_, &recipient_firstname_, &recipient_photo_, &g1, &g2, &g3, &g4, &g5)
			if err2 != nil {
				panic(err2)
			}
			fmt.Printf("rows_3: %v\n", rows_3)
			// Тут запрашиваем пустое ли поле Contracts
			var contract_completed bool = true
			if contract_name != "" {
				//Узнать user_id человека с определенным логином
				query := `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
				fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
				rows, err := db.Query(query, (*enter_struct).Login_parent)
				if err != nil {
					panic(err)
				}
				var user_id_parent int
				for rows.Next() {
					err2 := rows.Scan(&user_id_parent)
					if err2 != nil {
						panic(err2)
					}
				}
				rows.Close()
				fmt.Println("user_id_parent = ", user_id_parent)

				//Узнать user_id человека с определенным логином
				query = `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
				fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
				rows2, err := db.Query(query, "LOGIN_"+recipient_firstname_+"_"+recipient_surname_+"_"+(*enter_struct).Login_parent)
				if err != nil {
					panic(err)
				}
				var user_id_children int
				for rows2.Next() {

					err2 := rows2.Scan(&user_id_children)
					if err2 != nil {
						panic(err2)
					}
				}
				rows2.Close()
				fmt.Println("user_id_children = ", user_id_children)

				// Нужно реализовать в БД
				query = `SELECT "contracts"."completed" FROM "contracts" WHERE "payer" = $1 AND "recipient" = $2;`
				fmt.Println("# Запрашиваем поле активности определенного контракта в БД")
				rows3, err := db.Query(query, user_id_parent, user_id_children)
				if err != nil {
					panic(err)
				}

				for rows3.Next() {

					err2 := rows3.Scan(&contract_completed)
					if err2 != nil {
						panic(err2)
					}
				}
				rows3.Close()
			}

			if contract_completed == false {
				enterAnswer_struct.Contracts = append(enterAnswer_struct.Contracts,
					Contract{Firstname_children: recipient_firstname_, Secondname_children: recipient_surname_,
						Photo_children: recipient_photo_, Contract_name: contract_name, Mark_policy: Mark_Policy{One: strconv.Itoa(g1), Two: strconv.Itoa(g2), Three: strconv.Itoa(g3),
							Four: strconv.Itoa(g4), Five: strconv.Itoa(g5)}})
			}

		}
		defer rows_3.Close()

		return enterAnswer_struct, "OK"
	}
}

// Добавление нового пользователя
// ===========================================================================================

func SQL_Work_for_ProcessPostsAddChildren(add_children_struct *Add_children_struct) string {

	var enter_struct Enter_struct
	enter_struct.Login_parent = (*add_children_struct).Login_parent
	enter_struct.Password_parent = (*add_children_struct).Password_parent

	// Сохраняем в файл string-строку
	var buf bytes.Buffer
	b, _ := base64.StdEncoding.DecodeString((*add_children_struct).Photo_children)
	io.Copy(&buf, bytes.NewBuffer(b))
	var path_to_buf_file string = "./Photo_children.jpg"
	f1, err := os.Create(path_to_buf_file)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f1, &buf)
	if err != nil {
		panic(err.Error())
	}
	f1.Close()
	//defer os.Remove(path_to_buf_file)
	var auth_result string = Auth(&enter_struct)
	if auth_result == "AUTH_FALL" {
		return "AUTH_FALL"
	} else {
		// Отправляем в БД все данные для добавления ребенка
		query := `INSERT INTO "users" (
	"login",
	"hash",
	"surname",
	"firstname",
	"photo"
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
) returning user_id;`
		var lastInsertId int
		fmt.Println("# Отправляем в БД все данные для регистрации")
		err := db.QueryRow(query,
			"LOGIN_"+(*add_children_struct).Firstname_children+"_"+(*add_children_struct).Secondname_children+"_"+add_children_struct.Login_parent,
			"hash_str", // Абсолютно неважно это поле
			(*add_children_struct).Secondname_children,
			(*add_children_struct).Firstname_children,
			(*add_children_struct).Photo_children).Scan(&lastInsertId) //"bytea("+path_to_buf_file+")").Scan(&lastInsertId)
		if err != nil {
			panic(err)
		}

		query = `INSERT INTO "cards" (
			"user_id",
			"Name",
			"PersonalAcc",
			"BankName",
			"BIC",
			"CorrespAcc",
			"KPP",
			"PayeeINN"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		);`
		res, err := db.Exec(query,
			lastInsertId,
			(*add_children_struct).Children_CARD_info.Name,
			(*add_children_struct).Children_CARD_info.PersonalAcc,
			(*add_children_struct).Children_CARD_info.BankName,
			(*add_children_struct).Children_CARD_info.BIC,
			(*add_children_struct).Children_CARD_info.CorrespAcc,
			(*add_children_struct).Children_CARD_info.KPP,
			(*add_children_struct).Children_CARD_info.PayeeINN)
		if err != nil {
			panic(err)
		}
		fmt.Println("res = ", res)
		return "OK"
	}

}

// ===========================================================================================

// Заключение контракта
// ===========================================================================================

func SQL_Work_for_ProcessPostsSignContract(sign_contract_struct *Sign_contract_struct) string {

	// Сохраняем картинки-условия -----------------------------------------------
	buf := bytes.NewBuffer([]byte((*sign_contract_struct).Photo_condition_parent))
	var Photo_condition_parent_PATH string = "Photo_condition_parent.jpg"
	f1, err := os.Create(Photo_condition_parent_PATH)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f1, buf)
	if err != nil {
		panic(err.Error())
	}
	f1.Close()

	buf = bytes.NewBuffer([]byte((*sign_contract_struct).Photo_condition_children))
	var Photo_condition_children_PATH string = "Photo_condition_children.jpg"
	f1, err = os.Create(Photo_condition_children_PATH)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f1, buf)
	if err != nil {
		panic(err.Error())
	}
	f1.Close()
	// Сохраняем картинки-условия -----------------------------------------------

	var enter_struct Enter_struct
	enter_struct.Login_parent = sign_contract_struct.Login_parent
	enter_struct.Password_parent = sign_contract_struct.Password_parent
	var auth_result string = Auth(&enter_struct)
	if auth_result == "AUTH_FALL" {
		return "AUTH_FALL"
	} else {

		//Узнать user_id человека с определенным логином
		query := `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
		rows, err := db.Query(query, (*sign_contract_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		var user_id_parent int
		for rows.Next() {

			err2 := rows.Scan(&user_id_parent)
			if err2 != nil {
				panic(err2)
			}
		}
		rows.Close()
		fmt.Println("user_id_parent = ", user_id_parent)

		//Узнать user_id человека с определенным логином
		query = `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
		rows2, err := db.Query(query, "LOGIN_"+(*sign_contract_struct).Contract_.Firstname_children+"_"+(*sign_contract_struct).Contract_.Secondname_children+"_"+(*sign_contract_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		var user_id_children int
		for rows2.Next() {

			err2 := rows2.Scan(&user_id_children)
			if err2 != nil {
				panic(err2)
			}
		}
		rows2.Close()
		fmt.Println("user_id_children = ", user_id_children)
		var contract_id int
		query = `INSERT INTO "contracts" (
			"payer",
			"recipient",
			"contract_name",
			"checker_time"
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) returning contract_id;`

		err = db.QueryRow(query,
			user_id_parent,
			user_id_children,
			(*sign_contract_struct).Contract_.Contract_name,
			(*sign_contract_struct).Timer_interval_notisfaction).Scan(&contract_id)
		if err != nil {
			panic(err)
		}
		fmt.Println("contract_id = ", contract_id)

		for i := 1; i < 6; i++ {
			query = `INSERT INTO "prices" (
			"contract_id",
			"grade",
			"price"
		) VALUES (
			$1,
			$2,
			$3
		);`
			var price int
			if i == 1 {
				price, _ = strconv.Atoi((*sign_contract_struct).Contract_.Mark_policy.One)
			} else if i == 2 {
				price, _ = strconv.Atoi((*sign_contract_struct).Contract_.Mark_policy.Two)
			} else if i == 3 {
				price, _ = strconv.Atoi((*sign_contract_struct).Contract_.Mark_policy.Three)
			} else if i == 4 {
				price, _ = strconv.Atoi((*sign_contract_struct).Contract_.Mark_policy.Four)
			} else if i == 5 {
				price, _ = strconv.Atoi((*sign_contract_struct).Contract_.Mark_policy.Five)
			}
			res, err := db.Exec(query,
				contract_id,
				i,
				price)
			if err != nil {
				panic(err)
			}
			fmt.Println("res = ", res)
		}

		query = `INSERT INTO "diary_credentials" (
			"user_id",
			"diary_login",
			"diary_password"
		) VALUES (
			$1,
			$2,
			$3
		);`
		res1, err := db.Exec(query,
			user_id_children,
			(*sign_contract_struct).Diary_children_login,
			(*sign_contract_struct).Diary_children_password)
		if err != nil {
			panic(err)
		}
		fmt.Println("res1 = ", res1)

		return "OK"
	}
}

// ===========================================================================================

// Завершение контракта
// ===========================================================================================
func SQL_Work_for_ProcessPostsFinishContract(finish_contract_struct *Finish_contract_struct) string {
	// Сохраняем картинки-условия -----------------------------------------------
	buf := bytes.NewBuffer([]byte((*finish_contract_struct).Photo_condition_parent))
	var Photo_condition_parent_PATH string = "Photo_condition_parent.jpg"
	f1, err := os.Create(Photo_condition_parent_PATH)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f1, buf)
	if err != nil {
		panic(err.Error())
	}
	f1.Close()

	buf = bytes.NewBuffer([]byte((*finish_contract_struct).Photo_condition_children))
	var Photo_condition_children_PATH string = "Photo_condition_children.jpg"
	f1, err = os.Create(Photo_condition_children_PATH)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f1, buf)
	if err != nil {
		panic(err.Error())
	}
	f1.Close()
	// Сохраняем картинки-условия -----------------------------------------------

	var enter_struct Enter_struct
	enter_struct.Login_parent = (*finish_contract_struct).Login_parent
	enter_struct.Password_parent = (*finish_contract_struct).Password_parent
	var auth_result string = Auth(&enter_struct)
	if auth_result == "AUTH_FALL" {
		return "AUTH_FALL"
	} else {

		//Узнать user_id человека с определенным логином
		query := `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
		rows, err := db.Query(query, (*finish_contract_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		var user_id_parent int
		for rows.Next() {

			err2 := rows.Scan(&user_id_parent)
			if err2 != nil {
				panic(err2)
			}
		}
		rows.Close()
		fmt.Println("user_id_parent = ", user_id_parent)

		//Узнать user_id человека с определенным логином
		query = `SELECT "users"."user_id" FROM "users" WHERE "users"."login" = $1;`
		fmt.Println("# Запрашиваем из БД user_id конкретного родителя с определенным логином")
		rows2, err := db.Query(query, "LOGIN_"+(*finish_contract_struct).Contract_.Firstname_children+"_"+(*finish_contract_struct).Contract_.Secondname_children+"_"+(*finish_contract_struct).Login_parent)
		if err != nil {
			panic(err)
		}
		var user_id_children int
		for rows2.Next() {

			err2 := rows2.Scan(&user_id_children)
			if err2 != nil {
				panic(err2)
			}
		}
		rows2.Close()
		fmt.Println("user_id_children = ", user_id_children)

		// Нужно реализовать в БД
		query = `UPDATE "contracts" SET "completed" = $1 WHERE "payer" = $2 AND "recipient" = $3;`
		fmt.Println("# Обновляем поле активности определенного контракта в БД")
		rows3, err := db.Query(query, true, user_id_parent, user_id_children)
		if err != nil {
			panic(err)
		}
		rows3.Close()

		return "OK"
	}

}

// ===========================================================================================

// ===========================================================================================
// ===========================================================================================

// Registration
// ======================================================

type Registration_struct struct {
	Photo_parent      string `json:"Photo_parent"`
	Firstname_parent  string `json:"Firstname_parent"`
	Secondname_parent string `json:"Secondname_parent"`
	Login_parent      string `json:"Login_parent"`
	Password_parent   string `json:"Password_parent"`
}

type RegistrationAnswer_struct struct {
	Status string `json:"Status"`
}

// ======================================================

// Enter
// ======================================================

type Enter_struct struct {
	Login_parent    string `json:"Login_parent"`
	Password_parent string `json:"Password_parent"`
}

type Children_Card_info struct {
	Name        string `json:"Name"`
	PersonalAcc string `json:"PersonalAcc"`
	BankName    string `json:"BankName"`
	BIC         string `json:"BIC"`
	CorrespAcc  string `json:"CorrespAcc"`
	KPP         string `json:"KPP"`
	PayeeINN    string `json:"PayeeINN"`
}

type Mark_Policy struct {
	One   string `json:"One"`
	Two   string `json:"Two"`
	Three string `json:"Three"`
	Four  string `json:"Four"`
	Five  string `json:"Five"`
}

type Contract struct {
	Firstname_children  string      `json:"Firstname_children"`
	Secondname_children string      `json:"Secondname_children"`
	Photo_children      string      `json:"Photo_children"`
	Contract_name       string      `json:"Contract_name"`
	Mark_policy         Mark_Policy `json:"Mark_policy"`
}

type Children struct {
	Firstname_children  string             `json:"Firstname_children"`
	Secondname_children string             `json:"Secondname_children"`
	Photo_children      string             `json:"Photo_children"`
	Children_CARD_info  Children_Card_info `json:"Children_CARD_info"`
}

type EnterAnswer_struct struct {
	Status            string     `json:"Status"`
	Photo_parent      string     `json:"Photo_parent"`
	Firstname_parent  string     `json:"Firstname_parent"`
	Secondname_parent string     `json:"Secondname_parent"`
	Childrens         []Children `json:"Childrens"`
	Contracts         []Contract `json:"Contracts"`
}

// ======================================================

// Add_children
// ======================================================

type Add_children_struct struct {
	Login_parent            string             `json:"Login_parent"`
	Password_parent         string             `json:"Password_parent"`
	Firstname_children      string             `json:"Firstname_children"`
	Secondname_children     string             `json:"Secondname_children"`
	Photo_children          string             `json:"Photo_children"`
	Children_CARD_info      Children_Card_info `json:"Children_CARD_info"`
	Diary_children_login    string             `json:"Diary_children_login"`    // Заглушка
	Diary_children_password string             `json:"Diary_children_password"` // Заглушка
}

type Add_childrenAnswer_struct struct {
	Status string `json:"Status"`
}

// ======================================================

// Sign_contract
// ======================================================

type Sign_contract_struct struct {
	Login_parent                string   `json:"Login_parent"`
	Password_parent             string   `json:"Password_parent"`
	Contract_                   Contract `json:"Contract_"`
	Diary_children_login        string   `json:"Diary_children_login"`
	Diary_children_password     string   `json:"Diary_children_password"`
	Photo_condition_parent      string   `json:"Photo_condition_parent"`
	Photo_condition_children    string   `json:"Photo_condition_children"`
	Timer_interval_notisfaction string   `json:"Timer_interval_notisfaction"`
}

type Sign_contract_answer_struct struct {
	Status string `json:"Status"`
}

// ======================================================

// Finish_contract
// ======================================================

type Contract__ struct {
	Firstname_children  string `json:"Firstname_children"`
	Secondname_children string `json:"Secondname_children"`
	Contract_name       string `json:"Contract_name"`
}

type Finish_contract_struct struct {
	Login_parent             string     `json:"Login_parent"`
	Password_parent          string     `json:"Password_parent"`
	Contract_                Contract__ `json:"Contract_"`
	Photo_condition_parent   string     `json:"Photo_condition_parent"`
	Photo_condition_children string     `json:"Photo_condition_children"`
}

type Finish_contract_answer_struct struct {
	Status string `json:"Status"`
}

// ======================================================

func ProcessPostsRegistration(w http.ResponseWriter, r *http.Request) {

	// CORS ======================================================================
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// CORS ======================================================================

	var registration_struct Registration_struct = Registration_struct{}
	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &registration_struct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("registration_struct = %v\n", registration_struct)
	var registration_answer_struct RegistrationAnswer_struct = RegistrationAnswer_struct{}
	status := SQL_Work_for_ProcessPostsRegistration(&registration_struct)
	fmt.Println("status = ", status)
	registration_answer_struct.Status = status
	fmt.Println("registration_answer_struct.Status = ", registration_answer_struct.Status)
	json.NewEncoder(w).Encode(registration_answer_struct)
}

func ProcessPostsEnter(w http.ResponseWriter, r *http.Request) {

	// CORS ======================================================================
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// CORS ======================================================================

	var enter_struct Enter_struct
	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &enter_struct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("enter_struct = %v\n", enter_struct)
	var enter_answer_struct EnterAnswer_struct
	var Status string
	enter_answer_struct, Status = SQL_Work_for_ProcessPostsEnter(&enter_struct)
	// ... Проверка
	enter_answer_struct.Status = Status
	json.NewEncoder(w).Encode(enter_answer_struct)
}

func ProcessPostsAddChildren(w http.ResponseWriter, r *http.Request) {

	// CORS ======================================================================
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// CORS ======================================================================

	var add_children_struct Add_children_struct
	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &add_children_struct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("add_children_struct = %v\n", add_children_struct)
	var add_children_answer_struct Add_childrenAnswer_struct = Add_childrenAnswer_struct{}
	var Status string = SQL_Work_for_ProcessPostsAddChildren(&add_children_struct)
	add_children_answer_struct.Status = Status
	json.NewEncoder(w).Encode(add_children_answer_struct)
}

func ProcessPostsSignContract(w http.ResponseWriter, r *http.Request) {

	// CORS ======================================================================
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// CORS ======================================================================

	var sign_contract_struct Sign_contract_struct
	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &sign_contract_struct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sign_contract_struct = %v\n", sign_contract_struct)
	var sign_contract_answer_struct Sign_contract_answer_struct = Sign_contract_answer_struct{}
	var Status string = SQL_Work_for_ProcessPostsSignContract(&sign_contract_struct)
	sign_contract_answer_struct.Status = Status
	json.NewEncoder(w).Encode(sign_contract_answer_struct)
}

func ProcessPostsFinishContract(w http.ResponseWriter, r *http.Request) {

	// CORS ======================================================================
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// CORS ======================================================================

	var finish_contract_struct Finish_contract_struct
	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &finish_contract_struct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("finish_contract_struct = %v\n", finish_contract_struct)
	var finish_contract_answer_struct Finish_contract_answer_struct = Finish_contract_answer_struct{}
	var Status string = SQL_Work_for_ProcessPostsFinishContract(&finish_contract_struct)
	finish_contract_answer_struct.Status = Status
	json.NewEncoder(w).Encode(finish_contract_answer_struct)
}

/* func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Write([]byte("Hello, World!"))
} */

func RunREST_API(port string) {
	sm := http.NewServeMux()
	sm.HandleFunc("/registration/", ProcessPostsRegistration)
	sm.HandleFunc("/enter/", ProcessPostsEnter)
	sm.HandleFunc("/add_children/", ProcessPostsAddChildren)
	sm.HandleFunc("/sign_contract/", ProcessPostsSignContract)
	sm.HandleFunc("/finish_contract/", ProcessPostsFinishContract)
	l, err := net.Listen("tcp4", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", "Запущен http-сервер бекэнда для взаимодействия фронтенда и БД (Интерфейс - REST API)")
	time.Sleep(1 * time.Second)
	log.Fatal(http.Serve(l, sm))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	DB_USER     = "<db_user_name>"
	DB_PASSWORD = "<db_user_password>"
	DB_NAME     = "<db_name>"
)

var db *sql.DB

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage ./main <Port>")
		fmt.Println("For example: ./main  2021")
		os.Exit(1)
	}

	port := ":" + os.Args[1]
	go RunREST_API(port)
	time.Sleep(3 * time.Second)

	// Открытие и подключение к БД
	// ===========================================================================
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	DB, err := sql.Open("postgres", dbinfo)
	db = DB
	checkErr(err)
	fmt.Println("БД успешно открыта")
	defer db.Close()

	for {
	}

}
