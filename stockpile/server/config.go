/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/hashicorp/hcl2/gohcl"
  "github.com/hashicorp/hcl2/hcl"
  "github.com/hashicorp/hcl2/hclparse"
)

// defines the default address to listen on when none is given
const DefaultAddress = "0.0.0.0"

// defines the default port to listen on when none is given
const DefaultPort = 36623

// Represents a server configuration (typically parsed from one or more HCL files)
type Config struct {
  BindAddress     string         `hcl:"bind-address,attr"`
  EnableLegacyApi bool           `hcl:"legacy-api,attr"`
  Storage         *StorageConfig `hcl:"storage,block"`
  Ttl             *TtlConfig     `hcl:"ttl,block"`
}

// Represents a storage backend configuration
// The "type" parameter will simply equal the executable name within the plugin directory while all parameters are
// passed on upon startup
type StorageConfig struct {
  Type       string   `hcl:"type,label"`
  Parameters hcl.Body `hcl:",remain"`
}

// Represents the TTL (Time To Live) configuration (e.g. caching durations for various value types)
type TtlConfig struct {
  Name           time.Duration
  RawName        string `hcl:"name,attr"`
  NameHistory    time.Duration
  RawNameHistory string `hcl:"name-history,attr"`
  Profile        time.Duration
  RawProfile     string `hcl:"profile,attr"`
  Blacklist      time.Duration
  RawBlacklist   string `hcl:"blacklist,attr"`
}

// Creates an empty configuration
func EmptyConfig() *Config {
  return &Config{}
}

func DefaultConfig() *Config {
  cfg := &Config{
    BindAddress:     fmt.Sprintf("%s:%d", DefaultAddress, DefaultPort),
    EnableLegacyApi: false,
    Storage: &StorageConfig{
      Type: "mem",
    },
    Ttl: &TtlConfig{
      Name:        mojang.NameValidityPeriod,            // Full Mojang limit
      NameHistory: mojang.NameChangeRateLimitPeriod / 4, // 1/4th of the Mojang limit
      Profile:     time.Hour * 24 * 7,                   // 7 days
      Blacklist:   time.Hour * 24 * 7,                   // 7 days
    },
  }

  // since parse may be called on this config we'll have to copy the string representations as well
  ttl := cfg.Ttl
  ttl.RawName = ttl.Name.String()
  ttl.RawNameHistory = ttl.NameHistory.String()
  ttl.RawProfile = ttl.Profile.String()
  ttl.RawBlacklist = ttl.Blacklist.String()

  return cfg
}

func DevelopmentConfig() *Config {
  return DefaultConfig().Merge(&Config{
    BindAddress:     fmt.Sprintf("%s:%d", "127.0.0.1", DefaultPort),
    EnableLegacyApi: true,
  })
}

// Loads a file or directory
func LoadConfig(path string) (*Config, error) {
  file, err := os.Stat(path)
  if err != nil {
    return nil, err
  }

  if file.IsDir() {
    return LoadConfigDirectory(path)
  }

  base := DefaultConfig()
  cfg, err := LoadConfigFile(path)
  if err != nil {
    return nil, err
  }
  base.Merge(cfg)
  return base, nil
}

// Loads an entire directory of configuration files
func LoadConfigDirectory(path string) (*Config, error) {
  files, err := ioutil.ReadDir(path)
  if err != nil {
    return nil, err
  }

  base := DefaultConfig()
  for _, file := range files {
    if file.IsDir() {
      continue
    }

    if strings.HasSuffix(file.Name(), ".hcl") || strings.HasSuffix(file.Name(), ".json") {
      cfg, err := LoadConfigFile(filepath.Join(path, file.Name()))
      if err != nil {
        return nil, err
      }
      base.Merge(cfg)
    }
  }

  return base, nil
}

// Loads a single configuration file
func LoadConfigFile(path string) (*Config, error) {
  parser := hclparse.NewParser()
  file, diag := parser.ParseHCLFile(path)

  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to load configuration file \"%s\": %s", path, diag.Error())
  }

  cfg := EmptyConfig()
  diag = gohcl.DecodeBody(file.Body, nil, cfg)
  cfg.Parse()

  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to load configuration file \"%s\": %s", path, diag.Error())
  }
  return cfg, nil
}

// Merges two configuration instances into one
func (c *Config) Merge(other *Config) *Config {
  if other.BindAddress != "" {
    c.BindAddress = other.BindAddress
  }

  if other.EnableLegacyApi {
    c.EnableLegacyApi = true
  }

  if c.Storage == nil {
    c.Storage = other.Storage
  } else if other.Storage != nil {
    c.Storage.Merge(other.Storage)
  }

  if c.Ttl == nil {
    c.Ttl = other.Ttl
  } else if other.Ttl != nil {
    c.Ttl.Merge(other.Ttl)
  }

  return c
}

func (c *StorageConfig) Merge(other *StorageConfig) *StorageConfig {
  if other.Type != "" {
    c.Type = other.Type
    c.Parameters = other.Parameters
  }
  return c
}

func (c *TtlConfig) Merge(other *TtlConfig) *TtlConfig {
  if other.Name != 0 {
    c.Name = other.Name
  }
  if other.NameHistory != 0 {
    c.NameHistory = other.NameHistory
  }
  if c.Profile != 0 {
    c.Profile = other.Profile
  }
  if other.Blacklist != 0 {
    c.Blacklist = other.Blacklist
  }
  return c
}

func (c *TtlConfig) Parse() error {
  name, err := time.ParseDuration(c.RawName)
  if err != nil {
    return err
  }

  nameHistory, err := time.ParseDuration(c.RawNameHistory)
  if err != nil {
    return err
  }

  profile, err := time.ParseDuration(c.RawProfile)
  if err != nil {
    return err
  }

  blacklist, err := time.ParseDuration(c.RawBlacklist)
  if err != nil {
    return err
  }

  c.Name = name
  c.NameHistory = nameHistory
  c.Profile = profile
  c.Blacklist = blacklist
  return nil
}

func (c *Config) Parse() error {
  return c.Ttl.Parse()
}