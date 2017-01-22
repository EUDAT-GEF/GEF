package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import java.io.InputStream;
import java.net.MalformedURLException;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import java.text.DateFormat;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
@Path("jobs")
public class JobsHandler {
	private static final String apiUrl = "jobs";

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(JobsHandler.class);
	final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

	ReverseProxy rp;
	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	public JobsHandler() throws MalformedURLException {
		rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
	}

	@POST
	public InputStream postJob() throws Exception {
		return rp.forward(apiUrl, request, response);
	}

	@GET
	public InputStream listJobs() throws Exception {
		return rp.forward(apiUrl, request, response);
	}

	@GET
	@Path("{jobID}")
	public InputStream inspectJob(@PathParam("jobID") String jobID) throws Exception {
		return rp.forward(apiUrl + "/" + jobID, request, response);
	}
}
