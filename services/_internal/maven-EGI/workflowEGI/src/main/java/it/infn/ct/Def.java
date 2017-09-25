package it.infn.ct;

import java.io.File;
import java.io.FileInputStream;
import java.io.DataInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.FileWriter;
import java.io.FileReader;

import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;


public class Def

{
    public static void WriteJson (String vmId, String publicIP, String vmState, String jsonfile)
    {
        JSONParser parser = new JSONParser();
        JSONObject jsonObject = null;

        try
        {
            Object obj = parser.parse(new FileReader(jsonfile));
            jsonObject = (JSONObject) obj;
            jsonObject.put("vmId",vmId);
            jsonObject.put("publicIP",publicIP);
            jsonObject.put("vmState",vmState);
        } 
        catch (FileNotFoundException e) 
        {e.printStackTrace();}
<<<<<<< HEAD

        catch (IOException e) 
        {e.printStackTrace();} 

        catch (ParseException e) 
        {e.printStackTrace();}

       
=======

        catch (IOException e) 
        {e.printStackTrace();} 

        catch (ParseException e) 
        {e.printStackTrace();}


>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
        try (FileWriter file = new FileWriter(jsonfile)) 
        {
            file.write(jsonObject.toJSONString());
            file.flush();
        }
        catch (IOException e) 
        {e.printStackTrace();}

    }

    public static String[] readJson (String jsonfile)
    {
        JSONParser parser = new JSONParser();

        String endpoint = "", resource_tpl = "", os_tpl = "", publicKey = "", contextualisation = "", proxy = "";
        String[] egiList = {};
        try 
        {

            Object obj = parser.parse(new FileReader(jsonfile));

            JSONObject jsonObject = (JSONObject) obj;

            proxy = (String) jsonObject.get("proxy");
            endpoint = (String) jsonObject.get("endpoint");
            resource_tpl = (String) jsonObject.get("resource_tpl");
            os_tpl = (String) jsonObject.get("os_tpl");
            publicKey = (String) jsonObject.get("publicKey");
            contextualisation = (String) jsonObject.get("contextualisation");

            egiList = new String[]{proxy,endpoint,resource_tpl,os_tpl,publicKey,contextualisation};

        } 
        catch (FileNotFoundException e) 
        {e.printStackTrace();}

        catch (IOException e) 
        {e.printStackTrace();}

        catch (ParseException e) 
        {e.printStackTrace();}

        return egiList;
    }
}
