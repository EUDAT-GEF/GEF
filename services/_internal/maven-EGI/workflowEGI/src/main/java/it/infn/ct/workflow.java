package it.infn.ct;

import java.lang.Thread;

import org.json.simple.JSONObject;

public class workflow 
{

	public static String vmState = "inactive";
	public static String configEGI = "/home/configEGI.json";//"/Users/pivan/CERFACS/GEF/GEF_push_copy_20_09/services/_internal/maven-EGI/configEGI_local.json";//
	public static String[] vmFeatures = null;

	public static void main(String[] args) {
		// Get input from EGI
		Def def = new Def();
<<<<<<< HEAD
		JSONObject egiInput = def.getJson(configEGI);
=======
		String[] egiList = def.readJson(configEGI);
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c

		// Instantiate a new VM on EGI based on the input from configEGI.json
		InstantiateVM newVM = new InstantiateVM();
		String vmId = new String (newVM.instantiateVM(egiInput));
		System.out.println("Waiting for the Virtual Machine to be active...");
		// Describe the VM state - Once this part is successfull the VM is operational
		while (!vmState.equals("active")) 
		{
<<<<<<< HEAD
			try 
			{
				vmFeatures = DescribeVM.describe(vmId, egiInput);
				vmState = vmFeatures[1];
				Thread.sleep(8000);
				System.out.println("Virtual machine state is " + vmFeatures[1]);
=======

			try 
			{
				System.out.println("Waiting for the Virtual Machine to be active...");
				//Thread.sleep(8000);
				vmFeatures = DescribeVM.describe(vmId);
				vmState = vmFeatures[1];
				Thread.sleep(8000);
				System.out.println("Virtual machine is " + vmFeatures[1]);
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
			} catch(InterruptedException ex) {Thread.currentThread().interrupt();}
		}

		//Write the public IP and vmState on the JSON file
		def.WriteJson(vmId,vmFeatures[0],vmFeatures[1],configEGI);

	}
}