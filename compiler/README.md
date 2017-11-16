# Frugal Compiler

The Frugal compiler consists of two main pieces: the IDL parser, and the
code generator.

The parser takes a Frugal interface definition language (IDL) file as
input and returns a parsed representation of the file as a go struct. The
parser is defined using a [parsing expression
grammar](https://en.wikipedia.org/wiki/Parsing_expression_grammar) and
uses the [Pigeon](https://github.com/mna/pigeon) application to generate
go code given the PEG representation of the IDL.

The code generator takes in the parsed IDL as a go struct and outputs
required Frugal constructs in the language of choice.
