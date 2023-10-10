package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserHander handles restful calls for endpoint for user related actions
type UserHandler struct {
	CFG    *config.Config
	DB     *pg.DB
	Logger *logrus.Logger
}

type forgotEmailBody struct {
	Email string `json:"email"`
}

type newPasswordTokenBody struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

// ForgotPasswordHandler sends a new password out to the user with email
func (h *UserHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var emailBody forgotEmailBody
	err := decoder.Decode(&emailBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{}
	err = h.DB.Model(&user).Where("email = ?", emailBody.Email).Select()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ResetToken = h.generateResetToken(user.Password, user.GetID())
	err = h.DB.Update(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	domain := h.CFG.FrontEndURL
	link := fmt.Sprintf("http://%s/#/resetPassword?token=%s", domain, user.ResetToken)
	err = sendRequestReset(emailBody.Email, link, domain)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) generateResetToken(passwordHash []byte, login string) string {
	secret := []byte(h.CFG.ResetSecret)
	return passwordreset.NewToken(login, 12*time.Hour, passwordHash, secret)
}

func (h *UserHandler) getPasswordHash(login string) ([]byte, error) {
	// return password hash for the login,
	// or an error if there's no such user

	user := models.User{}
	err := h.DB.Model(&user).Where("id = ?", login).Select()
	if err != nil {
		return nil, err
	}
	return user.Password, nil
}

// NewPasswordWithResetToken takes in a reset token and validates it.  if valid sets the users pw to new value
func (h *UserHandler) NewPasswordWithResetToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body newPasswordTokenBody
	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	secret := []byte(h.CFG.ResetSecret)
	userID, err := passwordreset.VerifyToken(body.Token, h.getPasswordHash, secret)
	if err != nil {
		// verification failed, don't allow password reset
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{}
	err = h.DB.Model(&user).Where("reset_token = ? and id = ?", body.Token, repository.NewSQLCompatUUIDFromStr(userID)).Select()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Password = passwordHash
	user.ResetToken = ""
	err = h.DB.Update(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func sendRequestReset(email string, link string, domainname string) error {

	v := url.Values{}
	v.Set("from", "ESP Admin <espadmin@siimpl.io>")
	v.Add("to", email)
	v.Add("subject", "ESP Admin - Reset Password")
	v.Add("html", fmt.Sprintf("<html><b>Reset password</b> using this link: <a href=\"%s\">%s</a></html>", link, link))

	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/automailer.siimpl.io/messages", strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("api", "eef79f6b186795c78230a25b091754b9-52b0ea77-6a1c6b20")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Sending email failed: %s", resp.Status)
	} else {
		fmt.Printf("\nReset Password e-mail sent to :%s", email)
	}
	defer resp.Body.Close()
	return nil
}
