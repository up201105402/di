package model

type RefreshToken struct {
	ID           uint   `json:"-"`
	UID          uint   `json:"-"`
	SignedString string `json:"refreshToken"`
}

type IDToken struct {
	SignedString string `json:"idToken"`
}

type TokenPair struct {
	IDToken      `json:"accessToken"`
	RefreshToken `json:"refreshToken"`
}
