#!/usr/bin/env bash
set -e

# Wrap up package for pub
cd $FRUGAL_HOME
tar -C lib/dart -czf $FRUGAL_HOME/frugal.pub.tgz .

# Compile library code
cd $FRUGAL_HOME/lib/dart
timeout 5m pub get

#generate test runner
pub run dart_dev gen-test-runner

# Run the tests
pub run dart_dev test

pub run dart_dev format --check
