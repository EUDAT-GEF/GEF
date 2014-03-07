package eu.eudat.gef.rest.json;

import java.text.DateFormat;
import java.util.Date;

/**
 * @author edima
 */
public class WorkflowJson {
	final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

	public WorkflowJson(String id, String name, Date date) {
                this.id = id;
		this.name = name;
		this.date = dateFormatter.format(date);
	}
        public String id;
	public String name;
	public String date;
}
