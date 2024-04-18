package models

type UserRegister struct {
    Username        string `json:"username"`
    Email           string `json:"email"`
    PasswordHash    string `json:"password"`
    Salt            string `json:"salt"`
}

type UserLogin struct {
    Email           string `json:"email"`
    PasswordHash    string `json:"password"`
    Otp             string `json:"otp"`
}

type ArgonParams struct {
    Memory          uint32
    Iterations      uint32
    Parallelism     uint8
    SaltLength      uint32
    KeyLength       uint32
}