package eu.eudat.gef.interop;

import de.tuebingen.uni.sfs.epicpid.impl.Strings;
import eu.eudat.gef.service.PidService;
import eu.eudat.gef.interop.impl.B2ShareDatastreamReaderFactory;
import java.io.File;
import java.net.URL;
import java.util.HashSet;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

/**
 *
 * @author edima
 */
//@Ignore
public class DatastreamReaderManagerTest {

	String pid = "http://hdl.handle.net/11304/30f24e76-b988-11e3-8cd7-14feb57d12b9";
	PidService ps;

	public DatastreamReaderManagerTest() {
	}

	@org.junit.Test
	public void testListDatastream() throws Exception {
		System.out.println("listDatastream");

		DatastreamReaderManager instance = new DatastreamReaderManager();
		instance.registerDatastreamReaderFactory(new B2ShareDatastreamReaderFactory());

		URL pidURL = new URL(pid);
		DatastreamReader dr = instance.getDatastreamReader(pidURL);

		for (String dsname : dr.listDatastreams()) {
			System.out.println(pid + " has datastream: " + dsname);
		}
		assertEquals(dr.listDatastreams(),
				new HashSet<String>() {
					{
						add("32c9362c-82ad-11e3-9c41-005056943408.el");
						add("228c078a-82ad-11e3-8ef2-005056943408.el");
						add("42ec3dec-82ad-11e3-9c41-005056943408.el");
						add("30479830-82ad-11e3-bb5f-005056943408.aif");
						add("32a4b82e-82ad-11e3-ad80-005056943408.el");
						add("42b05cdc-82ad-11e3-b283-005056943408.aif");
					}
				});
	}

	@org.junit.Test
	public void getDatastream() throws Exception {
		System.out.println("getDatastream");

		DatastreamReaderManager instance = new DatastreamReaderManager();
		instance.registerDatastreamReaderFactory(new B2ShareDatastreamReaderFactory());

		URL pidURL = new URL(pid);
		DatastreamReader dr = instance.getDatastreamReader(pidURL);

		assert (dr.listDatastreams().size() > 0);
		for (String dsname : dr.listDatastreams()) {
			if (dsname.endsWith(".el")) {
				File f = File.createTempFile("gef-unit-test-get-datastream", "");
				System.out.println("saving datastream " + dsname + " to " + f.getAbsolutePath());
				dr.saveDatastream(dsname, f);
				System.out.println("done");
				assertTrue(f.length() > 0);
			}
		}
	}
}
