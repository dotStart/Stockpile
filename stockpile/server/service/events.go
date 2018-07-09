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
package service

import (
  "github.com/dotStart/Stockpile/rpc"
  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/golang/protobuf/ptypes/empty"
  "github.com/op/go-logging"
  "google.golang.org/grpc/peer"
)

type EventServiceImpl struct {
  logger *logging.Logger
  cache  *cache.Cache
}

func NewEventService(cache *cache.Cache) (*EventServiceImpl) {
  return &EventServiceImpl{
    logger: logging.MustGetLogger("event-srv"),
    cache:  cache,
  }
}

func (s *EventServiceImpl) StreamEvents(_ *empty.Empty, srv rpc.EventService_StreamEventsServer) error {
  p, ok := peer.FromContext(srv.Context())
  if ok {
    s.logger.Debugf("beginning to stream events to rpc client %s", p.Addr)
  }

  listener := s.cache.NewListener()
  defer listener.Close()

  for e := range listener.C {
    if ok {
      s.logger.Debugf("forwarding event of type %T (using key %T) to rpc client %s", e.Object, e.Key, p.Addr)
    }

    enc, err := rpc.EventToRpc(e)
    if err != nil {
      s.logger.Errorf("failed to encode event %v: %s", e, err)
      continue
    }
    srv.Send(enc)
  }

  if ok {
    s.logger.Debugf("event stream has ended - closing session with %s", p.Addr)
  }
  return nil
}
