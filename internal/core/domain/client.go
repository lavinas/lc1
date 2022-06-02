package domain

import (
	"errors"
	"strconv"
	"math"
	"regexp"
	"strings"
	"net"
	
	"github.com/google/uuid"
	"github.com/dongri/phonenumber"
)

const (
	cpf_min_length        = 8
	cpf_max_length        = 12
	cnpj_min_length       = 12
	cnpj_max_length       = 16
	email_regex           = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]" +
	"{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	default_county        = "BR"
)

type Client struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Document uint64 `json:"document"`
	Email string `json:"email"`
	Phone uint64 `json:"phone"`
	Password string `json:"password"`
}

func NewClient() *Client {
	id := uuid.NewString()
	return &Client{Id: id}
}

func (c *Client) ValidateId() error {
	if c.Id == "" {
		return errors.New("id should not be empty")
	}
	if _, err := uuid.Parse(c.Id); err != nil {
		return errors.New("id should be a valid uuid")
	}
	return nil
}

func (c *Client) ValidateName() error {
	if c.Name == "" {
		return errors.New("name should not be empty")
	}
	return nil
}

func (c *Client) ValidateDocument() error {
	// validate if it is not empty
	if c.Document == uint64(0) {
		return errors.New("document should no be empty")
	}
	// validate if it is valid CPF ou CNPJ
	if !c.IsDocumentCPF() && !c.IsDocumentCNPJ() {
		return errors.New("document should have a CPF or CNPJ number")
	}
	return nil
}

func (c *Client) ValidateEmail() error {
	// validate if it is not empty
	if c.Email == "" {
		return errors.New("email should not be empty")
	}
	// validate structure
	var emailRegex = regexp.MustCompile(email_regex)
	if !emailRegex.MatchString(c.Email) {
		return errors.New("email should have a valid email address format")
	}
	parts := strings.Split(c.Email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return errors.New("email should have a valid email address format")
	}
	return nil
}

// Validate phone number or/and country
func (c *Client) ValidatePhone() error {
	// check non nil number
	if c.Phone == uint64(0) {
		return errors.New("phone should not be empty")
	}
	country := c.GetPhoneCountry()
	if country == "" {
		return errors.New("phone should have defined country code")
	}
	sn := phonenumber.Parse(strconv.FormatUint(c.Phone, 10), country)
	_, err := strconv.ParseUint(sn, 10, 64)
	if err != nil {
		return errors.New("phone should have a valid number")
	}
	return nil
}

func (c *Client) IsDocumentCPF() bool {
	if c.Document == 0 {
		return false
	}
	// valid length
	len := len(strconv.FormatUint(c.Document, 10))
	if len < cpf_min_length || len > cpf_max_length {
		return false
	}
	// valid check digits (2 last digits)
	dig1 := int(c.Document % 100 / 10)
	dig2 := int(c.Document % 10)
	val1 := 0
	val2 := 0
	for i := 3; i <= len; i++ {
		x := int(math.Mod(float64(c.Document), math.Pow10(i)) / math.Pow10(i-1))
		val1 += x * (i - 1)
		val2 += x * i
	}
	val2 += dig1 * 2
	val1 = int(math.Mod(float64(val1*10), float64(11)))
	val2 = int(math.Mod(float64(val2*10), float64(11)))
	if val1 != dig1 || val2 != dig2 {
		return false
	}
	// Ok
	return true
}

func (c *Client) IsDocumentCNPJ() bool {
	// valid is not zero
	if c.Document == 0 {
		return false
	}
	// valid length
	len := len(strconv.FormatUint(c.Document, 10))
	if len < cnpj_min_length || len > cnpj_max_length {
		return false
	}
	// valid check digits (2 last digits)
	dig1 := int(c.Document % 100 / 10)
	dig2 := int(c.Document % 10)
	val1 := 0
	val2 := 0
	for i := 0; i <= len-3; i++ {
		x := int(math.Mod(float64(c.Document), math.Pow10(i+3)) / math.Pow10(i+2))
		val1 += x * (int(math.Mod(float64(i), float64(8))) + 2)
		val2 += x * (int(math.Mod(float64(i+1), float64(8))) + 2)
	}
	val2 += dig1 * 2
	val1 = int(math.Mod(float64(val1), float64(11)))
	val2 = int(math.Mod(float64(val2), float64(11)))
	if val1 < 2 {
		val1 = 0
	}
	if val2 < 2 {
		val2 = 0
	}
	val1 = 11 - val1
	val2 = 11 - val2
	if val1 != dig1 || val2 != dig2 {
		return false
	}
	// Ok
	return true
}

func (c *Client) GetPhoneCountry() (string) {
	if c.Phone == 0 {
		return ""
	}
	iso := phonenumber.GetISO3166ByNumber(strconv.FormatUint(c.Phone, 10), false)
	return iso.Alpha2
}
