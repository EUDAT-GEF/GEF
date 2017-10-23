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

import java.util.Arrays;
import java.util.Map;
import java.util.List;
import java.util.Properties;
import java.net.URI;

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
import cz.cesnet.cloud.occi.exception.AmbiguousIdentifierException;
import cz.cesnet.cloud.occi.exception.InvalidAttributeValueException;
import cz.cesnet.cloud.occi.exception.RenderingException;
import cz.cesnet.cloud.occi.parser.MediaType;

import org.json.simple.JSONObject;

public class DescribeVM
{
	// Describing cloud resources from provider
	// Available resources that can be described via API are the following:
	// - os_tpl = virtual appliance (aka VA) in the provider, 
	// - resource_tpl = template (aka flavour) resources, 
	// - compute = computing resources, 
	// - storage = storage resources,
	// - network = network resources.

	public static String vmState = "inactive";
	public static String publicIP = null;
	public static String[] vmFeatures = {};

	public static void doDescribe (Properties properties, Client client, Model model)
	{

		try 
		{
<<<<<<< HEAD
			if (properties.getProperty("OCCI_RESOURCE_ID").contains("compute"))
			{
				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));
				String[] Attributes = entities.get(0).toText().split(";");

				for (int i=0; i<Attributes.length; i++)
=======
<<<<<<< HEAD
			System.out.println();

			if (properties.getProperty("OCCI_RESOURCE_ID").contains("compute"))
			{
				System.out.println("[ VM DESCRIPTION ]");

				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));
				String[] Attributes = entities.get(0).toText().split(";");

				for (int i=0; i<Attributes.length; i++)
=======

			if properties.getProperty("OCCI_RESOURCE_ID").contains("compute")
			{
				System.out.println("[ VM DESCRIPTION ]");

				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));
				String[] Attributes = entities.get(0).toText().split(";");

				for int i=0; i<Attributes.length; i++
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				{
					if Attributes[i].contains("occi.networkinterface.address")
					{							
						publicIP = Attributes[i].replaceAll("[^\\d.]", "").substring(2);//Attributes[i].replace("occi.networkinterface.address=","");
					}
<<<<<<< HEAD
				}

				for (Entity entity : entities)
				{
					Map<Attribute, String> map = entity.getAttributes();

					for (Map.Entry<Attribute, String> entry : map.entrySet()) 
					{
						vmState = map.get(new Attribute("occi.compute.state"));
					}
				}
			}


			if (properties.getProperty("OCCI_RESOURCE_ID").contains("storage")) 
			{
				System.out.println("[ STORAGE DESCRIPTION ]");
				System.out.println("[[ " + properties.getProperty("OCCI_RESOURCE_ID") + " ]]");
				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));

				String title = "", size = "", state = "", ID = "";
				for (Entity entity : entities) 
				{
					ID = entity.getId();

=======
				}			
<<<<<<< HEAD

				for (Entity entity : entities)
				{
				Map<Attribute, String> map = entity.getAttributes();

				for (Map.Entry<Attribute, String> entry : map.entrySet()) 
				{
					vmState = map.get(new Attribute("occi.compute.state"));
				}

				if (vmState != null) System.out.println("occi.compute.state = " + vmState);
				}
			}


			if (properties.getProperty("OCCI_RESOURCE_ID").contains("storage")) 
			{
				System.out.println("[ STORAGE DESCRIPTION ]");
				System.out.println("[[ " + properties.getProperty("OCCI_RESOURCE_ID") + " ]]");
				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));

				String title = "", size = "", state = "", ID = "";
				for (Entity entity : entities) 
				{
					ID = entity.getId();

>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
					Map<Attribute, String> map = entity.getAttributes();

					for (Map.Entry<Attribute, String> entry : map.entrySet()) 
					{
						title = map.get(new Attribute("occi.core.title"));
						size = map.get(new Attribute("occi.storage.size"));
						state = map.get(new Attribute("occi.storage.state"));
<<<<<<< HEAD
					}

=======
						//System.out.println(entry.getKey() + " - " + entry.getValue());
					}

=======

				for Entity entity : entities
				{
					Map<Attribute, String> map = entity.getAttributes();

					for (Map.Entry<Attribute, String> entry : map.entrySet()) 
					{
						vmState = map.get(new Attribute("occi.compute.state"));
					}

					if (vmState != null) System.out.println("occi.compute.state = " + vmState);
				}
			}


			if properties.getProperty("OCCI_RESOURCE_ID").contains("storage") 
			{
				System.out.println("[ STORAGE DESCRIPTION ]");
				System.out.println("[[ " + properties.getProperty("OCCI_RESOURCE_ID") + " ]]");
				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));

				String title = "", size = "", state = "", ID = "";

				for Entity entity : entities
				{
					ID = entity.getId();

					Map<Attribute, String> map = entity.getAttributes();

					for Map.Entry<Attribute, String> entry : map.entrySet() 
					{
						title = map.get(new Attribute("occi.core.title"));
						size = map.get(new Attribute("occi.storage.size"));
						state = map.get(new Attribute("occi.storage.state"));
					}

>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
					System.out.println("\n>> location: " + "/storage/" + ID);
					System.out.println("occi.core.id = " + ID);
					System.out.println("occi.core.title = " + title);
					System.out.println("occi.storage.size = " + size);
					System.out.println("occi.storage.state = " + state);
				}
			}

<<<<<<< HEAD
			if (properties.getProperty("OCCI_RESOURCE_ID").contains("network")) 
=======
<<<<<<< HEAD
			if (properties.getProperty("OCCI_RESOURCE_ID").contains("network")) 
=======
			if properties.getProperty("OCCI_RESOURCE_ID").contains("network")
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
			{
				System.out.println("[ NETWORK DESCRIPTION ]");
				System.out.println("[[ " + properties.getProperty("OCCI_RESOURCE_ID") + " ]]");
				List<Entity> entities = client.describe(URI.create(properties.getProperty("OCCI_RESOURCE_ID")));

				String title = "", size = "", state = "", ID = "", summary = "", allocation = "";
				String address = "", vlan = "", netid = "", netvlan = "", phydev = "", bridge = "";

<<<<<<< HEAD
				for (Entity entity : entities) 
=======
<<<<<<< HEAD
				for (Entity entity : entities) 
=======
				for Entity entity : entities 
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				{
					ID = entity.getId();

					Map<Attribute, String> map = entity.getAttributes();
<<<<<<< HEAD
					for (Map.Entry<Attribute, String> entry : map.entrySet()) 
=======
<<<<<<< HEAD
					for (Map.Entry<Attribute, String> entry : map.entrySet()) 
=======
					for Map.Entry<Attribute, String> entry : map.entrySet()
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
					{
						title = map.get(new Attribute("occi.core.title"));
						size = map.get(new Attribute("occi.network.size"));
						summary = map.get(new Attribute("occi.core.summary"));
						vlan = map.get(new Attribute("occi.network.vlan"));
						state = map.get(new Attribute("occi.network.state"));
						allocation = map.get(new Attribute("occi.network.allocation"));
						address = map.get(new Attribute("occi.network.address"));
						netid = map.get(new Attribute("org.opennebula.network.id"));
						netvlan = map.get(new Attribute("org.opennebula.network.vlan"));
						phydev = map.get(new Attribute("org.opennebula.network.phydev"));
						bridge = map.get(new Attribute("org.opennebula.network.bridge"));
					}

					System.out.println("\n>> location: " + "/network/" + ID);
					if (ID != null) System.out.println("occi.core.id = " + ID);
					if (title != null) System.out.println("occi.core.title = " + title);
					if (summary != null) System.out.println("occi.core.summary = " + summary);
					if (vlan != null) System.out.println("occi.core.vlan = " + vlan);
					if (state != null) System.out.println("occi.network.state = " + state);
					if (allocation != null) System.out.println("occi.network.allocation = " + allocation);
					if (address != null) System.out.println("occi.network.address = " + address);
					if (netid != null) System.out.println("org.opennebula.network.id = " + netid);
					if (netvlan != null) System.out.println("org.opennebula.network.vlan = " + netvlan);
					if (phydev != null) System.out.println("org.opennebula.network.phydev = " + phydev);
					if (bridge != null) System.out.println("org.opennebula.network.bridge = " + bridge);
				}
			}

<<<<<<< HEAD
			if (properties.getProperty("RESOURCE").equals("os_tpl"))
=======
<<<<<<< HEAD
			if (properties.getProperty("RESOURCE").equals("os_tpl"))
=======
			if properties.getProperty("RESOURCE").equals("os_tpl")
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
			{
				System.out.println("[ OS_TPL DESCRIPTION ]");
				System.out.println("[[ " + properties.getProperty("OCCI_RESOURCE_ID") + " ]]");
				String TERM = properties.getProperty("OCCI_RESOURCE_ID")
				.substring(properties.getProperty("OCCI_RESOURCE_ID").lastIndexOf("#") + 1);
				List<Mixin> mixins = model.findRelatedMixins("os_tpl");

<<<<<<< HEAD
=======
<<<<<<< HEAD
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				for (Mixin mixin : mixins)
				{
					if ( ((mixin.getTerm()).contains(TERM)) ) 
					{
						System.out.println("title: \t\t" + mixin.getTitle());
						System.out.println("term: \t\t" + mixin.getTerm());
						System.out.println("location: \t" + mixin.getLocation());
					}
				}
			} 

			if (properties.getProperty("RESOURCE").equals("resource_tpl"))
			{
				// Getting the description of all the available template(s)
				System.out.println("[ RESOURCE_TPL DESCRIPTION ]");
				List<Mixin> mixins = model.findRelatedMixins(properties.getProperty("RESOURCE"));

				if (!mixins.isEmpty())
					for (Mixin mixin : mixins) 
					{
						if (mixin.getTerm().equals(properties.getProperty("OCCI_RESOURCE_ID"))) 
						{
						System.out.println("[[ " + mixin.getLocation() + " ]]");
						System.out.println("title: \t\t" + mixin.getTitle());
						System.out.println("term: \t\t" + mixin.getTerm());
						String locations = (mixin.getLocation()).toString();
						String segments[] = locations.split("/");
						System.out.println("location: \t" + "/" + segments[segments.length - 1] + "/");
						}
<<<<<<< HEAD
					}
=======
					}
=======
				for Mixin mixin : mixins
				{
					if mixin.getTerm()).contains(TERM)
					{
						System.out.println("title: \t\t" + mixin.getTitle());
						System.out.println("term: \t\t" + mixin.getTerm());
						System.out.println("location: \t" + mixin.getLocation());
					}
				}
			} 

			if properties.getProperty("RESOURCE").equals("resource_tpl")
			{
				// Getting the description of all the available template(s)
				System.out.println("[ RESOURCE_TPL DESCRIPTION ]");
				List<Mixin> mixins = model.findRelatedMixins(properties.getProperty("RESOURCE"));

				if !mixins.isEmpty()
					for Mixin mixin : mixins 
					{
						if mixin.getTerm().equals(properties.getProperty("OCCI_RESOURCE_ID"))
						{
							System.out.println("[[ " + mixin.getLocation() + " ]]");
							System.out.println("title: \t\t" + mixin.getTitle());
							System.out.println("term: \t\t" + mixin.getTerm());
							String locations = (mixin.getLocation()).toString();
							String segments[] = locations.split("/");
							System.out.println("location: \t" + "/" + segments[segments.length - 1] + "/");
						}
					}
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
			}
		}
		catch (CommunicationException | RenderingException | AmbiguousIdentifierException ex) 
		{throw new RuntimeException(ex);}
	}

	public static String[] describe (String vmID, JSONObject egiInput)
	{

<<<<<<< HEAD
		// Setting 

		String AUTH = (String) egiInput.get("auth");
		String OCCI_ENDPOINT_HOST = (String) egiInput.get("endpoint");
		String PROXY_PATH = (String) egiInput.get("proxyPath");
		String TRUSTED_CERT_REPOSITORY_PATH = (String) egiInput.get("trustedCertificatesPath");

		Boolean verbose = false;
=======
		// [ Setting preferences here! ]
		String AUTH = "x509"; 
		String OCCI_ENDPOINT_HOST = "https://carach5.ics.muni.cz:11443"; 

		String TRUSTED_CERT_REPOSITORY_PATH = "/etc/grid-security/certificates";
		String PROXY_PATH = "/tmp/x509up_u5040"; 

		Boolean verbose = true;
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c

		String ACTION = "describe";	

		// CESNET-MetaCloud
		// [ *Describing* available resources (e.g. os_tpl, resource_tpl, compute, storage and network) ]

		List<String> RESOURCE = Arrays.asList("compute",
		vmID); 

<<<<<<< HEAD
=======
<<<<<<< HEAD
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
		if (verbose) 
		{
			System.out.println();
			if (ACTION != null && !ACTION.isEmpty()) 
<<<<<<< HEAD
=======
=======
		if verbose 
		{
			System.out.println();
			if ACTION != null && !ACTION.isEmpty() 
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				System.out.println("[ACTION] = " + ACTION);
			else	
				System.out.println("[ACTION] = Get dump model");
				System.out.println("AUTH = " + AUTH);
<<<<<<< HEAD
=======
<<<<<<< HEAD
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
			if (OCCI_ENDPOINT_HOST != null && !OCCI_ENDPOINT_HOST.isEmpty()) 
				System.out.println("OCCI_ENDPOINT_HOST = " + OCCI_ENDPOINT_HOST);
			if (RESOURCE != null && !RESOURCE.isEmpty()) 
				System.out.println("RESOURCE = " + RESOURCE);
			if (TRUSTED_CERT_REPOSITORY_PATH != null && !TRUSTED_CERT_REPOSITORY_PATH.isEmpty()) 
				System.out.println("TRUSTED_CERT_REPOSITORY_PATH = " + TRUSTED_CERT_REPOSITORY_PATH);
			if (PROXY_PATH != null && !PROXY_PATH.isEmpty()) 
				System.out.println("PROXY_PATH = " + PROXY_PATH);
			if (verbose) 
<<<<<<< HEAD
=======
=======
			if OCCI_ENDPOINT_HOST != null && !OCCI_ENDPOINT_HOST.isEmpty() 
				System.out.println("OCCI_ENDPOINT_HOST = " + OCCI_ENDPOINT_HOST);
			if RESOURCE != null && !RESOURCE.isEmpty()
				System.out.println("RESOURCE = " + RESOURCE);
			if TRUSTED_CERT_REPOSITORY_PATH != null && !TRUSTED_CERT_REPOSITORY_PATH.isEmpty() 
				System.out.println("TRUSTED_CERT_REPOSITORY_PATH = " + TRUSTED_CERT_REPOSITORY_PATH);
			if PROXY_PATH != null && !PROXY_PATH.isEmpty()
				System.out.println("PROXY_PATH = " + PROXY_PATH);
			if verbose
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				System.out.println("Verbose = ON ");
			else System.out.println("Verbose = OFF ");
		}

		Properties properties = new Properties();
<<<<<<< HEAD
=======
<<<<<<< HEAD
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
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
				properties.setProperty("OCCI_RESOURCE_ID", RESOURCE.get(i));

			else 
			{ 
				properties.setProperty("RESOURCE", RESOURCE.get(i));
			}
		}
<<<<<<< HEAD
=======
=======
		if ACTION != null && !ACTION.isEmpty()
			properties.setProperty("ACTION", ACTION);
		if OCCI_ENDPOINT_HOST != null && !OCCI_ENDPOINT_HOST.isEmpty()
			properties.setProperty("OCCI_ENDPOINT_HOST", OCCI_ENDPOINT_HOST);

		if RESOURCE != null && !RESOURCE.isEmpty() 
			for (int i=0; i<RESOURCE.size(); i++)
			{
				if (!RESOURCE.get(i).equals("compute") && 
				!RESOURCE.get(i).equals("storage") &&
				!RESOURCE.get(i).equals("network")) &&
				!RESOURCE.get(i).equals("os_tpl") &&
				!RESOURCE.get(i).equals("resource_tpl")) 
					properties.setProperty("OCCI_RESOURCE_ID", RESOURCE.get(i));

				else 
				{ 
					properties.setProperty("RESOURCE", RESOURCE.get(i));
				}
			}
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c

		properties.setProperty("TRUSTED_CERT_REPOSITORY_PATH", TRUSTED_CERT_REPOSITORY_PATH);
		properties.setProperty("PROXY_PATH", PROXY_PATH);
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

<<<<<<< HEAD
			if  (ACTION.equals("describe")) 
=======
<<<<<<< HEAD
			if  (ACTION.equals("describe")) 
=======
			if  ACTION.equals("describe") 
>>>>>>> 1098d207ae9ed8e9e1670143fc84e89a2ba54dc6
>>>>>>> 291681728effedb5c5b45f3231aaa23b3d3b0d6c
				doDescribe(properties, client, model);
		} 
		catch (CommunicationException ex ) 
		{throw new RuntimeException(ex);}

		vmFeatures = new String[]{publicIP,vmState};
		return vmFeatures;
	}
}
