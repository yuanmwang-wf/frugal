package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	identifier     = regexp.MustCompile("^[A-Za-z]+[A-Za-z0-9]")
	prefixVariable = regexp.MustCompile("{\\w*}")
	defaultPrefix  = &ScopePrefix{String: "", Variables: make([]string, 0)}
)

type statementWrapper struct {
	comment   []string
	statement interface{}
}

type exception *Struct

type union *Struct

func newScopePrefix(prefix string) (*ScopePrefix, error) {
	variables := []string{}
	for _, variable := range prefixVariable.FindAllString(prefix, -1) {
		variable = variable[1 : len(variable)-1]
		if len(variable) == 0 || !identifier.MatchString(variable) {
			return nil, fmt.Errorf("parser: invalid prefix variable '%s'", variable)
		}
		variables = append(variables, variable)
	}
	return &ScopePrefix{String: prefix, Variables: variables}, nil
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func ifaceSliceToString(v interface{}) string {
	ifs := toIfaceSlice(v)
	b := make([]byte, len(ifs))
	for i, v := range ifs {
		b[i] = v.([]uint8)[0]
	}
	return string(b)
}

func rawCommentToDocStr(raw string) []string {
	rawLines := strings.Split(raw, "\n")
	comment := make([]string, len(rawLines))
	for i, line := range rawLines {
		comment[i] = strings.TrimLeft(line, "* ")
	}
	return comment
}

// toStruct converts a union to a struct with all fields optional.
func unionToStruct(u union) *Struct {
	st := (*Struct)(u)
	for _, f := range st.Fields {
		f.Modifier = Optional
	}
	return st
}

// toAnnotations converts an interface{} to an Annotation slice.
func toAnnotations(v interface{}) Annotations {
	if v == nil {
		return nil
	}
	return Annotations(v.([]*Annotation))
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 98, col: 1, offset: 2976},
			expr: &actionExpr{
				pos: position{line: 98, col: 12, offset: 2987},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 98, col: 12, offset: 2987},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 98, col: 12, offset: 2987},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 98, col: 15, offset: 2990},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 98, col: 26, offset: 3001},
								expr: &seqExpr{
									pos: position{line: 98, col: 28, offset: 3003},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 98, col: 28, offset: 3003},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 98, col: 38, offset: 3013},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 98, col: 45, offset: 3020},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 98, col: 45, offset: 3020},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 98, col: 51, offset: 3026},
									name: "SyntaxError",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SyntaxError",
			pos:  position{line: 163, col: 1, offset: 5375},
			expr: &actionExpr{
				pos: position{line: 163, col: 16, offset: 5390},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 163, col: 16, offset: 5390,
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 167, col: 1, offset: 5448},
			expr: &actionExpr{
				pos: position{line: 167, col: 14, offset: 5461},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 167, col: 14, offset: 5461},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 167, col: 14, offset: 5461},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 167, col: 21, offset: 5468},
								expr: &seqExpr{
									pos: position{line: 167, col: 22, offset: 5469},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 167, col: 22, offset: 5469},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 167, col: 32, offset: 5479},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 167, col: 37, offset: 5484},
							label: "statement",
							expr: &ruleRefExpr{
								pos:  position{line: 167, col: 47, offset: 5494},
								name: "FrugalStatement",
							},
						},
					},
				},
			},
		},
		{
			name: "FrugalStatement",
			pos:  position{line: 180, col: 1, offset: 5964},
			expr: &choiceExpr{
				pos: position{line: 180, col: 20, offset: 5983},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 180, col: 20, offset: 5983},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 30, offset: 5993},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 42, offset: 6005},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 50, offset: 6013},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 57, offset: 6020},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 67, offset: 6030},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 76, offset: 6039},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 88, offset: 6051},
						name: "Union",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 96, offset: 6059},
						name: "Service",
					},
					&ruleRefExpr{
						pos:  position{line: 180, col: 106, offset: 6069},
						name: "Scope",
					},
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 182, col: 1, offset: 6076},
			expr: &actionExpr{
				pos: position{line: 182, col: 12, offset: 6087},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 182, col: 12, offset: 6087},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 182, col: 12, offset: 6087},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 182, col: 22, offset: 6097},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 182, col: 24, offset: 6099},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 182, col: 29, offset: 6104},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 182, col: 37, offset: 6112},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 182, col: 39, offset: 6114},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 182, col: 51, offset: 6126},
								expr: &ruleRefExpr{
									pos:  position{line: 182, col: 51, offset: 6126},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 182, col: 68, offset: 6143},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 194, col: 1, offset: 6420},
			expr: &actionExpr{
				pos: position{line: 194, col: 14, offset: 6433},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 194, col: 14, offset: 6433},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 194, col: 14, offset: 6433},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 26, offset: 6445},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 28, offset: 6447},
							label: "scope",
							expr: &oneOrMoreExpr{
								pos: position{line: 194, col: 34, offset: 6453},
								expr: &charClassMatcher{
									pos:        position{line: 194, col: 34, offset: 6453},
									val:        "[*a-z.-]",
									chars:      []rune{'*', '.', '-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 44, offset: 6463},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 46, offset: 6465},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 194, col: 49, offset: 6468},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 60, offset: 6479},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 62, offset: 6481},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 74, offset: 6493},
								expr: &ruleRefExpr{
									pos:  position{line: 194, col: 74, offset: 6493},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 91, offset: 6510},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 202, col: 1, offset: 6696},
			expr: &actionExpr{
				pos: position{line: 202, col: 10, offset: 6705},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 202, col: 10, offset: 6705},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 202, col: 10, offset: 6705},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 18, offset: 6713},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 20, offset: 6715},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 24, offset: 6719},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 34, offset: 6729},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 36, offset: 6731},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 41, offset: 6736},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 52, offset: 6747},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 202, col: 54, offset: 6749},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 58, offset: 6753},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 60, offset: 6755},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 66, offset: 6761},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 77, offset: 6772},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 79, offset: 6774},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 202, col: 91, offset: 6786},
								expr: &ruleRefExpr{
									pos:  position{line: 202, col: 91, offset: 6786},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 108, offset: 6803},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 211, col: 1, offset: 6997},
			expr: &actionExpr{
				pos: position{line: 211, col: 9, offset: 7005},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 211, col: 9, offset: 7005},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 211, col: 9, offset: 7005},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 16, offset: 7012},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 18, offset: 7014},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 23, offset: 7019},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 34, offset: 7030},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 211, col: 37, offset: 7033},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 41, offset: 7037},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 44, offset: 7040},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 211, col: 51, offset: 7047},
								expr: &seqExpr{
									pos: position{line: 211, col: 52, offset: 7048},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 211, col: 52, offset: 7048},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 211, col: 62, offset: 7058},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 211, col: 67, offset: 7063},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 71, offset: 7067},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 73, offset: 7069},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 211, col: 85, offset: 7081},
								expr: &ruleRefExpr{
									pos:  position{line: 211, col: 85, offset: 7081},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 102, offset: 7098},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 235, col: 1, offset: 7760},
			expr: &actionExpr{
				pos: position{line: 235, col: 14, offset: 7773},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 235, col: 14, offset: 7773},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 235, col: 14, offset: 7773},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 235, col: 21, offset: 7780},
								expr: &seqExpr{
									pos: position{line: 235, col: 22, offset: 7781},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 235, col: 22, offset: 7781},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 235, col: 32, offset: 7791},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 235, col: 37, offset: 7796},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 235, col: 42, offset: 7801},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 53, offset: 7812},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 235, col: 55, offset: 7814},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 235, col: 61, offset: 7820},
								expr: &seqExpr{
									pos: position{line: 235, col: 62, offset: 7821},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 235, col: 62, offset: 7821},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 235, col: 66, offset: 7825},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 235, col: 68, offset: 7827},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 82, offset: 7841},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 235, col: 84, offset: 7843},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 235, col: 96, offset: 7855},
								expr: &ruleRefExpr{
									pos:  position{line: 235, col: 96, offset: 7855},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 235, col: 113, offset: 7872},
							expr: &ruleRefExpr{
								pos:  position{line: 235, col: 113, offset: 7872},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 251, col: 1, offset: 8270},
			expr: &actionExpr{
				pos: position{line: 251, col: 12, offset: 8281},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 251, col: 12, offset: 8281},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 251, col: 12, offset: 8281},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 251, col: 22, offset: 8291},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 251, col: 24, offset: 8293},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 251, col: 28, offset: 8297},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 251, col: 38, offset: 8307},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 251, col: 40, offset: 8309},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 251, col: 45, offset: 8314},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 251, col: 56, offset: 8325},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 251, col: 58, offset: 8327},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 251, col: 70, offset: 8339},
								expr: &ruleRefExpr{
									pos:  position{line: 251, col: 70, offset: 8339},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 251, col: 87, offset: 8356},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 259, col: 1, offset: 8528},
			expr: &actionExpr{
				pos: position{line: 259, col: 11, offset: 8538},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 259, col: 11, offset: 8538},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 259, col: 11, offset: 8538},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 259, col: 20, offset: 8547},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 259, col: 22, offset: 8549},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 259, col: 25, offset: 8552},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 260, col: 1, offset: 8592},
			expr: &actionExpr{
				pos: position{line: 260, col: 14, offset: 8605},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 260, col: 14, offset: 8605},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 260, col: 14, offset: 8605},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 26, offset: 8617},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 28, offset: 8619},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 31, offset: 8622},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Union",
			pos:  position{line: 261, col: 1, offset: 8673},
			expr: &actionExpr{
				pos: position{line: 261, col: 10, offset: 8682},
				run: (*parser).callonUnion1,
				expr: &seqExpr{
					pos: position{line: 261, col: 10, offset: 8682},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 261, col: 10, offset: 8682},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 18, offset: 8690},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 261, col: 20, offset: 8692},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 261, col: 23, offset: 8695},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 262, col: 1, offset: 8742},
			expr: &actionExpr{
				pos: position{line: 262, col: 15, offset: 8756},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 262, col: 15, offset: 8756},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 262, col: 15, offset: 8756},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 262, col: 20, offset: 8761},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 262, col: 31, offset: 8772},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 262, col: 34, offset: 8775},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 262, col: 38, offset: 8779},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 262, col: 41, offset: 8782},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 262, col: 48, offset: 8789},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 262, col: 58, offset: 8799},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 262, col: 62, offset: 8803},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 262, col: 64, offset: 8805},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 262, col: 76, offset: 8817},
								expr: &ruleRefExpr{
									pos:  position{line: 262, col: 76, offset: 8817},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 262, col: 93, offset: 8834},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 273, col: 1, offset: 9051},
			expr: &actionExpr{
				pos: position{line: 273, col: 14, offset: 9064},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 273, col: 14, offset: 9064},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 273, col: 21, offset: 9071},
						expr: &seqExpr{
							pos: position{line: 273, col: 22, offset: 9072},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 273, col: 22, offset: 9072},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 273, col: 28, offset: 9078},
									name: "__",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 282, col: 1, offset: 9259},
			expr: &actionExpr{
				pos: position{line: 282, col: 10, offset: 9268},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 282, col: 10, offset: 9268},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 282, col: 10, offset: 9268},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 282, col: 17, offset: 9275},
								expr: &seqExpr{
									pos: position{line: 282, col: 18, offset: 9276},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 282, col: 18, offset: 9276},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 282, col: 28, offset: 9286},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 282, col: 33, offset: 9291},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 36, offset: 9294},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 48, offset: 9306},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 282, col: 50, offset: 9308},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 54, offset: 9312},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 282, col: 56, offset: 9314},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 282, col: 60, offset: 9318},
								expr: &ruleRefExpr{
									pos:  position{line: 282, col: 60, offset: 9318},
									name: "FieldModifier",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 75, offset: 9333},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 282, col: 77, offset: 9335},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 81, offset: 9339},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 91, offset: 9349},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 282, col: 93, offset: 9351},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 98, offset: 9356},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 109, offset: 9367},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 282, col: 112, offset: 9370},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 282, col: 116, offset: 9374},
								expr: &seqExpr{
									pos: position{line: 282, col: 117, offset: 9375},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 282, col: 117, offset: 9375},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 282, col: 121, offset: 9379},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 282, col: 123, offset: 9381},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 282, col: 136, offset: 9394},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 282, col: 138, offset: 9396},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 282, col: 150, offset: 9408},
								expr: &ruleRefExpr{
									pos:  position{line: 282, col: 150, offset: 9408},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 282, col: 167, offset: 9425},
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 167, offset: 9425},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldModifier",
			pos:  position{line: 305, col: 1, offset: 9957},
			expr: &actionExpr{
				pos: position{line: 305, col: 18, offset: 9974},
				run: (*parser).callonFieldModifier1,
				expr: &choiceExpr{
					pos: position{line: 305, col: 19, offset: 9975},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 305, col: 19, offset: 9975},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 305, col: 32, offset: 9988},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 313, col: 1, offset: 10131},
			expr: &actionExpr{
				pos: position{line: 313, col: 12, offset: 10142},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 313, col: 12, offset: 10142},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 313, col: 12, offset: 10142},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 22, offset: 10152},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 313, col: 24, offset: 10154},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 313, col: 29, offset: 10159},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 40, offset: 10170},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 313, col: 42, offset: 10172},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 313, col: 50, offset: 10180},
								expr: &seqExpr{
									pos: position{line: 313, col: 51, offset: 10181},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 313, col: 51, offset: 10181},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 313, col: 61, offset: 10191},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 313, col: 64, offset: 10194},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 313, col: 75, offset: 10205},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 80, offset: 10210},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 313, col: 83, offset: 10213},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 87, offset: 10217},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 313, col: 90, offset: 10220},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 313, col: 98, offset: 10228},
								expr: &seqExpr{
									pos: position{line: 313, col: 99, offset: 10229},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 313, col: 99, offset: 10229},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 313, col: 108, offset: 10238},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 313, col: 114, offset: 10244},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 313, col: 114, offset: 10244},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 313, col: 120, offset: 10250},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 139, offset: 10269},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 313, col: 141, offset: 10271},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 313, col: 153, offset: 10283},
								expr: &ruleRefExpr{
									pos:  position{line: 313, col: 153, offset: 10283},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 313, col: 170, offset: 10300},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 330, col: 1, offset: 10741},
			expr: &actionExpr{
				pos: position{line: 330, col: 22, offset: 10762},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 330, col: 22, offset: 10762,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 334, col: 1, offset: 10831},
			expr: &actionExpr{
				pos: position{line: 334, col: 13, offset: 10843},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 334, col: 13, offset: 10843},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 334, col: 13, offset: 10843},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 334, col: 20, offset: 10850},
								expr: &seqExpr{
									pos: position{line: 334, col: 21, offset: 10851},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 334, col: 21, offset: 10851},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 334, col: 31, offset: 10861},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 334, col: 36, offset: 10866},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 334, col: 43, offset: 10873},
								expr: &seqExpr{
									pos: position{line: 334, col: 44, offset: 10874},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 334, col: 44, offset: 10874},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 334, col: 53, offset: 10883},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 334, col: 58, offset: 10888},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 334, col: 62, offset: 10892},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 75, offset: 10905},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 334, col: 78, offset: 10908},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 334, col: 83, offset: 10913},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 94, offset: 10924},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 334, col: 96, offset: 10926},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 100, offset: 10930},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 334, col: 103, offset: 10933},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 334, col: 113, offset: 10943},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 334, col: 123, offset: 10953},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 127, offset: 10957},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 334, col: 130, offset: 10960},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 334, col: 141, offset: 10971},
								expr: &ruleRefExpr{
									pos:  position{line: 334, col: 141, offset: 10971},
									name: "Throws",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 334, col: 149, offset: 10979},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 334, col: 151, offset: 10981},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 334, col: 163, offset: 10993},
								expr: &ruleRefExpr{
									pos:  position{line: 334, col: 163, offset: 10993},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 334, col: 180, offset: 11010},
							expr: &ruleRefExpr{
								pos:  position{line: 334, col: 180, offset: 11010},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 362, col: 1, offset: 11661},
			expr: &actionExpr{
				pos: position{line: 362, col: 17, offset: 11677},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 362, col: 17, offset: 11677},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 362, col: 22, offset: 11682},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 362, col: 22, offset: 11682},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 362, col: 31, offset: 11691},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 369, col: 1, offset: 11813},
			expr: &actionExpr{
				pos: position{line: 369, col: 11, offset: 11823},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 369, col: 11, offset: 11823},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 369, col: 11, offset: 11823},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 20, offset: 11832},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 369, col: 23, offset: 11835},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 27, offset: 11839},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 369, col: 30, offset: 11842},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 369, col: 41, offset: 11853},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 369, col: 51, offset: 11863},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 373, col: 1, offset: 11899},
			expr: &actionExpr{
				pos: position{line: 373, col: 14, offset: 11912},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 373, col: 14, offset: 11912},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 373, col: 19, offset: 11917},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 373, col: 19, offset: 11917},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 373, col: 30, offset: 11928},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 373, col: 46, offset: 11944},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 380, col: 1, offset: 12069},
			expr: &actionExpr{
				pos: position{line: 380, col: 13, offset: 12081},
				run: (*parser).callonBaseType1,
				expr: &seqExpr{
					pos: position{line: 380, col: 13, offset: 12081},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 380, col: 13, offset: 12081},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 380, col: 18, offset: 12086},
								name: "BaseTypeName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 380, col: 31, offset: 12099},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 380, col: 33, offset: 12101},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 380, col: 45, offset: 12113},
								expr: &ruleRefExpr{
									pos:  position{line: 380, col: 45, offset: 12113},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "BaseTypeName",
			pos:  position{line: 387, col: 1, offset: 12249},
			expr: &actionExpr{
				pos: position{line: 387, col: 17, offset: 12265},
				run: (*parser).callonBaseTypeName1,
				expr: &choiceExpr{
					pos: position{line: 387, col: 18, offset: 12266},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 387, col: 18, offset: 12266},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 27, offset: 12275},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 36, offset: 12284},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 44, offset: 12292},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 52, offset: 12300},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 60, offset: 12308},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 71, offset: 12319},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 387, col: 82, offset: 12330},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 391, col: 1, offset: 12377},
			expr: &actionExpr{
				pos: position{line: 391, col: 18, offset: 12394},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 391, col: 18, offset: 12394},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 391, col: 23, offset: 12399},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 391, col: 23, offset: 12399},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 391, col: 33, offset: 12409},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 391, col: 43, offset: 12419},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 395, col: 1, offset: 12454},
			expr: &actionExpr{
				pos: position{line: 395, col: 12, offset: 12465},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 395, col: 12, offset: 12465},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 395, col: 12, offset: 12465},
							expr: &ruleRefExpr{
								pos:  position{line: 395, col: 12, offset: 12465},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 395, col: 21, offset: 12474},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 28, offset: 12481},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 395, col: 31, offset: 12484},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 395, col: 35, offset: 12488},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 45, offset: 12498},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 395, col: 48, offset: 12501},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 52, offset: 12505},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 395, col: 55, offset: 12508},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 395, col: 61, offset: 12514},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 71, offset: 12524},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 395, col: 74, offset: 12527},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 395, col: 78, offset: 12531},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 395, col: 80, offset: 12533},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 395, col: 92, offset: 12545},
								expr: &ruleRefExpr{
									pos:  position{line: 395, col: 92, offset: 12545},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 404, col: 1, offset: 12743},
			expr: &actionExpr{
				pos: position{line: 404, col: 12, offset: 12754},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 404, col: 12, offset: 12754},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 404, col: 12, offset: 12754},
							expr: &ruleRefExpr{
								pos:  position{line: 404, col: 12, offset: 12754},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 404, col: 21, offset: 12763},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 404, col: 28, offset: 12770},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 404, col: 31, offset: 12773},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 404, col: 35, offset: 12777},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 404, col: 45, offset: 12787},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 404, col: 48, offset: 12790},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 404, col: 52, offset: 12794},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 404, col: 54, offset: 12796},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 404, col: 66, offset: 12808},
								expr: &ruleRefExpr{
									pos:  position{line: 404, col: 66, offset: 12808},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 412, col: 1, offset: 12970},
			expr: &actionExpr{
				pos: position{line: 412, col: 13, offset: 12982},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 412, col: 13, offset: 12982},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 412, col: 13, offset: 12982},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 412, col: 21, offset: 12990},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 412, col: 24, offset: 12993},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 412, col: 28, offset: 12997},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 412, col: 38, offset: 13007},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 412, col: 41, offset: 13010},
							val:        ">",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 412, col: 45, offset: 13014},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 412, col: 47, offset: 13016},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 412, col: 59, offset: 13028},
								expr: &ruleRefExpr{
									pos:  position{line: 412, col: 59, offset: 13028},
									name: "TypeAnnotations",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 420, col: 1, offset: 13191},
			expr: &actionExpr{
				pos: position{line: 420, col: 12, offset: 13202},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 420, col: 12, offset: 13202},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 420, col: 12, offset: 13202},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 420, col: 23, offset: 13213},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 420, col: 31, offset: 13221},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 424, col: 1, offset: 13258},
			expr: &choiceExpr{
				pos: position{line: 424, col: 15, offset: 13272},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 424, col: 15, offset: 13272},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 25, offset: 13282},
						name: "BoolConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 40, offset: 13297},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 57, offset: 13314},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 71, offset: 13328},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 82, offset: 13339},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 424, col: 94, offset: 13351},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "TypeAnnotations",
			pos:  position{line: 426, col: 1, offset: 13363},
			expr: &actionExpr{
				pos: position{line: 426, col: 20, offset: 13382},
				run: (*parser).callonTypeAnnotations1,
				expr: &seqExpr{
					pos: position{line: 426, col: 20, offset: 13382},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 426, col: 20, offset: 13382},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 426, col: 24, offset: 13386},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 426, col: 27, offset: 13389},
							label: "annotations",
							expr: &zeroOrMoreExpr{
								pos: position{line: 426, col: 39, offset: 13401},
								expr: &ruleRefExpr{
									pos:  position{line: 426, col: 39, offset: 13401},
									name: "TypeAnnotation",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 426, col: 55, offset: 13417},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "TypeAnnotation",
			pos:  position{line: 434, col: 1, offset: 13581},
			expr: &actionExpr{
				pos: position{line: 434, col: 19, offset: 13599},
				run: (*parser).callonTypeAnnotation1,
				expr: &seqExpr{
					pos: position{line: 434, col: 19, offset: 13599},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 434, col: 19, offset: 13599},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 434, col: 24, offset: 13604},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 434, col: 35, offset: 13615},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 434, col: 37, offset: 13617},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 434, col: 43, offset: 13623},
								expr: &actionExpr{
									pos: position{line: 434, col: 44, offset: 13624},
									run: (*parser).callonTypeAnnotation8,
									expr: &seqExpr{
										pos: position{line: 434, col: 44, offset: 13624},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 434, col: 44, offset: 13624},
												val:        "=",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 434, col: 48, offset: 13628},
												name: "__",
											},
											&labeledExpr{
												pos:   position{line: 434, col: 51, offset: 13631},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 434, col: 57, offset: 13637},
													name: "Literal",
												},
											},
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 434, col: 89, offset: 13669},
							expr: &ruleRefExpr{
								pos:  position{line: 434, col: 89, offset: 13669},
								name: "ListSeparator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 434, col: 104, offset: 13684},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "BoolConstant",
			pos:  position{line: 445, col: 1, offset: 13880},
			expr: &actionExpr{
				pos: position{line: 445, col: 17, offset: 13896},
				run: (*parser).callonBoolConstant1,
				expr: &choiceExpr{
					pos: position{line: 445, col: 18, offset: 13897},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 445, col: 18, offset: 13897},
							val:        "true",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 445, col: 27, offset: 13906},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 449, col: 1, offset: 13961},
			expr: &actionExpr{
				pos: position{line: 449, col: 16, offset: 13976},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 449, col: 16, offset: 13976},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 449, col: 16, offset: 13976},
							expr: &charClassMatcher{
								pos:        position{line: 449, col: 16, offset: 13976},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 449, col: 22, offset: 13982},
							expr: &ruleRefExpr{
								pos:  position{line: 449, col: 22, offset: 13982},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 453, col: 1, offset: 14046},
			expr: &actionExpr{
				pos: position{line: 453, col: 19, offset: 14064},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 453, col: 19, offset: 14064},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 453, col: 19, offset: 14064},
							expr: &charClassMatcher{
								pos:        position{line: 453, col: 19, offset: 14064},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 453, col: 25, offset: 14070},
							expr: &ruleRefExpr{
								pos:  position{line: 453, col: 25, offset: 14070},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 453, col: 32, offset: 14077},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 453, col: 36, offset: 14081},
							expr: &ruleRefExpr{
								pos:  position{line: 453, col: 36, offset: 14081},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 453, col: 43, offset: 14088},
							expr: &seqExpr{
								pos: position{line: 453, col: 45, offset: 14090},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 453, col: 45, offset: 14090},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 453, col: 52, offset: 14097},
										name: "IntConstant",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ConstList",
			pos:  position{line: 457, col: 1, offset: 14167},
			expr: &actionExpr{
				pos: position{line: 457, col: 14, offset: 14180},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 457, col: 14, offset: 14180},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 457, col: 14, offset: 14180},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 457, col: 18, offset: 14184},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 457, col: 21, offset: 14187},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 457, col: 28, offset: 14194},
								expr: &seqExpr{
									pos: position{line: 457, col: 29, offset: 14195},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 457, col: 29, offset: 14195},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 457, col: 40, offset: 14206},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 457, col: 43, offset: 14209},
											expr: &ruleRefExpr{
												pos:  position{line: 457, col: 43, offset: 14209},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 457, col: 58, offset: 14224},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 457, col: 63, offset: 14229},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 457, col: 66, offset: 14232},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 466, col: 1, offset: 14426},
			expr: &actionExpr{
				pos: position{line: 466, col: 13, offset: 14438},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 466, col: 13, offset: 14438},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 466, col: 13, offset: 14438},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 466, col: 17, offset: 14442},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 466, col: 20, offset: 14445},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 466, col: 27, offset: 14452},
								expr: &seqExpr{
									pos: position{line: 466, col: 28, offset: 14453},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 466, col: 28, offset: 14453},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 466, col: 39, offset: 14464},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 466, col: 42, offset: 14467},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 466, col: 46, offset: 14471},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 466, col: 49, offset: 14474},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 466, col: 60, offset: 14485},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 466, col: 64, offset: 14489},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 466, col: 64, offset: 14489},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 466, col: 70, offset: 14495},
													expr: &litMatcher{
														pos:        position{line: 466, col: 71, offset: 14496},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 466, col: 76, offset: 14501},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 466, col: 81, offset: 14506},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Scope",
			pos:  position{line: 486, col: 1, offset: 15056},
			expr: &actionExpr{
				pos: position{line: 486, col: 10, offset: 15065},
				run: (*parser).callonScope1,
				expr: &seqExpr{
					pos: position{line: 486, col: 10, offset: 15065},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 486, col: 10, offset: 15065},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 486, col: 17, offset: 15072},
								expr: &seqExpr{
									pos: position{line: 486, col: 18, offset: 15073},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 486, col: 18, offset: 15073},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 486, col: 28, offset: 15083},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 486, col: 33, offset: 15088},
							val:        "scope",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 41, offset: 15096},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 486, col: 44, offset: 15099},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 486, col: 49, offset: 15104},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 60, offset: 15115},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 486, col: 63, offset: 15118},
							label: "prefix",
							expr: &zeroOrOneExpr{
								pos: position{line: 486, col: 70, offset: 15125},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 70, offset: 15125},
									name: "Prefix",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 78, offset: 15133},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 486, col: 81, offset: 15136},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 85, offset: 15140},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 486, col: 88, offset: 15143},
							label: "operations",
							expr: &zeroOrMoreExpr{
								pos: position{line: 486, col: 99, offset: 15154},
								expr: &seqExpr{
									pos: position{line: 486, col: 100, offset: 15155},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 486, col: 100, offset: 15155},
											name: "Operation",
										},
										&ruleRefExpr{
											pos:  position{line: 486, col: 110, offset: 15165},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 486, col: 116, offset: 15171},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 486, col: 116, offset: 15171},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 486, col: 122, offset: 15177},
									name: "EndOfScopeError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 139, offset: 15194},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 486, col: 141, offset: 15196},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 486, col: 153, offset: 15208},
								expr: &ruleRefExpr{
									pos:  position{line: 486, col: 153, offset: 15208},
									name: "TypeAnnotations",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 486, col: 170, offset: 15225},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfScopeError",
			pos:  position{line: 508, col: 1, offset: 15822},
			expr: &actionExpr{
				pos: position{line: 508, col: 20, offset: 15841},
				run: (*parser).callonEndOfScopeError1,
				expr: &anyMatcher{
					line: 508, col: 20, offset: 15841,
				},
			},
		},
		{
			name: "Prefix",
			pos:  position{line: 512, col: 1, offset: 15908},
			expr: &actionExpr{
				pos: position{line: 512, col: 11, offset: 15918},
				run: (*parser).callonPrefix1,
				expr: &seqExpr{
					pos: position{line: 512, col: 11, offset: 15918},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 512, col: 11, offset: 15918},
							val:        "prefix",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 512, col: 20, offset: 15927},
							name: "__",
						},
						&ruleRefExpr{
							pos:  position{line: 512, col: 23, offset: 15930},
							name: "PrefixToken",
						},
						&zeroOrMoreExpr{
							pos: position{line: 512, col: 35, offset: 15942},
							expr: &seqExpr{
								pos: position{line: 512, col: 36, offset: 15943},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 512, col: 36, offset: 15943},
										val:        ".",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 512, col: 40, offset: 15947},
										name: "PrefixToken",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PrefixToken",
			pos:  position{line: 517, col: 1, offset: 16078},
			expr: &choiceExpr{
				pos: position{line: 517, col: 16, offset: 16093},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 517, col: 17, offset: 16094},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 517, col: 17, offset: 16094},
								val:        "{",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 517, col: 21, offset: 16098},
								name: "PrefixWord",
							},
							&litMatcher{
								pos:        position{line: 517, col: 32, offset: 16109},
								val:        "}",
								ignoreCase: false,
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 517, col: 39, offset: 16116},
						name: "PrefixWord",
					},
				},
			},
		},
		{
			name: "PrefixWord",
			pos:  position{line: 519, col: 1, offset: 16128},
			expr: &oneOrMoreExpr{
				pos: position{line: 519, col: 15, offset: 16142},
				expr: &charClassMatcher{
					pos:        position{line: 519, col: 15, offset: 16142},
					val:        "[^\\r\\n\\t\\f .{}]",
					chars:      []rune{'\r', '\n', '\t', '\f', ' ', '.', '{', '}'},
					ignoreCase: false,
					inverted:   true,
				},
			},
		},
		{
			name: "Operation",
			pos:  position{line: 521, col: 1, offset: 16160},
			expr: &actionExpr{
				pos: position{line: 521, col: 14, offset: 16173},
				run: (*parser).callonOperation1,
				expr: &seqExpr{
					pos: position{line: 521, col: 14, offset: 16173},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 521, col: 14, offset: 16173},
							label: "docstr",
							expr: &zeroOrOneExpr{
								pos: position{line: 521, col: 21, offset: 16180},
								expr: &seqExpr{
									pos: position{line: 521, col: 22, offset: 16181},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 521, col: 22, offset: 16181},
											name: "DocString",
										},
										&ruleRefExpr{
											pos:  position{line: 521, col: 32, offset: 16191},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 521, col: 37, offset: 16196},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 521, col: 42, offset: 16201},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 521, col: 53, offset: 16212},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 521, col: 55, offset: 16214},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 521, col: 59, offset: 16218},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 521, col: 62, offset: 16221},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 521, col: 66, offset: 16225},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 521, col: 76, offset: 16235},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 521, col: 78, offset: 16237},
							label: "annotations",
							expr: &zeroOrOneExpr{
								pos: position{line: 521, col: 90, offset: 16249},
								expr: &ruleRefExpr{
									pos:  position{line: 521, col: 90, offset: 16249},
									name: "TypeAnnotations",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 521, col: 107, offset: 16266},
							expr: &ruleRefExpr{
								pos:  position{line: 521, col: 107, offset: 16266},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 538, col: 1, offset: 16826},
			expr: &actionExpr{
				pos: position{line: 538, col: 12, offset: 16837},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 538, col: 13, offset: 16838},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 538, col: 14, offset: 16839},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 538, col: 14, offset: 16839},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 538, col: 18, offset: 16843},
									expr: &choiceExpr{
										pos: position{line: 538, col: 19, offset: 16844},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 538, col: 19, offset: 16844},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 538, col: 26, offset: 16851},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 538, col: 33, offset: 16858},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 538, col: 41, offset: 16866},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 538, col: 41, offset: 16866},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 538, col: 46, offset: 16871},
									expr: &choiceExpr{
										pos: position{line: 538, col: 47, offset: 16872},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 538, col: 47, offset: 16872},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 538, col: 54, offset: 16879},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 538, col: 61, offset: 16886},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 547, col: 1, offset: 17172},
			expr: &actionExpr{
				pos: position{line: 547, col: 15, offset: 17186},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 547, col: 15, offset: 17186},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 547, col: 15, offset: 17186},
							expr: &choiceExpr{
								pos: position{line: 547, col: 16, offset: 17187},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 547, col: 16, offset: 17187},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 547, col: 25, offset: 17196},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 547, col: 31, offset: 17202},
							expr: &choiceExpr{
								pos: position{line: 547, col: 32, offset: 17203},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 547, col: 32, offset: 17203},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 547, col: 41, offset: 17212},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 547, col: 49, offset: 17220},
										val:        "[._]",
										chars:      []rune{'.', '_'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListSeparator",
			pos:  position{line: 551, col: 1, offset: 17275},
			expr: &charClassMatcher{
				pos:        position{line: 551, col: 18, offset: 17292},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 552, col: 1, offset: 17297},
			expr: &charClassMatcher{
				pos:        position{line: 552, col: 11, offset: 17307},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 553, col: 1, offset: 17316},
			expr: &charClassMatcher{
				pos:        position{line: 553, col: 10, offset: 17325},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 555, col: 1, offset: 17332},
			expr: &anyMatcher{
				line: 555, col: 15, offset: 17346,
			},
		},
		{
			name: "DocString",
			pos:  position{line: 556, col: 1, offset: 17348},
			expr: &actionExpr{
				pos: position{line: 556, col: 14, offset: 17361},
				run: (*parser).callonDocString1,
				expr: &seqExpr{
					pos: position{line: 556, col: 14, offset: 17361},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 556, col: 14, offset: 17361},
							val:        "/**@",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 556, col: 21, offset: 17368},
							expr: &seqExpr{
								pos: position{line: 556, col: 23, offset: 17370},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 556, col: 23, offset: 17370},
										expr: &litMatcher{
											pos:        position{line: 556, col: 24, offset: 17371},
											val:        "*/",
											ignoreCase: false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 556, col: 29, offset: 17376},
										name: "SourceChar",
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 556, col: 43, offset: 17390},
							val:        "*/",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 562, col: 1, offset: 17570},
			expr: &choiceExpr{
				pos: position{line: 562, col: 12, offset: 17581},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 562, col: 12, offset: 17581},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 562, col: 31, offset: 17600},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 563, col: 1, offset: 17618},
			expr: &seqExpr{
				pos: position{line: 563, col: 21, offset: 17638},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 563, col: 21, offset: 17638},
						expr: &ruleRefExpr{
							pos:  position{line: 563, col: 22, offset: 17639},
							name: "DocString",
						},
					},
					&litMatcher{
						pos:        position{line: 563, col: 32, offset: 17649},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 563, col: 37, offset: 17654},
						expr: &seqExpr{
							pos: position{line: 563, col: 39, offset: 17656},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 563, col: 39, offset: 17656},
									expr: &litMatcher{
										pos:        position{line: 563, col: 40, offset: 17657},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 563, col: 45, offset: 17662},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 563, col: 59, offset: 17676},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 564, col: 1, offset: 17681},
			expr: &seqExpr{
				pos: position{line: 564, col: 37, offset: 17717},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 564, col: 37, offset: 17717},
						expr: &ruleRefExpr{
							pos:  position{line: 564, col: 38, offset: 17718},
							name: "DocString",
						},
					},
					&litMatcher{
						pos:        position{line: 564, col: 48, offset: 17728},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 564, col: 53, offset: 17733},
						expr: &seqExpr{
							pos: position{line: 564, col: 55, offset: 17735},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 564, col: 55, offset: 17735},
									expr: &choiceExpr{
										pos: position{line: 564, col: 58, offset: 17738},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 564, col: 58, offset: 17738},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 564, col: 65, offset: 17745},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 564, col: 71, offset: 17751},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 564, col: 85, offset: 17765},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 565, col: 1, offset: 17770},
			expr: &choiceExpr{
				pos: position{line: 565, col: 22, offset: 17791},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 565, col: 23, offset: 17792},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 565, col: 23, offset: 17792},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 565, col: 28, offset: 17797},
								expr: &seqExpr{
									pos: position{line: 565, col: 30, offset: 17799},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 565, col: 30, offset: 17799},
											expr: &ruleRefExpr{
												pos:  position{line: 565, col: 31, offset: 17800},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 565, col: 35, offset: 17804},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 565, col: 53, offset: 17822},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 565, col: 53, offset: 17822},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 565, col: 57, offset: 17826},
								expr: &seqExpr{
									pos: position{line: 565, col: 59, offset: 17828},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 565, col: 59, offset: 17828},
											expr: &ruleRefExpr{
												pos:  position{line: 565, col: 60, offset: 17829},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 565, col: 64, offset: 17833},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 567, col: 1, offset: 17849},
			expr: &zeroOrMoreExpr{
				pos: position{line: 567, col: 7, offset: 17855},
				expr: &choiceExpr{
					pos: position{line: 567, col: 9, offset: 17857},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 567, col: 9, offset: 17857},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 567, col: 22, offset: 17870},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 567, col: 28, offset: 17876},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 568, col: 1, offset: 17887},
			expr: &zeroOrMoreExpr{
				pos: position{line: 568, col: 6, offset: 17892},
				expr: &choiceExpr{
					pos: position{line: 568, col: 8, offset: 17894},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 568, col: 8, offset: 17894},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 568, col: 21, offset: 17907},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 569, col: 1, offset: 17943},
			expr: &zeroOrMoreExpr{
				pos: position{line: 569, col: 7, offset: 17949},
				expr: &ruleRefExpr{
					pos:  position{line: 569, col: 7, offset: 17949},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 571, col: 1, offset: 17962},
			expr: &charClassMatcher{
				pos:        position{line: 571, col: 15, offset: 17976},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 572, col: 1, offset: 17984},
			expr: &litMatcher{
				pos:        position{line: 572, col: 8, offset: 17991},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 573, col: 1, offset: 17996},
			expr: &choiceExpr{
				pos: position{line: 573, col: 8, offset: 18003},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 573, col: 8, offset: 18003},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 573, col: 8, offset: 18003},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 573, col: 11, offset: 18006},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 573, col: 17, offset: 18012},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 573, col: 17, offset: 18012},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 573, col: 19, offset: 18014},
								expr: &ruleRefExpr{
									pos:  position{line: 573, col: 19, offset: 18014},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 573, col: 38, offset: 18033},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 573, col: 44, offset: 18039},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 573, col: 44, offset: 18039},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 573, col: 47, offset: 18042},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 575, col: 1, offset: 18047},
			expr: &notExpr{
				pos: position{line: 575, col: 8, offset: 18054},
				expr: &anyMatcher{
					line: 575, col: 9, offset: 18055,
				},
			},
		},
	},
}

func (c *current) onGrammar1(statements interface{}) (interface{}, error) {
	stmts := toIfaceSlice(statements)
	frugal := &Frugal{
		Scopes:         []*Scope{},
		ParsedIncludes: make(map[string]*Frugal),
		Includes:       []*Include{},
		Namespaces:     []*Namespace{},
		Typedefs:       []*TypeDef{},
		Constants:      []*Constant{},
		Enums:          []*Enum{},
		Structs:        []*Struct{},
		Exceptions:     []*Struct{},
		Unions:         []*Struct{},
		Services:       []*Service{},
		typedefIndex:   make(map[string]*TypeDef),
		namespaceIndex: make(map[string]*Namespace),
	}

	for _, st := range stmts {
		wrapper := st.([]interface{})[0].(*statementWrapper)
		switch v := wrapper.statement.(type) {
		case *Namespace:
			frugal.Namespaces = append(frugal.Namespaces, v)
			frugal.namespaceIndex[v.Scope] = v
		case *Constant:
			v.Comment = wrapper.comment
			frugal.Constants = append(frugal.Constants, v)
		case *Enum:
			v.Comment = wrapper.comment
			frugal.Enums = append(frugal.Enums, v)
		case *TypeDef:
			v.Comment = wrapper.comment
			frugal.Typedefs = append(frugal.Typedefs, v)
			frugal.typedefIndex[v.Name] = v
		case *Struct:
			v.Type = StructTypeStruct
			v.Comment = wrapper.comment
			frugal.Structs = append(frugal.Structs, v)
		case exception:
			strct := (*Struct)(v)
			strct.Type = StructTypeException
			strct.Comment = wrapper.comment
			frugal.Exceptions = append(frugal.Exceptions, strct)
		case union:
			strct := unionToStruct(v)
			strct.Type = StructTypeUnion
			strct.Comment = wrapper.comment
			frugal.Unions = append(frugal.Unions, strct)
		case *Service:
			v.Comment = wrapper.comment
			v.Frugal = frugal
			frugal.Services = append(frugal.Services, v)
		case *Include:
			frugal.Includes = append(frugal.Includes, v)
		case *Scope:
			v.Comment = wrapper.comment
			v.Frugal = frugal
			frugal.Scopes = append(frugal.Scopes, v)
		default:
			return nil, fmt.Errorf("parser: unknown value %#v", v)
		}
	}
	return frugal, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["statements"])
}

func (c *current) onSyntaxError1() (interface{}, error) {
	return nil, errors.New("parser: syntax error")
}

func (p *parser) callonSyntaxError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSyntaxError1()
}

func (c *current) onStatement1(docstr, statement interface{}) (interface{}, error) {
	wrapper := &statementWrapper{statement: statement}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		wrapper.comment = rawCommentToDocStr(raw)
	}
	return wrapper, nil
}

func (p *parser) callonStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement1(stack["docstr"], stack["statement"])
}

func (c *current) onInclude1(file, annotations interface{}) (interface{}, error) {
	name := filepath.Base(file.(string))
	if ix := strings.LastIndex(name, "."); ix > 0 {
		name = name[:ix]
	}
	return &Include{
		Name:        name,
		Value:       file.(string),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonInclude1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInclude1(stack["file"], stack["annotations"])
}

func (c *current) onNamespace1(scope, ns, annotations interface{}) (interface{}, error) {
	return &Namespace{
		Scope:       ifaceSliceToString(scope),
		Value:       string(ns.(Identifier)),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonNamespace1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespace1(stack["scope"], stack["ns"], stack["annotations"])
}

func (c *current) onConst1(typ, name, value, annotations interface{}) (interface{}, error) {
	return &Constant{
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Value:       value,
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonConst1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst1(stack["typ"], stack["name"], stack["value"], stack["annotations"])
}

func (c *current) onEnum1(name, values, annotations interface{}) (interface{}, error) {
	vs := toIfaceSlice(values)
	en := &Enum{
		Name:        string(name.(Identifier)),
		Values:      make([]*EnumValue, len(vs)),
		Annotations: toAnnotations(annotations),
	}
	// Assigns numbers in order. This will behave badly if some values are
	// defined and other are not, but I think that's ok since that's a silly
	// thing to do.
	next := 0
	for idx, v := range vs {
		ev := v.([]interface{})[0].(*EnumValue)
		if ev.Value < 0 {
			ev.Value = next
		}
		if ev.Value >= next {
			next = ev.Value + 1
		}
		en.Values[idx] = ev
	}
	return en, nil
}

func (p *parser) callonEnum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnum1(stack["name"], stack["values"], stack["annotations"])
}

func (c *current) onEnumValue1(docstr, name, value, annotations interface{}) (interface{}, error) {
	ev := &EnumValue{
		Name:        string(name.(Identifier)),
		Value:       -1,
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		ev.Comment = rawCommentToDocStr(raw)
	}
	if value != nil {
		ev.Value = int(value.([]interface{})[2].(int64))
	}
	return ev, nil
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1(stack["docstr"], stack["name"], stack["value"], stack["annotations"])
}

func (c *current) onTypeDef1(typ, name, annotations interface{}) (interface{}, error) {
	return &TypeDef{
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonTypeDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeDef1(stack["typ"], stack["name"], stack["annotations"])
}

func (c *current) onStruct1(st interface{}) (interface{}, error) {
	return st.(*Struct), nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["st"])
}

func (c *current) onException1(st interface{}) (interface{}, error) {
	return exception(st.(*Struct)), nil
}

func (p *parser) callonException1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onException1(stack["st"])
}

func (c *current) onUnion1(st interface{}) (interface{}, error) {
	return union(st.(*Struct)), nil
}

func (p *parser) callonUnion1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnion1(stack["st"])
}

func (c *current) onStructLike1(name, fields, annotations interface{}) (interface{}, error) {
	st := &Struct{
		Name:        string(name.(Identifier)),
		Annotations: toAnnotations(annotations),
	}
	if fields != nil {
		st.Fields = fields.([]*Field)
	}
	return st, nil
}

func (p *parser) callonStructLike1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructLike1(stack["name"], stack["fields"], stack["annotations"])
}

func (c *current) onFieldList1(fields interface{}) (interface{}, error) {
	fs := fields.([]interface{})
	flds := make([]*Field, len(fs))
	for i, f := range fs {
		flds[i] = f.([]interface{})[0].(*Field)
	}
	return flds, nil
}

func (p *parser) callonFieldList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldList1(stack["fields"])
}

func (c *current) onField1(docstr, id, mod, typ, name, def, annotations interface{}) (interface{}, error) {
	f := &Field{
		ID:          int(id.(int64)),
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		f.Comment = rawCommentToDocStr(raw)
	}
	if mod != nil {
		f.Modifier = mod.(FieldModifier)
	} else {
		f.Modifier = Default
	}

	if def != nil {
		f.Default = def.([]interface{})[2]
	}
	return f, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["docstr"], stack["id"], stack["mod"], stack["typ"], stack["name"], stack["def"], stack["annotations"])
}

func (c *current) onFieldModifier1() (interface{}, error) {
	if bytes.Equal(c.text, []byte("required")) {
		return Required, nil
	} else {
		return Optional, nil
	}
}

func (p *parser) callonFieldModifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldModifier1()
}

func (c *current) onService1(name, extends, methods, annotations interface{}) (interface{}, error) {
	ms := methods.([]interface{})
	svc := &Service{
		Name:        string(name.(Identifier)),
		Methods:     make([]*Method, len(ms)),
		Annotations: toAnnotations(annotations),
	}
	if extends != nil {
		svc.Extends = string(extends.([]interface{})[2].(Identifier))
	}
	for i, m := range ms {
		mt := m.([]interface{})[0].(*Method)
		svc.Methods[i] = mt
	}
	return svc, nil
}

func (p *parser) callonService1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onService1(stack["name"], stack["extends"], stack["methods"], stack["annotations"])
}

func (c *current) onEndOfServiceError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of service")
}

func (p *parser) callonEndOfServiceError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfServiceError1()
}

func (c *current) onFunction1(docstr, oneway, typ, name, arguments, exceptions, annotations interface{}) (interface{}, error) {
	m := &Method{
		Name:        string(name.(Identifier)),
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		m.Comment = rawCommentToDocStr(raw)
	}
	t := typ.(*Type)
	if t.Name != "void" {
		m.ReturnType = t
	}
	if oneway != nil {
		m.Oneway = true
	}
	if arguments != nil {
		m.Arguments = arguments.([]*Field)
	}
	if exceptions != nil {
		m.Exceptions = exceptions.([]*Field)
		for _, e := range m.Exceptions {
			e.Modifier = Optional
		}
	}
	return m, nil
}

func (p *parser) callonFunction1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunction1(stack["docstr"], stack["oneway"], stack["typ"], stack["name"], stack["arguments"], stack["exceptions"], stack["annotations"])
}

func (c *current) onFunctionType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(*Type); ok {
		return t, nil
	}
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonFunctionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionType1(stack["typ"])
}

func (c *current) onThrows1(exceptions interface{}) (interface{}, error) {
	return exceptions, nil
}

func (p *parser) callonThrows1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onThrows1(stack["exceptions"])
}

func (c *current) onFieldType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(Identifier); ok {
		return &Type{Name: string(t)}, nil
	}
	return typ, nil
}

func (p *parser) callonFieldType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldType1(stack["typ"])
}

func (c *current) onBaseType1(name, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        name.(string),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonBaseType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseType1(stack["name"], stack["annotations"])
}

func (c *current) onBaseTypeName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonBaseTypeName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseTypeName1()
}

func (c *current) onContainerType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonContainerType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContainerType1(stack["typ"])
}

func (c *current) onMapType1(key, value, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "map",
		KeyType:     key.(*Type),
		ValueType:   value.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonMapType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapType1(stack["key"], stack["value"], stack["annotations"])
}

func (c *current) onSetType1(typ, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "set",
		ValueType:   typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonSetType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetType1(stack["typ"], stack["annotations"])
}

func (c *current) onListType1(typ, annotations interface{}) (interface{}, error) {
	return &Type{
		Name:        "list",
		ValueType:   typ.(*Type),
		Annotations: toAnnotations(annotations),
	}, nil
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["typ"], stack["annotations"])
}

func (c *current) onCppType1(cppType interface{}) (interface{}, error) {
	return cppType, nil
}

func (p *parser) callonCppType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCppType1(stack["cppType"])
}

func (c *current) onTypeAnnotations1(annotations interface{}) (interface{}, error) {
	var anns []*Annotation
	for _, ann := range annotations.([]interface{}) {
		anns = append(anns, ann.(*Annotation))
	}
	return anns, nil
}

func (p *parser) callonTypeAnnotations1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotations1(stack["annotations"])
}

func (c *current) onTypeAnnotation8(value interface{}) (interface{}, error) {
	return value, nil
}

func (p *parser) callonTypeAnnotation8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotation8(stack["value"])
}

func (c *current) onTypeAnnotation1(name, value interface{}) (interface{}, error) {
	var optValue string
	if value != nil {
		optValue = value.(string)
	}
	return &Annotation{
		Name:  string(name.(Identifier)),
		Value: optValue,
	}, nil
}

func (p *parser) callonTypeAnnotation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeAnnotation1(stack["name"], stack["value"])
}

func (c *current) onBoolConstant1() (interface{}, error) {
	return string(c.text) == "true", nil
}

func (p *parser) callonBoolConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolConstant1()
}

func (c *current) onIntConstant1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIntConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntConstant1()
}

func (c *current) onDoubleConstant1() (interface{}, error) {
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonDoubleConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleConstant1()
}

func (c *current) onConstList1(values interface{}) (interface{}, error) {
	valueSlice := values.([]interface{})
	vs := make([]interface{}, len(valueSlice))
	for i, v := range valueSlice {
		vs[i] = v.([]interface{})[0]
	}
	return vs, nil
}

func (p *parser) callonConstList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstList1(stack["values"])
}

func (c *current) onConstMap1(values interface{}) (interface{}, error) {
	if values == nil {
		return nil, nil
	}
	vals := values.([]interface{})
	kvs := make([]KeyValue, len(vals))
	for i, kv := range vals {
		v := kv.([]interface{})
		kvs[i] = KeyValue{
			Key:   v[0],
			Value: v[4],
		}
	}
	return kvs, nil
}

func (p *parser) callonConstMap1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstMap1(stack["values"])
}

func (c *current) onScope1(docstr, name, prefix, operations, annotations interface{}) (interface{}, error) {
	ops := operations.([]interface{})
	scope := &Scope{
		Name:        string(name.(Identifier)),
		Operations:  make([]*Operation, len(ops)),
		Prefix:      defaultPrefix,
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		scope.Comment = rawCommentToDocStr(raw)
	}
	if prefix != nil {
		scope.Prefix = prefix.(*ScopePrefix)
	}
	for i, o := range ops {
		op := o.([]interface{})[0].(*Operation)
		scope.Operations[i] = op
	}
	return scope, nil
}

func (p *parser) callonScope1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onScope1(stack["docstr"], stack["name"], stack["prefix"], stack["operations"], stack["annotations"])
}

func (c *current) onEndOfScopeError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of scope")
}

func (p *parser) callonEndOfScopeError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfScopeError1()
}

func (c *current) onPrefix1() (interface{}, error) {
	prefix := strings.TrimSpace(strings.TrimPrefix(string(c.text), "prefix"))
	return newScopePrefix(prefix)
}

func (p *parser) callonPrefix1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefix1()
}

func (c *current) onOperation1(docstr, name, typ, annotations interface{}) (interface{}, error) {
	o := &Operation{
		Name:        string(name.(Identifier)),
		Type:        typ.(*Type),
		Annotations: toAnnotations(annotations),
	}
	if docstr != nil {
		raw := docstr.([]interface{})[0].(string)
		o.Comment = rawCommentToDocStr(raw)
	}
	return o, nil
}

func (p *parser) callonOperation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperation1(stack["docstr"], stack["name"], stack["typ"], stack["annotations"])
}

func (c *current) onLiteral1() (interface{}, error) {
	if len(c.text) != 0 && c.text[0] == '\'' {
		intermediate := strings.Replace(string(c.text[1:len(c.text)-1]), `\'`, `'`, -1)
		return strconv.Unquote(`"` + strings.Replace(intermediate, `"`, `\"`, -1) + `"`)
	}

	return strconv.Unquote(string(c.text))
}

func (p *parser) callonLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral1()
}

func (c *current) onIdentifier1() (interface{}, error) {
	return Identifier(string(c.text)), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onDocString1() (interface{}, error) {
	comment := string(c.text)
	comment = strings.TrimPrefix(comment, "/**@")
	comment = strings.TrimSuffix(comment, "*/")
	return strings.TrimSpace(comment), nil
}

func (p *parser) callonDocString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDocString1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
