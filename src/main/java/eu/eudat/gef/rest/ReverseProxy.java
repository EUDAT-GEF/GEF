package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEFConfig;
import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.ProtocolException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 *
 * @author edima
 */
public class ReverseProxy {
	GEFConfig.GefDocker cfg;

	public ReverseProxy(GEFConfig.GefDocker cfg) {
		this.cfg = cfg;
	}

	public void forward(String path, HttpServletRequest request, HttpServletResponse response) throws ProtocolException, IOException {
		URL url = new URL(cfg.url, path);
		String query = request.getQueryString();
		if (query != null && !query.isEmpty()) {
			debug("Adding query string: " + query);
			url = new URL(url.toExternalForm() + "?" + query);
			debug("New url: " + url);
		}

		HttpURLConnection conn = (HttpURLConnection) url.openConnection();

		conn.setReadTimeout(cfg.timeout);
		conn.setConnectTimeout(cfg.timeout);
		conn.setRequestMethod(request.getMethod());
		conn.setDoInput(true);
		conn.setDoOutput(true);

		for (String header : toList(request.getHeaderNames())) {
			List<String> values = toList(request.getHeaders(header));
			conn.setRequestProperty(header, join(values, ","));
		}

		copyAndClose(request.getInputStream(), conn.getOutputStream());

		response.setStatus(conn.getResponseCode());
		for (Map.Entry<String, List<String>> header : conn.getHeaderFields().entrySet()) {
			if (header.getKey() != null) {
				String values = join(header.getValue(), ",");
				debug("Set response header: " + header.getKey() + "=" + values);
				response.setHeader(header.getKey(), values);
			}
		}
		copyAndClose(conn.getInputStream(), response.getOutputStream());
		conn.disconnect();
		response.flushBuffer();
	}

	private static void copyAndClose(InputStream inputStream, OutputStream outputStream) throws IOException {
		BufferedInputStream is = new BufferedInputStream(inputStream);
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		BufferedOutputStream os = new BufferedOutputStream(outputStream);

		byte[] b = new byte[64 * 1024];
		for (int n = is.read(b); n > 0; n = is.read(b)) {
			os.write(b, 0, n);
			baos.write(b, 0, n);
		}
		is.close();
		os.flush();
		os.close();
		debug("\n --- copy streams START");
		debug(baos.toString("UTF-8"));
		debug("\n --- copy streams END");
	}

	private static <X> List<X> toList(Enumeration<X> enm) {
		List<X> list = new ArrayList<X>();
		while (enm.hasMoreElements()) {
			list.add(enm.nextElement());
		}
		return list;
	}

	private static String join(Iterable<String> iterable, String connector) {
		StringBuilder sb = new StringBuilder();
		Iterator<String> it = iterable.iterator();
		while (it.hasNext()) {
			sb.append(it.next());
			sb.append(connector);
		}
		return sb.substring(0, sb.length() - connector.length());
	}

	private static void debug(String s) {
//		System.out.println(s);
	}
}
