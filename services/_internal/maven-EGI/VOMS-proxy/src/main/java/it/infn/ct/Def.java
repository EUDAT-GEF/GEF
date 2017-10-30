package it.infn.ct;

import java.io.File;
import java.io.FileInputStream;
import java.io.DataInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.FileWriter;
import java.io.FileReader;

import java.util.ArrayList;

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

		catch (IOException e) 
		{e.printStackTrace();} 

		catch (ParseException e) 
		{e.printStackTrace();}

		try (FileWriter file = new FileWriter(jsonfile)) 
		{
			file.write(jsonObject.toJSONString());
			file.flush();
		}
		catch (IOException e) 
		{e.printStackTrace();}
	}

	public static JSONObject getJson (String jsonfile)
	{
		JSONParser parser = new JSONParser();
		JSONObject jsonObject = new JSONObject();
		try 
		{
			Object obj = parser.parse(new FileReader(jsonfile));
			jsonObject = (JSONObject) obj;
		} 
		catch (FileNotFoundException e) 
		{e.printStackTrace();}

		catch (IOException e) 
		{e.printStackTrace();}

		catch (ParseException e) 
		{e.printStackTrace();}

		return jsonObject;
	}
}
