package eu.eudat.gef.app;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.dropwizard.Configuration;
import java.net.MalformedURLException;
import java.net.URL;
import org.hibernate.validator.constraints.NotEmpty;

public class GEFConfig extends Configuration {

	public static class Params {

		@NotEmpty
		@JsonProperty
		Irods irods;

		@NotEmpty
		@JsonProperty
		Pid pid;

		@NotEmpty
		@JsonProperty
		public GefDocker gefDocker;
	}

	public static class Irods {

		@NotEmpty
		@JsonProperty
		String server;

		@NotEmpty
		@JsonProperty
		int port;

		@NotEmpty
		@JsonProperty
		String username;

		@NotEmpty
		@JsonProperty
		String password;

		@NotEmpty
		@JsonProperty
		String path;

		@NotEmpty
		@JsonProperty
		String resource;

	}

	public static class Pid {

		@NotEmpty
		@JsonProperty
		String epicServerUrl;

		@NotEmpty
		@JsonProperty
		String localPrefix;

		@NotEmpty
		@JsonProperty
		String user;

		@NotEmpty
		@JsonProperty
		String pass;
	}

	public static class GefDocker {
		@NotEmpty
		@JsonProperty
		public URL url;

		@NotEmpty
		@JsonProperty
		public int timeout = 2000; // 2 secs
	}

	public Params gefParams = new Params();

	public static Pid makePid(String epicServerUrl, String localPrefix, String user, String pass) {
		Pid pid = new Pid();
		pid.epicServerUrl = epicServerUrl;
		pid.localPrefix = localPrefix;
		pid.user = user;
		pid.pass = pass;
		return pid;
	}

	public static Irods makeIrods(String server, int port, String username, String password, String path, String resource) {
		Irods irods = new Irods();
		irods.server = server;
		irods.port = port;
		irods.username = username;
		irods.password = password;
		irods.path = path;
		irods.resource = resource;
		return irods;
	}

	public static GefDocker makeDocker(String url) throws MalformedURLException {
		GefDocker gefDocker = new GefDocker();
		gefDocker.url = new URL(url);
		return gefDocker;
	}
}
