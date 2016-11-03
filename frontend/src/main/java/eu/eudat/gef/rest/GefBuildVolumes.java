package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import org.slf4j.LoggerFactory;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.text.DateFormat;

/**
 * Created by wqiu on 05/10/16.
 */
@Path("buildVolumes")
public class GefBuildVolumes {
    private static final String gefDockerBuildApi = "buildVolumes";

    private static final org.slf4j.Logger log = LoggerFactory.getLogger(GefBuildImages.class);
    final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

    ReverseProxy rp;
    @Context
    HttpServletRequest request;
    @Context
    HttpServletResponse response;

    public GefBuildVolumes () throws MalformedURLException {
        rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
    }

    @POST
    public InputStream newBuild() throws Exception {
        return rp.forward(gefDockerBuildApi, request, response);
    }

    @POST
    @Path("{buildID}")
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    public InputStream doBuild(@PathParam("buildID") String buildID) throws Exception {
        return rp.forward(gefDockerBuildApi + "/" + buildID, request, response);
    }
}
