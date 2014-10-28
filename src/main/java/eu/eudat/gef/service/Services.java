package eu.eudat.gef.service;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.mockimpl.PidMockImpl;
import eu.eudat.gef.irodslink.IrodsAccessConfig;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.impl.jargon.JargonConnection;
import eu.eudat.gef.rest.DataSets;
import eu.eudat.gef.rest.Jobs;
import eu.eudat.gef.rest.Workflows;
import javax.servlet.ServletContextEvent;
import javax.servlet.ServletContextListener;
import org.picocontainer.DefaultPicoContainer;
import org.picocontainer.MutablePicoContainer;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
public class Services implements ServletContextListener {

	static org.slf4j.Logger log = LoggerFactory.getLogger(Services.class);
	private static MutablePicoContainer pico = new DefaultPicoContainer();

	@Override
	public void contextInitialized(ServletContextEvent sce) {
		init();
	}

	@Override
	public void contextDestroyed(ServletContextEvent sce) {
	}

	public static void init() {
		initPidService();
		initIrodsService();
	}

	public static void initPidService() {
		if (get(PidServerConfig.class) != null) {
			return;
		}
		PidServerConfig pidConfig = new PidServerConfig();
		pidConfig.epicServerUrl = WebAppConfig.get("/config/pid/epicServerUrl");
		pidConfig.localPrefix = WebAppConfig.get("/config/pid/localPrefix");
		pidConfig.username = WebAppConfig.get("/config/pid/user");
		pidConfig.password = WebAppConfig.get("/config/pid/pass");
		Services.register(pidConfig);
		Services.register(PidMockImpl.class);
//		Services.register(PidServerImpl.class);
	}

	public static void initIrodsService() {
		if (get(IrodsAccessConfig.class) != null) {
			return;
		}
		IrodsAccessConfig irodsConfig = new IrodsAccessConfig();
		irodsConfig.server = WebAppConfig.get("/config/irods/server");
		irodsConfig.port = WebAppConfig.getInt("/config/irods/port");
		irodsConfig.username = WebAppConfig.get("/config/irods/username");
		irodsConfig.password = WebAppConfig.get("/config/irods/password");
		irodsConfig.path = WebAppConfig.get("/config/irods/path");
		irodsConfig.resource = WebAppConfig.get("/config/irods/resource");

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
}
