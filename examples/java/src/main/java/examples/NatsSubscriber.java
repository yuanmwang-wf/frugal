package examples;

import com.workiva.frugal.protocol.FProtocolFactory;
import com.workiva.frugal.provider.FScopeProvider;
import com.workiva.frugal.transport.FNatsPublisherTransport;
import com.workiva.frugal.transport.FNatsSubscriberTransport;
import com.workiva.frugal.transport.FPublisherTransportFactory;
import com.workiva.frugal.transport.FSubscriberTransportFactory;
import io.nats.client.Connection;
import io.nats.client.Nats;
import io.nats.client.Options;
import org.apache.thrift.TException;
import org.apache.thrift.protocol.TBinaryProtocol;
import v1.music.AlbumWinnersSubscriber;

import java.io.IOException;
import java.util.Properties;
import java.util.concurrent.TimeoutException;

/**
 * Create a NATS PubSub subscriber.
 */
public class NatsSubscriber {

    public static void main(String[] args) throws TException, IOException, TimeoutException {
        // Specify the protocol used for serializing requests.
        // The protocol stack must match the protocol stack of the publisher.
        FProtocolFactory protocolFactory = new FProtocolFactory(new TBinaryProtocol.Factory());

        // Create a NATS client (using default options for local dev)
        Properties properties = new Properties();
        properties.put(Options.PROP_URL, Options.DEFAULT_URL);
        Options.Builder optionsBuilder = new Options.Builder(properties);
        Connection conn = Nats.connect(optionsBuilder.build());

        // Create the pubsub scope provider, given the NATs connection and protocol
        FPublisherTransportFactory publisherFactory = new FNatsPublisherTransport.Factory(conn);
        FSubscriberTransportFactory subscriberFactory = new FNatsSubscriberTransport.Factory(conn);
        FScopeProvider provider = new FScopeProvider(publisherFactory, subscriberFactory, protocolFactory);

        // Subscribe to winner announcements
        AlbumWinnersSubscriber.Iface subscriber = new AlbumWinnersSubscriber.Client(provider);
        subscriber.subscribeWinner((ctx, album) -> System.out.println("You won! " + album));
        subscriber.subscribeContestStart((ctx, albums) -> System.out.println("Contest started, available albums: " + albums));
        System.out.println("Subscriber started...");
    }
}
