package models

import (
    "github.com/golang-jwt/jwt"
)

type UserRegister struct {
    Username        string `json:"username"`
    Email           string `json:"email"`
    PasswordHash    string `json:"password"`
    Salt            string `json:"salt"`
    PublicKey       string `json:"public_key"`
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
    Otp             bool   `json:"otp"`
}

type PasswordRequestSearch struct {
    Domain          string `json:"domain"`
}

type PasswordRequestInfo struct {
    PasswordId      string `json:"passwordId"`
    Domain          string `json:"domain"`
    Otp             string `json:"otp"`
}

type PasswordLeakCheck struct {
    PasswordPartialHash        string `json:"password"`
}

type ArgonParams struct {
    Memory          uint32
    Iterations      uint32
    Parallelism     uint8
    SaltLength      uint32
    KeyLength       uint32
}

type CustomJWTClaims struct {
	jwt.StandardClaims
	Ip string
}
