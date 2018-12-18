/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

part of frugal.src.frugal;

final TSerializer _serializer =
    new TSerializer(protocolFactory: new TJsonProtocolFactory());

/// Convert the given frugal object to string.
String fObjToJson(Object obj) {
  if (obj is TBase) {
    return new String.fromCharCodes(_serializer.write(obj));
  }
  if (obj is FContext) {
    return json.encode(obj.requestHeaders());
  }
  return json.encode(obj);
}
