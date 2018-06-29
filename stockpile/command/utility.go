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

// evaluates whether a given field is hidden (specifically whether it is merely used to represent
// the internal state of a protobuf message)
func isHiddenField(field *reflect.StructField) bool {
  return strings.HasPrefix(field.Name, "XXX_")
}

// pretty prints an arbitrary value
func printValue(val reflect.Value) []string {
  if !val.IsValid() {
    return []string{"unset"}
  }

  convMethod := val.MethodByName("String")
  if convMethod.IsValid() {
    retValues := convMethod.Call([]reflect.Value{})
    return []string{retValues[0].String()}
  }

  switch val.Kind() {
  case reflect.Slice:
    fallthrough
  case reflect.Array:
    encoded := make([]string, 0)

    for i := 0; i < val.Len(); i++ {
      encodedVal := printValue(val.Index(i))

      for j, str := range encodedVal {
        if j == 0 {
          encoded = append(encoded, "-> "+str)
        } else {
          encoded = append(encoded, "   "+str)
        }
      }

      if i+1 != val.Len() {
        encoded = append(encoded, "")
      }
    }
    return encoded
  case reflect.Struct:
    encoded := make([]string, 0)

    for i := 0; i < val.NumField(); i++ {
      field := val.Type().Field(i)
      fieldValue := val.Field(i)

      if isHiddenField(&field) {
        continue
      }

      encodedField := printValue(fieldValue)
      for j, str := range encodedField {
        if j == 0 {
          encoded = append(encoded, field.Name+": "+str)
        } else {
          encoded = append(encoded, strings.Repeat(" ", len(field.Name)+2)+str)
        }
      }
    }
    return encoded
  case reflect.Ptr:
    return printValue(val.Elem())
  }

  return []string{fmt.Sprintf("%v", val)}
}

// prints a struct in a human readable table format
func writeTable(writer io.Writer, data interface{}) {
  val := reflect.ValueOf(data)
  typ := reflect.TypeOf(data)

  headerCellLength := 3
  valueCellLength := 5
  fieldCount := typ.NumField()

  valueMap := make(map[string][]string)
  for i := 0; i < typ.NumField(); i++ {
    fieldDef := typ.Field(i)
    field := val.Field(i)

    if isHiddenField(&fieldDef) {
      fieldCount--
      continue
    }

    valueMap[fieldDef.Name] = printValue(field)

    keyLength := len(fieldDef.Name)

    if keyLength > headerCellLength {
      headerCellLength = keyLength
    }
    for _, str := range valueMap[fieldDef.Name] {
      l := len(str)

      if l > valueCellLength {
        valueCellLength = l
      }
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

    if isHiddenField(&fieldDef) {
      continue
    }

    io.WriteString(writer, fieldDef.Name)
    io.WriteString(writer, strings.Repeat(" ", headerCellLength-len(fieldDef.Name)))
    io.WriteString(writer, " | ")

    for j, str := range valueMap[fieldDef.Name] {
      if j != 0 {
        io.WriteString(writer, strings.Repeat(" ", headerCellLength))
        io.WriteString(writer, " | ")
      }

      io.WriteString(writer, str)
      io.WriteString(writer, "\n")
    }

    if i+1 != fieldCount {
      io.WriteString(writer, strings.Repeat("-", headerCellLength))
      io.WriteString(writer, "-+-")
      io.WriteString(writer, strings.Repeat("-", valueCellLength))
      io.WriteString(writer, "\n")
    }
  }
}
