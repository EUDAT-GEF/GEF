package eu.eudat.gef.service;

import java.io.File;
import javax.servlet.ServletContextEvent;
import javax.servlet.ServletContextListener;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.xpath.XPathConstants;
import javax.xml.xpath.XPathFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.w3c.dom.Document;
import org.w3c.dom.Node;

/**
 * @author edima
 */
public class WebAppConfig implements ServletContextListener {
	final static String CONFIG_FILE_LOCATION = "/WEB-INF/config.xml";
	final static Logger log = LoggerFactory.getLogger(Services.class);
	static Document xdoc = null;
	static XPathFactory xpf = null;

	public void contextInitialized(ServletContextEvent sce) {
		try {
			String configFilePath = sce.getServletContext().getRealPath(CONFIG_FILE_LOCATION);
			if (configFilePath == null) {
				throw new RuntimeException(
						"cannot find the configuration file, it should be located at $"
						+ CONFIG_FILE_LOCATION);
			}
			File file = new File(configFilePath);
			xdoc = DocumentBuilderFactory.newInstance().newDocumentBuilder().parse(file);
			xpf = XPathFactory.newInstance();
		} catch (Exception ex) {
			log.error("exception while initializing GEF servlet: " + ex.getMessage(), ex);
			throw new RuntimeException(ex);
		}
	}

	public static String get(String xpath) {
		try {
			Node n = (Node) xpf.newXPath().evaluate(xpath, xdoc, XPathConstants.NODE);
			return n.getTextContent();
		} catch (Exception e) {
			return null;
		}
	}

	public static int getInt(String xpath) {
		String x = get(xpath);
		try {
			return Integer.parseInt(x);
		} catch (NumberFormatException ex) {
			log.error("number expected: " + ex.getMessage(), ex);
			throw ex;
		}
	}

	public void contextDestroyed(ServletContextEvent sce) {
	}
}
