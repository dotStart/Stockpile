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
package command

import (
  "fmt"
  "io"
  "reflect"
  "strings"
)

// prints a struct in a human readable table format
func writeTable(writer io.Writer, data interface{}) {
  val := reflect.ValueOf(data)
  typ := reflect.TypeOf(data)

  headerCellLength := 3
  valueCellLength := 5

  for i := 0; i < typ.NumField(); i++ {
    fieldDef := typ.Field(i)
    field := val.Field(i)

    if strings.HasPrefix(fieldDef.Name, "XXX_") {
      continue
    }

    keyLength := len(fieldDef.Name)
    val := fmt.Sprintf("%v", field)
    valueLength := len(val)

    if keyLength > headerCellLength {
      headerCellLength = keyLength
    }
    if valueLength > valueCellLength {
      valueCellLength = valueLength
    }
  }

  io.WriteString(writer, "Key")
  io.WriteString(writer, strings.Repeat(" ", headerCellLength-3))
  io.WriteString(writer, " | ")
  io.WriteString(writer, "Value\n")

  io.WriteString(writer, strings.Repeat("-", headerCellLength))
  io.WriteString(writer, "-+-")
  io.WriteString(writer, strings.Repeat("-", valueCellLength))
  io.WriteString(writer, "\n")

  for i := 0; i < typ.NumField(); i++ {
    fieldDef := typ.Field(i)
    field := val.Field(i)

    if strings.HasPrefix(fieldDef.Name, "XXX_") {
      continue
    }

    io.WriteString(writer, fieldDef.Name)
    io.WriteString(writer, strings.Repeat(" ", headerCellLength-len(fieldDef.Name)))
    io.WriteString(writer, " | ")

    io.WriteString(writer, fmt.Sprintf("%v", field))
    io.WriteString(writer, "\n")
  }
}
