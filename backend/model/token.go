package model

type RefreshToken struct {
	ID           uint   `json:"-"`
	UID          uint   `json:"-"`
	SignedString string `json:"signedString"`
}

type IDToken struct {
	SignedString string `json:"signedString"`
}

type TokenPair struct {
	IDToken      `json:"accessToken"`
	RefreshToken `json:"refreshToken"`
}
