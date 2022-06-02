package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
)

func TestNewClient(t *testing.T){
	c := NewClient()
	_, err := uuid.Parse(c.Id)
	assert.Nil(t, err)
	assert.Equal(t, "", c.Name)
	assert.Equal(t, uint64(0), c.Document)
	assert.Equal(t, "", c.Email)
	assert.Equal(t, uint64(0), c.Phone)
	assert.Equal(t, "", c.Password)
}

func TestValidateId(t *testing.T) {
	// test empty id
	c := Client{}
	err := c.ValidateId()
	assert.NotNil(t, err)
	assert.Equal(t, "id should not be empty", err.Error())
	// test valid id with dash
	c = Client{Id: "cf357e70-7dc9-4e73-8323-f9ae2be36f4a"}
	err = c.ValidateId()
	assert.Nil(t, err)
	// test valid id without dash
	c = Client{Id: "cf357e707dc94e738323f9ae2be36f4a"}
	err = c.ValidateId()
	assert.Nil(t, err)
	// test invalid uuid
	c = Client{Id: "cf357e70-7dc9-4e73-8323-f9ae2be36f"}
	err = c.ValidateId()
	assert.NotNil(t, err)
	assert.Equal(t, "id should be a valid uuid", err.Error())
	// test invalid uuid
	c = Client{Id: "cf357e707dc9-4e73-8323-f9ae2be36f"}
	err = c.ValidateId()
	assert.NotNil(t, err)
	assert.Equal(t, "id should be a valid uuid", err.Error())
}

func TestValidateName(t *testing.T){
	// test empty
	c := Client{}
	err := c.ValidateName()
	assert.NotNil(t, err)
	assert.Equal(t, "name should not be empty", err.Error())
	// test valid name
	c.Name = "test"
	err = c.ValidateName()
	assert.Nil(t, err)
}

func TestGetPhoneCountry(t *testing.T){
	// Brasil cell test
	c := Client{Phone: 5511999999999}
	p := c.GetPhoneCountry()
	assert.Equal(t, "BR", p)	
	// Brasil test
	c.Phone = 551199999999
	p = c.GetPhoneCountry()
	assert.Equal(t, "BR", p)
	// USA test
	c.Phone = 12129240446
	p = c.GetPhoneCountry()
	assert.Equal(t, "US", p)
	// Brasil Phone Error
	c.Phone = 559919899999
	p = c.GetPhoneCountry()
	assert.Equal(t, "", p)
	// Brasil Phone Error
	c.Phone = 99899999
	p = c.GetPhoneCountry()
	assert.Equal(t, "", p)
}





