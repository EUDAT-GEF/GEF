package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import java.io.InputStream;
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
public class BuildsHandler {
	private static final String apiUrl = "builds";

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(BuildsHandler.class);
	final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

	ReverseProxy rp;
	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	public BuildsHandler() throws MalformedURLException {
		rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
	}

	@POST
	public InputStream newBuild() throws Exception {
		return rp.forward(apiUrl, request, response);
	}

	@POST
	@Path("{buildID}")
	@Consumes(MediaType.MULTIPART_FORM_DATA)
	public InputStream doBuild(@PathParam("buildID") String buildID) throws Exception {
		return rp.forward(apiUrl + "/" + buildID, request, response);
	}
}
