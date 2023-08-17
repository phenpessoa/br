package br

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

var (
	// ErrInvalidUF is returned when an invalid state (UF) code is passed.
	ErrInvalidUF = errors.New("br: invalid uf passed")
)

// UF stands for Unidade Federativa and represents a Brazilian state.
type UF uint8

// NewUF creates a UF instance from a given state code.
func NewUF(codigo int) (UF, error) {
	if codigo < 11 || codigo > 53 {
		return 0, ErrInvalidUF
	}

	uf := UF(uint8(codigo))
	if uf.String() == "" {
		return 0, ErrInvalidUF
	}

	return uf, nil
}

// NewUFFromStr creates a UF instance from a string representation of the
// state's name or abbreviation.
func NewUFFromStr(uf string) (UF, error) {
	if len(uf) < 2 || len(uf) > 19 {
		return 0, ErrInvalidUF
	}

	switch strings.ToLower(uf) {
	case "ro", "rondônia", "rondonia":
		return UF(11), nil
	case "ac", "acre":
		return UF(12), nil
	case "am", "amazonas":
		return UF(13), nil
	case "rr", "roraima":
		return UF(14), nil
	case "pa", "pará", "para":
		return UF(15), nil
	case "ap", "amapá", "amapa":
		return UF(16), nil
	case "to", "tocantins":
		return UF(17), nil
	case "ma", "maranhão", "maranhao":
		return UF(21), nil
	case "pi", "piauí", "piaui":
		return UF(22), nil
	case "ce", "ceará", "ceara":
		return UF(23), nil
	case "rn", "rio grande do norte", "riograndedonorte":
		return UF(24), nil
	case "pb", "paraíba", "paraiba":
		return UF(25), nil
	case "pe", "pernambuco":
		return UF(26), nil
	case "al", "alagoas":
		return UF(27), nil
	case "se", "sergipe":
		return UF(28), nil
	case "ba", "bahia":
		return UF(29), nil
	case "mg", "minas gerais", "minasgerais":
		return UF(31), nil
	case "es", "espírito santo", "espiritosanto":
		return UF(32), nil
	case "rj", "rio de janeiro", "riodejaneiro":
		return UF(33), nil
	case "sp", "são paulo", "sao paulo", "sãopaulo", "saopaulo":
		return UF(35), nil
	case "pr", "paraná", "parana":
		return UF(41), nil
	case "sc", "santa catarina", "santacatarina":
		return UF(42), nil
	case "rs", "rio grande do sul", "riograndedosul":
		return UF(43), nil
	case "ms", "mato grosso do sul", "matogrossodosul":
		return UF(50), nil
	case "mt", "mato grosso", "matogrosso":
		return UF(51), nil
	case "go", "goiás", "goias":
		return UF(52), nil
	case "df", "distrito federal", "distritofederal":
		return UF(53), nil
	default:
		return 0, ErrInvalidUF
	}
}

// String returns the abbreviation of the UF, such as RJ and SP.
func (uf UF) String() string {
	switch uf {
	case 11:
		return "RO"
	case 12:
		return "AC"
	case 13:
		return "AM"
	case 14:
		return "RR"
	case 15:
		return "PA"
	case 16:
		return "AP"
	case 17:
		return "TO"
	case 21:
		return "MA"
	case 22:
		return "PI"
	case 23:
		return "CE"
	case 24:
		return "RN"
	case 25:
		return "PB"
	case 26:
		return "PE"
	case 27:
		return "AL"
	case 28:
		return "SE"
	case 29:
		return "BA"
	case 31:
		return "MG"
	case 32:
		return "ES"
	case 33:
		return "RJ"
	case 35:
		return "SP"
	case 41:
		return "PR"
	case 42:
		return "SC"
	case 43:
		return "RS"
	case 50:
		return "MS"
	case 51:
		return "MT"
	case 52:
		return "GO"
	case 53:
		return "DF"
	default:
		return ""
	}
}

// UnidadeFederativa returns the full name of the UF, such as Rio de Janeiro.
func (uf UF) UnidadeFederativa() string {
	switch uf {
	case 11:
		return "Rondônia"
	case 12:
		return "Acre"
	case 13:
		return "Amazonas"
	case 14:
		return "Roraima"
	case 15:
		return "Pará"
	case 16:
		return "Amapá"
	case 17:
		return "Tocantins"
	case 21:
		return "Maranhão"
	case 22:
		return "Piauí"
	case 23:
		return "Ceará"
	case 24:
		return "Rio Grande do Norte"
	case 25:
		return "Paraíba"
	case 26:
		return "Pernambuco"
	case 27:
		return "Alagoas"
	case 28:
		return "Sergipe"
	case 29:
		return "Bahia"
	case 31:
		return "Minas Gerais"
	case 32:
		return "Espírito Santo"
	case 33:
		return "Rio de Janeiro"
	case 35:
		return "São Paulo"
	case 41:
		return "Paraná"
	case 42:
		return "Santa Catarina"
	case 43:
		return "Rio Grande do Sul"
	case 50:
		return "Mato Grosso do Sul"
	case 51:
		return "Mato Grosso"
	case 52:
		return "Goiás"
	case 53:
		return "Distrito Federal"
	default:
		return ""
	}
}

// Codigo returns the numeric code of the UF.
func (uf UF) Codigo() int {
	return int(uf)
}

// UnmarshalJSON implements the json.Unmarshaler interface for UF.
func (uf *UF) UnmarshalJSON(b []byte) error {
	str := unsafe.String(unsafe.SliceData(b), len(b))
	if strings.Contains(str, `"`) {
		str = strings.ReplaceAll(str, `"`, "")
		_uf, err := NewUFFromStr(str)
		if err != nil {
			return fmt.Errorf("can not unmarshal %s into uf", str)
		}
		*uf = _uf
		return nil
	}

	parsedCode, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("can not unmarshal %s into uf", str)
	}

	_uf, err := NewUF(int(parsedCode))
	if err != nil {
		return fmt.Errorf("can not unmarshal %d into uf", parsedCode)
	}

	*uf = _uf
	return nil
}

// MarshalJSON implements the json.Marshaler interface for UF.
func (uf UF) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uf.String() + `"`), nil
}

// Scan implements the sql.Scanner interface for UF.
func (uf *UF) Scan(value any) error {
	i64, ok := value.(int64)
	if !ok {
		return fmt.Errorf("br: unknown type passed to UF Scan: %T", value)
	}

	_uf, err := NewUF(int(i64))
	if err != nil {
		return fmt.Errorf("br: can not convert %d to UF", i64)
	}

	*uf = _uf
	return nil
}

// Value implements the driver.Valuer interface for UF.
func (uf UF) Value() (driver.Value, error) {
	return int64(uf), nil
}
