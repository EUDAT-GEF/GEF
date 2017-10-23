package it.infn.ct;

import java.lang.Thread;

import org.json.simple.JSONObject;

public class workflow 
{

	public static String vmState = "inactive";
	public static String configEGI = "/home/configEGI.json";
	public static String[] vmFeatures = null;

	public static void main(String[] args) {
		// Get input from EGI
		Def def = new Def();
		JSONObject egiInput = def.getJson(configEGI);

		// Instantiate a new VM on EGI based on the input from configEGI.json
		InstantiateVM newVM = new InstantiateVM();
		String vmId = new String (newVM.instantiateVM(egiInput));
		System.out.println("Waiting for the Virtual Machine to be active...");
		
		// Describe the VM state - Once this part is successfull the VM is operational
		while (!vmState.equals("active")) 
		{
			try 
			{
				vmFeatures = DescribeVM.describe(vmId, egiInput);
				vmState = vmFeatures[1];
				Thread.sleep(8000);
				System.out.println("Virtual machine state: " + vmFeatures[1]);
			} catch(InterruptedException ex) {Thread.currentThread().interrupt();}
		}

		//Write the public IP and vmState on the JSON file
		def.WriteJson(vmId,vmFeatures[0],vmFeatures[1],configEGI);

	}
}