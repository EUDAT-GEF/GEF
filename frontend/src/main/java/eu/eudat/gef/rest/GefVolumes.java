package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEF;
import org.slf4j.LoggerFactory;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.Context;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.text.DateFormat;

/**
 * Created by wqiu on 05/10/16.
 */


@Path("volumes")
public class GefVolumes {
    private static final String gefDockerVolumesApi = "volumes";

    private static final org.slf4j.Logger log = LoggerFactory.getLogger(GefVolumes.class);
    final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

    ReverseProxy rp;
    @Context
    HttpServletRequest request;
    @Context
    HttpServletResponse response;

    public GefVolumes() throws MalformedURLException {
        rp = new ReverseProxy(GEF.getInstance().config.gefParams.gefDocker);
    }

    @POST
    public InputStream postVolume() throws Exception {
        return rp.forward(gefDockerVolumesApi, request, response);
    }

    @GET
    public InputStream listVolumes() throws Exception {
        return rp.forward(gefDockerVolumesApi, request, response);
    }

    @GET
    @Path("{volumeID}/{path: .*}")
    public InputStream volumeContentHandler(@PathParam("volumeID") String volumeID, @PathParam("path") String path) throws Exception {
        return rp.forward(gefDockerVolumesApi + "/" + volumeID + "/" + path, request, response);
    }
}
