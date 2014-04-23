package eu.eudat.gef.interop;

import java.io.IOException;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;

/**
 *
 * @author edima
 */
public class DatastreamReaderManager {
	List<DatastreamReaderFactory > l = new ArrayList<DatastreamReaderFactory>();
	
	public void registerDatastreamReaderFactory(DatastreamReaderFactory f) {
		l.add(f);
	}
	
	public DatastreamReader getDatastreamReader(URL pid) throws IOException {
		for (DatastreamReaderFactory f :l)
			try {
				DatastreamReader dr = f.make(pid);
				if (dr != null) {
					return dr;
				}
			} catch (Exception xc) {
				//ignore
			}
		return null;
	}
}
