package eu.eudat.gef.service;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.impl.PidServerImpl;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
public class PidService {

	static org.slf4j.Logger log = LoggerFactory.getLogger(IrodsService.class);

	static {
		PidServerConfig pidConfig = new PidServerConfig();
		pidConfig.epicServerUrl = WebAppConfig.get("/config/pid/epicServerUrl");
		pidConfig.localPrefix = WebAppConfig.get("/config/pid/localPrefix");
		pidConfig.username = WebAppConfig.get("/config/pid/user");
		pidConfig.password = WebAppConfig.get("/config/pid/pass");
		Services.register(pidConfig);
		Services.register(PidServerImpl.class);
	}
}
