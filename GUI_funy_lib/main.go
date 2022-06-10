package main

import (
	"GUI_funy_lib/tutorials"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var port string = ""

var IP_address string = ""

var URL_ENTER string
var URL_REGISTRACTION string
var URL_ADD_CHILDREN string
var URL_SIGN_CONTRACT string
var URL_FINISH_CONTRACT string

// Значения по умолчанию - dark_theme
var (
	LOG_IN_ICON     string = "./icons/dark_theme/log_in_icon.png"
	REG_ICON        string = "./icons/dark_theme/reg_icon.png"
	ADD_PHOTO_ICON  string = "./icons/dark_theme/add_photo_icon.png"
	FINISH_CONTRACT string = "./icons/dark_theme/finish_contract_icon.png"
	ADD_CONTRACT    string = "./icons/dark_theme/add_contract_icon.png"
	ADD_CHILDREN    string = "./icons/dark_theme/add_children_icon.png"
	APP_ICON        string = "./icons/dark_theme/app_icon.png"
)

var Add_photo_ptr *widget.Button
var Registration_button_ptr *widget.Button
var Sign_contract_button_ptr *widget.Button
var Add_contract_button_ptr *widget.Button
var Add_children_button_ptr *widget.Button

type My_Image struct {
	name    string
	content []byte
}

func (m *My_Image) Content() []byte {
	return (*m).content
}

func (m *My_Image) Name() string {
	return (*m).name
}

var topWindow fyne.Window

const preferenceCurrentTutorial = "currentTutorial"

var applicat *fyne.App
var picture_registration *My_Image
var w_enter *fyne.Window
var w_reg *fyne.Window
var w_main *fyne.Window
var w_add_children *fyne.Window
var w_add_contract *fyne.Window
var Picture canvas.Image
var Add_contract_button *widget.Button
var CurrentChildrenName string
var CurrentChildrenSurname string
var ParentName string
var ParentSurname string
var CurrentName string
var GLOBAL_current_index_child *int

// Необходимые переменные для регистрации --------------------------
var Photo_Parent *canvas.Image
var Parent_name *widget.Entry
var Parent_surname *widget.Entry
var Parent_login *widget.Entry
var Parent_password *widget.Entry

// -----------------------------------------------------------------

// Необходимые переменные для входа --------------------------
var Parent_login_enter *widget.Entry
var Parent_password_enter *widget.Entry

// -----------------------------------------------------------------

// -----------------------------------------------------------------

// Необходимые переменные для добавления ребенка  --------------------------
var Children_name *widget.Entry
var Children_surname *widget.Entry
var Photo_Children *canvas.Image
var _Name *widget.Entry
var _PersonalAcc *widget.Entry
var _BankName *widget.Entry
var _BIC *widget.Entry
var _CorrespAcc *widget.Entry
var _KPP *widget.Entry
var _PayeeINN *widget.Entry

// -----------------------------------------------------------------

// Необходимые переменные для заключения контракта  --------------------------
var Children_name_contract *widget.Entry
var Children_surname_contract *widget.Entry
var Contract_Name *widget.Entry
var _Diary_children_login *widget.Entry
var _Diary_children_password *widget.Entry
var Timer_interval_notisfaction *numEntry
var text_obj1 *widget.Entry
var text_obj2 *widget.Entry
var text_obj3 *widget.Entry
var text_obj4 *widget.Entry
var text_obj5 *widget.Entry

// -----------------------------------------------------------------

// Необходимые переменные для завершения контракта  --------------------------
var Children_name_contract_finish *widget.Label
var Children_surname_contract_finish *widget.Label
var _Tree *widget.Tree
var Contract_Name_finish *string

// -----------------------------------------------------------------

// Глобальные объекты, которые заполняются по ходу выполнения программы, и которые ДОЛЖНЫ быть заполнены к моменту вызова функций, которые будут их отправлять
var registration_struct Registration_struct
var registrationAnswer_struct RegistrationAnswer_struct
var enter_struct Enter_struct
var enterAnswer_struct EnterAnswer_struct
var add_children_struct Add_children_struct
var add_childrenAnswer_struct Add_childrenAnswer_struct
var sign_contract_struct Sign_contract_struct
var sign_contract_answer_struct Sign_contract_answer_struct
var finish_contract_struct Finish_contract_struct
var finish_contract_answer_struct Finish_contract_answer_struct

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

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage ./main <IP_address of backand> <Port>")
		fmt.Println("For example: ./main 127.0.0.1 2021 for local connect")
		fmt.Println("For example: ./main 46.72.127.251 2021 for remote connect")
		os.Exit(1)
	}

	IP_address = os.Args[1]
	port = os.Args[2]
	fmt.Printf("IP_address of backand: |%s|\n", IP_address)

	URL_ENTER = "http://" + IP_address + ":" + port + "/enter/"
	URL_REGISTRACTION = "http://" + IP_address + ":" + port + "/registration/"
	URL_ADD_CHILDREN = "http://" + IP_address + ":" + port + "/add_children/"
	URL_SIGN_CONTRACT = "http://" + IP_address + ":" + port + "/sign_contract/"
	URL_FINISH_CONTRACT = "http://" + IP_address + ":" + port + "/finish_contract/"

	picture_registration = &My_Image{}

	image_data_before, err := ioutil.ReadFile(APP_ICON)
	if err != nil {
		panic(err.Error())
	}
	// Эмуляция того, что я вытащил строку из JSON
	var my_str_image string = base64.StdEncoding.EncodeToString(image_data_before)
	// Декодирование строки base64 в []byte
	image_data_after, err := base64.StdEncoding.DecodeString(my_str_image)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var my_image My_Image
	my_image = My_Image{name: "app_icon.jpg", content: image_data_after}
	a := app.NewWithID("Social_contract")
	a.SetIcon(&my_image)
	applicat = &a
	w_enter_ := (*applicat).NewWindow("ENTRANCE")
	w_reg_ := (*applicat).NewWindow("REGISTRATION")
	w_main_ := (*applicat).NewWindow("SOCIAL CONTRACT")
	w_add_children_ := (*applicat).NewWindow("ADDING CHILDREN")
	w_add_contract_ := (*applicat).NewWindow("SIGNING CONTRACT")
	Add_contract_button = &widget.Button{}
	w_enter = &w_enter_
	w_reg = &w_reg_
	w_main = &w_main_
	w_add_children = &w_add_children_
	w_add_contract = &w_add_contract_
	topWindow = *w_enter
	Enter_Window()
}

func SendPostRequest(URL string, json_struct interface{}) string {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	var req *http.Request
	switch json_struct.(type) {
	case *Enter_struct:
		body, err := json.Marshal(json_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("json string:\n\t%s\n", string(body))
		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error happend", err)
			panic(err)
		}
		defer resp.Body.Close() // важный пункт!
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		err = json.Unmarshal(buf.Bytes(), &enterAnswer_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("enter_struct = %v\n", enterAnswer_struct)
		return enterAnswer_struct.Status
	case *Registration_struct:
		body, err := json.Marshal(json_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("json string:\n\t%s\n", string(body))
		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		// NOTE this !!

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error happend", err)
			panic(err)
		}
		defer resp.Body.Close() // важный пункт!
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)

		/* json.NewDecoder(resp.Body).Decode(&registrationAnswer_struct) */
		err = json.Unmarshal(buf.Bytes(), &registrationAnswer_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("registrationAnswer_struct = %v\n", registrationAnswer_struct)
		return registrationAnswer_struct.Status
	case *Add_children_struct:
		body, err := json.Marshal(json_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("json string:\n\t%s\n", string(body))
		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error happend", err)
			panic(err)
		}
		defer resp.Body.Close() // важный пункт!
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		err = json.Unmarshal(buf.Bytes(), &add_childrenAnswer_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("enter_struct = %v\n", add_childrenAnswer_struct)
		return add_childrenAnswer_struct.Status
	case *Sign_contract_struct:
		body, err := json.Marshal(json_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("json string:\n\t%s\n", string(body))
		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error happend", err)
			panic(err)
		}
		defer resp.Body.Close() // важный пункт!
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		err = json.Unmarshal(buf.Bytes(), &sign_contract_answer_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("enter_struct = %v\n", sign_contract_answer_struct)
		return sign_contract_answer_struct.Status
	case *Finish_contract_struct:
		body, err := json.Marshal(json_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("json string:\n\t%s\n", string(body))
		req, err = http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error happend", err)
			panic(err)
		}
		defer resp.Body.Close() // важный пункт!
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		err = json.Unmarshal(buf.Bytes(), &finish_contract_answer_struct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("enter_struct = %v\n", finish_contract_answer_struct)
		return finish_contract_answer_struct.Status
	default:
		fmt.Println("Переданный в функцию пустой интерфейс не является определенной в проекте json-структурой!")
		return "NO_SUCH_JSON_STRUCT"
	}
}

func Main_Window() {
	(*w_reg).Hide()
	(*w_main).SetMaster()
	topWindow = (*w_main)
	content := container.NewMax()
	title := widget.NewLabel("Warning!") //Component name
	intro := widget.NewLabel("Only contracted children are displayed in the list. \nYou need to add a child and then add a contract for that child. \nOnce both of these actions are completed, \nthe child will be displayed in the list on the left.")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t tutorials.Tutorial) {
		title.SetText(t.Title)
		intro.SetText(t.Intro)
		content.Objects = []fyne.CanvasObject{t.View((*w_main))}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		(*w_main).SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		(*w_main).SetContent(split)
	}
	(*w_main).Resize(fyne.NewSize(520, 620))
	(*w_main).SetFixedSize(true)
	(*w_main).CenterOnScreen()
	(*w_main).Show()
}

func Log_in() {
	if Parent_login_enter.Text == "" {
		info_msg := "Поле 'Login' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_enter))
		fmt.Println("Error in Enter_request!")
		return
	}
	if Parent_password_enter.Text == "" {
		info_msg := "Поле 'Password' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_enter))
		fmt.Println("Error in Enter_request!")
		return
	}

	enter_struct.Login_parent = Parent_login_enter.Text
	enter_struct.Password_parent = Parent_password_enter.Text
	var Status string = SendPostRequest(URL_ENTER, &enter_struct)
	fmt.Println("Вывод структуры для входа:\n\t", enter_struct)
	fmt.Println("Status = ", Status)
	if Status != "OK" {
		if Status == "AUTH_FALL" {
			Parent_login_enter.Text = ""
			Parent_login_enter.Refresh()
			Parent_password_enter.Text = ""
			Parent_password_enter.Refresh()
			err := errors.New("Authentication error")
			dialog.ShowError(err, (*w_enter))
			fmt.Println("Error in Enter_request!")
			return
		}
	}

	(*w_enter).Hide()
	Refresh_main_window()
}

func Refresh_main_window() {
	// Заполняем дерево
	if TutorialIndex_my != nil {
		for k := range *TutorialIndex_my {
			delete(*TutorialIndex_my, k)
		}
	}
	TutorialIndex_my_buf := make(map[string][]string)
	Tutorials_my_buf := make(map[string]tutorials.Tutorial)

	//Формируем Список всех детей в виде ("Имя Фамилия")
	var childrens []string
	for _, val := range enterAnswer_struct.Childrens {
		childrens = append(childrens, val.Firstname_children+" "+val.Secondname_children)
		Tutorials_my_buf[val.Firstname_children+" "+val.Secondname_children] = tutorials.Tutorial{Title: val.Firstname_children + " " + val.Secondname_children, Intro: "", View: widgetChildrenInformation}
	}
	TutorialIndex_my_buf[""] = childrens
	// Формируем список контрактов для каждого ребенка
	for _, child := range childrens {
		var contrs []string
		for _, contr := range enterAnswer_struct.Contracts {
			if contr.Firstname_children+" "+contr.Secondname_children == child {
				contrs = append(contrs, contr.Contract_name)
				Tutorials_my_buf[contr.Contract_name] = tutorials.Tutorial{contr.Contract_name, "", widgetContractInformation}
			}
		}
		TutorialIndex_my_buf[child] = contrs
	}
	TutorialIndex_my = &TutorialIndex_my_buf
	Tutorials_my = &Tutorials_my_buf
	(*w_main).Content().Refresh()
	Main_Window()
}

// Нужно передавать в подобные (формирующие JSON-структур) функции адреса объектов!!!
func Add_registration_info_to_server() {

	if Parent_login.Text == "" {
		info_msg := "Поле 'Login' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_reg))
		fmt.Println("Error in Registration_request!")
		return
	}

	if Parent_password.Text == "" {
		info_msg := "Поле 'Password' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_reg))
		fmt.Println("Error in Registration_request!")
		return
	}

	if len(Photo_Parent.Resource.Content()) == 0 {
		info_msg := "Добавление фотографии обязательно."
		dialog.ShowInformation("Information", info_msg, (*w_reg))
		fmt.Println("Error in Registration_request!")
		return
	}

	if Parent_name.Text == "" {
		info_msg := "Поле 'Firstname' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_reg))
		fmt.Println("Error in Registration_request!")
		return
	}

	if Parent_surname.Text == "" {
		info_msg := "Поле 'Surname' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_reg))
		fmt.Println("Error in Registration_request!")
		return
	}

	// Тут заполняем JSON-структуру
	registration_struct.Photo_parent = base64.StdEncoding.EncodeToString(Photo_Parent.Resource.Content())
	registration_struct.Firstname_parent = Parent_name.Text
	registration_struct.Secondname_parent = Parent_surname.Text
	registration_struct.Login_parent = Parent_login.Text
	registration_struct.Password_parent = Parent_password.Text
	enter_struct.Login_parent = Parent_login.Text
	enter_struct.Password_parent = Parent_password.Text

	var Status string = SendPostRequest(URL_REGISTRACTION, &registration_struct)
	fmt.Println("Вывод структуры ответа регистрации:\n\t", registrationAnswer_struct)
	fmt.Println("Status = ", Status)
	if Status != "OK" {
		if Status == "LOGIN_BUSY" {
			Parent_login.Text = ""
			Parent_login.Refresh()
			info_msg := "This login is busy. Try another one."
			dialog.ShowInformation("Information", info_msg, (*w_reg))
			fmt.Println("Error in Registration_request!")
			return
		}
	}

	var Status_ string = SendPostRequest(URL_ENTER, &enter_struct)
	fmt.Println("Вывод структуры для входа:\n\t", enter_struct)

	fmt.Println("Status_ = ", Status_)
	if Status != "OK" {
		fmt.Println("Error in Enter_request!")
		// Вывести окошко ошибки!
	}

	(*w_reg).Hide()
	Refresh_main_window()
}

func Registration_window() {
	(*w_enter).Hide()
	(*w_reg).SetMaster()
	image_data, err := ioutil.ReadFile(ADD_PHOTO_ICON)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "add_photo.jpg", content: image_data}

	image_data_, err := ioutil.ReadFile(REG_ICON)
	if err != nil {
		panic(err.Error())
	}
	var my_image_ My_Image
	my_image_ = My_Image{name: "reg.jpg", content: image_data_}
	topWindow = *w_reg
	content := container.NewMax()
	title := widget.NewLabel("Registration")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true
	login := widget.NewEntry()
	login.SetPlaceHolder("Login")
	login.FocusGained()
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	firstname := widget.NewEntry()
	firstname.SetPlaceHolder("Firstname")
	surname := widget.NewEntry()
	surname.SetPlaceHolder("Surname")
	Picture_ := canvas.NewImageFromResource((*&picture_registration))
	Picture = *Picture_
	Photo_Parent = &Picture
	Parent_name = firstname
	Parent_surname = surname
	Parent_login = login
	Parent_password = password
	Add_photo := widget.NewButton("Add photo", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, *w_reg)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			imageOpened(reader)
		}, *w_reg)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
	Add_photo_ptr = Add_photo
	Add_photo.Alignment = widget.ButtonAlignCenter
	Add_photo.SetIcon(&my_image)
	Registration_button := widget.NewButton("Registration", Add_registration_info_to_server)
	Registration_button.Alignment = widget.ButtonAlignCenter
	Registration_button.IconPlacement = widget.ButtonIconLeadingText
	Registration_button_ptr = Registration_button
	Registration_button.SetIcon(&my_image_)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true
	nc := container.NewVBox(&Picture)
	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), login, password, nc), container.NewVBox(Add_photo, firstname, surname, Registration_button), nil, nil, content)
	(*w_reg).SetContent(tutorial)
	(*w_reg).Resize(fyne.NewSize(400, 600))
	(*w_reg).SetFixedSize(true)
	(*w_reg).CenterOnScreen()
	(*w_reg).Show()
}

func get_ver_space(number_lines int) *widget.Label {
	var vert_str string = ""
	for i := 0; i < number_lines-1; i++ {
		vert_str += "\n"
	}
	vert_space := widget.NewLabel(vert_str)
	return vert_space
}

func Enter_Window() {
	// Окно входа в систему
	image_data, err := ioutil.ReadFile(LOG_IN_ICON)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "log_in.jpg", content: image_data}

	image_data_, err := ioutil.ReadFile(REG_ICON)
	if err != nil {
		panic(err.Error())
	}
	var my_image_ My_Image
	my_image_ = My_Image{name: "reg.jpg", content: image_data_}

	topWindow = (*w_enter)
	content := container.NewMax()
	title := widget.NewLabel("Authorization")

	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	login := widget.NewEntry()
	login.SetPlaceHolder("Login")
	login.FocusGained()
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	Parent_login_enter = login
	Parent_password_enter = password

	Log_in_button := widget.NewButton("Log in", Log_in)
	Log_in_button.Alignment = widget.ButtonAlignCenter
	Log_in_button.SetIcon(&my_image)
	Registration_button := widget.NewButton("Registration", Registration_window)
	Registration_button.Alignment = widget.ButtonAlignCenter
	Registration_button.IconPlacement = widget.ButtonIconLeadingText

	Registration_button.SetIcon(&my_image_)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), login, password), container.NewVBox(Log_in_button, Registration_button), nil, nil, content)

	(*w_enter).SetContent(tutorial)

	(*w_enter).Resize(fyne.NewSize(350, 215))
	(*w_enter).SetFixedSize(true)
	(*w_enter).CenterOnScreen()
	(*w_enter).ShowAndRun()
}

var TutorialIndex_my *map[string][]string
var Tutorials_my *map[string]tutorials.Tutorial

func widgetChildrenInformation(_ fyne.Window) fyne.CanvasObject {
	var current_index_child int
	for indx, val := range enterAnswer_struct.Childrens {
		if val.Firstname_children+" "+val.Secondname_children == CurrentName {
			current_index_child = indx
			break
		}
	}
	// Нужно сохранить этот ИНДЕКС
	GLOBAL_current_index_child = &current_index_child
	image_data_after, err := base64.StdEncoding.DecodeString(enterAnswer_struct.Childrens[current_index_child].Photo_children)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	var my_image My_Image
	my_image = My_Image{name: "photo_child.jpg", content: image_data_after}
	Pict := canvas.NewImageFromResource((&my_image))
	Pict.FillMode = canvas.ImageFillStretch
	Pict.Resize(fyne.NewSize(200, 50))
	Pict.Move(fyne.NewPos(10, 25))
	Pict.Refresh()
	title_1 := widget.NewLabel("Children_bank_card_info")
	title_1.Alignment = fyne.TextAlignCenter
	// Необходимо узнать, чью информацию (какого именно ребенка считывать)
	ac_2 := widget.NewAccordion(
		widget.NewAccordionItem("Name", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.Name)),
		widget.NewAccordionItem("PersonalAcc", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.PersonalAcc)),
		widget.NewAccordionItem("BankName", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.BankName)),
		widget.NewAccordionItem("BIC", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.BIC)),
		widget.NewAccordionItem("CorrespAcc", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.CorrespAcc)),
		widget.NewAccordionItem("KPP", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.KPP)),
		widget.NewAccordionItem("PayeeINN", widget.NewLabel(enterAnswer_struct.Childrens[current_index_child].Children_CARD_info.PayeeINN)))
	cont := container.NewVScroll(container.NewVBox(fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(90, 90)), canvas.NewImageFromResource((&My_Image{})), Pict), get_ver_space(1), title_1, widget.NewSeparator(), ac_2))
	return cont
}

func widgetContractInformation(w fyne.Window) fyne.CanvasObject {
	title_ := widget.NewLabel("Diary_children_credits")
	title_.Alignment = fyne.TextAlignCenter
	var current_indx_contract int
	for indx, val := range enterAnswer_struct.Contracts {
		if val.Firstname_children+" "+val.Secondname_children == enterAnswer_struct.Childrens[*GLOBAL_current_index_child].Firstname_children+" "+enterAnswer_struct.Childrens[*GLOBAL_current_index_child].Secondname_children {
			if val.Contract_name == CurrentName {
				current_indx_contract = indx
				break
			}
		}
	}
	wl_1 := widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Firstname_children)
	wl_2 := widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Secondname_children)
	Children_name_contract_finish = wl_1
	Children_surname_contract_finish = wl_2
	ac_ := widget.NewAccordion(
		widget.NewAccordionItem("Children_name", wl_1),
		widget.NewAccordionItem("Children_surname", wl_2),
	)
	title_1 := widget.NewLabel("Marks policy")
	title_1.Alignment = fyne.TextAlignCenter
	ac_1 := widget.NewAccordion(
		widget.NewAccordionItem("One", widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Mark_policy.One)),
		widget.NewAccordionItem("Two", widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Mark_policy.Two)),
		&widget.AccordionItem{
			Title:  "Three",
			Detail: widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Mark_policy.Three),
		},
	)
	ac_1.Append(widget.NewAccordionItem("Four", widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Mark_policy.Four)))
	ac_1.Append(widget.NewAccordionItem("Five", widget.NewLabel(enterAnswer_struct.Contracts[current_indx_contract].Mark_policy.Five)))
	image_data, err := ioutil.ReadFile(FINISH_CONTRACT)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "finish_contract.jpg", content: image_data}
	Remove_contract_button := widget.NewButton("Finish contract", Finish_contract)
	Remove_contract_button.Alignment = widget.ButtonAlignCenter
	Remove_contract_button.SetIcon(&my_image)
	Remove_contract_button.Disable()
	cont := container.NewVScroll(container.NewVBox(title_, widget.NewSeparator(), ac_, get_ver_space(1), title_1, widget.NewSeparator(), ac_1, widget.NewSeparator(), get_ver_space(1), widget.NewSeparator(), Remove_contract_button))
	return cont
}

func Finish_contract() {
	finish_contract_struct.Login_parent = enter_struct.Login_parent
	finish_contract_struct.Password_parent = enter_struct.Password_parent
	finish_contract_struct.Contract_.Firstname_children = Children_name_contract_finish.Text
	finish_contract_struct.Contract_.Secondname_children = Children_surname_contract_finish.Text
	finish_contract_struct.Contract_.Contract_name = CurrentName //*Contract_Name_finish

	// Заглушка
	finish_contract_struct.Photo_condition_parent = base64.StdEncoding.EncodeToString([]byte("pomp"))
	finish_contract_struct.Photo_condition_children = base64.StdEncoding.EncodeToString([]byte("pomp"))

	var Status string = SendPostRequest(URL_FINISH_CONTRACT, &finish_contract_struct)

	// Вывод структуры
	fmt.Println("finish_contract_struct = \n\t", finish_contract_struct)

	fmt.Println("Status = ", Status)
	if Status != "OK" {
		fmt.Println("Error in Finish_contract_request!")
		// Вывести окошко ошибки!
	}

	var Status_ string = SendPostRequest(URL_ENTER, &enter_struct)
	fmt.Println("Вывод структуры для входа:\n\t", enter_struct)

	fmt.Println("Status_ = ", Status_)
	if Status != "OK" {
		fmt.Println("Error in Enter_request!")
		// Вывести окошко ошибки!
	}

	(*w_main).Hide()
	Refresh_main_window()
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	image_data_after, err := base64.StdEncoding.DecodeString(enterAnswer_struct.Photo_parent)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	var my_image_2 My_Image
	my_image_2 = My_Image{name: "photo_parent.jpg", content: image_data_after}
	Pict := canvas.NewImageFromResource((&my_image_2))
	Pict.FillMode = canvas.ImageFillStretch
	title_0 := widget.NewLabel(enterAnswer_struct.Firstname_parent)
	title_0.TextStyle = fyne.TextStyle{Bold: true, Italic: true, Monospace: false}
	title_0.Alignment = fyne.TextAlignCenter
	title_00 := widget.NewLabel(enterAnswer_struct.Secondname_parent)
	title_00.TextStyle = fyne.TextStyle{Bold: true, Italic: true, Monospace: false}
	title_00.Alignment = fyne.TextAlignCenter
	title_and_picture_cont := container.NewVBox(container.NewHBox(fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(90, 90)), Pict), widget.NewSeparator(), container.NewVBox(title_0, title_00)), widget.NewSeparator())
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return (*TutorialIndex_my)[uid] //TutorialIndex_[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := (*TutorialIndex_my)[uid] //TutorialIndex_[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := (*Tutorials_my)[uid] //Tutorials_[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := (*Tutorials_my)[uid]; ok { //Tutorials_[uid]; ok {
				CurrentName = uid
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}
	if loadPrevious {
		if len(enterAnswer_struct.Childrens) > 0 && len(enterAnswer_struct.Contracts) > 0 {
			tree.Select(enterAnswer_struct.Childrens[0].Firstname_children + " " + enterAnswer_struct.Childrens[0].Secondname_children)
		}
	}
	themes := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			LOG_IN_ICON = "./icons/dark_theme/log_in_icon.png"
			REG_ICON = "./icons/dark_theme/reg_icon.png"
			ADD_PHOTO_ICON = "./icons/dark_theme/add_photo_icon.png"
			FINISH_CONTRACT = "./icons/dark_theme/finish_contract_icon.png"
			ADD_CONTRACT = "./icons/dark_theme/add_contract_icon.png"
			ADD_CHILDREN = "./icons/dark_theme/add_children_icon.png"
			APP_ICON = "./icons/dark_theme/app_icon.png"
			Refresh_main_window()
			a.Settings().SetTheme(theme.DarkTheme())
			// Меняем иконки
		}),
		widget.NewButton("Light", func() {
			LOG_IN_ICON = "./icons/light_theme/log_in_icon.png"
			REG_ICON = "./icons/light_theme/reg_icon.png"
			ADD_PHOTO_ICON = "./icons/light_theme/add_photo_icon.png"
			FINISH_CONTRACT = "./icons/light_theme/finish_contract_icon.png"
			ADD_CONTRACT = "./icons/light_theme/add_contract_icon.png"
			ADD_CHILDREN = "./icons/light_theme/add_children_icon.png"
			APP_ICON = "./icons/light_theme/app_icon.png"
			Refresh_main_window()
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)
	image_data, err := ioutil.ReadFile(ADD_CONTRACT)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "add_contract.jpg", content: image_data}
	image_data_, err := ioutil.ReadFile(ADD_CHILDREN)
	if err != nil {
		panic(err.Error())
	}
	var my_image_ My_Image
	my_image_ = My_Image{name: "add_children.jpg", content: image_data_}
	Add_contract_button = widget.NewButton("Add contract", Add_contract)
	Add_contract_button.Alignment = widget.ButtonAlignCenter
	Add_contract_button.SetIcon(&my_image)
	Add_contract_button_ptr = Add_contract_button
	Add_children_button := widget.NewButton("Add children", Add_children)
	Add_children_button.Alignment = widget.ButtonAlignCenter
	Add_children_button.SetIcon(&my_image_)
	Add_children_button_ptr = Add_children_button
	title := widget.NewLabel("↑ Childrens | Themes ↓")
	title.TextStyle = fyne.TextStyle{Bold: false, Italic: true, Monospace: true}
	title.Alignment = fyne.TextAlignCenter
	title_1 := widget.NewLabel("Childrens")
	title_1.Alignment = fyne.TextAlignCenter
	return container.NewBorder(title_and_picture_cont, container.NewVBox(Add_contract_button, Add_children_button, widget.NewSeparator(), title, widget.NewSeparator(), themes), nil, nil, container.NewVScroll(tree))
}

func Add_children() {
	(*w_main).Hide()
	(*w_add_children).SetMaster()
	topWindow = (*w_add_children)
	image_data, err := ioutil.ReadFile(ADD_PHOTO_ICON)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "add_photo.jpg", content: image_data}

	content := container.NewMax()
	title := widget.NewLabel("Children information")

	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	title_1 := widget.NewLabel("Children CARD information")

	title_1.Alignment = fyne.TextAlignCenter
	title_1.TextStyle.Bold = true

	firstname := widget.NewEntry()
	firstname.SetPlaceHolder("Firstname")
	surname := widget.NewEntry()
	surname.SetPlaceHolder("Surname")

	Children_name = firstname
	Children_surname = surname

	Picture_ := canvas.NewImageFromResource((*&picture_registration))
	Picture = *Picture_

	Photo_Children = &Picture

	Add_photo := widget.NewButton("Add photo", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, *w_add_children)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			imageOpened(reader)
		}, *w_add_children)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	Add_photo.Alignment = widget.ButtonAlignCenter
	Add_photo.SetIcon(&my_image)

	Add_children_button := widget.NewButton("Add_children", Add_children_info_to_server)
	Add_children_button.Alignment = widget.ButtonAlignCenter
	Add_children_button.IconPlacement = widget.ButtonIconLeadingText

	Add_children_button.SetIcon(&my_image)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	nc := container.NewVBox(&Picture)

	Name := widget.NewEntry()
	Name.SetPlaceHolder("Name")
	PersonalAcc := widget.NewEntry()
	PersonalAcc.SetPlaceHolder("PersonalAcc")
	BankName := widget.NewEntry()
	BankName.SetPlaceHolder("BankName")
	BIC := widget.NewEntry()
	BIC.SetPlaceHolder("BIC")
	CorrespAcc := widget.NewEntry()
	CorrespAcc.SetPlaceHolder("CorrespAcc")
	KPP := widget.NewEntry()
	KPP.SetPlaceHolder("KPP")
	PayeeINN := widget.NewEntry()
	PayeeINN.SetPlaceHolder("PayeeINN")

	_Name = Name
	_PersonalAcc = PersonalAcc
	_BankName = BankName
	_BIC = BIC
	_CorrespAcc = CorrespAcc
	_KPP = KPP
	_PayeeINN = PayeeINN

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), nc), container.NewVBox(get_ver_space(10), widget.NewSeparator(), Add_photo, firstname, surname, get_ver_space(1), title_1, widget.NewSeparator(),
			Name, PersonalAcc, BankName, BIC, CorrespAcc, KPP, PayeeINN, widget.NewSeparator(), Add_children_button), nil, nil, content) // content

	(*w_add_children).SetContent(tutorial)

	(*w_add_children).Resize(fyne.NewSize(400, 900))
	(*w_add_children).CenterOnScreen()
	(*w_add_children).Show()
}

func Add_children_info_to_server() {
	if len(Photo_Children.Resource.Content()) == 0 {
		info_msg := "Добавление фотографии обязательно."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if Children_name.Text == "" {
		info_msg := "Поле 'Firstname' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if Children_surname.Text == "" {
		info_msg := "Поле 'Surname' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _Name.Text == "" {
		info_msg := "Поле 'Name' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _PersonalAcc.Text == "" {
		info_msg := "Поле 'PersonalAcc' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _BankName.Text == "" {
		info_msg := "Поле 'BankName' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _BIC.Text == "" {
		info_msg := "Поле 'BIC' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _CorrespAcc.Text == "" {
		info_msg := "Поле 'CorrespAcc' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _KPP.Text == "" {
		info_msg := "Поле 'KPP' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	if _PayeeINN.Text == "" {
		info_msg := "Поле 'PayeeINN' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_children))
		fmt.Println("Error in Add_children_request!")
		return
	}

	// Заполнение JSON-структуры
	add_children_struct.Login_parent = enter_struct.Login_parent
	add_children_struct.Password_parent = enter_struct.Password_parent
	add_children_struct.Firstname_children = Children_name.Text
	add_children_struct.Secondname_children = Children_surname.Text
	add_children_struct.Photo_children = base64.StdEncoding.EncodeToString(Photo_Children.Resource.Content())
	add_children_struct.Children_CARD_info.Name = _Name.Text
	add_children_struct.Children_CARD_info.PersonalAcc = _PersonalAcc.Text
	add_children_struct.Children_CARD_info.BankName = _BankName.Text
	add_children_struct.Children_CARD_info.BIC = _BIC.Text
	add_children_struct.Children_CARD_info.CorrespAcc = _CorrespAcc.Text
	add_children_struct.Children_CARD_info.KPP = _KPP.Text
	add_children_struct.Children_CARD_info.PayeeINN = _PayeeINN.Text

	// Заглушка
	add_children_struct.Diary_children_login = "pomp"
	add_children_struct.Diary_children_password = "pomp"

	var Status string = SendPostRequest(URL_ADD_CHILDREN, &add_children_struct)

	// Вывод структуры
	fmt.Println("add_children_struct = \n\t", add_children_struct)

	fmt.Println("Status = ", Status)
	if Status != "OK" {
		fmt.Println("Error in Add_children_request!")
		// Вывести окошко ошибки!
	}

	var Status_ string = SendPostRequest(URL_ENTER, &enter_struct)
	fmt.Println("Вывод структуры для входа:\n\t", enter_struct)

	fmt.Println("Status_ = ", Status_)
	if Status != "OK" {
		fmt.Println("Error in Enter_request!")
		// Вывести окошко ошибки!
	}

	(*w_add_children).Hide()
	Refresh_main_window()
}

func Add_contract() {
	(*w_main).Hide()
	(*w_add_contract).SetMaster()
	topWindow = (*w_add_contract)

	image_data, err := ioutil.ReadFile(ADD_CONTRACT)
	if err != nil {
		panic(err.Error())
	}
	var my_image My_Image
	my_image = My_Image{name: "sign_contract_icon.jpg", content: image_data}

	content := container.NewMax()
	title := widget.NewLabel("Contract information")

	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	title_1 := widget.NewLabel("Marks policy")

	title_1.Alignment = fyne.TextAlignCenter
	title_1.TextStyle.Bold = true

	firstname := widget.NewEntry()
	firstname.SetPlaceHolder("Firstname children")
	surname := widget.NewEntry()
	surname.SetPlaceHolder("Surname children")

	Children_name_contract = firstname
	Children_surname_contract = surname

	Contract_Name_enter := widget.NewEntry()
	Contract_Name_enter.SetPlaceHolder("Contract Name (Must be unique)")
	Contract_Name = Contract_Name_enter

	Diary_login := widget.NewEntry()
	Diary_login.SetPlaceHolder("Diary login (Must be unique)")
	_Diary_children_login = Diary_login

	Diary_password := widget.NewPasswordEntry()
	Diary_password.SetPlaceHolder("Diary password")
	_Diary_children_password = Diary_password

	Timer_ := newNumEntry()
	Timer_.SetPlaceHolder("Time interval of notisfaction in days (only digits)")

	Timer_interval_notisfaction = Timer_

	paint_mark_policy := func() fyne.CanvasObject {

		text_obj1 = &widget.Entry{Text: "0"}
		text_obj2 = &widget.Entry{Text: "0"}
		text_obj3 = &widget.Entry{Text: "40"}
		text_obj4 = &widget.Entry{Text: "80"}
		text_obj5 = &widget.Entry{Text: "100"}

		obj1 := widget.NewAccordionItem("One", text_obj1)
		obj2 := widget.NewAccordionItem("Two", text_obj2)
		obj3 := widget.NewAccordionItem("Three", text_obj3)
		obj4 := widget.NewAccordionItem("Four", text_obj4)
		obj5 := widget.NewAccordionItem("Five", text_obj5)

		ac := widget.NewAccordion(obj1, obj2, obj3, obj4, obj5)

		return ac
	}
	mark_policy_tree := paint_mark_policy()

	Sign_contract_button := widget.NewButton("Sign_contract", Sign_contract_info_to_server)
	Sign_contract_button.Alignment = widget.ButtonAlignCenter
	Sign_contract_button.IconPlacement = widget.ButtonIconLeadingText
	Sign_contract_button.SetIcon(&my_image)
	Sign_contract_button_ptr = Sign_contract_button
	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), Contract_Name_enter, firstname, surname, Diary_login, Diary_password, Timer_, widget.NewSeparator(), get_ver_space(2), title_1, widget.NewSeparator(), mark_policy_tree), container.NewVBox(get_ver_space(1), widget.NewSeparator(), Sign_contract_button), nil, nil, content) // content
	(*w_add_contract).SetContent(tutorial)
	(*w_add_contract).Resize(fyne.NewSize(400, 600))
	(*w_add_children).SetFixedSize(true)
	(*w_add_contract).CenterOnScreen()
	(*w_add_contract).Show()
}

func Sign_contract_info_to_server() {
	if Contract_Name.Text == "" {
		info_msg := "Поле 'Contract Name' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if Children_name_contract.Text == "" {
		info_msg := "Поле 'Firstname children' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if Children_surname_contract.Text == "" {
		info_msg := "Поле 'Surname children' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if _Diary_children_login.Text == "" {
		info_msg := "Поле 'Diary login' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if _Diary_children_password.Text == "" {
		info_msg := "Поле 'Diary password' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if Timer_interval_notisfaction.Text == "" {
		info_msg := "Поле 'Time interval ...' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if text_obj1.Text == "" {
		info_msg := "Поле 'Marks policy.One' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if text_obj2.Text == "" {
		info_msg := "Поле 'Marks policy.Two' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if text_obj3.Text == "" {
		info_msg := "Поле 'Marks policy.Three' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if text_obj4.Text == "" {
		info_msg := "Поле 'Marks policy.Four' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}

	if text_obj5.Text == "" {
		info_msg := "Поле 'Marks policy.Five' не должно быть пустым."
		dialog.ShowInformation("Information", info_msg, (*w_add_contract))
		fmt.Println("Error in Sign_contract_request!")
		return
	}
	sign_contract_struct.Login_parent = enter_struct.Login_parent
	sign_contract_struct.Password_parent = enter_struct.Password_parent
	sign_contract_struct.Diary_children_login = _Diary_children_login.Text
	sign_contract_struct.Diary_children_password = _Diary_children_password.Text
	sign_contract_struct.Timer_interval_notisfaction = Timer_interval_notisfaction.Text
	
	// Заглушки
	sign_contract_struct.Photo_condition_parent = base64.RawStdEncoding.EncodeToString([]byte("pomp"))
	sign_contract_struct.Photo_condition_children = base64.RawStdEncoding.EncodeToString([]byte("pomp"))
	sign_contract_struct.Contract_.Photo_children = base64.RawStdEncoding.EncodeToString([]byte("pomp"))
	
	sign_contract_struct.Contract_.Firstname_children = Children_name_contract.Text
	sign_contract_struct.Contract_.Secondname_children = Children_surname_contract.Text
	sign_contract_struct.Contract_.Contract_name = Contract_Name.Text
	
	sign_contract_struct.Contract_.Mark_policy.One = text_obj1.Text
	sign_contract_struct.Contract_.Mark_policy.Two = text_obj2.Text
	sign_contract_struct.Contract_.Mark_policy.Three = text_obj3.Text
	sign_contract_struct.Contract_.Mark_policy.Four = text_obj4.Text
	sign_contract_struct.Contract_.Mark_policy.Five = text_obj5.Text

	// Send data_to_server, обработка ответа
	var Status string = SendPostRequest(URL_SIGN_CONTRACT, &sign_contract_struct)

	// Вывод структуры
	fmt.Println("sign_contract_struct = \n\t", sign_contract_struct)

	fmt.Println("Status = ", Status)
	if Status != "OK" {
		fmt.Println("Error in Sign_contract_request!")
		// Вывести окошко ошибки!
	}

	var Status_ string = SendPostRequest(URL_ENTER, &enter_struct)
	fmt.Println("Вывод структуры для входа:\n\t", enter_struct)

	fmt.Println("Status_ = ", Status_)
	if Status != "OK" {
		fmt.Println("Error in Enter_request!")
		// Вывести окошко ошибки!
	}

	(*w_add_contract).Hide()
	Refresh_main_window()
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func imageOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()
	showImage(f)
}

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillContain

	// Вот как брать данные из загруженного изображения !!!
	img_bytes := img.Resource.Content()
	(*picture_registration).content = img_bytes
	(*picture_registration).name = f.URI().Name()
	Picture_ := canvas.NewImageFromResource((*&picture_registration))
	Picture = *Picture_
	Picture.FillMode = canvas.ImageFillStretch
	Picture.Resize(fyne.NewSize(200, 250))
	Picture.Move(fyne.NewPos(100, 25))
	Picture.Refresh()
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)
	return canvas.NewImageFromResource(res)
}

type numEntry struct {
	widget.Entry
}

func newNumEntry() *numEntry {
	e := &numEntry{}
	e.ExtendBaseWidget(e)
	e.Validator = validation.NewRegexp(`\d`, "Must contain a number")
	return e
}
