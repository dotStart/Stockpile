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
  "fmt"
  "regexp"
  "strings"

  "github.com/google/uuid"
)

// this is just a rudimentary regex for standardized and Mojang UUIDs and will not verify whether
// the version and variant bits are set correctly
// TODO: This is only used by the legacy api
var idRegex, _ = regexp.Compile("^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})|([a-fA-F0-9]{8}[a-fA-F0-9]{4}[a-fA-F0-9]{4}[a-fA-F0-9]{4}[a-fA-F0-9]{12})$")

// Parses an identifier (regardless of whether it is supplied in its Mojang or standard format)
func ParseId(id string) (uuid.UUID, error) {
  if IsMojangId(id) {
    return ToStandardId(id)
  }

  return uuid.Parse(id)
}

// evaluates whether the passed string is a profileId
// TODO: This is only used by the legacy API
func IsProfileId(id string) bool {
  return idRegex.MatchString(id)
}

// Evaluates whether the passed ID is a mojang identifier
func IsMojangId(id string) bool {
  return !strings.Contains(id, "-")
}

// Converts a standard UUID into its Mojang format
func ToMojangId(id uuid.UUID) string {
  return strings.Replace(id.String(), "-", "", -1)
}

// Converts a Mojang UUID into its RFC format
func ToStandardId(id string) (uuid.UUID, error) {
  encoded, err := ToStandardIdString(id)
  if err != nil {
    return uuid.Nil, err
  }

  return uuid.Parse(encoded)
}

// Converts a Mojang UUID into its RFC format
func ToStandardIdString(id string) (string, error) {
  if len(id) != 32 {
    return uuid.Nil.String(), fmt.Errorf("illegal Mojang identifier length: expected 32 characters but got %d", len(id))
  }

  return fmt.Sprintf("%s-%s-%s-%s-%s", id[0:8], id[8:12], id[12:16], id[16:20], id[20:32]), nil
}
