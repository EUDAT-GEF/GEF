package eu.eudat.gef.rest;

import com.google.gson.Gson;
import com.sun.jersey.core.header.FormDataContentDisposition;
import com.sun.jersey.multipart.FormDataBodyPart;
import com.sun.jersey.multipart.FormDataParam;
import de.tuebingen.uni.sfs.epicpid.Pid;
import de.tuebingen.uni.sfs.epicpid.PidServer;
import eu.eudat.gef.app.Services;
import eu.eudat.gef.irodslink.IrodsCollection;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.IrodsFile;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.io.File;
import java.io.InputStream;
import java.net.URI;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.slf4j.LoggerFactory;

/**
 * @author edima
 */
@Path("datasets")
@Produces(MediaType.APPLICATION_JSON)
public class DataSets {

	static org.slf4j.Logger log = LoggerFactory.getLogger(DataSets.class);

	public static final String DATA_DIR = "datasets";

	@Context
	HttpServletRequest request;
	@Context
	HttpServletResponse response;

	@POST
	@Consumes(MediaType.MULTIPART_FORM_DATA)
	public Response uploadMultiple(@FormDataParam("file") List<FormDataBodyPart> fdbpList) throws Exception {
		log.info("upload multiple: " + fdbpList.size());

		Pid pid = Services.get(PidServer.class).makePid("", null, null);
		IrodsConnection conn = Services.get(IrodsConnection.class);
		String newColl = conn.getInitialPath() + "/" + DATA_DIR + "/" + pid.getId();
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
			IrodsConnection conn, IrodsCollection coll) throws Exception {
		log.info("upload: " + fileDetail.getType() + "; " + fileDetail.getName() + "; " + fileDetail.getFileName());

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
		while (name.length() < 3) {
			name += "_";
		}
		File f = File.createTempFile(name, ext);
		f.delete();
		Files.copy(inputStream, f.toPath());

		IrodsFile ifile = conn.getObject(coll.getFullPath() + "/" + name + ext).asFile();
		log.info("upload from " + f.getPath() + " to " + ifile.getFullPath());
		ifile.uploadFromLocalFile(f);
		f.delete();
		return coll + "/" + ifile.getName();
	}

	public static class DatasetJson {

		public DatasetJson(String id, String name, int size, long date) {
			this.id = id;
			this.name = name;
			this.size = size;
			this.date = date;
		}
		public String id;
		public String name;
		public long date;
		public int size;
	}

	@GET
	public Response getObjects() throws Exception {
		IrodsConnection conn = Services.get(IrodsConnection.class);
		String dataCollPath = conn.getInitialPath() + "/" + DATA_DIR;
		IrodsCollection coll = conn.getObject(dataCollPath).asCollection();
		if (!coll.exists()) {
			coll.create();
		}

		List<DatasetJson> list = new ArrayList<>();
		for (IrodsCollection c : coll.listCollections()) {
			for (IrodsCollection c2 : c.listCollections()) {
				int size = 0;
				Collection<IrodsFile> files = c2.listFiles();
				for (IrodsFile f : files) {
					size += f.getSize();
				}
				String id = c.getName() + "/" + c2.getName();
				String name = files.isEmpty() ? "" : files.iterator().next().getName();
				list.add(new DatasetJson(id, name, size, c2.getDate().getTime()));
			}
		}

		Map<String, Object> map = new HashMap<String, Object>();
		map.put("datasets", list);
		String json = new Gson().toJson(map);

		return Response.ok().entity(json).build();
	}

	private void println(String string) {
		System.out.println(string);
	}
}
