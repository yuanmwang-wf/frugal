package com.workiva.frugal.server;

import java.util.Map;

/**
 * Provides an interface with which to handle events from a frugal server.
 *
 * An FServerEventHandler should serve a distinct purpose from middleware.
 * A set of middleware, with a processor, should be able to be used on any kind
 * of frugal server. Some use cases for this could be monitoring how long a
 * request sits before being processed, or adding keys to MDC before processing
 * a message.
 *
 * It is preferred to use middleware if either solution works, as middleware is
 * more portable between different servers.
 */
public interface FServerEventHandler {
    /**
     * Called when a request is first received. For some async servers, this
     * can be different than when a request is processed. For many synchronous
     * servers, {@link FServerEventHandler#onRequestStarted(Map)} will be
     * called at the same time.
     *
     * @param ephemeralProperties
     */
    void onRequestReceived(Map<Object, Object> ephemeralProperties);

    /**
     * Called when a request is about to be processed. For some async servers,
     * this can be different than when a request is received. For many
     * synchronous servers, {@link FServerEventHandler#onRequestReceived(Map)}
     * will be called at the same time.
     *
     * @param ephemeralProperties
     */
    void onRequestStarted(Map<Object, Object> ephemeralProperties);

    /**
     * Called when a request is done being processed.
     *
     * @param ephemeralProperties
     */
    void onRequestEnded(Map<Object, Object> ephemeralProperties);
}
