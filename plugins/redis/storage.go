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
package main

import (
  "fmt"
  "time"

  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/go-redis/redis"
  "github.com/hashicorp/hcl2/gohcl"
  "github.com/op/go-logging"
)

var defaultPassword = ""

type redisStorageBackendInterface struct {
  logger *logging.Logger
  cfg    *RedisStorageBackendConfig
  client *redis.Client
}

type RedisStorageBackendConfig struct {
  Address    string  `hcl:"address,attr"`
  Password   *string `hcl:"password,attr"`
  DatabaseId int     `hcl:"database,attr"`
}

func NewRedisStorageBackend(cfg *server.Config) (storage.StorageBackend, error) {
  redisCfg := &RedisStorageBackendConfig{}
  diag := gohcl.DecodeBody(cfg.Storage.Parameters, nil, redisCfg)
  if diag.HasErrors() {
    return nil, fmt.Errorf("illegal backend configuration: %s", diag.Error())
  }

  if redisCfg.Password == nil {
    redisCfg.Password = &defaultPassword
  }

  client := redis.NewClient(&redis.Options{
    Addr:     redisCfg.Address,
    Password: *redisCfg.Password,
    DB:       redisCfg.DatabaseId,
  })

  _, err := client.Ping().Result()
  if err != nil {
    return nil, fmt.Errorf("cannot reach configured redis server: %s", err)
  }

  return storage.NewEncodedStorageBackend(cfg, &redisStorageBackendInterface{
    logger: logging.MustGetLogger("redis"),
    cfg:    redisCfg,
    client: client,
  }), nil
}

func (f *redisStorageBackendInterface) GetCacheEntry(category string, key string, ttl time.Duration) ([]byte, error) {
  enc, err := f.client.Get(fmt.Sprintf("%s_%s", category, key)).Bytes()
  if err == redis.Nil {
    return nil, nil
  }
  return enc, err
}

func (f *redisStorageBackendInterface) PutCacheEntry(category string, key string, data []byte, ttl time.Duration) error {
  return f.client.Set(fmt.Sprintf("%s_%s", category, key), data, ttl).Err()
}

func (f *redisStorageBackendInterface) PurgeCacheEntry(category string, key string) error {
  return f.client.Del(fmt.Sprintf("%s_%s", category, key)).Err()
}

func (f *redisStorageBackendInterface) Close() error {
  return f.client.Close()
}
