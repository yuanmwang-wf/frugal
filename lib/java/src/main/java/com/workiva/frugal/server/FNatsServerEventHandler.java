package com.workiva.frugal.server;

import java.util.Map;

/**
 * Provides an interface with which to handle events from an FNatsServer.
 *
 * An FNatsServerEventHandler should serve a distinct purpose from middleware.
 * A set of middleware, with a processor, should be able to be used on any kind
 * of frugal server. Conversely, this should only be used for behaviour and
 * events for an FNatsServer not applicable to other frugal servers.
 *
 * It is preferred to use middleware if either solution works, as middleware is
 * more portable between different servers.
 */
public interface FNatsServerEventHandler {
    /**
     * Called when a request is received, but before it is put onto a work
     * queue.
     *
     * @param ephemeralProperties
     */
    void onRequestReceived(Map<Object, Object> ephemeralProperties);

    /**
     * Called when a request is retrieved from a work queue to begin being
     * processed.
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
