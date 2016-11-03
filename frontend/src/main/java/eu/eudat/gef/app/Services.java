package eu.eudat.gef.app;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.mockimpl.PidMockImpl;
import org.picocontainer.DefaultPicoContainer;
import org.picocontainer.MutablePicoContainer;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
public class Services {

	static org.slf4j.Logger log = LoggerFactory.getLogger(Services.class);
	private static MutablePicoContainer pico = new DefaultPicoContainer();

	public static void init(GEFConfig cfg) {
		initPidService(cfg.gefParams.pid);
	}

	public static void initPidService(GEFConfig.Pid cfg) {
		if (getSilent(PidServerConfig.class) != null) {
			return;
		}
		PidServerConfig pidConfig = new PidServerConfig();
		pidConfig.epicServerUrl = cfg.epicServerUrl;
		pidConfig.localPrefix = cfg.localPrefix;
		pidConfig.username = cfg.user;
		pidConfig.password = cfg.pass;
		Services.register(pidConfig);
		Services.register(PidMockImpl.class);
	}

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

	public static <T> T getSilent(Class<T> klass) {
		return pico.getComponent(klass);
	}
}
