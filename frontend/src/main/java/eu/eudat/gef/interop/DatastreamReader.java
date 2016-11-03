package eu.eudat.gef.interop;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.util.Set;

/**
 *
 * @author edima
 */
public interface DatastreamReader {

	Set<String> listDatastreams() throws DatastreamReaderException;

	void saveDatastream(String datastream, File destination) throws DatastreamReaderException, IOException;

	InputStream streamDatastream(String datastream) throws DatastreamReaderException, IOException;
}
