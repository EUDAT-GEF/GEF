package eu.eudat.gef;

import de.tuebingen.uni.sfs.epicpid.PidServerConfig;
import de.tuebingen.uni.sfs.epicpid.impl.PidServerImpl;
import eu.eudat.gef.irodslink.IrodsAccessConfig;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.impl.jargon.JargonConnection;
import eu.eudat.gef.rest.DataSets;
import eu.eudat.gef.rest.Jobs;
import eu.eudat.gef.rest.Workflows;
import javax.servlet.http.HttpSession;
import javax.ws.rs.core.Context;
import org.picocontainer.DefaultPicoContainer;
import org.picocontainer.MutablePicoContainer;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
public class Services {

    static org.slf4j.Logger log = LoggerFactory.getLogger(Services.class);
    private static MutablePicoContainer pico = null;

    private static void init() {
        pico = new DefaultPicoContainer();
        IrodsAccessConfig irodsConfig = new IrodsAccessConfig();
        irodsConfig.server = WebAppConfig.get("/config/irods/server");
        irodsConfig.port = WebAppConfig.getInt("/config/irods/port");
        irodsConfig.username = WebAppConfig.get("/config/irods/username");
        irodsConfig.password = WebAppConfig.get("/config/irods/password");
        irodsConfig.path = WebAppConfig.get("/config/irods/path");
        irodsConfig.resource = WebAppConfig.get("/config/irods/resource");

        pico.addComponent(irodsConfig)
                .addComponent(JargonConnection.class);

        PidServerConfig pidConfig = new PidServerConfig();
        pidConfig.username = WebAppConfig.get("/config/pid/user");
        pidConfig.password = WebAppConfig.get("/config/pid/pass");
        pico.addComponent(pidConfig)
                .addComponent(PidServerImpl.class);
    }
    @Context
    HttpSession session;

    public static <T> T get(Class<T> klass) {
        if (pico == null) {
            init();
            try {
                IrodsConnection ic = get(IrodsConnection.class);
                ic.getObject(ic.getInitialPath()).asCollection().create();
                ic.getObject(ic.getInitialPath()+"/"+ DataSets.DATA_DIR).asCollection().create();
                ic.getObject(ic.getInitialPath()+"/"+ Workflows.WORKFLOW_DIR).asCollection().create();
                ic.getObject(ic.getInitialPath()+"/"+ Jobs.JOBS_DIR).asCollection().create();
            } catch (Exception xc) {
            }
        }
        return pico.getComponent(klass);
    }
}
