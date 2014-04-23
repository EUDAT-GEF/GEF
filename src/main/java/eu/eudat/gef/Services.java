package eu.eudat.gef;

import org.picocontainer.DefaultPicoContainer;
import org.picocontainer.MutablePicoContainer;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
public class Services {

	static org.slf4j.Logger log = LoggerFactory.getLogger(Services.class);
	private static MutablePicoContainer pico = new DefaultPicoContainer();

	public static void register(Object o) {
		pico.addComponent(o);
	}

	public static <T> T get(Class<T> klass) {
		return pico.getComponent(klass);
	}
}
