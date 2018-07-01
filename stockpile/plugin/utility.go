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
package plugin

import (
  "crypto/sha1"
  "encoding/hex"
  "strings"
  "time"
)

// provides a primitive wrapper object which handles expiration in the memory storage backend
type expirationWrapper struct {
  content   interface{}
  createdAt time.Time
}

// evaluates whether a particular entry is still considered valid
func (w *expirationWrapper) isValid(ttl time.Duration) bool {
  return time.Since(w.createdAt) <= ttl
}

// calculates a unified cache for a given input value (typically for primitive built-in cache types)
func calculateHash(input string) string {
  enc := sha1.Sum([]byte(strings.ToLower(input)))
  return hex.EncodeToString(enc[:])
}
