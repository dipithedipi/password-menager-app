package models

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

type PasswordSet struct {
    Domain          string `json:"domain"`
    Username        string `json:"username"`
    Description     string `json:"description"`
    Password        string `json:"password"`
}

type PasswordRequestSearch struct {
    Domain          string `json:"domain"`
}

type PasswordRequestInfo struct {
    PasswordId      string `json:"passwordId"`
    Domain          string `json:"domain"`
}

type ArgonParams struct {
    Memory          uint32
    Iterations      uint32
    Parallelism     uint8
    SaltLength      uint32
    KeyLength       uint32
}