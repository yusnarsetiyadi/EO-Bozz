package service

import (
	"capstone-alta1/features/user"
	"capstone-alta1/utils/helper"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository user.RepositoryInterface
	validate       *validator.Validate
}

func New(repo user.RepositoryInterface) user.ServiceInterface {
	return &userService{
		userRepository: repo,
		validate:       validator.New(),
	}
}

// Create implements user.ServiceInterface
func (service *userService) Create(input user.Core, c echo.Context) (err error) {
	// validasi input
	if errValidate := service.validate.Struct(input); errValidate != nil {
		return errValidate
	}

	// validasi email harus unik

	_, errFindEmail := service.userRepository.FindUser(input.Email)

	if errFindEmail != nil && !strings.Contains(errFindEmail.Error(), "found") {
		return helper.ServiceErrorMsg(errFindEmail)
	}

	bytePass, errEncrypt := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if errEncrypt != nil {
		log.Error(errEncrypt.Error())
		return err
	}

	input.Password = string(bytePass)
	input.Role = "Admin"
	errCreate := service.userRepository.Create(input)
	if errCreate != nil {
		log.Error(errCreate.Error())
		return err
	}

	return nil
}

func (service *userService) GetAll(query string) (data []user.Core, err error) {
	data, err = service.userRepository.GetAll()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return data, err
}

func (service *userService) GetById(id uint) (data user.Core, err error) {
	data, err = service.userRepository.GetById(id)
	if err != nil {
		log.Error(err.Error())
		return user.Core{}, err
	}

	return data, err

}

func (service *userService) Update(input user.Core, id uint) error {
	if input.Password != "" {
		generate, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
		input.Password = string(generate)
	}
	err := service.userRepository.Update(input, id)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (service *userService) Delete(id uint) error {
	err := service.userRepository.Delete(id)
	if err != nil {
		log.Error(err.Error())
		return helper.ServiceErrorMsg(err)
	}
	return nil
}

func (service *userService) UpdatePassword(input user.Core, id uint) error {
	err := service.userRepository.Update(input, id)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
