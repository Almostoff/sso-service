package model

import (
	"time"
)

type ClientUuid struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type Client struct {
	ClientUuid       string          `json:"client_uuid" db:"client_uuid"`
	ID               int64           `json:"-" db:"id"`
	Language         string          `json:"language"`
	Nickname         string          `json:"nickname" db:"nickname"`
	RegistrationDate string          `json:"registration_date" db:"registration_date"`
	LastActivity     string          `json:"last_activity" db:"last_activity"`
	LastLogin        string          `json:"last_login" db:"last_login"`
	Contacts         *ClientContacts `json:"contacts,omitempty"`
	AuthLevel        *AuthLevel      `json:"auth_level,omitempty"`
	Credential       *Credential     `json:"-"`
}

type ClientRepoTest struct {
	ClientUuid       string         `json:"client_uuid" db:"client_uuid"`
	ID               int64          `json:"-" db:"id"`
	Language         string         `json:"language"`
	Nickname         string         `json:"nickname" db:"nickname"`
	RegistrationDate string         `json:"registration_date" db:"registration_date"`
	LastActivity     string         `json:"last_activity" db:"last_activity"`
	LastLogin        string         `json:"last_login" db:"last_login"`
	Contacts         ClientContacts `json:"contacts,omitempty"`
	AuthLevel        AuthLevel      `json:"auth_level,omitempty"`
	Credential       Credential     `json:"-"`
}

type ClientDB struct {
	ClientUuid       string          `json:"client_uuid" db:"client_uuid"`
	ID               int64           `json:"-" db:"id"`
	Nickname         string          `json:"nickname" db:"nickname"`
	RegistrationDate string          `json:"registration_date" db:"registration_date"`
	LastActivity     string          `json:"last_activity" db:"last_activity"`
	LastLogin        string          `json:"last_login" db:"last_login"`
	Contacts         *ClientContacts `json:"contacts,omitempty" db:"contacts"`
	AuthLevel        *AuthLevel      `json:"auth_level,omitempty" db:"auth_level"`
	Credential       *Credential     `json:"-" db:"credential"`
}

type ClientContacts struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phone" db:"phone"`
	Tg         string `json:"tg" db:"tg"`
}

type AuthLevel struct {
	ClientUuid     string `json:"client_uuid" db:"client_uuid"`
	KYC            bool   `json:"kyc" db:"kyc"`
	Email          bool   `json:"email" db:"email"`
	Phone          bool   `json:"phone" db:"phone"`
	Totp           bool   `json:"totp" db:"totp"`
	StrongPassword bool   `json:"strong_password" db:"strong_password"`
	TG             bool   `json:"tg" db:"tg"`
	ResolvedIp     bool   `json:"resolved_ip" db:"resolved_ip"`
}

type AuthCode struct {
	ClientUuid  string    `json:"client_uuid" db:"client_uuid"`
	Type        string    `json:"type_code" db:"type_code"`
	CodeNeed    string    `json:"code_need" db:"code_need"`
	Date        time.Time `json:"create_time,omitempty" db:"create_time"`
	Destination string    `json:"destination" db:"destination"`
}

type Session struct {
	Id         int64     `json:"id" db:"id"`
	ClientUuid string    `json:"client_uuid" db:"client_uuid"`
	UA         string    `json:"ua" db:"ua"`
	IP         string    `json:"ip" db:"ip"`
	LoginTime  time.Time `json:"login_time" db:"login_time"`
	LogoutTime time.Time `json:"logout_time" db:"logout_time"`
	IsLogout   bool      `json:"is_logout" db:"is_logout"`
}

type Credential struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Password   string `json:"-" db:"password"`
	TotpSecret string `json:"-" db:"totp_secret"`
	TgId       int64  `json:"-" db:"tg_id"`
}

type HistoryPasswords struct {
	ID         int    `json:"id" db:"id"`
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Password   string `json:"-" db:"password"`
	ChangeTime string `json:"-" db:"change_time"`
}

type NicknameHistory struct {
	ClientUuid  string    `json:"client_uuid" db:"client_uuid"`
	OldNickname string    `json:"old_nickname" db:"old_nickname"`
	ChangeTime  time.Time `json:"change_time" db:"change_time"`
}

type Access struct {
	Access string `json:"Access"`
}

type ResponseConfirmKycModel struct {
	AuthToken string `json:"authToken"`
}
