package eu.eudat.gef.service;

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
		T ret = pico.getComponent(klass);
		if (ret == null) {
			log.error("null reference when retrieving object of class: " + klass);
		}
		return ret;
	}
}
