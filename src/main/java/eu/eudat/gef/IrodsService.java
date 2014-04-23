package eu.eudat.gef;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.impl.PidServerImpl;
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
public class IrodsService {

	static org.slf4j.Logger log = LoggerFactory.getLogger(IrodsService.class);

	static {
		init();
	}

	private static void init() {
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
		}

	}
}
