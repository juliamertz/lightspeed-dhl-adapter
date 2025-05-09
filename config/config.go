package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Dhl struct {
	UserId    string `json:"userId"`
	ApiKey    string `json:"key"`
	AccountId string `json:"-"`
}

type Lightspeed struct {
	Key       string
	Secret    string
	Frontend  string
	Cluster   string
	ShopId    string
	ClusterId string
}

type CompanyInfo struct {
	Name         string
	Street       string
	City         string
	PostalCode   string
	CountryCode  string
	Number       string
	Addition     string
	Email        string
	PhoneNumber  string
	PersonalNote *string
}

type Options struct {
	DryRun          *bool
	Port            *int
	PollingInterval *int
	Environment     *string
	Debug           *bool
}

type Secrets struct {
	Dhl         Dhl
	Lightspeed  Lightspeed
	CompanyInfo CompanyInfo
	Options     *Options
}

func LoadSecrets(path string) (*Secrets, error) {
	var secrets Secrets
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}

	if secrets.Options.DryRun == nil {
		dryRun := false
		secrets.Options.DryRun = &dryRun
	}

	if secrets.Options.Port == nil {
		port := 8080
		secrets.Options.Port = &port
	}

	if secrets.Options.Debug == nil {
		debug := false
		secrets.Options.Debug = &debug
	}

	if secrets.Options.Environment == nil {
		env := "production"
		secrets.Options.Environment = &env
	} else if *secrets.Options.Environment != "production" && *secrets.Options.Environment != "development" {
		panic("Invalid environment specified in config.toml, must be either 'production' or 'development'")
	}

	if secrets.Options.PollingInterval == nil {
		interval := 15
		secrets.Options.PollingInterval = &interval
	}

	if secrets.Dhl.UserId == "" || secrets.Dhl.ApiKey == "" || secrets.Dhl.AccountId == "" {
		return nil, errors.New("DHL secrets are incomplete")
	}

	if secrets.Lightspeed.Key == "" || secrets.Lightspeed.Secret == "" || secrets.Lightspeed.Cluster == "" || secrets.Lightspeed.Frontend == "" || secrets.Lightspeed.ClusterId == "" || secrets.Lightspeed.ShopId == "" {
		return nil, errors.New("Lightspeed secrets are incomplete")
	}

	if secrets.CompanyInfo.Name == "" || secrets.CompanyInfo.Street == "" || secrets.CompanyInfo.City == "" || secrets.CompanyInfo.PostalCode == "" || secrets.CompanyInfo.CountryCode == "" || secrets.CompanyInfo.Number == "" || secrets.CompanyInfo.Addition == "" || secrets.CompanyInfo.Email == "" || secrets.CompanyInfo.PhoneNumber == "" {
		return nil, errors.New("Company info is incomplete")
	}
	return &secrets, nil
}
