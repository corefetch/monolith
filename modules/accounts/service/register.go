package service

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"learnt.io/core"
	"learnt.io/modules/accounts/store"
	"learnt.io/modules/emails"
)

type CreateAcountData struct {
	Role core.UserRole `json:"role" validate:"required"`

	Login         string `json:"login" validate:"required"`
	Password      string `json:"password" validate:"required"`
	PasswordMatch string `json:"password_match" validate:"required"`
	Referrer      string `json:"referrer"`

	Names     []string            `json:"names" binding:"required"`
	Telephone string              `json:"telephone"`
	Birthday  time.Time           `json:"birthday"`
	Location  *store.UserLocation `json:"location"`
}

func Register(create CreateAcountData) (account *store.Account, err error) {

	validate := validator.New()
	err = validate.Struct(create)

	if err != nil {
		return nil, err
	}

	if err := ValidatePassword(create.Password); err != nil {
		return nil, err
	}

	if create.Password != create.PasswordMatch {
		return nil, errors.New("password does not match")
	}

	password, err := bcrypt.GenerateFromPassword(
		[]byte(create.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create password")
	}

	account = &store.Account{
		ID:       primitive.NewObjectID(),
		Role:     create.Role,
		Login:    create.Login,
		Password: string(password),
		Profile: store.Profile{
			Names:     create.Names,
			Birthday:  create.Birthday,
			Telephone: create.Telephone,
		},
		CreatedAt: time.Now(),
	}

	err = account.Save()

	if err != nil {
		return nil, err
	}

	// Send email after register
	switch account.Role {
	case core.RoleStudent:
		emails.AfterRegisterStudent(account, emails.TemplateParams{
			"NAME": strings.Join(account.Profile.Names, " "),
		})
	case core.RoleTutor:
		emails.AfterRegisterTutor(account, emails.TemplateParams{
			"NAME": strings.Join(account.Profile.Names, " "),
		})
	}

	return
}

func ValidateStudent(create CreateAcountData) (err error) {
	return nil
}

func ValidatePassword(password string) error {

	matchLower, _ := regexp.MatchString("[a-z]", password)
	matchUpper, _ := regexp.MatchString("[A-Z]", password)

	matchDigit, _ := regexp.MatchString("[0-9]", password)

	if !matchLower || !matchUpper || !matchDigit {
		return errors.New("password must contain at least a lowercase character, an uppercase character, and a digit")
	}

	if len(password) < 8 {
		return errors.New("password must be bigger than 8 characters")
	}

	return nil
}
