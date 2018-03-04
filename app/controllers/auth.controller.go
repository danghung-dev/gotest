package controllers

import (
	"log"
	"net/http"
	"strconv"
	 "gotest/app"
	 "gotest/app/models"
	 "gotest/app/services"
	 "gotest/app/utils"

	"golang.org/x/oauth2"
	"fmt"
	"net/url"
	"strings"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/tools/go/gcimporter15/testdata"
)

type AuthController struct {
	App *app.App
	*models.UserHelper
	jwtService services.JWTAuthService
	ggConfig *oauth2.Config
}
var (
	oauthConf = &oauth2.Config{
		ClientID:     "324311177982130",
		ClientSecret: "65a6389320e7ee747278598fd84d224e",
		RedirectURL:  "http://localhost:3002/api/v1/auth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	a = &oauth2.Config{}
	oauthStateString = "dung de kiem tra state"
)
func NewAuthController(a *app.App, us *models.UserHelper, jwtService services.JWTAuthService) *AuthController {
	//oauthConf = &oauth2.Config{
	//
	//}
	ggConfig := &oauth2.Config{
		ClientID:     a.Config.Google.ClientID,
		ClientSecret: a.Config.Google.ClientSecret,
		RedirectURL:  a.Config.Google.RedirectURL,
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &AuthController{a, us, jwtService, ggConfig}
}

func (ac *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	ggoauthConf := ac.ggConfig
	Url, err := url.Parse(ggoauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", ggoauthConf.ClientID)
	parameters.Add("scope", strings.Join(ggoauthConf.Scopes, " "))
	parameters.Add("redirect_uri", ggoauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ac *AuthController) GoogleCallBack(w http.ResponseWriter, r *http.Request) {
	ggoauthConf := ac.ggConfig
	code := r.FormValue("code")

	token, err := ggoauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	client := ggoauthConf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := GetJSON(resp.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}

	email, err := response.GetString("email")
	//var user models.User
	user, err := ac.UserHelper.FindByEmail(email)

	if err != nil && err.Error() == "record not found" {
		user = &models.User{ Email: email}
		result := ac.App.Database.Create(&user)
		if result.Error != nil {
			NewAPIError(&APIError{false, "Cannot create user in db", http.StatusBadRequest}, w)
			return
		}
	}

	tokens, err := ac.jwtService.GenerateTokens(user)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}
	authUser := &models.AuthUser{
		User:  user,
		Admin: user.Admin,
	}

	data := struct {
		Tokens *services.Tokens `json:"tokens"`
		User   *models.AuthUser `json:"user"`
	}{
		tokens,
		authUser,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: data}, w, http.StatusOK)
}

func (ac *AuthController) FaceBookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ac *AuthController) FaceBookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//defer resp.Body.Close()

	//response, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("ReadAll: %s\n", err)
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}

	response, err := GetJSON(resp.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	id, err := response.GetString("id")
	name, err := response.GetString("name")
	//var user models.User
	user, err := ac.UserHelper.FindByFacebookId(id)

	if err != nil && err.Error() == "record not found" {
		user = &models.User{ FacebookId: id, Name: name}
		result := ac.App.Database.Create(&user)
		if result.Error != nil {
			NewAPIError(&APIError{false, "Cannot create user in db", http.StatusBadRequest}, w)
			return
		}
	}

	tokens, err := ac.jwtService.GenerateTokens(user)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}
	authUser := &models.AuthUser{
		User:  user,
		Admin: user.Admin,
	}

	data := struct {
		Tokens *services.Tokens `json:"tokens"`
		User   *models.AuthUser `json:"user"`
	}{
		tokens,
		authUser,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: data}, w, http.StatusOK)

}

func (ac *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {
	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	email, err := j.GetString("email")
	if err != nil {
		NewAPIError(&APIError{false, "Email is required", http.StatusBadRequest}, w)
		return
	}
	if ok := utils.IsEmail(email); !ok {
		NewAPIError(&APIError{false, "You must provide a valid email address", http.StatusBadRequest}, w)
		return
	}
	u, err := ac.UserHelper.FindByEmail(email)
	if err != nil {
		NewAPIError(&APIError{false, "Incorrect email or password", http.StatusBadRequest}, w)
		return
	}

	pw, err := j.GetString("password")
	if err != nil {
		NewAPIError(&APIError{false, "Password is required", http.StatusBadRequest}, w)
		return
	}

	if ok := u.CheckPassword(pw); !ok {
		NewAPIError(&APIError{false, "Incorrect email or password", http.StatusBadRequest}, w)
		return
	}

	tokens, err := ac.jwtService.GenerateTokens(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	authUser := &models.AuthUser{
		User:  u,
		Admin: u.Admin,
	}

	data := struct {
		Tokens *services.Tokens `json:"tokens"`
		User   *models.AuthUser `json:"user"`
	}{
		tokens,
		authUser,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: data}, w, http.StatusOK)
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString, err := services.GetTokenFromRequest(&ac.App.Config, r)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	tokenHash, err := services.ExtractTokenHash(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	/*jti, err := services.ExtractJti(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}*/

	keys := ac.App.Redis.Keys("*" + tokenHash + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
			NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
			return
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Logout successful"}, w, http.StatusOK)

}

func (ac *AuthController) LogoutAll(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	userId := strconv.Itoa(uid)
	keys := ac.App.Redis.Keys("*." + userId + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Logout successful"}, w, http.StatusOK)
}

func (ac *AuthController) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	tokenString, err := services.GetRefreshTokenFromRequest(&ac.App.Config, r)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	tokenHash, err := services.ExtractRefreshTokenHash(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}
	u, err := ac.UserHelper.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusBadRequest}, w)
		return
	}
	tokens, err := ac.jwtService.GenerateTokens(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	keys := ac.App.Redis.Keys("*" + tokenHash + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
			NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
			return
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Refresh successful", Data: tokens}, w, http.StatusOK)
}
