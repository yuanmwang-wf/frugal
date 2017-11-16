---
authors: Kevin Sookocheff <kevin@sookocheff.com>
state: draft
---

# RFD 1 Annotations

## Introduction

The Frugal IDL is feature rich for the use cases of remote procedure calls
and pub/sub messaging. To allow extension of this feature set to custom
use cases not defined or anticipated by the core Frugal library, this
document proposes the general purpose concept of `annotations`.
Annotations are presented in conjunction with a compiler plugin system
that can be leveraged by any developer to augment the functionality of the
Frugal compiler.

### Reasoning

Annotations allow third-parties to augment the Frugal compiler to generate
code specific to their application or need. Examples include generating
a RESTful HTTP server from a Frugal IDL file, or generating custom
documentation for their service, without having to modify the core Frugal
library.

### Key requirements

With the above in mind, we lay out the following core requirements:

1. Syntax for annotating any Frugal IDL statements.
2. A compiler plugin system for generating custom code from an IDL.

## Annotations Syntax

Frugal currently supports two existing annotations: `deprecated` and
`vendor`. Each annotation is specified using the same syntax:

```
<statement> (<annotation-name>="<annotation-value")
```

As a concrete example, the syntax for a `vendor` annotation is as follows:

```
namespace go bar (vendor="github.com/Workiva/my-repo/gen-go/bar")
```

where,

* `statement` = `namespace go bar`
* `annotation-name` = `vendor`
* `annotation-value` = `github.com/Workiva/my-repo/gen-go/bar`

And the syntax for a `deprecated` annotation is as follows:

```
bool enterAlbumGiveaway( 1: string email, 2: string name ) (deprecated="use something else")
```

where,

* `statement` = `bool enterAlbumGiveaway( 1: string email, 2: string name )`
* `annotation-name` = `deprecated`
* `annotation-value` = `use something else`

Annotations are specified using PEG as

```
TypeAnnotations <- '(' __ annotations:TypeAnnotation* ')' {
    var anns []*Annotation
    for _, ann := range annotations.([]interface{}) {
        anns = append(anns, ann.(*Annotation))
    }
    return anns, nil
}

TypeAnnotation <- name:Identifier _ value:('=' __ value:Literal { return value, nil })? ListSeparator? __ {
    var optValue string
    if value != nil {
        optValue = value.(string)
    }
    return &Annotation{
        Name:  string(name.(Identifier)),
        Value: optValue,
    }, nil
}
```


## Annotations Semantics

## Compiler Plugin

## Prototype

