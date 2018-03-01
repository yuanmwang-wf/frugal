# Gateway

Gateway is a Frugal compiler option and associated library. It reads IDL
service definitions, and generates a reverse-proxy server which translates
a RESTful JSON API into Frugal service requests. This server can be customized
using annotations in your Frugal IDL definition.

It helps you to provide your APIs in both Frugal and RESTful style at the same time.

## Getting Started

1. Write your IDL as usual.
2. Compile the gateway. The HTTP gateway depends on compiled Go code. You must
   compile your IDL with both the Go and Gateway options.
    * `frugal --gen go gateway_test.frugal && frugal --gen gateway
      gateway_test.frugal`
3. Run your server that implements the IDL. This step depends on your
implementation language.
4. Run the HTTP gateway.
5. Send HTTP+JSON requests to the gateway.

## Configuring the Gateway

### Default JSON Mapping

| frugal   | JSON          | JSON example                     | Notes                                                                                                                                                                               |
| ------   | ----          | ------------                     | -----                                                                                                                                                                               |
| bool     | boolean       | true, false                      |                                                                                                                                                                                     |
| byte     | base64 string | "Y"                              | JSON value will be the data encoded as a string using standard base64 encoding with paddings. Either standard or URL-safe base64 encoding with/without paddings are accepted.       |
| i16, i32 | number        | 1, -10, 0                        | JSON value will be a decimal number. Either numbers or strings are accepted.                                                                                                        |
| i64      | string        | 1, -10, 0                        | JSON value will be a decimal string. Either numbers or strings are accepted.                                                                                                        |
| double   | number        | 1.1, -10.0, 0, "NaN", "Infinity" | JSON value will be a number or one of the special string values "NaN", "Infinity", and "-Infinity". Either numbers or strings are accepted. Exponent notation is also accepted.     |
| binary   | base64 string | "YWJjMTIzIT8kKiYoKSctPUBY+"      | JSON value will be the data encoded as a string using standard base64 encoding with paddings. Either standard or URL-safe base64 encoding with/without paddings are accepted.       |
| string   | string        | "Hello World!"                   |                                                                                                                                                                                     |
| list<V>  | array         | [v, …]                           | null is accepted as the empty list [].                                                                                                                                              |
| map<K,V> | object        | {"k": v, …}                      | All keys are converted to strings. Values are converted as per rules in this table.                                                                                                 |
| struct   | object        | {"fBar": v, "g": null, …}        | Generates JSON objects. Message field names are mapped to lowerCamelCase and become JSON object keys. A null value is treated as the default value of the corresponding field type. |
| enum     | string        | "FOO_BAR"                        | The enum value’s name as specified in the IDL.                                                                                                                                      |

### Custom JSON Mapping

You can customize the mapping for your types using annotations in the
Frugal IDL.

To initially limit the scope of the proxy, the currently supported
annotation syntax will be restricted to support RPC methods that accept
a single struct parameter in the request, and return either a single
struct or void in the response. For example, the following service RPC
could be described using the annotations listed here:

```
TweetSearchResult searchTweets(1:TweetSearchRequest query)
```

Whereas the following service RPC could not — at least currently — be
described using the annotations listed here:

```
string searchTweets(1:string query, 2:int32 tweetType)
```

The reason for this limitation is to provide a clear mapping for RPC
services to be described as HTTP endpoints that accept and return JSON. To
extend this proposal to include all services would require agreeing on the
syntax and semantics for the annotations. This has purposefully been left
out of the initial design to avoid losing the forest for the
trees.

With traditional Frugal syntax, the name of a field in a struct or
parameter is not important and can be changed at any time. However, when
using a struct directly as an interface to an HTTP endpoint this
flexibility is no longer allowed. The approach taken to maintaining
backwards compatibility with Frugal is to allow annotations to structs
that specify how they are to be interpreted as JSON requests or responses.

##### JSON Annotations

| Annotation        | Description                                                                                                                                                              | Example                       |
| ----------        | -----------                                                                                                                                                              | -------                       |
| http.jsonProperty | The property name to use when converting this struct to JSON.

```
struct Tweet {
  1: required i32 userId;
  2: required string userName;
  3: required string text;
  4: optional Location loc; (http.jsonProperty="location")
  16: optional string language = "english";
}
```

### Custom RPC Mapping

RPC calls are mapped to a URL path with query parameters. The default HTTP
method is GET. These can be customized using the annotations described in
this section.

#### RPC Annotations

| Annotation        | Description                                                                                                                                                              | Example                       |
| ----------        | -----------                                                                                                                                                              | -------                       |
| http.method       | The HTTP verb to use when calling this procedure.                                                                                                                        | get                           |
| http.pathTemplate | The HTTP path template to use when calling this procedure.                                                                                                               | /v1/twitter/tweets/{tweet_id} |
| http.query        | For HTTP endpoints requiring query parameters. This field specifies the fields in the request parameter of the RPC call to use as query parameters to the HTTP endpoint. | sheet_id                      |

The annotations specify how different portions of the Frugal method
parameter are mapped to URL path and query parameters, and the HTTP request
body.

#### Interpreting Variables

The annotation syntax allows for two methods of variable interpolation:
path templating and query parameters. The behaviour of this interpolation
is described as follows.

##### http.pathTemplate

The http.pathTemplate annotation maps fields in a Frugal struct to URL
path parameters. For example, with the following annotations

```
struct LookupAlbumRequest {
  1: string ASIN (http.jsonProperty="asin")
}

service Store {
  Album lookupAlbum(1: LookupAlbumRequest request) (http.pathTemplate="/v1/store/album/{asin}")
}
```

incoming HTTP+JSON requests on the URL `/v1/store/album/asdfg` will populate
`LookupAlbumRequest` as `LookupAlbumRequest.ASIN = asdfg`.

##### http.query

The http.query annotation maps fields in a Frugal struct to URL query
parameters. For example, with the following annotations

```
struct LookupAlbumRequest {
  1: string ASIN (http.jsonProperty="asin")
}

service Store {
  Album lookupAlbum(1: LookupAlbumRequest request) (http.pathTemplate="/v1/store/album/", http.query="asin")
}
```

incoming HTTP+JSON requests on the URL `/v1/store/album?asin=asdfg` will
populate `LookupAlbumRequest` as `LookupAlbumRequest.ASIN = asdfg`.

## Known Limitations

* For duplicate query parameters, currently only the first is respected.
* Frugal allows map keys to be any type (including structs and
  exceptions). For maximum compatibility, only base types are allowed for
  map keys in an HTTP proxy.
* The Frugal `binary` type is not currently supported.
* Sets cannot be serialized with the default marshaler, if you use sets,
  another marshaler must be used (not yet developed).

## Todo

* Consistent error handling (with customizable response envelope)
  * Map known error types to HTTP status codes
  * Consistent handling of missing required fields, incorrect data types,
    etc. See:
    https://github.com/Workiva/endpoints/blob/master/docs/guidelines.md#382-error-responses
* Consistent success handling (with customizable response envelope). See:
  https://github.com/Workiva/endpoints/blob/master/docs/guidelines.md#383-success-responses
* Forward headers to FContext (request_id, tracing data, etc.). See:
  https://github.com/Workiva/endpoints/blob/master/docs/guidelines.md#36-custom-headers
* Authorization
