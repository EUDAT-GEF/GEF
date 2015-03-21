package eu.eudat.gef.app;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.mockimpl.PidMockImpl;
import eu.eudat.gef.irodslink.IrodsAccessConfig;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.impl.jargon.JargonConnection;
import eu.eudat.gef.rest.DataSets;
import eu.eudat.gef.rest.Jobs;
import eu.eudat.gef.rest.Workflows;
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
		initIrodsService(cfg.gefParams.irods);
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
//		Services.register(PidServerImpl.class);
	}

	public static void initIrodsService(GEFConfig.Irods cfg) {
		if (getSilent(IrodsAccessConfig.class) != null) {
			return;
		}
		IrodsAccessConfig irodsConfig = new IrodsAccessConfig();
		irodsConfig.server = cfg.server;
		irodsConfig.port = cfg.port;
		irodsConfig.username = cfg.username;
		irodsConfig.password = cfg.password;
		irodsConfig.path = cfg.path;
		irodsConfig.resource = cfg.resource;

		Services.register(irodsConfig);
		Services.register(JargonConnection.class);

		try {
			IrodsConnection ic = Services.get(IrodsConnection.class);
			ic.getObject(ic.getInitialPath()).asCollection().create();
			ic.getObject(ic.getInitialPath() + "/" + DataSets.DATA_DIR).asCollection().create();
			ic.getObject(ic.getInitialPath() + "/" + Workflows.WORKFLOW_DIR).asCollection().create();
			ic.getObject(ic.getInitialPath() + "/" + Jobs.JOBS_DIR).asCollection().create();
		} catch (Exception xc) {
			log.error(xc.getMessage(), xc);
			// ignore this one
		}
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
