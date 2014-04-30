package eu.eudat.gef.irods;

import com.google.common.base.Charsets;
import com.google.common.base.Joiner;
import com.google.common.io.CharStreams;
import com.google.common.io.Files;
import eu.eudat.gef.service.ResourceSupplier;
import eu.eudat.gef.service.Services;
import eu.eudat.gef.irodslink.IrodsConnection;
import eu.eudat.gef.irodslink.IrodsException;
import eu.eudat.gef.irodslink.IrodsFile;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.math.BigInteger;
import java.net.URI;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.concurrent.Callable;

/**
 * @author edima
 */
public class IrodsWso implements Callable<String> {
	static final String WSO_RUNS_DIR = "wso/runs";
	private SecureRandom random = new SecureRandom();
	File stageDir;
	IrodsFile dataFile;
	String wsoMpfPrefix;
	String parameters;
	URI runDirUri;

	public IrodsWso(IrodsFile dataFile, String parameters) throws IOException, IrodsException {
		this.dataFile = dataFile;
		this.parameters = parameters;
	}

	public String randomId() {
		return new BigInteger(40, random).toString(16);
	}

	public String str(InputStream is) throws IOException {
		try {
			return CharStreams.toString(new InputStreamReader(is, Charsets.UTF_8));
		} finally {
			is.close();
		}
	}

	public String call() throws Exception {
		String text = CharStreams.toString(CharStreams.newReaderSupplier(
				new ResourceSupplier("/wso-template.mpf"), Charsets.UTF_8));

		stageDir = new File("/tmp/wsostage" + randomId());

		text = text.replaceAll("\\$StageDir", escape(stageDir.getPath()));
		text = text.replaceAll("\\$FileName", escape(dataFile.getName()));
		text = text.replaceAll("\\$Filter", escape(parameters));
		text = text.replaceAll("\\$IrodsFilePathAndName", escape(dataFile.getFullPath()));

		File tmp = Files.createTempDir();
		String wsoinstance = "wso-" + randomId();

		File wsoMpf = new File(tmp, wsoinstance + ".mpf");
		File wsoRun = new File(tmp, wsoinstance + ".run");
		Files.write(text, wsoMpf, Charsets.UTF_8);

		IrodsConnection irodsConn = Services.get(IrodsConnection.class);
		String irodsWsoInst = irodsConn.getInitialPath() + "/" + WSO_RUNS_DIR + "/";

		//IrodsFile irodsWsoMpf = irodsConn.getObject(irodsWsoMpfPath).asFile();
		//irodsWsoMpf.uploadFromLocalFile(wsompf);
		String command[] = new String[]{"/usr/local/bin/iput", wsoMpf.getPath(), irodsWsoInst};
		System.out.println(Arrays.asList(command));
		Process p = Runtime.getRuntime().exec(command);
		System.out.println("" + p.waitFor());
		System.out.println(str(p.getInputStream()));
		System.out.println(str(p.getErrorStream()));

		String irodsWsoRunDir = irodsWsoInst + wsoRun.getName() + "Dir";
		runDirUri = irodsConn.makeUri(irodsConn.getObject(irodsWsoRunDir));

		//IrodsFile irodsWsoRun = irodsConn.getObject(irodsWsoRunPath).asFile();

//		while (!irodsWsoRun.exists()) {
//			Thread.sleep(100);
//		}
//
//		irodsWsoRun.downloadToLocalFile(wsorun);

		command = new String[]{"/usr/local/bin/iget", irodsWsoInst + wsoRun.getName(), "-"};
		System.out.println(Arrays.asList(command));
		p = Runtime.getRuntime().exec(command);
		System.out.println("" + p.waitFor());
		String out = str(p.getInputStream());
		System.out.println(out);
		System.out.println(str(p.getErrorStream()));
		
		return out;
		//return Joiner.on("\n").join(Files.readLines(wsorun, Charsets.UTF_8));
	}

	private String escape(String str) {
		return str.replaceAll("\"", "'");
	}

	public URI getRunDirUri() {
		return runDirUri;
	}
}
