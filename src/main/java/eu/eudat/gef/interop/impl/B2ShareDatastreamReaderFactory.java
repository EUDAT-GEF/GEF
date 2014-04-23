package eu.eudat.gef.interop.impl;

import com.google.common.io.Files;
import com.google.common.io.Resources;
import eu.eudat.gef.interop.DatastreamReader;
import eu.eudat.gef.interop.DatastreamReaderException;
import eu.eudat.gef.interop.DatastreamReaderFactory;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;

/**
 *
 * @author edima
 */
public class B2ShareDatastreamReaderFactory implements DatastreamReaderFactory {

	public static final String SELECTOR_FOR_FILES = "#files #collapseTwo table tbody td a";
	public static final int TIMEOUT = 10 * 1000;

	Map<String, URL> fileMap = new HashMap<String, URL>();

	public DatastreamReader make(URL pidUrl) throws IOException {
		Document doc = Jsoup.connect(pidUrl.toExternalForm())
				.timeout(TIMEOUT).get();
		Elements files = doc.select(SELECTOR_FOR_FILES);
		for (Element e : files) {
			try {
				String href = e.attr("href");
				URL u = new URL(href);
				String[] segments = u.getPath().split("/");
				if (segments.length == 0) {
					continue;
				}
				String name = segments[segments.length - 1];
				fileMap.put(name, u);
			} catch (MalformedURLException xc) {
				// if it's not a valid URL, no good
			}
		}

		return new DatastreamReader() {

			public Set<String> listDatastreams() throws DatastreamReaderException {
				return fileMap.keySet();
			}

			public void saveDatastream(String datastream, File destination) throws DatastreamReaderException, IOException {
				URL url = fileMap.get(datastream);
				if (url == null) {
					throw new DatastreamReaderException("inexistent datastream name: " + datastream);
				}
				Resources.asByteSource(url).copyTo(Files.asByteSink(destination));
			}

			public InputStream streamDatastream(String datastream) throws DatastreamReaderException, IOException {
				URL url = fileMap.get(datastream);
				if (url == null) {
					throw new DatastreamReaderException("inexistend datastream name: " + datastream);
				}
				return Resources.asByteSource(url).openStream();
			}
		};
	}
}
