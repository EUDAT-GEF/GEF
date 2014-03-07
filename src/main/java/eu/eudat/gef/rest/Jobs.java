package eu.eudat.gef.rest;

import com.google.common.base.Charsets;
import com.google.common.io.Files;
import com.google.gson.Gson;
import com.sun.jersey.multipart.FormDataBodyPart;
import com.sun.jersey.multipart.FormDataMultiPart;
import de.tuebingen.uni.sfs.epicpid.PidServer;
import eu.eudat.gef.Services;
import eu.eudat.gef.irodslink.IrodsCollection;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.IrodsFile;
import eu.eudat.gef.irodslink.IrodsObject;
import eu.eudat.gef.rest.json.JobJson;
import java.io.File;
import java.io.FileInputStream;
import java.net.URI;
import java.net.URLConnection;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.util.Map;

/**
 * @author edima
 */
@Path("jobs")
public class Jobs {

    public static final String JOBS_DIR = "jobs";
    public static final String SUFFIX = ".gefcommand";
    @Context
    HttpServletRequest request;
    @Context
    HttpServletResponse response;

    @POST
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    @Produces(MediaType.APPLICATION_JSON)
    public Response post(FormDataMultiPart multiPart) throws Exception {
        System.out.println("--------- " + request.getPathInfo());
        Map<String, String> params = new LinkedHashMap<String, String>();
        for (Map.Entry<String, List<FormDataBodyPart>> e : multiPart.getFields().entrySet()) {
            for (FormDataBodyPart fdbp : e.getValue()) {
                params.put(fdbp.getName(), fdbp.getValue());
            }
        }
        String workflowPid = params.get("workflowPid");
        if (workflowPid == null || workflowPid.isEmpty()) {
            return Response.status(Response.Status.BAD_REQUEST).entity("dataset PID needed!").build();
        }

        IrodsConnection conn = Services.get(IrodsConnection.class);
        Map<String, String> newparams = new LinkedHashMap<String, String>();
        try {
            for (String key : params.keySet()) {
                String value = params.get(key);
                if (value != null && !value.isEmpty()) {
                    if (key.endsWith("Pid")) {
                        value = resolvePid(conn, value);
                        key = key.substring(0, key.length() - "Pid".length()) + "IrodsPath";
                        newparams.put(key, value);
                    } else {
                        newparams.put(key, params.get(key));
                    }
                }
            }
            params = newparams;
        } catch (Exception xc) {
            xc.printStackTrace();
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR).entity(xc.getMessage()).build();
        }

        File f = File.createTempFile("wkf", SUFFIX);
        StringBuilder content = new StringBuilder();
        for (Map.Entry<String, String> e : params.entrySet()) {
            content.append(e.getKey()).append("=").append(e.getValue()).append("\n");
        }
        Files.write(content, f, Charsets.UTF_8);

        String jobId = f.getName().substring(0, f.getName().length() - SUFFIX.length());
        String jobCollPath = conn.getInitialPath() + "/" + JOBS_DIR + "/" + jobId;
        IrodsCollection jobColl = conn.getObject(jobCollPath).asCollection();
        jobColl.create();
        IrodsFile jobparams = conn.getObject(jobCollPath + "/" + SUFFIX).asFile();
        jobparams.uploadFromLocalFile(f);
        f.delete();

        String json = new Gson().toJson(conn.makeUri(jobColl).toASCIIString());
        return Response.created(conn.makeUri(jobColl)).entity(jobId).build();
    }

    String resolvePid(IrodsConnection conn, String pid) throws Exception {
        String irodsCollPath;
        try {
            PidServer pider = Services.get(PidServer.class);
            String workflowUrl = pider.fromString(pid).retrieveUrl();
            irodsCollPath = new URI(workflowUrl).getPath();
        } catch (Exception xc) {
            System.out.println("cannot resolve PID: "+ xc);
            String testDataPid = conn.getInitialPath() + "/" + DataSets.DATA_DIR + "/" + pid;
            String testWflwPid = conn.getInitialPath() + "/" + Workflows.WORKFLOW_DIR + "/" + pid;
            if (conn.getObject(testDataPid).exists())
                irodsCollPath = testDataPid;
            else if (conn.getObject(testWflwPid).exists())
                irodsCollPath = testWflwPid;
            else throw xc;
        }
        try {
            return conn.getObject(irodsCollPath).asCollection()
                    .listFiles().iterator().next().getFullPath();
        } catch (Exception xc) {
            throw new Exception("cannot find PID resources", xc);
        }
    }

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public Response get0() throws Exception {
        return getAny();
    }

    @GET
    @Path("{a}")
    @Produces(MediaType.APPLICATION_JSON)
    public Response get1() throws Exception {
        return getAny();
    }

    @GET
    @Path("{a}/{b}")
    @Produces(MediaType.APPLICATION_OCTET_STREAM)
    public Response get2() throws Exception {
        return getAny();
    }

    public Response getAny() throws Exception {
        String JOBS = "/jobs";
        String path = request.getPathInfo();
        int idx = path.indexOf(JOBS);
        if (idx < 0) {
            return Response.status(400).entity("No /jobs/ in path, bad routing").build();
        }
        path = path.substring(idx + JOBS.length());
        System.out.println("--------- " + path);

        IrodsConnection conn = Services.get(IrodsConnection.class);
        IrodsObject o = conn.getObject(conn.getInitialPath() + "/" + JOBS_DIR + path);
        if (!o.exists()) {
            return Response.status(404).entity("Object " + o.getFullPath() + " not found").build();
        } else if (o.isFile()) {
            IrodsFile irf = o.asFile();
            String name = irf.getName(), extension = "";
            int i = name.lastIndexOf('.');
            if (i >= 0) {
                name = name.substring(0, i);
                extension = name.substring(i);
            }
            File f = File.createTempFile(name + "-wfk-", extension);
            f.delete();
            irf.downloadToLocalFile(f);

            return Response.ok()
                    .type(URLConnection.guessContentTypeFromName(f.getName()))
                    .entity(new FileInputStream(f)).build();
        } else {
            IrodsCollection jobColl = o.asCollection();
            List<JobJson> jobs = new ArrayList<JobJson>();
            for (IrodsCollection c : jobColl.listCollections()) {
                jobs.add(new JobJson(c.getName(), c.getDate()));
            }
            for (IrodsFile f : jobColl.listFiles()) {
                jobs.add(new JobJson(f.getName(), f.getDate()));
            }

            String json = new Gson().toJson(jobs);
            return Response.ok().entity(json).build();
        }
    }
}
