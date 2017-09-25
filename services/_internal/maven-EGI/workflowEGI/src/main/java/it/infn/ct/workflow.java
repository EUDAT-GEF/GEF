package it.infn.ct;

import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.io.FileWriter;

import java.lang.Thread;

import org.json.simple.JSONObject;
import org.json.simple.JSONArray;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;

public class workflow 
{

	public static String vmState = "inactive";
	public static String configEGI = "/home/workflowEGI/configEGI.json";
	public static String[] vmFeatures = null;

	public static void main(String[] args) {

		// Get input from EGI
		Def def = new Def();
		String[] egiList = def.readJson(configEGI);

		// Instantiate a new VM on EGI based on the input from configEGI.json
		InstantiateVM newVM = new InstantiateVM();
		String vmId = new String (newVM.instantiateVM(egiList));
		System.out.println("New Virtual Machine ID: " + vmId);

		// Describe the VM state - Once this part is successfull the VM is operational
		while (!vmState.equals("active")) 
		{

			try 
			{
				System.out.println("Waiting for the Virtual Machine to be active...");
				//Thread.sleep(8000);
				vmFeatures = DescribeVM.describe(vmId);
				vmState = vmFeatures[1];
				Thread.sleep(8000);
				System.out.println("Virtual machine is " + vmFeatures[1]);
			} catch(InterruptedException ex) {Thread.currentThread().interrupt();}
		}

		//Write the public IP and vmState on the JSON file
		def.WriteJson(vmId,vmFeatures[0],vmFeatures[1],configEGI);

	}
}