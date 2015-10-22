package eu.eudat.gef.rest;

import eu.eudat.gef.app.GEFConfig;
import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.Closeable;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.ProtocolException;
import java.net.SocketTimeoutException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import org.apache.commons.io.input.AutoCloseInputStream;
import org.slf4j.LoggerFactory;

/**
 *
 * @author edima
 */
public class ReverseProxy {

	private static final org.slf4j.Logger log = LoggerFactory.getLogger(ReverseProxy.class);
	GEFConfig.GefDocker cfg;

	public ReverseProxy(GEFConfig.GefDocker cfg) {
		this.cfg = cfg;
	}

	public InputStream forward(String path, HttpServletRequest request, HttpServletResponse response) throws MalformedURLException, ProtocolException, IOException {
		String method = request.getMethod();
		URL url = new URL(cfg.url, path);
		String query = request.getQueryString();
		if (query != null && !query.isEmpty()) {
			url = new URL(url.toExternalForm() + "?" + query);
			log.debug("url with query string: " + url);
		}

		final HttpURLConnection conn;
		try {
			conn = (HttpURLConnection) url.openConnection();
		} catch (IOException ex) {
			log.error("exception while opening remote url", ex);
			throw ex;
		}

		conn.setReadTimeout(cfg.timeout);
		conn.setConnectTimeout(cfg.timeout);
		conn.setRequestMethod(method);
		conn.setDoInput(true);
		conn.setDoOutput(true);

		log.debug(request.getMethod() + " " + url);

		for (String header : toList(request.getHeaderNames())) {
			List<String> values = toList(request.getHeaders(header));
			log.debug("request header: " + header + "=" + values);
			conn.setRequestProperty(header, join(values, ","));
		}

		if (method != "GET") { // if we try this on GET we get a 404 (?!)
			try {
				copyAndClose(request.getInputStream(), conn.getOutputStream());
			} catch (IOException ex) {
				log.error("exception while forwarding body", ex);
				throw ex;
			}
		}

		int responseCode;
		InputStream entityStream;
		try {
			responseCode = conn.getResponseCode();
			entityStream = (200 <= responseCode && responseCode < 400) ? conn.getInputStream() : conn.getErrorStream();
		} catch (SocketTimeoutException ex) {
			log.error("timeout from " + url);
			response.setStatus(504); // gateway timeout
			return new ByteArrayInputStream(new byte[0]);
		} catch (IOException ex) {
			log.error("exception while setting status from " + url, ex);
			response.setStatus(502); // gateway error
			return new ByteArrayInputStream(new byte[0]);
		}

		response.setStatus(responseCode);
		log.info("" + url + " : " + responseCode);
		for (Map.Entry<String, List<String>> header : conn.getHeaderFields().entrySet()) {
			if (header.getKey() != null) {
				String values = join(header.getValue(), ",");
				log.debug("response header: " + header.getKey() + "=" + values);
				response.setHeader(header.getKey(), values);
			}
		}

		try {
			response.flushBuffer();
		} catch (IOException ex) {
			log.error("exception while flushing", ex);
			// continue returning request
		}

		return new ExtraAutoCloseInputStream(entityStream, new Closeable() {
			@Override
			public void close() throws IOException {
				conn.disconnect();
			}
		});
	}

	public static class ExtraAutoCloseInputStream extends AutoCloseInputStream {

		Closeable closeable;

		public ExtraAutoCloseInputStream(InputStream is, Closeable closeable) {
			super(is);
			this.closeable = closeable;
		}

		@Override
		public void close() throws IOException {
			super.close();
			closeable.close();
		}
	}

	static void copyAndClose(InputStream inputStream, OutputStream outputStream) throws IOException {
		BufferedInputStream is = new BufferedInputStream(inputStream);
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		BufferedOutputStream os = new BufferedOutputStream(outputStream);

		byte[] b = new byte[64 * 1024];
		for (int n = is.read(b); n > 0; n = is.read(b)) {
			os.write(b, 0, n);
			if (log.isDebugEnabled()) {
				baos.write(b, 0, n);
			}
		}
		is.close();
		os.flush();
		os.close();
		if (log.isDebugEnabled()) {
			log.debug("\n --- copy streams START\n"
					+ baos.toString("UTF-8")
					+ "\n --- copy streams END");
		}
	}

	static <X> List<X> toList(Enumeration<X> enm) {
		List<X> list = new ArrayList<X>();
		while (enm.hasMoreElements()) {
			list.add(enm.nextElement());
		}
		return list;
	}

	static String join(Iterable<String> iterable, String connector) {
		StringBuilder sb = new StringBuilder();
		Iterator<String> it = iterable.iterator();
		while (it.hasNext()) {
			sb.append(it.next());
			sb.append(connector);
		}
		return sb.substring(0, sb.length() - connector.length());
	}
}
