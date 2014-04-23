package eu.eudat.gef.rest;

import com.google.common.io.ByteStreams;
import com.google.gson.Gson;
import com.sun.jersey.core.header.FormDataContentDisposition;
import com.sun.jersey.multipart.FormDataBodyPart;
import com.sun.jersey.multipart.FormDataParam;
import de.tuebingen.uni.sfs.epicpid.Pid;
import de.tuebingen.uni.sfs.epicpid.PidServer;
import eu.eudat.gef.Services;
import eu.eudat.gef.irodslink.IrodsCollection;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.IrodsFile;
import eu.eudat.gef.rest.json.DatasetJson;
import eu.eudat.gef.rest.json.WorkflowJson;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.URI;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

/**
 * @author edima
 */
@Path("workflows")
@Produces(MediaType.APPLICATION_JSON)
public class Workflows {
	public static final String WORKFLOW_DIR = "workflows";
	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	@POST
	@Consumes(MediaType.MULTIPART_FORM_DATA)
	public Response uploadMultiple(@FormDataParam("file") List<FormDataBodyPart> fdbpList) throws Exception {
		println("upload multiple: " + fdbpList.size());

		Pid pid = Services.get(PidServer.class).makePid("", null, null);
		IrodsConnection conn = Services.get(IrodsConnection.class);
		String newColl = conn.getInitialPath() + "/" + WORKFLOW_DIR + "/" + pid.getId();
		IrodsCollection collWfl = conn.getObject(newColl).asCollection();
		collWfl.create();
		URI collUri = conn.makeUri(collWfl);
		pid.changeUrlTo(collUri);

		for (FormDataBodyPart fdbp : fdbpList) {
			InputStream is = fdbp.getEntityAs(InputStream.class);
			FormDataContentDisposition cd = fdbp.getFormDataContentDisposition();
			uploadFile(is, cd, conn, collWfl);
		}

		String json = new Gson().toJson(collUri);
		return Response.created(collUri).entity(json).build();
	}

 public String uploadFile(InputStream inputStream, FormDataContentDisposition fileDetail,
			IrodsConnection conn, IrodsCollection collWfl) throws Exception {
		println("upload: " + fileDetail.getType() + "; " + fileDetail.getName() + "; " + fileDetail.getFileName());

		String name = fileDetail.getFileName();
		if (name == null || name.isEmpty()) {
			name = fileDetail.getName();
		}
		int idx = name.lastIndexOf(".");
		String ext = "";
		if (idx > 0) {
			ext = name.substring(idx);
			name = name.substring(0, name.length() - ext.length());
		}

		File f = File.createTempFile(name, ext);
		final OutputStream out = new FileOutputStream(f);
		try {
			ByteStreams.copy(inputStream, out);
		} finally {
			out.close();
		}

		IrodsFile ifile = conn.getObject(collWfl.getFullPath() + "/" +  name + ext).asFile();
		System.out.println("upload from " + f.getPath() + " to " + ifile.getFullPath());
		ifile.uploadFromLocalFile(f);
		f.delete();
		return collWfl + "/" + ifile.getName();
	}

	@GET
	public Response getObjects() throws Exception {
		IrodsConnection conn = Services.get(IrodsConnection.class);
		String dataCollPath = conn.getInitialPath() + "/" + WORKFLOW_DIR;
		IrodsCollection coll = conn.getObject(dataCollPath).asCollection();
		if (!coll.exists()) {
			coll.create();
		}

		List<WorkflowJson> list = new ArrayList<WorkflowJson>();
		for (IrodsCollection c : coll.listCollections()) {
			for (IrodsCollection c2 : c.listCollections()) {
                                Collection<IrodsFile> files = c2.listFiles();
                                String id = c.getName() + "/" + c2.getName();
                                String name = files.isEmpty() ? "" : files.iterator().next().getName();
				list.add(new WorkflowJson(id, name, c2.getDate()));
			}
		}

		String json = new Gson().toJson(list);
		return Response.ok().entity(json).build();
	}

	private void println(String string) {
		System.out.println(string);
	}
}
