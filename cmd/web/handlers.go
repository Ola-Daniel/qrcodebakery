package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"errors"

	"github.com/Ola-Daniel/qrcodebakery/internal/request"
	"github.com/Ola-Daniel/qrcodebakery/internal/response"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/crypto/bcrypt"
)


var ImageFile string

var ImageFileUploadPath string

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
        UsernameOrEmail string `form:"Username"`
        Password string `form:"Password"`
    }//  

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
		}//


		//At this point, the user is successfully authenticated
		//Redirect the user to the rpotected page or perform any other action
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} 
}


func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	type createUserForm struct {
        Username string `form:"Username"`
        Password string `form:"Password"`
		Email    string `form:"Email"`
    }///  

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

	//Clear Cookies to sign out!!!!!!!!!!



	//
	//
	//

	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}


func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.DashboardPage(w, http.StatusOK, data, "pages/dashboard.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) admin(w http.ResponseWriter, r *http.Request) {

}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {

    type response struct {
		DataString string `form:"dataString"`
		DataType   string `form:"dataType"` //Field for radio button value
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
		if err := validateURL(form.DataString);
		    err != nil {
				app.badRequest(w, r, err)     
			}
		
        //Validate Vcard
	case "Contact": 


	    if err := validateVcard(form.DataString);
		    err != nil {
				app.badRequest(w, r, err)
			}
       //Validate Wifi Connection String
	case "WiFi":

		if err := validateWifi(form.DataString);
		    err != nil {
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


	ImageFile = "generated-qrcode-"+random+".jpeg"

	ImageFileUploadPath = "./files/generated/"+ImageFile     
	
 
	
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