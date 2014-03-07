package eu.eudat.gef.rest.json;
    
import java.text.DateFormat;
import java.util.Date;

/**
 * @author edima
 */
public class JobJson {

    final static DateFormat dateFormatter = DateFormat.getDateTimeInstance(DateFormat.DEFAULT, DateFormat.SHORT);

    public JobJson(String id, Date date) {
        this.id = id;
        this.date = dateFormatter.format(date);
    }
    public String id;
    public String date;
}
