package models

import (
	"github.com/golang-jwt/jwt"
)

type UserRegister struct {
    Username        string `json:"username"`
    Email           string `json:"email"`
    PasswordHash    string `json:"password"`
    Salt            string `json:"salt"`
    PublicKey       string `json:"publicKey"`
}

type UserVerify struct {
    Email           string `json:"email"`
    Otp             string `json:"otp"`
}

type UserRegisterUsernameCheck struct {
    Username        string `json:"username"`
}

type UserLogin struct {
    Email           string `json:"email"`
    PasswordHash    string `json:"password"`
    Otp             string `json:"otp"`
}

type UserSaltLogin struct {
    Email           string `json:"email"`
}

type PasswordSet struct {
    Domain          string `json:"domain"`
    Username        string `json:"username"`
    Description     string `json:"description"`
    Password        string `json:"password"`
    Category        string `json:"category"`
    Otp             bool   `json:"otp"`
}

type PasswordRequestSearch struct {
    Domain          string `json:"domain"`
    Category        []string `json:"category"`
}

type PasswordRequestInfo struct {
    PasswordId      string `json:"passwordId"`
    Domain          string `json:"domain"`
    Otp             string `json:"otp"`
}

type PasswordDelete struct {
    PasswordId      string `json:"passwordId"`
    Otp             string `json:"otp"`
}

type PasswordUpdate struct {
    PasswordId      string `json:"passwordId"`
    NewPassword     string `json:"newPassword"`
    NewUsername     string `json:"NewUsername"`
    NewDescription  string `json:"NewDescription"`
    NewCategory     string `json:"newCategory"`
    Otp             string `json:"otp"`
    OtpProtected    bool   `json:"otpProtected"`
}

type PasswordLeakCheck struct {
    PasswordPartialHash        string `json:"password"`
}

type EventRequest struct {
    StartDateTime   string `json:"start"`
    EndDateTime     string `json:"end"`
}

type CategoryCreate struct {
    Name            string `json:"name"`
}

type CategoryUpdate struct {
    Name      string `json:"name"`
    NewName   string `json:"newName"`
}

type CategoryDelete struct {
    Name      string `json:"name"`
    Otp       string `json:"otp"`
}

type ArgonParams struct {
    Memory          uint32
    Iterations      uint32
    Parallelism     uint8
    SaltLength      uint32
    KeyLength       uint32
}

type SessionModelResponse struct {
    DatabaseElemID string
    IpAddress   string
    CreatedAt   string
    LastUse     string
    UserAgent   string
    CurrentUser bool
}

type SessionDeleteRequest struct {
    DatabaseElemID string `json:"id"`
    Otp           string  `json:"otp"`
}

type CustomJWTClaims struct {
	jwt.StandardClaims
	Ip string
}
