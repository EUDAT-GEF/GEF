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

import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.io.FileWriter;

import org.joda.time.LocalTime;

import org.json.simple.JSONObject;
import org.json.simple.JSONArray;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;

import org.italiangrid.voms.clients.VomsProxyInit;
import org.italiangrid.voms.clients.VomsProxyInfo;

import org.apache.log4j.Logger;
import java.util.Properties;

import org.json.simple.JSONObject;


public class VOMSProxyInit

{
	private static Logger log = Logger.getLogger(VOMSProxyInit.class);
  
        public static boolean isEmpty(String str)
        {
		if (str != null && !str.isEmpty()) return false; 
		else return true; 
        }
	
 	public static void proxyGeneration (JSONObject egiInput)
	{			
		String VONAME = (String) egiInput.get("auth");; 
		String VOMS_PROXY_FILEPATH = (String) egiInput.get("proxyPath");
		String VOMS_LIFETIME = "12:00";
		String VOMSES_DIR = (String) egiInput.get("vomsDir");
		String CERT_DIR = (String) egiInput.get("trustedCertificatesPath");
		Boolean ENABLE_RFC = true;
		
		try {
			if (isEmpty(VONAME) && 
                           (isEmpty(VOMS_PROXY_FILEPATH)) &&
			   (isEmpty(VOMS_LIFETIME)) &&
			   (isEmpty(VOMSES_DIR)) &&
			   (isEmpty(CERT_DIR))) 
                                throw new Exception ("[ ARGUMENTS EXCEPTION ]");
			
			Properties p = new Properties(System.getProperties());
 			p.setProperty("X509_USER_PROXY", VOMS_PROXY_FILEPATH);
			System.setProperties(p);
			
			if (ENABLE_RFC)
				VomsProxyInit.main(new String[]{
					"-voms", VONAME,
					"-vomses", VOMSES_DIR,
					"-out", VOMS_PROXY_FILEPATH,
					"-certdir", CERT_DIR,
					"-vomslife", VOMS_LIFETIME,
					"-ignorewarn",
					"-limited",
					"-rfc",
					"-debug"
				});
			else
				VomsProxyInit.main(new String[]{
                                        "-voms", VONAME,
                                        "-vomses", VOMSES_DIR,
                                        "-out", VOMS_PROXY_FILEPATH,
                                        "-certdir", CERT_DIR,
                                        "-vomslife", VOMS_LIFETIME,
					"-ignorewarn",
					"-limited",
                                        "-debug"
                                }); } catch (Exception exc){ System.out.println (exc.toString()); }
	}
}
		

