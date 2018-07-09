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
package client

import (
  "context"
  "io"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/golang/protobuf/ptypes/empty"
)

// provides an alias for a function which's sole purpose is to respond to error cases in async
// functions
type ErrorFunc = func(error)

// creates an event channel which will be notified about cache events as they occur
func (s *Stockpile) EventChannel(errorHandler ErrorFunc) (chan *entity.Event, error) {
  eventClient, err := s.eventService.StreamEvents(context.Background(), &empty.Empty{})
  if err != nil {
    return nil, err
  }

  outputChannel := make(chan *entity.Event)
  go func() {
    for true {
      evt, err := eventClient.Recv()
      if err != nil {
        if err == io.EOF {
          return
        }

        s.Logger.Printf("Failed to poll for events: %s", err)
        if errorHandler != nil {
          errorHandler(err)
        }
        return
      }

      parsed, err := rpc.EventFromRpc(evt)
      if err != nil {
        s.Logger.Printf("Failed to decode event: %s", err)
        continue
      }

      outputChannel <- parsed
    }
  }()
  return outputChannel, nil
}
