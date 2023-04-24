package model

type RefreshToken struct {
	ID  uint   `json:"-"`
	UID uint   `json:"-"`
	SS  string `json:"refreshToken"`
}

type IDToken struct {
	SS string `json:"idToken"`
}

type TokenPair struct {
	IDToken
	RefreshToken
}
