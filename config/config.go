package config

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	path string

	FishPi   *FishPi   `yaml:"fishPi"`
	Settings *Settings `yaml:"settings"`
	Ice      *Ice      `yaml:"ice"`
	Elves    *Elves    `yaml:"elves"`
}

type FishPi struct {
	ApiBase     string `yaml:"apiBase"`
	UserAgent   string `yaml:"userAgent"`
	ApiKey      string `yaml:"apiKey"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	PasswordMd5 string `yaml:"passwordMd5"`
	MfaCode     string `yaml:"mfaCode"`
}

type Settings struct {
	WsInterval  int `yaml:"wsInterval"`
	MsgCacheNum int `yaml:"msgCacheNum"`
}

type Ice struct {
	Url      string `yaml:"url"`
	Ck       string `yaml:"ck"`
	Username string `yaml:"username"`
	Uid      string `yaml:"uid"`
}

type Elves struct {
	Token string `yaml:"token"`
}

func NewConfig(path string) (*Config, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err = yaml.Unmarshal(body, &c); err != nil {
		return nil, err
	}
	c.FishPi.Init()
	c.path = path

	return &c, err
}

func (c *Config) UpdateApiKey(apiKey string) error {
	c.FishPi.ApiKey = apiKey

	return c.save()
}

func (c *Config) UpdateCK(ck string) error {
	c.Ice.Ck = ck

	return c.save()
}

func (c *Config) save() error {
	body, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path, body, os.ModePerm)
}

func (f *FishPi) Init() {
	if f.PasswordMd5 != "" {
		return
	}
	f.PasswordMd5 = f.md5(f.Password)
}

func (f *FishPi) md5(source string) string {
	newSig := md5.Sum([]byte(source))
	newArr := fmt.Sprintf("%x", newSig)
	return strings.ToLower(newArr)
}
