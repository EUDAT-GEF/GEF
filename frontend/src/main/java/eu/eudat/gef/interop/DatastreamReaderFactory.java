package eu.eudat.gef.interop;

import java.io.IOException;
import java.net.URL;

/**
 *
 * @author edima
 */
public interface DatastreamReaderFactory {
	DatastreamReader make(URL url) throws IOException;
}
