package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	//"time"

	"github.com/Ola-Daniel/qrcodebakery/internal/cookies"
	"github.com/Ola-Daniel/qrcodebakery/internal/request"
	"github.com/Ola-Daniel/qrcodebakery/internal/response"
	"github.com/Ola-Daniel/qrcodebakery/internal/validator"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/crypto/bcrypt"
)

var ImageFile string

var DynamicImageFile string

var ImageFileUploadPath string

var DynamicImageFileUploadPath string

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data["QRCodeImagePath"] = ImageFileUploadPath

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) privacypolicy(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/privacypolicy.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) disclaimer(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/disclaimer.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) tos(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/termsofservice.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	type loginForm struct {
		UsernameOrEmail string              `form:"Username"`
		Password        string              `form:"Password"`
		Validator       validator.Validator `form:"-"`
	} //

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)

		err := response.Page(w, http.StatusOK, data, "pages/login.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		var form loginForm
		err := request.DecodePostForm(r, &form)

		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		user, err := app.db.GetUser(form.UsernameOrEmail)

		if err != nil {
			//Handle the error (e.g., user not found)//
			app.invalidCredentials(w, r, err)
			return
		}

		//Compare the hashed password from the database with the password provided in the form

		err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(form.Password))
		if err != nil {
			//Handle the error (e.g., invalid password)
			app.invalidCredentials(w, r, err)
			return
		} //
		//authentication cookie
		cookie := http.Cookie{
			Name:     "isAuthenticated",
			Value:    "yes",
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		err = cookies.WriteEncrypted(w, cookie, app.config.cookie.secretKey)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
        uid := strconv.Itoa(user.ID)
		//uid := user.ID

		//authentication cookie
		cookie = http.Cookie{
			Name:     "userid",
			Value:    uid,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		err = cookies.WriteEncrypted(w, cookie, app.config.cookie.secretKey)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		//At this point, the user is successfully authenticated
		//Redirect the user to the rpotected page or perform any other action
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	type createUserForm struct {
		Username  string              `form:"Username"`
		Password  string              `form:"Password"`
		Email     string              `form:"Email"`
		Validator validator.Validator `form:"-"`
	} ///

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)

		err := response.Page(w, http.StatusOK, data, "pages/register.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		var form createUserForm

		err := request.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		//Validate the username
		if len(form.Username) < 8 {
			//If the username is too short , return an error
			app.badRequest(w, r, errors.New("username must be at least 8 characters long"))
			return
		}

		// Hash the password using bcrypt//

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		fmt.Println(hashedPassword)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		//err = bcrypt.CompareHashAndPassword([]byte(string(hashedPassword)), []byte(string(form.Password)))
		//if err != nil {
		//      // Handle the error (e.g., invalid password)
		//     app.invalidCredentials(w, r, err)
		//     return
		// }

		err = app.db.NewUser(string(form.Username), string(hashedPassword), string(form.Email))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
   //is authenticated?
	cookie := http.Cookie{
		Name:     "isAuthenticated",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	 

	http.SetCookie(w, &cookie)

	//user id 
	userCookie := http.Cookie{
		Name:     "userid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &userCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) viewQRCodes(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadEncrypted(r, "isAuthenticated", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequest(w, r, err)
		case errors.Is(err, cookies.ErrInvalidValue):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Println(value)


	value, err = cookies.ReadEncrypted(r, "userid", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequest(w, r, err)
		case errors.Is(err, cookies.ErrInvalidValue):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	    }
	

	userid, err := strconv.Atoi(value)
	if err != nil { 
		app.badRequest(w, r, err)   
	}


	//type QRCode struct {
	//	QrcodeID  int
	//	UserID    int
	//	Data      string
	//	ImagePath string
	//	CreatedAt time.Time
	//}







	qrcodes, err := app.db.GetAllQRCodesByUserID(userid) 
	if err != nil {
		app.badRequest(w, r, err)
	}

	fmt.Println(qrcodes)

	for _, code := range qrcodes {
        fmt.Printf("QR Code ID: %d, User ID: %d, Data: %s, Image Path: %s, Created At: %s\n",
            code.QrcodeID, code.UserID, code.Data, code.ImagePath, code.CreatedAt)
    }

	

	data := app.newTemplateData(r)

	data["QRCodeList"] = qrcodes

	err = response.DashboardPage(w, http.StatusOK, data, "pages/get_all_user_code.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) deleteQRCode(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadEncrypted(r, "isAuthenticated", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequest(w, r, err)
		case errors.Is(err, cookies.ErrInvalidValue):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Println(value)

}

func (app *application) createQRCode(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadEncrypted(r, "isAuthenticated", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequest(w, r, err)
		case errors.Is(err, cookies.ErrInvalidValue):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Println(value)


	value, err = cookies.ReadEncrypted(r, "userid", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequest(w, r, err)
		case errors.Is(err, cookies.ErrInvalidValue):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	    }
	

	userid, err := strconv.Atoi(value)
	if err != nil { 
		app.badRequest(w, r, err)   
	}

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)

		data["DynamicQRCodeImagePath"] = "../../" + DynamicImageFileUploadPath 

		err = response.DashboardPage(w, http.StatusOK, data, "pages/create_dynamic_code.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		type response struct {
			DataString string              `form:"dataString"`
			DataType   string              `form:"dataType"` //Field for radio button value
			Validator  validator.Validator `form:"-"`
		}
		var form response

		err := request.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		if form.DataType == "" {
			form.DataType = "URL"
		}

		form.Validator.CheckField(form.DataString != "", "Data", "Input data cannot be empty")

		if form.DataString == "" {
			app.badRequest(w, r, errors.New("input data cannot be empty"))
		}

		switch form.DataType {
		case "URL":

			//Validate URL and domain
			if err := validateURL(form.DataString); err != nil {
				app.badRequest(w, r, err)
			}

			//Validate Vcard
		case "Contact":

			if err := validateVcard(form.DataString); err != nil {
				app.badRequest(w, r, err)
			}
			//Validate Wifi Connection String
		case "WiFi":

			if err := validateWifi(form.DataString); err != nil {
				app.badRequest(w, r, err)
			}

		default:

		}
		qrc, err := qrcode.New(form.DataString)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		random := uuid.New().String()

		DynamicImageFile = "dynamic-qrcode-" + random + ".jpeg"

		DynamicImageFileUploadPath = "./files/generated/" + DynamicImageFile

		wr, err := standard.New(DynamicImageFileUploadPath)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// save file
		if err = qrc.Save(wr); err != nil {
			app.serverError(w, r, err)
			return
		}


		

		_, err = app.db.CreateQRCode(userid, form.DataString, DynamicImageFileUploadPath)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// Redirect back to homepage after successful generation
		http.Redirect(w, r, "/dashboard/create-qr-code", http.StatusSeeOther) 
 
	}

}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {

	value, err := cookies.ReadEncrypted(r, "isAuthenticated", app.config.cookie.secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		default:
			app.serverError(w, r, err)
			return
		}

	}
	fmt.Println(value)

	data := app.newTemplateData(r)

	err = response.DashboardPage(w, http.StatusOK, data, "pages/dashboard.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) admin(w http.ResponseWriter, r *http.Request) {

}

func (app *application) editQRCode(w http.ResponseWriter, r *http.Request) {
   
}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {

	type response struct {
		DataString string              `form:"dataString"`
		DataType   string              `form:"dataType"` //Field for radio button value
		Validator  validator.Validator `form:"-"`
	}
	var form response

	err := request.DecodePostForm(r, &form)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if form.DataType == "" {
		form.DataType = "URL"
	}

	form.Validator.CheckField(form.DataString != "", "Data", "Input data cannot be empty")

	if form.DataString == "" {
		app.badRequest(w, r, errors.New("input data cannot be empty"))
	}

	// Check if the "dynamic" parameter is present
	if r.FormValue("dynamic") != "" {
		// Redirect to login page if dynamic parameter is present
		http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		return
	}

	switch form.DataType {
	case "URL":

		//Validate URL and domain
		if err := validateURL(form.DataString); err != nil {
			app.badRequest(w, r, err)
		}

		//Validate Vcard
	case "Contact":

		if err := validateVcard(form.DataString); err != nil {
			app.badRequest(w, r, err)
		}
		//Validate Wifi Connection String
	case "WiFi":

		if err := validateWifi(form.DataString); err != nil {
			app.badRequest(w, r, err)
		}

	default:

	}
	qrc, err := qrcode.New(form.DataString)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	random := uuid.New().String()

	ImageFile = "generated-qrcode-" + random + ".jpeg"

	ImageFileUploadPath = "./files/generated/" + ImageFile

	wr, err := standard.New(ImageFileUploadPath)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// save file
	if err = qrc.Save(wr); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect back to homepage after successful generation
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

/*
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("isAuthenticated")
		if err != nil || cookie.Value != "yes" {
			// If cookie is not present or has invalid value, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
} */

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}

func validateURL(inputURL string) error {
	u, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return err
	}

	if u.Scheme == "" || u.Host == "" {
		return err
	}

	//Check if the host name is a valid domain
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(u.Host) {
		return err
	}

	return nil
}

func validateVcard(input string) error {
	vCardRegex := `BEGIN:VCARD*END:VCARD`
	match, err := regexp.MatchString(vCardRegex, input)
	if err != nil {
		return err
	}

	if !match {
		return err
	}
	return nil
}

func validateWifi(input string) error {

	wifiRegex := `WIFI:T:.*;s;.*;p:.*;;`
	match, err := regexp.MatchString(wifiRegex, input)
	if err != nil {
		return err
	}

	if !match {
		return err
	}
	return nil
}
