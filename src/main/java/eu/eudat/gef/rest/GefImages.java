package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import java.io.InputStream;
import java.net.MalformedURLException;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import java.text.DateFormat;
import javax.ws.rs.core.MediaType;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
@Path("images")
public class GefImages {
	private static final String gefDockerImagesApi = "images";

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(GefImages.class);
	final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

	ReverseProxy rp;
	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	public GefImages() throws MalformedURLException {
		rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
	}

	@GET
	public InputStream listImages() throws Exception {
		return rp.forward(gefDockerImagesApi, request, response);
	}

	@GET
	@Path("{imageID}")
	public InputStream inspectImage(@PathParam("imageID") String imageID) throws Exception {
		return rp.forward(gefDockerImagesApi + "/" + imageID, request, response);
	}
}
