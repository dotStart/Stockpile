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
package entity

import (
  "crypto/sha1"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "regexp"
  "strings"

  "golang.org/x/text/encoding/charmap"
)

var ipPattern, _ = regexp.Compile("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")

// represents a server blacklist
type Blacklist struct {
  Hashes []string
}

// creates a new blacklist from the supplied list of hashes
func NewBlacklist(hashes []string) (*Blacklist, error) {
  for i := 0; i < len(hashes); {
    hash := hashes[i]

    if hash == "" {
      hashes = append(hashes[:i], hashes[i+1:]...)
      continue // skip extras
    }

    if len(hash) != 40 {
      return nil, fmt.Errorf("encountered malformed hash \"%s\": must be exactly 40 characters long", hash)
    }

    i++
  }

  return &Blacklist{
    Hashes: hashes,
  }, nil
}

func (b *Blacklist) Serialize() ([]byte, error) {
  return json.Marshal(b.Hashes)
}

func (b *Blacklist) Deserialize(enc []byte) error {
  hashes := make([]string, 0)
  err := json.Unmarshal(enc, &hashes)
  if err != nil {
    return err
  }

  b.Hashes = hashes
  return nil
}

// evaluates whether a certain hash is part of a blacklist
func (b *Blacklist) Contains(hash string) bool {
  for _, blacklistedHash := range b.Hashes {
    if blacklistedHash == hash {
      return true
    }
  }

  return false
}

// evaluates whether the passed hostname has been blacklisted
func (b *Blacklist) IsBlacklisted(addr string) (bool, error) {
  hash, err := calculateHash(addr)
  if err != nil {
    return false, err
  }
  if b.Contains(hash) {
    return true, nil
  }

  if isIpAddress(addr) {
    return b.IsBlacklistedIP(addr)
  }

  return b.IsBlacklistedDomain(addr)
}

// evaluates whether a given IPv4 address has been blacklisted
func (b *Blacklist) IsBlacklistedIP(ip string) (bool, error) {
  elements := strings.Split(ip, ".") // TODO: does Minecraft support IPv6 blacklisting?
  for i := 3; i > 0; i-- {
    addr := strings.Join(elements[:i], ".") + ".*"
    hash, err := calculateHash(addr)
    if err != nil {
      return false, err
    }

    if b.Contains(hash) {
      return true, nil
    }
  }

  return false, nil
}

// evaluates whether a given hostname has been blacklisted
func (b *Blacklist) IsBlacklistedDomain(hostname string) (bool, error) {
  elements := strings.Split(hostname, ".")
  length := len(elements)

  for i := 1; i < length; i++ {
    addr := "*." + strings.Join(elements[i:], ".")
    hash, err := calculateHash(addr)
    if err != nil {
      return false, err
    }

    if b.Contains(hash) {
      return true, nil
    }
  }

  return false, nil
}

// calculates a blacklist compatible hash
func calculateHash(input string) (string, error) {
  encoder := charmap.ISO8859_10.NewEncoder()

  encoded, err := encoder.String(input)
  if err != nil {
    return "", err
  }

  hash := sha1.Sum([]byte(encoded))
  return hex.EncodeToString(hash[:]), nil
}

// evaluates whether the given address is an IPv4 address
func isIpAddress(address string) bool {
  return ipPattern.MatchString(address)
}
