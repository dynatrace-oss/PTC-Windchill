# PTC-Windchill
This repo contains files that are required to create assets in Dynatrace tenant including an application, metrics and dashboards to monitor PTC Windchill.
In order to create these entities in Dynatrace, you would need 2 items:
Name | Description
------------ | -------------
Dynatrace tenant url | `Managed` https://{your-domain}/e/{your-environment-id}  <br/>`SaaS` https://{your-environment-id}.live.dynatrace.com
API Token | You need the Write configuration (WriteConfig) permission assigned to your API token  

The API Token needs to have these minimum permissions:
* Write Configuration
* Capture Request Data
* Write Entities
* Write Extensions

![GitHub Logo](/images/TokenPermissions.png)

All the relevant files required to create Dynatrace objects are in **assets** directory. Under this directory, subfolders are available to create the relevant objects in Dynatrace. For example, all the files to create different dashboards are under Dashboards directory. Each subfolder is prefixed by a numeric value to show the order in which these objects need to be created. Since ManagementZones are created first, it is named as 01-ManagementZones. 

There are 2 ways to create PTC related components to your Dynatrace tenant
### Utility Tool

You can [Download](https://github.com/dynatrace-oss/PTC-Windchill/releases/latest) the utility to upload all the conponents. You will require Dynatrace tenant URL and the token for this utility to work.

### Manually upload the components
* Please use any tool of your convenience (Postman, curl etc.) to make POST REST call to your Dynatrace tenant. Please use the tool of your choice to apply REST calls for all the files in a given subfolder.

  
Method | Endpoint | Subfolder Name | Notes
------------| ----------------------------------- | --------------- | -----------------------
POST | /api/config/v1/managementZones | 01-ManagementZones |
POST | /api/config/v1/applications/web | 02-Application |
POST | /api/config/v1/applicationDetectionRules | 03-DetectionRules | Replace the application id in the json with the application id 
POST | /api/config/v1/autoTags | 04-AutoTags | 
POST | /api/config/v1/service/requestAttributes | 05-RequestAttributes |  
POST | /api/config/v1/calculatedMetrics/service | 06-MetricsService |  
POST | /api/config/v1/service/customServices/java?position=APPEND | 07-CustomServices |
POST | /api/config/v1/extensions | 08-Extensions | Extensions require uploading of the zip file. 
POST | /api/config/v1/dashboards | 11-Dashboards
