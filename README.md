# PTC-Windchill
This repo contains files that are required to create an application, metrics and dashboards to monitor PTC Windchill with Dynatrace.  
In order to create these entities in Dynatrace, you would need 2 items:
Name | Description
------------ | -------------
Dynatrace tenant url | `Managed` https://{your-domain}/e/{your-environment-id}  <br/>`SaaS` https://{your-environment-id}.live.dynatrace.com
API Token | You need the Write configuration (WriteConfig) permission assigned to your API token  

* Please use any tool of your convenience (Postman, curl etc.) to make POST REST call to your Dynatrace tenant. 
  
Method | Endpoint | Filename | Notes
------------| ----------------------------------- | --------------- | -----------------------
POST | /api/config/v1/extensions | custom.jmx.ActiveUsers.zip |
POST | /api/config/v1/extensions | custom.jmx.Queues.zip |
POST | /api/config/v1/applications/web | 1-application.json |
GET  | /applications/web | Run this command to get the application id | Note the application id
POST | /applicationDetectionRules | 2-applicationDetectionRules.json | Replace the application id in the json with the id received from the GET command above  
POST | /api/config/v1/managementZones | 3-mgmtZone.json  
POST | /api/config/v1/calculatedMetrics/service | 4-exceptionCount.json  
POST | /api/config/v1/calculatedMetrics/service | 5-topSqlStatements.json  
POST | /api/config/v1//dashboards | dashboards/*.json
