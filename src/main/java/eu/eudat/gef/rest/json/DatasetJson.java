package eu.eudat.gef.rest.json;

import java.text.DateFormat;
import java.util.Date;

/**
 * @author edima
 */
public class DatasetJson {

    final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

    public DatasetJson(String id, String name, int size, Date date) {
        this.id = id;
        this.name = name;
        this.size = size;
        this.date = dateFormatter.format(date);
    }
    public String id;
    public String name;
    public String date;
    public int size;
}
