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
package metadata

import "time"

// Identifies the application brand (typically "vanilla")
func Brand() string {
  return brand
}

// Evaluates whether this version is a custom build
func IsCustomBuild() bool {
  return Brand() != "vanilla"
}

// Retrieves the application version
func Version() string {
  return version
}

// Retrieves the full application version (including its build identifier)
func VersionFull() string {
  versionExtension := "+dev"
  if commitHash != "" {
    versionExtension = "+git-" + commitHash
  }

  return version + versionExtension
}

// Retrieves the commit hash from which this version was built
func CommitHash() string {
  return commitHash
}

// Retrieves the date and time at which this version was built
func Timestamp() time.Time {
  return timestampParsed
}
