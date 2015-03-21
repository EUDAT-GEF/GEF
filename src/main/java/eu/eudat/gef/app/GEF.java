package eu.eudat.gef.app;

import com.fasterxml.jackson.databind.ObjectMapper;
import eu.eudat.gef.rest.DataSets;
import eu.eudat.gef.rest.Jobs;
import eu.eudat.gef.rest.Workflows;
import io.dropwizard.Application;
import io.dropwizard.assets.AssetsBundle;
import io.dropwizard.setup.Bootstrap;
import io.dropwizard.setup.Environment;
import java.io.IOException;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 *
 * TODO: UI: good basic UI for execution and inspection of results
 *
 * TODO: PID dataset inspection: visualization, deep inspection of a data set as
 * a special case of execution via GEF
 *
 */
public class GEF extends Application<GEFConfig> {

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(GEF.class);

	private static GEF instance;

	public static void main(String[] args) throws Exception {
		new GEF().run(args);
	}

	@Override
	public String getName() {
		return "GEF";
	}

	@Override
	public void initialize(Bootstrap<GEFConfig> bootstrap) {
		bootstrap.addBundle(new AssetsBundle("/assets", "/", "index.html"));
	}

	@Override
	public void run(GEFConfig config, Environment environment) throws Exception {
		log.info("GEF initialization started.");
		instance = this;

		System.out.println("Using parameters: ");
		try {
			System.out.println(new ObjectMapper().writerWithDefaultPrettyPrinter().
					writeValueAsString(config.gefParams));
		} catch (IOException xc) {
		}

		try {
			environment.jersey().setUrlPattern("/rest/*");
			environment.jersey().register(DataSets.class);
			environment.jersey().register(Jobs.class);
			environment.jersey().register(Workflows.class);

			Services.init(config);
		} catch (Exception ex) {
			log.error("INIT EXCEPTION", ex);
			throw ex; // force exit
		}
		log.info("GEF initialization finished.");
	}

	public static GEF getInstance() {
		return instance;
	}
}
