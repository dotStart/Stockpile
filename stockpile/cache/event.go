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
package cache

import (
  "errors"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

// represents an even which has occurred within the cache
type Event struct {
  Type   EventType
  Key    interface{}
  Object interface{}
}

func (e *Event) ProfileIdPayload() (*entity.ProfileId, error) {
  if e.Type != ProfileIdEvent {
    return nil, errors.New("cannot convert event payload to ProfileId")
  }

  return e.Object.(*entity.ProfileId), nil
}

func (e *Event) NameChangeHistoryPayload() (*entity.NameChangeHistory, error) {
  if e.Type != NameHistoryEvent {
    return nil, errors.New("cannot convert event payload to NameChangeHistory")
  }

  return e.Object.(*entity.NameChangeHistory), nil
}

func (e *Event) ProfilePayload() (*entity.Profile, error) {
  if e.Type != ProfileEvent {
    return nil, errors.New("cannot convert event payload to Profile")
  }

  return e.Object.(*entity.Profile), nil
}

func (e *Event) BlacklistPayload() (*entity.Blacklist, error) {
  if e.Type != BlacklistEvent {
    return nil, errors.New("cannot convert event payload to Blacklist")
  }

  return e.Object.(*entity.Blacklist), nil
}

func (e *Event) ProfileIdKey() (*ProfileIdKey, error) {
  if e.Type != ProfileIdEvent {
    return nil, errors.New("cannot convert event key to ProfileIdKey")
  }

  return e.Key.(*ProfileIdKey), nil
}

func (e *Event) IdKey() (*uuid.UUID, error) {
  if e.Type != NameHistoryEvent && e.Type != ProfileEvent {
    return nil, errors.New("cannot convert event key to UUID")
  }

  return e.Key.(*uuid.UUID), nil
}

// indicates which type of data is changed as part of an event
type EventType int32

const (
  ProfileIdEvent   EventType = 0
  NameHistoryEvent EventType = 1
  ProfileEvent     EventType = 2
  BlacklistEvent   EventType = 3
)

type ProfileIdKey struct {
  Name string
  At   time.Time
}
