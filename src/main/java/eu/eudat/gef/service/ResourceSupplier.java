package eu.eudat.gef.service;

import com.google.common.io.InputSupplier;
import java.io.IOException;
import java.io.InputStream;

/**
 * @author edima
 */
public class ResourceSupplier implements InputSupplier<InputStream> {
	String resource;

	public ResourceSupplier(String resource) {
		this.resource = resource;
	}

	public InputStream getInput() throws IOException {
		return getClass().getResourceAsStream(resource);
	}
}
