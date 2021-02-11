# PTC-Windchill
This repo contains files that are required to create an application, metrics and dashboards to monitor PTC Windchill with Dynatrace.  
In order to create these entities in Dynatrace, you would need 2 items:
Name | Description
------------ | -------------
Dynatrace tenant url | `Managed` https://{your-domain}/e/{your-environment-id}  <br/>`SaaS` https://{your-environment-id}.live.dynatrace.com
API Token | You need the Write configuration (WriteConfig) permission assigned to your API token  

* Using Dynatrace UI, upload the custom JMX metrics file, **custom.jmx.Queues-1.zip**, to the Dynatrace tenant
  
* Please use any tool of your convenience (Postman, curl etc.) to make POST REST call to your Dynatrace tenant. 
  
Method | Endpoint | Filename
------------| ----------------------------------- | ---------------  
POST | /api/config/v1/applications/web | 1-application.json  
POST | /api/config/v1/managementZones | 2-mgmtZone.json  
POST | /api/config/v1/calculatedMetrics/service | 3-exceptionCount.json  
POST | /api/config/v1/calculatedMetrics/service | 4-topSqlStatements.json  
POST | /api/config/v1//dashboards | PTCWindchillUserData.json
POST | /api/config/v1//dashboards | dashboards/*.json
  
* Using Dynatrace UI, edit the application **PTC Windchill** to the following application detection rule:  
![Application Detection Rule](/images/ApplicationDetectionRule.png)
