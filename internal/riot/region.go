package riot

import "errors"

type Region uint8

var (
	ErrRegionUnknown = errors.New("unknown region")
)

const (
	// NA1 is the North American server
	NA1 Region = iota
	// EUW1 is the EU West server
	EUW1
	// EUN1 is the EU Nordic & East server
	EUN1
	// BR1 is the Brazilian server
	BR1
	// JP1 is the Japanese server
	JP1
	// KR is the Korean server
	KR
	// LA1 is the Northern Latin American server
	LA1
	// LA2 is the Southern Latin American server
	LA2
	// OC1 is the Oceania server
	OC1
	// RU is the Russian server
	RU
	// TR1 is the Turkish server
	TR1
)

var realms = map[Region]string{
	BR1:  "br1",
	EUN1: "eun1",
	EUW1: "euw1",
	JP1:  "jp1",
	KR:   "kr",
	LA1:  "la1",
	LA2:  "la2",
	NA1:  "na1",
	OC1:  "oc1",
	RU:   "ru",
	TR1:  "tr1",
}

var continents = map[Region]string{
	BR1:  "americas",
	EUN1: "europe",
	EUW1: "europe",
	JP1:  "asia",
	KR:   "asia",
	LA1:  "americas",
	LA2:  "americas",
	NA1:  "americas",
	OC1:  "asia",
	RU:   "europe",
	TR1:  "europe",
}

func (r Region) Realm() (string, bool) {
	realm, isValid := realms[r]
	return realm, isValid
}

func (r Region) Continent() (string, bool) {
	continent, isValid := continents[r]
	return continent, isValid
}
