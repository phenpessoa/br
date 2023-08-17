package br

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/phenpessoa/gutils/cache"
	"github.com/phenpessoa/gutils/unsafex"
)

var (
	// ErrInvalidUF is returned when an invalid state (UF) code is passed.
	ErrInvalidUF = errors.New("br: invalid uf passed")

	// ErrInvalidCEP is returned when an invalid CEP code is passed.
	ErrInvalidCEP = errors.New("br: invalid cep passed")

	// ErrInvalidSerializedAddress is returned when trying to deserialize an
	// invalid string into an Address.
	ErrInvalidSerializedAddress = errors.New(
		"br: invalid serialized address",
	)
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
	str := unsafex.String(b)
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

// CEP represents a Brazilian postal code (Código de Endereçamento Postal).
type CEP string

// NewCEP creates a CEP instance from a string representation and removes any
// punctuation characters.
//
// It verifies the length of the CEP, the digits, and checks the cache for known
// invalid CEPs. It does not make requests to APIs to validate the corresponding
// address.
//
// To check if this CEP actually represents an Address, call the CEP.ToAddress
// method.
func NewCEP(s string) (CEP, error) {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ".", "")
	cep := CEP(s)
	if !cep.IsValid() {
		return "", ErrInvalidCEP
	}
	return cep, nil
}

// String returns the formatted CEP string as XXXXX-XXX.
func (cep CEP) String() string {
	if !cep.IsValid() {
		return ""
	}

	out := make([]byte, 9)
	for i := range cep {
		switch {
		case i < 5:
			out[i] = cep[i]
		case i == 5:
			out[i] = '-'
			out[i+1] = cep[i]
		default:
			out[i+1] = cep[i]
		}
	}

	return unsafex.String(out)
}

// IsValid checks whether the provided CEP is valid based on its length, digits
// and check the caches for known invalid CEPs.
//
// It does not make requests to APIs to validate the corresponding address.
//
// To check if this CEP actually represents an Address, call the CEP.ToAddress
// method.
func (cep CEP) IsValid() bool {
	if len(cep) != 8 {
		return false
	}

	for i := range cep {
		if cep[i] < '0' || cep[i] > '9' {
			return false
		}
	}

	return !invalidCEPs.Contains(string(cep))
}

// ToAddress converts a CEP into an Address instance, retrieving address
// information associated with the CEP.
//
// This method may perform requests to external APIs to fetch address details
// based on the CEP.
func (cep CEP) ToAddress() (addr Address, err error) {
	if !cep.IsValid() {
		return Address{}, ErrInvalidCEP
	}

	if addr, ok := addresses.Get(string(cep)); ok {
		return addr, nil
	}

	defer func() {
		if err != nil && errors.Is(err, ErrInvalidCEP) {
			invalidCEPs.Set(string(cep), unit{})
		}
	}()

	addr, err = buscaCEP(cep)
	if err != nil {
		if errors.Is(err, ErrInvalidCEP) {
			return Address{}, ErrInvalidCEP
		}

		var err2 error
		addr, err2 = viaCEP(cep)
		if err2 != nil {
			if errors.Is(err2, ErrInvalidCEP) {
				return Address{}, ErrInvalidCEP
			}

			return Address{}, fmt.Errorf(
				"failed to validate cep: %w",
				errors.Join(err, err2),
			)
		}
	}

	addresses.Set(string(cep), addr)
	return addr, nil
}

func buscaCEP(cep CEP) (Address, error) {
	type buscaCEPResponse struct {
		Erro     bool      `json:"erro"`
		Mensagem string    `json:"mensagem"`
		Total    int       `json:"total"`
		Dados    []Address `json:"dados"`
	}

	const buscaCEPURL = "https://buscacepinter.correios.com.br/app/endereco/carrega-cep-endereco.php"
	form := url.Values{}
	form.Set("pagina", "/app/endereco/index.php")
	form.Set("endereco", string(cep))
	form.Set("tipoCEP", "ALL")
	res, err := http.Post(
		buscaCEPURL,
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return Address{}, fmt.Errorf(
			"failed to make post request to buscaCEP: %w", err,
		)
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return Address{}, fmt.Errorf(
			"buscaCEP status code not ok: %d\ndata: %s",
			res.StatusCode, string(data),
		)
	}

	var dados buscaCEPResponse
	if err := json.Unmarshal(data, &dados); err != nil {
		return Address{}, fmt.Errorf(
			"failed to unmarshal buscaCEP json: %w\ndata: %s",
			err, string(data),
		)
	}

	if dados.Erro {
		return Address{}, fmt.Errorf(
			"busca cep server returned error: %s", dados.Mensagem,
		)
	}

	if len(dados.Dados) == 0 {
		return Address{}, ErrInvalidCEP
	}

	addr := dados.Dados[0]
	logradouro := addr.Logradouro

	logradouroParts := strings.Split(logradouro, ",")
	if len(logradouroParts) == 1 {
		logradouroParts = strings.Split(logradouro, "-")
	}

	var complemento string
	if len(logradouroParts) > 1 {
		logradouro = strings.TrimSpace(logradouroParts[0])
		complemento = strings.TrimSpace(logradouroParts[1])
	}

	addr.Logradouro = logradouro
	addr.Complemento = complemento
	addr.CEP = cep
	return addr, nil
}

func viaCEP(cep CEP) (Address, error) {
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	res, err := http.Get(viaCEPURL)
	if err != nil {
		return Address{}, fmt.Errorf(
			"failed to make get request to viaCEP: %w", err,
		)
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return Address{}, fmt.Errorf(
			"viaCEP status code not ok: %d\ndata: %s",
			res.StatusCode, string(data),
		)
	}

	type viaCEPResponse struct {
		Address
		Erro bool `json:"erro"`
	}

	var dados viaCEPResponse
	if err := json.Unmarshal(data, &dados); err != nil {
		return Address{}, fmt.Errorf(
			"failed to unmarshal viaCEP json: %w\ndata: %s", err, string(data),
		)
	}

	if dados.Erro || dados.Address.CEP == "" {
		return Address{}, ErrInvalidCEP
	}

	return dados.Address, nil
}

type unit struct{}

var (
	addresses   *cache.Cache[string, Address]
	invalidCEPs *cache.Cache[string, unit]
)

const (
	cacheAddressesAndCEPsFor = 2 * time.Hour
)

func init() {
	addresses = cache.New[string, Address](cacheAddressesAndCEPsFor)
	invalidCEPs = cache.New[string, unit](cacheAddressesAndCEPsFor)
}

// Address represents an address associated with a Brazilian CEP.
type Address struct {
	UF          UF     `json:"uf"`
	CEP         CEP    `json:"cep"`
	Localidade  string `json:"localidade"`
	Logradouro  string `json:"logradouroDNEC"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	NomeUnidade string `json:"nomeUnidade"`
}

func (addr Address) serializedSize() int {
	return 2 + 1 + // uf
		len(addr.Localidade) + 1 +
		len(addr.Logradouro) + 1 +
		len(addr.Complemento) + 1 +
		len(addr.Bairro) + 1 +
		len(addr.NomeUnidade) + 1 +
		len(addr.CEP)
}

// Serialize converts the Address instance into a serialized string
// representation.
//
// This can be used to store the address as a string on a database, for example.
func (addr Address) Serialize() string {
	var buf bytes.Buffer
	buf.Grow(addr.serializedSize())
	buf.WriteString(addr.UF.String())
	buf.WriteRune(';')
	buf.WriteString(addr.Localidade)
	buf.WriteRune(';')
	buf.WriteString(addr.Logradouro)
	buf.WriteRune(';')
	buf.WriteString(addr.Complemento)
	buf.WriteRune(';')
	buf.WriteString(addr.Bairro)
	buf.WriteRune(';')
	buf.WriteString(addr.NomeUnidade)
	buf.WriteRune(';')
	buf.WriteString(string(addr.CEP))
	return buf.String()
}

// Deserialize parses a serialized string and populates the Address fields.
func (addr *Address) Deserialize(str string) error {
	parts := strings.Split(str, ";")
	if len(parts) != 7 {
		return fmt.Errorf(
			"%w: invalid length: %d",
			ErrInvalidSerializedAddress, len(parts),
		)
	}

	uf, err := NewUFFromStr(parts[0])
	if err != nil {
		return fmt.Errorf(
			"%w: unknown uf: %s",
			ErrInvalidSerializedAddress, parts[0],
		)
	}

	cep, err := NewCEP(parts[6])
	if err != nil {
		return fmt.Errorf(
			"%w: invalid CEP: %s",
			ErrInvalidSerializedAddress, parts[6],
		)
	}

	addr.UF = uf
	addr.Localidade = parts[1]
	addr.Logradouro = parts[2]
	addr.Complemento = parts[3]
	addr.Bairro = parts[4]
	addr.NomeUnidade = parts[5]
	addr.CEP = cep
	return nil
}

// Scan implements the sql.Scanner interface for Address.
func (addr *Address) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf(
			"br: unknown type passed to Address Scan: %T",
			value,
		)
	}

	if err := addr.Deserialize(str); err != nil {
		return fmt.Errorf("br: invalid serialized Address: %w", err)
	}

	return nil
}

// Value implements the driver.Valuer interface for Address.
func (addr Address) Value() (driver.Value, error) {
	return addr.Serialize(), nil
}

// Validate fetches additional address information based on the associated CEP
// and updates the Address fields.
//
// This method may perform requests to external APIs to retrieve address
// details.
func (addr *Address) Validate() error {
	parsedAddr, err := addr.CEP.ToAddress()
	if err != nil {
		return err
	}
	*addr = parsedAddr
	return nil
}
