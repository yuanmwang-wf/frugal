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

package test

import (
	"bufio"
	"flag"
	"io"
	"os"
	"testing"
)

const (
	outputDir               = "out"
	delim                   = "."
	validFile               = "idl/valid.frugal"
	invalidFile             = "idl/invalid.frugal"
	duplicateServices       = "idl/duplicate_services.frugal"
	duplicateScopes         = "idl/duplicate_scopes.frugal"
	duplicateMethods        = "idl/duplicate_methods.frugal"
	duplicateOperations     = "idl/duplicate_operations.frugal"
	duplicateMethodArgIds   = "idl/duplicate_arg_ids.frugal"
	duplicateStructFieldIds = "idl/duplicate_field_ids.frugal"
	frugalGenFile           = "idl/variety.frugal"
	badNamespace            = "idl/bad_namespace.frugal"
	includeVendor           = "idl/include_vendor.frugal"
	includeVendorNoPath     = "idl/include_vendor_no_path.frugal"
	vendorNamespace         = "idl/vendor_namespace.frugal"
)

var copyFiles bool

func init() {
	copyFilesPtr := flag.Bool("copy-files", false, "")
	flag.Parse()
	copyFiles = *copyFilesPtr
}

type FileComparisonPair struct {
	ExpectedPath  string
	GeneratedPath string
}

func compareFiles(t *testing.T, expectedPath, generatedPath string) {
	expected, err := os.Open(expectedPath)
	if err != nil {
		t.Fatalf("Failed to open file %s", expectedPath)
	}
	defer expected.Close()
	generated, err := os.Open(generatedPath)
	if err != nil {
		t.Fatalf("Failed to open file %s", generatedPath)
	}
	defer generated.Close()

	expectedScanner := bufio.NewScanner(expected)
	generatedScanner := bufio.NewScanner(generated)
	line := 1
	for expectedScanner.Scan() {
		generatedScanner.Scan()
		expectedLine := expectedScanner.Text()
		generatedLine := generatedScanner.Text()
		if expectedLine != generatedLine {
			t.Fatalf("\nExpected line \n<%s> (%s)\n generated line\n <%s> (%s) at line %d",
				expectedLine, expectedPath, generatedLine, generatedPath, line)
		}
		line++
	}

	if generatedScanner.Scan() {
		t.Fatalf("Generated has more lines than expected: exp %s gen %s", expectedPath, generatedPath)
	}
}

func compareAllFiles(t *testing.T, pairs []FileComparisonPair) {
	for _, pair := range pairs {
		compareFiles(t, pair.ExpectedPath, pair.GeneratedPath)
	}
}

func copyAllFiles(t *testing.T, pairs []FileComparisonPair) {
	if !copyFiles {
		return
	}

	for _, pair := range pairs {
		if err := copyFilePair(pair); err != nil {
			t.Fatal(err)
		}
	}
}

func copyFilePair(pair FileComparisonPair) error {
	// TODO automatically create a missing expected file?

	generatedFile, err := os.Open(pair.GeneratedPath)
	if err != nil {
		return err
	}
	defer generatedFile.Close()

	expectedFile, err := os.Create(pair.ExpectedPath)
	if err != nil {
		return err
	}
	defer expectedFile.Close()

	_, err = io.Copy(expectedFile, generatedFile)
	return err
}

func assertFilesNotExist(t *testing.T, filePaths []string) {
	for _, fileThatShouldNotExist := range filePaths {
		if _, err := os.Stat(fileThatShouldNotExist); !os.IsNotExist(err) {
			if err != nil {
				t.Errorf("Unexpected error checking for existence on %q: %s", fileThatShouldNotExist, err)
			} else {
				t.Errorf("Expected %q not to exist, but it did", fileThatShouldNotExist)
			}
		}
	}
}
