package eu.eudat.gef.interop.impl;

import eu.eudat.gef.app.GEFConfig;
import eu.eudat.gef.interop.DatastreamReader;
import eu.eudat.gef.app.Services;
import java.net.URL;
import static org.junit.Assert.*;

/**
 *
 * @author edima
 */
public class B2ShareDatastreamReaderFactoryTest {

	public B2ShareDatastreamReaderFactoryTest() {
		GEFConfig.Pid pid = GEFConfig.makePid("", "", "", "");
		Services.initPidService(pid);
	}

	/**
	 * Test of make method, of class B2ShareDatastreamReaderFactory.
	 */
	@org.junit.Test
	public void testMakeDatastreamReader() throws Exception {
		System.out.println("makeDatastreamReader");

		// a B2SHARE pid we use for testing
		String pid = "http://hdl.handle.net/11304/2eb6781c-b988-11e3-8cd7-14feb57d12b9";
		URL pidURL = new URL(pid);

		DatastreamReader dr = new B2ShareDatastreamReaderFactory().make(pidURL);
		assert (dr.listDatastreams().size() == 1);
		String first = dr.listDatastreams().iterator().next();
		assertEquals(first, "dfe036b2-8294-11e3-9c41-005056943408.zip");
	}
}
