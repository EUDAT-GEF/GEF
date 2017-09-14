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

import org.italiangrid.voms.clients.VomsProxyInit;

import org.apache.log4j.Logger;
import java.util.Properties;

public class VOMSProxyInit
{
	private static Logger log = Logger.getLogger(VOMSProxyInit.class);
  
        public static boolean isEmpty(String str)
        {
		if (str != null && !str.isEmpty()) return false; 
		else return true; 
        }
	
	/* M	A	I	N */
 	public static void main (String[] args)
	{			
		String VONAME = "fedcloud.egi.eu"; // <= Change here!
		String VOMS_PROXY_FILEPATH = "/tmp/x509up_u5040"; // <= Change here!
		String VOMS_LIFETIME = "24:00";
		String VOMSES_DIR = "/etc/vomses/";
		String X509_CERT_DIR = "/etc/grid-security/certificates/";
		Boolean ENABLE_RFC = true;
		
		try {
			if (isEmpty(VONAME) && 
                           (isEmpty(VOMS_PROXY_FILEPATH)) &&
			   (isEmpty(VOMS_LIFETIME)) &&
			   (isEmpty(VOMSES_DIR)) &&
			   (isEmpty(X509_CERT_DIR))) 
                                throw new Exception ("[ ARGUMENTS EXCEPTION ]");
			
			Properties p = new Properties(System.getProperties());
 			p.setProperty("X509_USER_PROXY", VOMS_PROXY_FILEPATH);
			System.setProperties(p);
			
			if (ENABLE_RFC)
				VomsProxyInit.main(new String[]{
					"-voms", VONAME,
					"-vomses", VOMSES_DIR,
					"-out", VOMS_PROXY_FILEPATH,
					"-certdir", X509_CERT_DIR,
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
                                        "-certdir", X509_CERT_DIR,
                                        "-vomslife", VOMS_LIFETIME,
					"-ignorewarn",
					"-limited",
                                        "-debug"
                                });

		        //VomsProxyInfo.main(new String[]{"--all"});
		} catch (Exception exc){ System.out.println (exc.toString()); }		
	}
}
