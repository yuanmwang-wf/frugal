// Autogenerated by Frugal Compiler (2.0.0-RC3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library v1_music.src.f_albumwinners_scope;

import 'dart:async';

import 'package:thrift/thrift.dart' as thrift;
import 'package:frugal/frugal.dart' as frugal;

import 'package:v1_music/v1_music.dart' as t_v1_music;


const String delimiter = '.';

/// Scopes are a Frugal extension to the IDL for declaring PubSub
/// semantics. Subscribers to this scope will be notified if they win a contest.
/// Scopes must have a prefix.
class AlbumWinnersPublisher {
  frugal.FPublisherTransport transport;
  frugal.FProtocolFactory protocolFactory;
  Map<String, frugal.FMethod> _methods;
  AlbumWinnersPublisher(frugal.FScopeProvider provider, [List<frugal.Middleware> middleware]) {
    transport = provider.publisherTransportFactory.getTransport();
    protocolFactory = provider.protocolFactory;
    this._methods = {};
    this._methods['Winner'] = new frugal.FMethod(this._publishWinner, 'AlbumWinners', 'publishWinner', middleware);
  }

  Future open() {
    return transport.open();
  }

  Future close() {
    return transport.close();
  }

  Future publishWinner(frugal.FContext ctx, t_v1_music.Album req) {
    return this._methods['Winner']([ctx, req]);
  }

  Future _publishWinner(frugal.FContext ctx, t_v1_music.Album req) async {
    var op = "Winner";
    var prefix = "v1.music.";
    var topic = "${prefix}AlbumWinners${delimiter}${op}";
    var memoryBuffer = new frugal.TMemoryOutputBuffer(transport.publishSizeLimit);
    var oprot = protocolFactory.getProtocol(memoryBuffer);
    var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, 0);
    oprot.writeRequestHeader(ctx);
    oprot.writeMessageBegin(msg);
    req.write(oprot);
    oprot.writeMessageEnd();
    await transport.publish(topic, memoryBuffer.writeBytes);
  }
}


/// Scopes are a Frugal extension to the IDL for declaring PubSub
/// semantics. Subscribers to this scope will be notified if they win a contest.
/// Scopes must have a prefix.
class AlbumWinnersSubscriber {
  final frugal.FScopeProvider provider;
  final List<frugal.Middleware> _middleware;

  AlbumWinnersSubscriber(this.provider, [this._middleware]) {}

  Future<frugal.FSubscription> subscribeWinner(dynamic onAlbum(frugal.FContext ctx, t_v1_music.Album req)) async {
    var op = "Winner";
    var prefix = "v1.music.";
    var topic = "${prefix}AlbumWinners${delimiter}${op}";
    var transport = provider.subscriberTransportFactory.getTransport();
    await transport.subscribe(topic, _recvWinner(op, provider.protocolFactory, onAlbum));
    return new frugal.FSubscription(topic, transport);
  }

  frugal.FAsyncCallback _recvWinner(String op, frugal.FProtocolFactory protocolFactory, dynamic onAlbum(frugal.FContext ctx, t_v1_music.Album req)) {
    frugal.FMethod method = new frugal.FMethod(onAlbum, 'AlbumWinners', 'subscribeAlbum', this._middleware);
    callbackWinner(thrift.TTransport transport) {
      var iprot = protocolFactory.getProtocol(transport);
      var ctx = iprot.readRequestHeader();
      var tMsg = iprot.readMessageBegin();
      if (tMsg.name != op) {
        thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
        iprot.readMessageEnd();
        throw new thrift.TApplicationError(
        thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
      }
      var req = new t_v1_music.Album();
      req.read(iprot);
      iprot.readMessageEnd();
      method([ctx, req]);
    }
    return callbackWinner;
  }
}

