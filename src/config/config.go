package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
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
	DryRun          bool
	Port            int
	PollingInterval int
}

type Secrets struct {
	Dhl         Dhl
	Lightspeed  Lightspeed
	CompanyInfo CompanyInfo
	Options     Options
}

func LoadSecrets(path string) Secrets {
	logger := log.With().Str("secrets_path", path).Logger()

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to read file")
	}

	var secrets Secrets
	err = toml.Unmarshal(data, &secrets)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to unmarshal config")
	}

	if secrets.Dhl.UserId == "" || secrets.Dhl.ApiKey == "" || secrets.Dhl.AccountId == "" {
		logger.Fatal().Msg("DHL secrets are incomplete")
	}

	if secrets.Lightspeed.Key == "" || secrets.Lightspeed.Secret == "" || secrets.Lightspeed.Cluster == "" || secrets.Lightspeed.Frontend == "" || secrets.Lightspeed.ClusterId == "" || secrets.Lightspeed.ShopId == "" {
		logger.Fatal().Msg("Lightspeed secrets are incomplete")
	}

	if secrets.CompanyInfo.Name == "" || secrets.CompanyInfo.Street == "" || secrets.CompanyInfo.City == "" || secrets.CompanyInfo.PostalCode == "" || secrets.CompanyInfo.CountryCode == "" || secrets.CompanyInfo.Number == "" || secrets.CompanyInfo.Addition == "" || secrets.CompanyInfo.Email == "" || secrets.CompanyInfo.PhoneNumber == "" {
		logger.Fatal().Msg("Company info is incomplete")
	}

	return secrets
}
