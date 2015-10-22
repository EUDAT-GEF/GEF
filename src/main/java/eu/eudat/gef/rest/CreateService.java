package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import java.net.MalformedURLException;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import java.text.DateFormat;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
@Path("builds")
public class CreateService {
	private static final String gefDockerBuildApi = "builds";

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(CreateService.class);
	final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

	ReverseProxy rp;
	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	public CreateService() throws MalformedURLException {
		rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
	}

	@POST
	public void newBuild() throws Exception {
		rp.forward(gefDockerBuildApi, request, response);
	}

	@POST
	@Path("{buildID}")
	@Consumes(MediaType.MULTIPART_FORM_DATA)
	public void doBuild(@PathParam("buildID") String buildID) throws Exception {
		rp.forward(gefDockerBuildApi + "/" + buildID, request, response);
	}

	private void println(String string) {
		System.out.println(string);
	}
}
