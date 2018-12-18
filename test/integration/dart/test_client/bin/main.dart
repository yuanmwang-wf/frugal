/// Licensed to the Apache Software Foundation (ASF) under one
/// or more contributor license agreements. See the NOTICE file
/// distributed with this work for additional information
/// regarding copyright ownership. The ASF licenses this file
/// to you under the Apache License, Version 2.0 (the
/// 'License'); you may not use this file except in compliance
/// with the License. You may obtain a copy of the License at
///
/// http://www.apache.org/licenses/LICENSE-2.0
///
/// Unless required by applicable law or agreed to in writing,
/// software distributed under the License is distributed on an
/// 'AS IS' BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
/// KIND, either express or implied. See the License for the
/// specific language governing permissions and limitations
/// under the License.

import 'dart:async';
import 'dart:io';

import 'dart:typed_data';
import 'package:args/args.dart';
import 'package:thrift/thrift.dart';
import 'package:frugal/frugal.dart';
import 'package:frugal_test_client/test_cases.dart';
import 'package:frugal_test/frugal_test.dart';
import 'package:w_transport/w_transport.dart' as wt;
// ignore: deprecated_member_use
import 'package:w_transport/w_transport_vm.dart' show configureWTransportForVM;

List<FTest> _tests;
FFrugalTestClient client;
bool verbose;
var middleware_called = false;

main(List<String> args) async {
  configureWTransportForVM();
  ArgResults results = _parseArgs(args);

  if (results == null) {
    exit(1);
  }

  verbose = results['verbose'] == true;

  await _initTestClient(
      host: results['host'],
      port: int.parse(results['port']),
      transportType: results['transport'],
      protocolType: results['protocol']).catchError((e) {
    stdout.writeln('Error:');
    stdout.writeln('$e');
    if (e is Error) {
      stdout.writeln('${e.stackTrace}');
    }
    exit(1);
  });

  // run tests
  int result = 0;
  _tests = createTests(client);

  for (FTest test in _tests) {
    if (verbose) stdout.write('${test.name}... ');
    try {
      await test.func();
      if (verbose) stdout.writeln('success!');
    } catch (e) {
      stdout.writeln(e.toString());
      result = result | test.errorCode;
    }
  }

  if (middleware_called) {
    stdout.writeln("Middleware successfully called.");
  } else {
    stdout.writeln("Middleware never called!");
    result = 1;
  }

  exit(result);
}

ArgResults _parseArgs(List<String> args) {
  var parser = new ArgParser();
  parser.addOption('host', defaultsTo: 'localhost', help: 'The server host');
  parser.addOption('port', defaultsTo: '9090', help: 'The port to connect to');
  parser.addOption('transport',
      defaultsTo: 'http',
      allowed: ['http'],
      help: 'The transport name',
      allowedHelp: {
        'http': 'http transport'
      });
  parser.addOption('protocol',
      defaultsTo: 'binary',
      allowed: ['binary', 'compact', 'json'],
      help: 'The protocol name',
      allowedHelp: {
        'binary': 'TBinaryProtocol',
        'compact': 'TCompactProtocol',
        'json': 'TJsonProtocol'
      });
  parser.addFlag('verbose', defaultsTo: false);

  ArgResults results;
  try {
    results = parser.parse(args);
  } catch (e) {
    stdout.writeln('$e\n');
  }

  if (results == null) stdout.write(parser.usage);

  return results;
}

TProtocolFactory getProtocolFactory(String protocolType) {
  if (protocolType == 'binary') {
    return new TBinaryProtocolFactory();
  } else if (protocolType == 'compact') {
    return new TCompactProtocolFactory();
  } else if (protocolType == 'json') {
    return new TJsonProtocolFactory();
  }

  throw new ArgumentError.value(protocolType);
}

Middleware clientMiddleware() {
  return (InvocationHandler next) {
    return (String serviceName, String methodName, List args) {
      if (args.length > 1 && args[1].runtimeType == Uint8List){
          stdout.write(methodName + "(" + args[1].length.toString() + ")"
              " = ");
      } else {
        stdout.write(methodName + "(" + args.sublist(1).toString() + ") = ");
      }
      middleware_called = true;
      return next(serviceName, methodName, args).then((result) {
        stdout.write(result.toString() + '\n');
        return result;
      }).catchError((e) {
        stdout.write(e.toString() + '\n');
        throw e;
      });
    };

  };
}

Future _initTestClient(
    {String host, int port, String transportType, String protocolType}) async {

  FProtocolFactory fProtocolFactory = null;
  FTransport transport = null;

//  Nats is not available without the SDK in dart, so HTTP is the only transport we can test
  var uri = Uri.parse('http://$host:$port');
// Set request and response size limit to 1mb
  var maxSize = 1048576;
  transport = new FHttpTransport(new wt.HttpClient(), uri, requestSizeLimit: maxSize, responseSizeLimit: maxSize);
  await transport.open();

  fProtocolFactory = new FProtocolFactory(getProtocolFactory(protocolType));
  client = new FFrugalTestClient(new FServiceProvider(transport, fProtocolFactory), [clientMiddleware()]);
}
