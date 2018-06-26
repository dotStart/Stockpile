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
package mojang

import "time"

// defines the total amount of time a name can be safely associated with a given profile
const NameValidityPeriod = time.Hour * 24 * 37

// defines the total amount of time that has to pass before a user can choose a new name again
const NameChangeRateLimitPeriod = time.Hour * 24 * 30

// calculates the beginning of a theoretical grace period
func CalculateNameGracePeriodBeginning(end time.Time) time.Time {
  return end.Add(-NameValidityPeriod)
}

// calculates the end of a theoretical grace period
func CalculateNameGracePeriodEnd(start time.Time) time.Time {
  return start.Add(NameValidityPeriod)
}

// evaluates whether a given name association is still considered valid (e.g. no other user was able to claim the name
// since it was last encountered)
func IsNameAssociationValidAt(at time.Time, lastSeen time.Time) bool {
  return at.Sub(lastSeen) < NameValidityPeriod
}

// evaluates whether a given name association is still considered valid (e.g. no other user was able to claim the name
// since it was last encountered)
func IsNameAssociationValid(lastSeen time.Time) bool {
  return IsNameAssociationValidAt(time.Now(), lastSeen)
}
