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

/// FAsyncTransport is an extension of FTransport that asynchronous frameworks
/// can implement. Implementations need only implement [flush] to send request
/// data and call [handleResponse] when asynchronous responses are received.
abstract class FAsyncTransport extends FTransport {
  final Logger _log = Logger('FAsyncTransport');
  Map<int, Completer<Uint8List>> _handlers = {};

  /// Instantiate an [FAsyncTransport].
  FAsyncTransport({int requestSizeLimit})
      : super(requestSizeLimit: requestSizeLimit);

  /// Flush the payload to the server. Implementations must be threadsafe.
  Future<Null> flush(Uint8List payload);

  @override
  Future<Null> oneway(FContext ctx, Uint8List payload) async {
    _preflightRequestCheck(payload);
    await flush(payload).timeout(ctx.timeout, onTimeout: () {
      throw TTransportError(FrugalTTransportErrorType.TIMED_OUT,
          'request timed out after ${ctx.timeout}');
    });
  }

  @override
  Future<TTransport> request(FContext ctx, Uint8List payload) async {
    _preflightRequestCheck(payload);

    Completer<Uint8List> resultCompleter = Completer();

    if (_handlers.containsKey(ctx._opId)) {
      throw StateError("frugal: context already registered");
    }
    _handlers[ctx._opId] = resultCompleter;
    Completer<Uint8List> closedCompleter = Completer();
    StreamSubscription<Object> closedSub = onClose.listen((_) {
      closedCompleter
          .completeError(TTransportError(FrugalTTransportErrorType.NOT_OPEN));
    });

    try {
      await flush(payload);
      Future<Uint8List> resultFuture =
          resultCompleter.future.timeout(ctx.timeout);

      // Bail early if the transport is closed
      Uint8List response =
          await Future.any([resultFuture, closedCompleter.future]);
      return TMemoryTransport.fromUint8List(response);
    } on TimeoutException catch (_) {
      throw TTransportError(FrugalTTransportErrorType.TIMED_OUT,
          "request timed out after ${ctx.timeout}");
    } finally {
      _handlers.remove(ctx._opId);

      // don't wait until this is disposed to cancel these
      await closedSub.cancel();
      if (!closedCompleter.isCompleted) {
        closedCompleter.complete();
      }
      if (!resultCompleter.isCompleted) {
        resultCompleter.complete();
      }
    }
  }

  /// Handles a frugal frame response. NOTE: this frame must NOT include the
  /// frame size. Implementations should call this when asynchronous responses
  /// are received from the server.
  void handleResponse(Uint8List frame) {
    var headers = Headers.decodeFromFrame(frame);
    var opId;
    try {
      opId = int.parse(headers[_opidHeader]);
    } catch (e) {
      _log.warning("frugal: invalid protocol frame: op id not a uint64", e);
      return;
    }

    Completer<Uint8List> handler = _handlers[opId];
    if (handler == null) {
      _log.warning("frugal: no handler found for message, dropping message");
      return;
    }

    if (handler.isCompleted) {
      _log.warning(
          "frugal: handler already called for message, dropping message");
    }
    handler.complete(frame);
  }
}
