package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	"github.com/bakare-dev/simple-bank-api/pkg/util"
)

type UserService struct {
	userRepo    userrepository.UserRepository
	profileRepo userrepository.ProfileRepository
}

func NewUserService(userRepo userrepository.UserRepository, profileRepo userrepository.ProfileRepository) *UserService {
	return &UserService{userRepo: userRepo, profileRepo: profileRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, int, error) {
	userExists, err := s.userRepo.Get(ctx, map[string]any{"email": user.Email})

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists != nil {
		return nil, http.StatusConflict, fmt.Errorf("user with this email already exists")
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err)
	}

	otp := util.RandomNumberString(6)

	fmt.Println(otp)

	//to do add mailer

	otpData := map[string]interface{}{
		"otp":    otp,
		"userId": user.ID,
	}

	otpDataJSON, _ := json.Marshal(otpData)

	_ = util.SetKey(context.Background(), fmt.Sprintf("otp:%s", otp), otpDataJSON, 300)

	createdUser.Password = ""

	return createdUser, http.StatusCreated, nil
}

func (s *UserService) CreateProfile(ctx context.Context, profile *model.Profile) (*model.Profile, int, error) {

	user, err := s.userRepo.FindByID(ctx, profile.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if user == nil {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}

	arg := map[string]any{
		"user_id": user.ID,
	}

	profileExists, err := s.profileRepo.Get(ctx, arg)

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if profileExists != nil {
		return nil, http.StatusConflict, fmt.Errorf("profile existed")
	}

	createdProfile, err := s.profileRepo.Create(ctx, profile)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create profile: %v", err)
	}

	return createdProfile, http.StatusCreated, nil
}

func (s *UserService) SignIn(ctx context.Context, user *model.User) (*string, int, map[string]interface{}, *bool, error) {
	arg := map[string]interface{}{
		"email": user.Email,
	}

	userExists, err := s.userRepo.Get(ctx, arg)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, nil, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return nil, http.StatusNotFound, nil, nil, fmt.Errorf("invalid email")
	}

	if userExists.Status != "activated" {
		return &userExists.ID, http.StatusUnauthorized, nil, nil, fmt.Errorf("please activate account to signin")
	}

	if err := util.CheckPassword(userExists.Password, user.Password); err != nil {
		return nil, http.StatusUnauthorized, nil, nil, fmt.Errorf("invalid password")
	}

	tokenMap, err := util.GenerateJWTToken(userExists.ID, string(userExists.Role))
	if err != nil {
		return nil, http.StatusInternalServerError, nil, nil, fmt.Errorf("try again later")
	}

	profilearg := map[string]interface{}{
		"user_id": userExists.ID,
	}

	profile, err := s.profileRepo.Get(ctx, profilearg)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, nil, fmt.Errorf("try again later")
	}

	profileExisted := profile != nil

	return &userExists.Email, http.StatusOK, map[string]interface{}{
		"accessToken": tokenMap["accessToken"],
		"expiresAt":   tokenMap["expiresAt"],
	}, &profileExisted, nil
}

func (s *UserService) GetProfile(ctx context.Context, user *model.User) (interface{}, int, error) {
	userExists, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}

	arg := map[string]any{
		"user_id": userExists.ID,
	}

	profile, err := s.profileRepo.Get(ctx, arg)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	data := map[string]interface{}{
		"user_id":   userExists.ID,
		"email":     userExists.Email,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"dob":       profile.DateOfBirth,
		"phoneno":   profile.PhoneNumber,
		"activated": userExists.Status,
	}

	return data, http.StatusOK, nil
}

func (s *UserService) ActivateUser(ctx context.Context, user *model.User, otp string) (int, error) {
	otpData, _ := util.GetKey(context.Background(), fmt.Sprintf("otp:%s", otp))

	if otpData == "" {
		return http.StatusForbidden, fmt.Errorf("otp invalid/expired")
	}

	var otpMap map[string]interface{}
	if err := json.Unmarshal([]byte(otpData), &otpMap); err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("otp invalid/expired")
	}

	otpValue := otpMap["otp"].(string)
	userIdValue := otpMap["userId"].(string)

	if otpValue != otp {
		return http.StatusForbidden, fmt.Errorf("otp invalid/expired")
	}

	if userIdValue != user.ID {
		return http.StatusForbidden, fmt.Errorf("otp invalid/expired")
	}

	userExists, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}

	userExists.Status = "activated"

	if err := s.userRepo.Update(ctx, userExists.ID, userExists); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to activate user")
	}

	return http.StatusOK, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, profile *model.Profile) (*string, int, error) {
	userExists, err := s.userRepo.FindByID(ctx, profile.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}

	arg := map[string]any{
		"user_id": userExists.ID,
	}

	profileExists, err := s.profileRepo.Get(ctx, arg)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if profileExists == nil {
		return nil, http.StatusConflict, fmt.Errorf("profile doesn't exist")
	}

	updatedData := make(map[string]interface{})

	if profile.LastName != "" {
		updatedData["last_name"] = profile.LastName
	}
	if profile.FirstName != "" {
		updatedData["first_name"] = profile.FirstName
	}
	if profile.PhoneNumber != "" {
		updatedData["phone_number"] = profile.PhoneNumber
	}

	if len(updatedData) == 0 {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("no valid fields to update")
	}

	if err := s.profileRepo.PartialUpdate(ctx, profileExists.ID, updatedData); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to update profile")
	}

	return &userExists.Email, http.StatusOK, nil
}

func (s *UserService) ResendOtp(ctx context.Context, userId string) (*string, int, error) {
	userExists, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}

	otp := util.RandomNumberString(6)

	fmt.Println(otp)

	//to do add mailer

	otpData := map[string]interface{}{
		"otp":    otp,
		"userId": userExists.ID,
	}

	otpDataJSON, _ := json.Marshal(otpData)

	_ = util.SetKey(context.Background(), fmt.Sprintf("otp:%s", otp), otpDataJSON, 300)

	return &userExists.Email, http.StatusOK, nil
}

func (s *UserService) ChangePassword(ctx context.Context, password string, userId string) (int, error) {
	userExists, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to hash password")
	}

	updatedData := map[string]interface{}{
		"password": hashedPassword,
	}

	if err := s.userRepo.PartialUpdate(ctx, userExists.ID, updatedData); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to update password")
	}

	return http.StatusOK, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, userId string) (int, error) {
	userExists, err := s.userRepo.FindByID(ctx, userId)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("try again later")
	}

	if userExists == nil {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}

	userExists.Status = "deactivated"

	if err := s.userRepo.Update(ctx, userExists.ID, userExists); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to deactivated user")
	}

	return http.StatusOK, nil
}
