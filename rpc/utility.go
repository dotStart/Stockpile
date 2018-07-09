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
package rpc

import (
  "errors"
  "fmt"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/golang/protobuf/proto"
  "github.com/golang/protobuf/ptypes"
  "github.com/golang/protobuf/ptypes/any"
  "github.com/google/uuid"
)

//go:generate protoc -I . --go_out=plugins=grpc:. common.proto events.proto profile.proto server.proto system.proto

const MessageTypeBaseUrl = "github.com/dotStart/Stockpile/stockpile/server/rpc/"

// Converts a profileId into its fully parsed representation
func ProfileIdFromRpc(rpc *ProfileId) (*entity.ProfileId, error) {
  if !rpc.IsPopulated() {
    return nil, nil
  }

  id, err := uuid.Parse(rpc.Id)
  if err != nil {
    return nil, err
  }

  return &entity.ProfileId{
    Id:          id,
    Name:        rpc.Name,
    FirstSeenAt: time.Unix(rpc.FirstSeenAt, 0),
    LastSeenAt:  time.Unix(rpc.LastSeenAt, 0),
    ValidUntil:  time.Unix(rpc.ValidUntil, 0),
  }, nil
}

// Converts an array of profileIds into their fully parsed representation
func ProfileIdsFromRpcArray(rpc []*ProfileId) ([]*entity.ProfileId, error) {
  arr := make([]*entity.ProfileId, len(rpc))
  for i, encoded := range rpc {
    profileId, err := ProfileIdFromRpc(encoded)
    if err != nil {
      return nil, err
    }
    if profileId == nil {
      return nil, errors.New("encountered one or more unpopulated profiles in response")
    }
    arr[i] = profileId
  }
  return arr, nil
}

// Converts a profileId into its rpc representation
func ProfileIdToRpc(profileId *entity.ProfileId) *ProfileId {
  return &ProfileId{
    Id:          profileId.Id.String(),
    Name:        profileId.Name,
    FirstSeenAt: profileId.FirstSeenAt.Unix(),
    LastSeenAt:  profileId.LastSeenAt.Unix(),
    ValidUntil:  profileId.ValidUntil.Unix(),
  }
}

// Converts an array of profileIds into their rpc representation
func ProfileIdsToRpcArray(profileIds []*entity.ProfileId) []*ProfileId {
  arr := make([]*ProfileId, len(profileIds))
  for i, profileId := range profileIds {
    arr[i] = ProfileIdToRpc(profileId)
  }
  return arr
}

// converts a name change into its rpc representation
func NameChangeToRpc(change *entity.NameChange) *NameHistoryEntry {
  return &NameHistoryEntry{
    Name:        change.Name,
    ChangedToAt: change.ChangedToAt.Unix(),
    ValidUntil:  change.ValidUntil.Unix(),
  }
}

// converts an array of name changes into their rpc representation
func NameChangesToRpcArray(changes []*entity.NameChange) []*NameHistoryEntry {
  arr := make([]*NameHistoryEntry, len(changes))
  for i, change := range changes {
    arr[i] = NameChangeToRpc(change)
  }
  return arr
}

// converts a name change from its rpc representation
func NameChangeFromRpc(rpc *NameHistoryEntry) *entity.NameChange {
  return &entity.NameChange{
    Name:        rpc.Name,
    ChangedToAt: time.Unix(rpc.ChangedToAt, 0),
    ValidUntil:  time.Unix(rpc.ValidUntil, 0),
  }
}

// converts an array of name changes from their rpc representation
func NameChangesFromRpcArray(rpc []*NameHistoryEntry) []*entity.NameChange {
  arr := make([]*entity.NameChange, len(rpc))
  for i, change := range rpc {
    arr[i] = NameChangeFromRpc(change)
  }
  return arr
}

// converts a name history element into its rpc representation
func NameHistoryToRpc(history *entity.NameChangeHistory) *NameHistory {
  if history == nil || len(history.History) == 0 {
    return &NameHistory{}
  }

  return &NameHistory{
    History: NameChangesToRpcArray(history.History),
  }
}

// converts a name history from its rpc representation
func NameHistoryFromRpc(history *NameHistory) *entity.NameChangeHistory {
  if !history.IsPopulated() {
    return nil
  }

  return &entity.NameChangeHistory{
    History: NameChangesFromRpcArray(history.History),
  }
}

// converts the result of a bulk id resolve operation into its rpc representation
func BulkIdsToRpc(ids []*entity.ProfileId) *BulkIdResponse {
  if ids == nil || len(ids) == 0 {
    return &BulkIdResponse{}
  }

  return &BulkIdResponse{
    Ids: ProfileIdsToRpcArray(ids),
  }
}

// converts the result of a bulk id resolve operation from its rpc representation
func BulkIdsFromRpc(rpc *BulkIdResponse) ([]*entity.ProfileId, error) {
  if !rpc.IsPopulated() {
    return nil, nil
  }

  return ProfileIdsFromRpcArray(rpc.Ids)
}

// converts a profile into its rpc representation
func ProfileToRpc(profile *entity.Profile) *Profile {
  var tex *ProfileTextures = nil
  if profile.Textures != nil {
    tex = ProfileTexturesToRpc(profile.Textures)
  }

  return &Profile{
    Id:         profile.Id.String(),
    Name:       profile.Name,
    Properties: ProfilePropertiesToRpcArray(profile.Properties),
    Textures:   tex,
  }
}

// converts a profile from its rpc representation
func ProfileFromRpc(rpc *Profile) (*entity.Profile, error) {
  id, err := uuid.Parse(rpc.Id)
  if err != nil {
    return nil, err
  }

  var tex *entity.ProfileTextures
  if rpc.Textures != nil {
    tex, err = ProfileTexturesFromRpc(rpc.Textures)
    if err != nil {
      return nil, err
    }
  }

  return &entity.Profile{
    Id:         id,
    Name:       rpc.Name,
    Properties: ProfilePropertiesFromRpcArray(rpc.Properties),
    Textures:   tex,
  }, nil
}

// converts a map of profile properties into their rpc representation
func ProfilePropertiesToRpcArray(props map[string]*entity.ProfileProperty) []*ProfileProperty {
  arr := make([]*ProfileProperty, 0)
  for _, prop := range props {
    arr = append(arr, ProfilePropertyToRpc(prop))
  }
  return arr
}

// converts an array of profile properties from their rpc representation
func ProfilePropertiesFromRpcArray(rpc []*ProfileProperty) map[string]*entity.ProfileProperty {
  arr := make(map[string]*entity.ProfileProperty)
  for _, prop := range rpc {
    arr[prop.Name] = ProfilePropertyFromRpc(prop)
  }
  return arr
}

// converts a profile property into its rpc representation
func ProfilePropertyToRpc(prop *entity.ProfileProperty) *ProfileProperty {
  return &ProfileProperty{
    Name:      prop.Name,
    Value:     prop.Value,
    Signature: prop.Signature,
  }
}

// converts a profile property from its rpc representation
func ProfilePropertyFromRpc(rpc *ProfileProperty) *entity.ProfileProperty {
  return &entity.ProfileProperty{
    Name:      rpc.Name,
    Value:     rpc.Value,
    Signature: rpc.Signature,
  }
}

// converts a profile textures attribute into its rpc representation
func ProfileTexturesToRpc(tex *entity.ProfileTextures) *ProfileTextures {
  return &ProfileTextures{
    ProfileId:   tex.ProfileId.String(),
    ProfileName: tex.ProfileName,
    SkinUrl:     tex.Textures["SKIN"], // TODO: this sucks
    CapeUrl:     tex.Textures["CAPE"],
    Timestamp:   tex.Timestamp.Unix(),
  }
}

func ProfileTexturesFromRpc(rpc *ProfileTextures) (*entity.ProfileTextures, error) {
  id, err := uuid.Parse(rpc.ProfileId)
  if err != nil {
    return nil, err
  }

  textures := make(map[string]string)
  textures["SKIN"] = rpc.SkinUrl
  textures["CAPE"] = rpc.CapeUrl

  return &entity.ProfileTextures{
    Timestamp:   time.Unix(rpc.Timestamp, 0),
    ProfileId:   id,
    ProfileName: rpc.ProfileName,
    Textures:    textures,
  }, nil
}

func BlacklistToRpc(blacklist *entity.Blacklist) *Blacklist {
  return &Blacklist{
    Hashes: blacklist.Hashes,
  }
}

func BlacklistFromRpc(blacklist *Blacklist) (*entity.Blacklist, error) {
  return entity.NewBlacklist(blacklist.Hashes)
}

// converts an arbitrary event into its rpc representation
func EventToRpc(event *entity.Event) (*Event, error) {
  key, err := EventKeyToRpc(event.Key)
  if err != nil {
    return nil, err
  }

  obj, err := EventPayloadToRpc(event.Object)
  if err != nil {
    return nil, err
  }

  enc, err := ptypes.MarshalAny(obj)
  if err != nil {
    return nil, err
  }

  return &Event{
    Type:   EventTypeToRpc(event.Type),
    Key:    key,
    Object: enc,
  }, nil
}

func EventFromRpc(event *Event) (*entity.Event, error) {
  typ, err := EventTypeFromRpc(event.Type)
  if err != nil {
    return nil, err
  }

  key, err := EventKeyFromRpc(event.Key)
  if err != nil {
    return nil, err
  }

  payload, err := EventPayloadFromRpc(event.Object)
  if err != nil {
    return nil, err
  }

  return &entity.Event{
    Type:   typ,
    Key:    key,
    Object: payload,
  }, nil
}

// converts an event type into its rpc representation
func EventTypeToRpc(typ entity.EventType) EventType {
  switch typ {
  case entity.ProfileIdEvent:
    return EventType_PROFILE_ID
  case entity.NameHistoryEvent:
    return EventType_NAME_HISTORY
  case entity.ProfileEvent:
    return EventType_PROFILE
  case entity.BlacklistEvent:
    return EventType_BLACKLIST
  default:
    return -1 // TODO: Unknown?
  }
}

// converts an event type from its rpc representation
func EventTypeFromRpc(typ EventType) (entity.EventType, error) {
  switch typ {
  case EventType_PROFILE_ID:
    return entity.ProfileIdEvent, nil
  case EventType_NAME_HISTORY:
    return entity.NameHistoryEvent, nil
  case EventType_PROFILE:
    return entity.ProfileEvent, nil
  case EventType_BLACKLIST:
    return entity.BlacklistEvent, nil
  default:
    return -1, fmt.Errorf("illegal event type: %d", typ)
  }
}

// encodes an arbitrary key type into its rpc representation
func EventKeyToRpc(key interface{}) (*any.Any, error) {
  // nil key is passed as is as there is no identifying information there
  if key == nil {
    return nil, nil
  }

  profileId, ok := key.(*entity.ProfileIdKey)
  if ok {
    return ptypes.MarshalAny(&ProfileIdKey{
      Name: profileId.Name,
      At:   profileId.At.Unix(),
    })
  }

  id, ok := key.(*uuid.UUID)
  if ok {
    return ptypes.MarshalAny(&IdKey{
      Id: id.String(),
    })
  }

  return nil, fmt.Errorf("unknown key type: %v", key)
}

// decodes an arbitrary key type from its rpc representation
func EventKeyFromRpc(key *any.Any) (interface{}, error) {
  if key == nil {
    return nil, nil
  }

  obj := &ptypes.DynamicAny{}
  err := ptypes.UnmarshalAny(key, obj)
  if err != nil {
    return nil, err
  }

  profileId, ok := obj.Message.(*ProfileIdKey)
  if ok {
    return &entity.ProfileIdKey{
      Name: profileId.Name,
      At:   time.Unix(profileId.At, 0),
    }, nil
  }

  id, ok := obj.Message.(*IdKey)
  if ok {
    i, err := uuid.Parse(id.Id)
    return &i, err
  }

  return nil, fmt.Errorf("unknown key type: %v", obj.Message)
}

// converts an event payload into its rpc format
func EventPayloadToRpc(payload interface{}) (proto.Message, error) {
  if payload == nil {
    return nil, errors.New("payload cannot be nil")
  }

  profileId, ok := payload.(*entity.ProfileId)
  if ok {
    return ProfileIdToRpc(profileId), nil
  }

  history, ok := payload.(*entity.NameChangeHistory)
  if ok {
    return NameHistoryToRpc(history), nil
  }

  profile, ok := payload.(*entity.Profile)
  if ok {
    return ProfileToRpc(profile), nil
  }

  blacklist, ok := payload.(*entity.Blacklist)
  if ok {
    return BlacklistToRpc(blacklist), nil
  }

  return nil, fmt.Errorf("illegal payload value: %v", payload)
}

// converts an event payload from its rpc format
func EventPayloadFromRpc(payload *any.Any) (interface{}, error) {
  if payload == nil {
    return nil, errors.New("payload cannot be nil")
  }

  obj := &ptypes.DynamicAny{}
  err := ptypes.UnmarshalAny(payload, obj)
  if err != nil {
    return nil, err
  }

  profileId, ok := obj.Message.(*ProfileId)
  if ok {
    return ProfileIdFromRpc(profileId)
  }

  history, ok := obj.Message.(*NameHistory)
  if ok {
    return NameHistoryFromRpc(history), nil
  }

  profile, ok := obj.Message.(*Profile)
  if ok {
    return ProfileFromRpc(profile)
  }

  blacklist, ok := obj.Message.(*Blacklist)
  if ok {
    return BlacklistFromRpc(blacklist)
  }

  return nil, fmt.Errorf("illegal payload value: %v", payload)
}

func PluginMetadataListToRpc(list []*plugin.Metadata) *PluginList {
  enc := make([]*Plugin, len(list))
  for i, metadata := range list {
    enc[i] = PluginMetadataToRpc(metadata)
  }
  return &PluginList{
    Plugins: enc,
  }
}

func PluginMetadataListFromRpc(list *PluginList) []*plugin.Metadata {
  decoded := make([]*plugin.Metadata, len(list.Plugins))
  for i, metadata := range list.Plugins {
    decoded[i] = PluginMetadataFromRpc(metadata)
  }
  return decoded
}

func PluginMetadataToRpc(metadata *plugin.Metadata) *Plugin {
  return &Plugin{
    Name:    metadata.Name,
    Version: metadata.Version,
    Authors: metadata.Authors,
    Website: metadata.Website,
  }
}

func PluginMetadataFromRpc(metadata *Plugin) *plugin.Metadata {
  return &plugin.Metadata{
    Name:    metadata.Name,
    Version: metadata.Version,
    Authors: metadata.Authors,
    Website: metadata.Website,
  }
}

// evaluates whether the message has been populated with actual data (e.g. whether it is not empty)
func (p *ProfileId) IsPopulated() bool {
  return p.Id != ""
}

// evaluates whether the message has been populated with actual data (e.g. whether it is not empty)
func (r *NameHistory) IsPopulated() bool {
  return r.History != nil && len(r.History) != 0
}

// evaluates whether the message has been populated with actual data (e.g. whether it is not empty)
func (r *BulkIdResponse) IsPopulated() bool {
  return r.Ids != nil && len(r.Ids) != 0
}

// evaluates whether the message has been populated with actual data (e.g. whether it is not empty)
func (p *Profile) IsPopulated() bool {
  return p.Id != ""
}
