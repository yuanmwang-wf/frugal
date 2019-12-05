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

/// Wraps a [TSocketTransport] to produce an [FTransport] which uses the given
/// socket for send/callback operations in a way that is compatible with Frugal.
/// Messages received on the [TSocket] (i.e. Frugal frames) are routed to the
/// [FAsyncTransport]'s handleResponse method.
class FAdapterTransport extends FAsyncTransport {
  final Logger _adapterTransportLog = Logger('FAdapterTransport');
  _TFramedTransport _framedTransport;

  StreamSubscription<_FrameWrapper> _onFrameSub;

  /// Create an [FAdapterTransport] with the given [TSocketTransport].
  FAdapterTransport(TSocketTransport transport)
      : _framedTransport = _TFramedTransport(transport.socket),
        super() {
    // If there is an error on the socket, close the transport pessimistically.
    // This error is already logged upstream in TSocketTransport.
    listenToStream(transport.socket.onError, (e) => close(e));
    // Forward state changes on to the transport monitor.
    // Note: Just forwarding OPEN on for the time-being.
    listenToStream(transport.socket.onState, (state) {
      if (state == TSocketState.OPEN) _monitor?.signalOpen();
    });

    manageDisposable(_framedTransport);
  }

  @override
  bool get isOpen => _framedTransport.isOpen;

  @override
  Future open() async {
    await _framedTransport.open();

    _onFrameSub = _framedTransport.onFrame.listen(_handleFrame);
  }

  void _handleFrame(_FrameWrapper frame) {
    try {
      handleResponse(frame.frameBytes);
    } catch (e) {
      // Fatal error. Close the transport.
      _adapterTransportLog.severe(
          "FAsyncCallback had a fatal error ${e.toString()}." +
              "Closing transport.");
      close(e);
    }
  }

  @override
  Future close([Error error]) async {
    await _framedTransport?.close();
    await _onFrameSub?.cancel();
    await super.close(error);
  }

  @override
  Future<Null> flush(Uint8List payload) {
    _framedTransport.socket.send(payload);
    return Future.value();
  }

  @override
  Future<Null> onDispose() async {
    await close();
    await super.onDispose();
  }
}
