/*
 *  Copyright 2016 EGI Foundation
 * 
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

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

import java.util.regex.Pattern;
import java.util.regex.Matcher;
import java.util.Properties;
import java.util.Set;
import java.util.Arrays;
import java.util.ArrayList;
import java.util.List;

import java.net.URI;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.net.URISyntaxException;

import cz.cesnet.cloud.occi.Model;
import cz.cesnet.cloud.occi.api.Client;
import cz.cesnet.cloud.occi.api.EntityBuilder;
import cz.cesnet.cloud.occi.api.exception.CommunicationException;
import cz.cesnet.cloud.occi.api.exception.EntityBuildingException;
import cz.cesnet.cloud.occi.api.http.HTTPClient;
import cz.cesnet.cloud.occi.api.http.auth.HTTPAuthentication;
import cz.cesnet.cloud.occi.api.http.auth.VOMSAuthentication;
import cz.cesnet.cloud.occi.core.Action;
import cz.cesnet.cloud.occi.core.ActionInstance;
import cz.cesnet.cloud.occi.core.Attribute;
import cz.cesnet.cloud.occi.core.Entity;
import cz.cesnet.cloud.occi.core.Mixin;
import cz.cesnet.cloud.occi.core.Resource;
import cz.cesnet.cloud.occi.core.Kind;
import cz.cesnet.cloud.occi.exception.AmbiguousIdentifierException;
import cz.cesnet.cloud.occi.exception.InvalidAttributeValueException;
import cz.cesnet.cloud.occi.exception.RenderingException;
import cz.cesnet.cloud.occi.infrastructure.Compute;
import cz.cesnet.cloud.occi.parser.MediaType;
import cz.cesnet.cloud.occi.core.Link;
import cz.cesnet.cloud.occi.infrastructure.NetworkInterface;
import cz.cesnet.cloud.occi.infrastructure.IPNetworkInterface;
import cz.cesnet.cloud.occi.infrastructure.Storage;

import org.apache.commons.codec.binary.Base64;

import org.json.simple.JSONObject;

public class InstantiateVM
{
	// Create cloud resources in the selected cloud provider.
	// Available resources that can be created via API are the following:
	// - compute = computing resource, 
	// - storage = storage resources.

	//  Action and type of action
	public static String ACTION = "create";
	public static List<String> RESOURCE = Arrays.asList("compute"); 

	// Authentication and VM title 
	public static String AUTH = null;
	public static String OCCI_PUBLICKEY_NAME = "egieudat";
	public static String TRUSTED_CERT_REPOSITORY_PATH = null;
	public static String PROXY_PATH = null; 
	public static List<String> ATTRIBUTES = Arrays.asList("occi.core.title=EUDAT");


	// Instantiation of all EGI input
	public static String OCCI_ENDPOINT_HOST = null;
	public static String RES_TPL = null;
	public static String OS_TPL = null;
	public static String PUBLIC_KEY_PATH = null;
	public static String CONTEXT_PATH = null;
	public static String instantiatedVmId = null;
	public static List<String> MIXIN = null;
	public static List<String> CONTEXT = null;

	public static Boolean verbose = true;

	// Creating a new VM in the OCCI_ENDPOINT_HOST cloud resource

	public static String doCreate (Properties properties, EntityBuilder eb, Model model, Client client, JSONObject egiInput)
	{
	
		URI uri_location = null;
		String networkInterfaceLocation = "";
		String networkInterfaceLocation_stripped = "";
		Resource vm_resource = null;

		try 
		{

			if (properties.getProperty("RESOURCE").equals("compute")) 
			{

				String segments[] = properties.getProperty("OCCI_OS_TPL").split("#");
				String OCCI_OS_TPL = segments[segments.length - 1];

				String segments2[] = properties.getProperty("OCCI_RESOURCE_TPL").split("#");
				String OCCI_RESOURCE_TPL = segments2[segments2.length - 1];

				System.out.println("[+] Creating a new compute Virtual Machine (VM)");

				// Creating a compute instance
				Resource compute = eb.getResource("compute");
				Mixin mixin = model.findMixin(OCCI_OS_TPL);
					compute.addMixin(mixin);
					compute.addMixin(model.findMixin(OCCI_OS_TPL, "os_tpl"));
					compute.addMixin(model.findMixin(OCCI_RESOURCE_TPL, "resource_tpl"));

				// Checking the context
				if (properties.getProperty("PUBLIC_KEY_FILE") != null && 
					!properties.getProperty("PUBLIC_KEY_FILE").isEmpty()) 
				{				
					String _public_key_file = 
						properties.getProperty("PUBLIC_KEY_FILE").substring(properties.getProperty("PUBLIC_KEY_FILE").lastIndexOf(":") + 1);

					File f = new File(_public_key_file);

					FileInputStream fis = new FileInputStream(f);
					DataInputStream dis = new DataInputStream(fis);
					byte[] keyBytes = new byte[(int) f.length()];
					dis.readFully(keyBytes);
					dis.close();
					String _publicKey = new String (keyBytes).trim();

					// Add SSH public key
					compute.addMixin(model.findMixin(URI.create("http://schemas.openstack.org/instance/credentials#public_key")));
					compute.addAttribute("org.openstack.credentials.publickey.data", _publicKey);
			
					// Add the name for the public key	
					if (OCCI_PUBLICKEY_NAME != null && !OCCI_PUBLICKEY_NAME.isEmpty()) 
						compute.addAttribute("org.openstack.credentials.publickey.name",
						properties.getProperty("OCCI_PUBLICKEY_NAME"));
				} 

				if (properties.getProperty("USER_DATA") != null && 
					!properties.getProperty("USER_DATA").isEmpty()) 
				{
					String _user_data =
						properties.getProperty("USER_DATA").substring(properties.getProperty("USER_DATA").lastIndexOf(":") + 1);

						File f = new File(_user_data);
						FileInputStream fis = new FileInputStream(f);
						DataInputStream dis = new DataInputStream(fis);
						byte[] keyBytes = new byte[(int) f.length()];
						dis.readFully(keyBytes);
						dis.close();
						byte[] data = Base64.encodeBase64(keyBytes);
						String user_data = new String (data);

					compute.addMixin(model.findMixin(URI.create("http://schemas.openstack.org/compute/instance#user_data")));
			
					compute.addAttribute("org.openstack.compute.user_data", user_data);
				}

				// Set VM title
				compute.setTitle(properties.getProperty("OCCI_CORE_TITLE"));
				URI location = client.create(compute);

				return location.toString();		

			} 
			
			if (properties.getProperty("RESOURCE").equals("storage")) 
			{
	 			System.out.println("[+] Creating a volume storage");

				// Creating a storage instance
				Storage storage = eb.getStorage();
	 			storage.setTitle(properties.getProperty("OCCI_CORE_TITLE"));
				storage.setSize(properties.getProperty("OCCI_STORAGE_SIZE"));

				URI storageLocation = client.create(storage);
				
				List<URI> list = client.list("storage");
				List<URI> storageURIs = new ArrayList<URI>();

				for (URI uri : list) 
				{
					if (uri.toString().contains("storage")) 
						storageURIs.add(uri);
				}
						
				System.out.println("URI = " + storageLocation);
			} 

		} 

		catch (FileNotFoundException ex) 
		{throw new RuntimeException(ex);}

		catch (IOException ex) 
		{throw new RuntimeException(ex);}

		catch (EntityBuildingException | AmbiguousIdentifierException |
			InvalidAttributeValueException | CommunicationException ex) 
		{throw new RuntimeException(ex);}

		return "";
	}

	public static String instantiateVM(JSONObject egiInput)
	{
		OCCI_ENDPOINT_HOST = (String) egiInput.get("endpoint");
		RES_TPL = (String) egiInput.get("resourceTpl");
		OS_TPL = (String) egiInput.get("osTpl");
		PUBLIC_KEY_PATH = (String) egiInput.get("publicKey");
		CONTEXT_PATH = (String) egiInput.get("contextualisation");
		AUTH = (String) egiInput.get("auth");
		TRUSTED_CERT_REPOSITORY_PATH = (String) egiInput.get("trustedCertificatesPath");
		PROXY_PATH = (String) egiInput.get("proxyPath");

		MIXIN = Arrays.asList(RES_TPL, 
		OS_TPL);

		CONTEXT = Arrays.asList("public_key="+PUBLIC_KEY_PATH, 
		"user_data="+CONTEXT_PATH); 

		Boolean result = false;
		String networkInterfaceLocation = "";
		String networkInterfaceLocation_stripped = "";
		Resource vm_resource = null;
		URI uri_location = null;

		if (verbose) 
		{
			System.out.println();
			if (ACTION != null && !ACTION.isEmpty()) 
				System.out.println("[ACTION] = " + ACTION);
			else	
				System.out.println("[ACTION] = Get dump model");
			System.out.println("AUTH = " + AUTH);
			if (OCCI_ENDPOINT_HOST != null && !OCCI_ENDPOINT_HOST.isEmpty()) 
				System.out.println("OCCI_ENDPOINT_HOST = " + OCCI_ENDPOINT_HOST);
			if (RESOURCE != null && !RESOURCE.isEmpty()) 
				System.out.println("RESOURCE = " + RESOURCE);
			if (MIXIN != null && !MIXIN.isEmpty()) 
				System.out.println("MIXIN = " + MIXIN);
			if (TRUSTED_CERT_REPOSITORY_PATH != null && !TRUSTED_CERT_REPOSITORY_PATH.isEmpty()) 
				System.out.println("TRUSTED_CERT_REPOSITORY_PATH = " + TRUSTED_CERT_REPOSITORY_PATH);
			if (PROXY_PATH != null && !PROXY_PATH.isEmpty()) 
				System.out.println("PROXY_PATH = " + PROXY_PATH);
			if (CONTEXT != null && !CONTEXT.isEmpty()) 
				System.out.println("CONTEXT = " + CONTEXT);
			if (OCCI_PUBLICKEY_NAME != null && !OCCI_PUBLICKEY_NAME.isEmpty()) 
				System.out.println("OCCI_PUBLICKEY_NAME = " + OCCI_PUBLICKEY_NAME);
			if (ATTRIBUTES != null && !ATTRIBUTES.isEmpty()) 
				System.out.println("ATTRIBUTES = " + ATTRIBUTES);
			if (verbose) System.out.println("Verbose = True ");
			else System.out.println("Verbose = False ");
		}

		Properties properties = new Properties();

		if (ACTION != null && !ACTION.isEmpty())
			properties.setProperty("ACTION", ACTION);

		if (OCCI_ENDPOINT_HOST != null && !OCCI_ENDPOINT_HOST.isEmpty())
			properties.setProperty("OCCI_ENDPOINT_HOST", OCCI_ENDPOINT_HOST);

		if (RESOURCE != null && !RESOURCE.isEmpty()) 
			for (int i=0; i<RESOURCE.size(); i++) 
			{
				if ((!RESOURCE.get(i).equals("compute")) && 
					(!RESOURCE.get(i).equals("storage")) &&
					(!RESOURCE.get(i).equals("network")) &&
					(!RESOURCE.get(i).equals("os_tpl")) &&
					(!RESOURCE.get(i).equals("resource_tpl")))
					properties.setProperty("OCCI_VM_RESOURCE_ID", RESOURCE.get(i));
				else 
				{ 
				properties.setProperty("RESOURCE", RESOURCE.get(i));
				properties.setProperty("OCCI_VM_RESOURCE_ID", "empty");
				}
			}
			
		if (MIXIN != null && !MIXIN.isEmpty()) 
			for (int i=0; i<MIXIN.size(); i++) 
			{
				if (MIXIN.get(i).contains("template") || 
					MIXIN.get(i).contains("os_tpl")) 
					properties.setProperty("OCCI_OS_TPL", MIXIN.get(i));

				if (MIXIN.get(i).contains("resource_tpl")) 
					properties.setProperty("OCCI_RESOURCE_TPL", MIXIN.get(i));
			}

		if (ATTRIBUTES != null && !ATTRIBUTES.isEmpty())
			for (int i=0; i<ATTRIBUTES.size(); i++) 
			{
				if (ATTRIBUTES.get(i).contains("occi.core.title")) 
				{
					String _OCCI_CORE_TITLE = ATTRIBUTES.get(i)
					.substring(ATTRIBUTES.get(i).lastIndexOf("=") + 1);

					properties.setProperty("OCCI_CORE_TITLE", _OCCI_CORE_TITLE);
				}

				if (ATTRIBUTES.get(i).contains("occi.storage.size")) 
				{
					String _OCCI_STORAGE_SIZE = ATTRIBUTES.get(i)
					.substring(ATTRIBUTES.get(i).lastIndexOf("=") + 1);

					properties.setProperty("OCCI_STORAGE_SIZE", _OCCI_STORAGE_SIZE);
				}
			}

		properties.setProperty("TRUSTED_CERT_REPOSITORY_PATH", TRUSTED_CERT_REPOSITORY_PATH);
		properties.setProperty("PROXY_PATH", PROXY_PATH);

		if (CONTEXT != null && !CONTEXT.isEmpty()) 
		{
			for (int i=0; i<CONTEXT.size(); i++) 
			{
				if (CONTEXT.get(i).contains("public_key")) 
				properties.setProperty("PUBLIC_KEY_FILE", CONTEXT.get(i));

				if (CONTEXT.get(i).contains("user_data")) 
				properties.setProperty("USER_DATA", CONTEXT.get(i));
			}
		}

		if (OCCI_PUBLICKEY_NAME != null && !OCCI_PUBLICKEY_NAME.isEmpty())
			properties.setProperty("OCCI_PUBLICKEY_NAME", OCCI_PUBLICKEY_NAME);
		properties.setProperty("OCCI_AUTH", AUTH);

		try 
		{
			HTTPAuthentication authentication = new VOMSAuthentication(PROXY_PATH);

			authentication.setCAPath(TRUSTED_CERT_REPOSITORY_PATH);

			Client client = new HTTPClient(URI.create(OCCI_ENDPOINT_HOST),
			authentication, MediaType.TEXT_PLAIN, false);

			//Connect client
			client.connect();

			Model model = client.getModel();
			EntityBuilder eb = new EntityBuilder(model);

			instantiatedVmId = doCreate(properties, eb, model, client, egiInput);
			return instantiatedVmId;
		}
		catch (CommunicationException ex) 
		{throw new RuntimeException(ex);}

	}
}
