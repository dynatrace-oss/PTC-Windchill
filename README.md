# PTC-Windchill
This repo contains files that are required to create assets in a Dynatrace tenant including an application, metrics, dashboards and synthetic to monitor PTC Windchill.
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

All the relevant files required to create Dynatrace objects are in the **assets** directory. Under this directory, subfolders are available to create the relevant objects in Dynatrace. For example, all the files to create different dashboards are under Dashboards directory. Each subfolder is prefixed by a numeric value to show the order in which these objects need to be created. Since ManagementZones are created first, it is named as 01-ManagementZones. Note that some dashboards may require minor configuration after deployment - such as filtering a tile to a particular process group.

Example synthetic monitors are available in the **synthetic** directory, please refer to the [README-Synthetic.md](https://github.com/dynatrace-oss/PTC-Windchill/tree/main/synthetic/README-Synthetic.md) - examples include remote file server, queue and server status monitors.

There are 2 ways to create PTC related components to your Dynatrace tenant
### 1. Utility Tool

A utility tool is available to upload all the assets to the Dynatrace tenant. In order to use the utility:
* Clone or download this repo. 
* [Download](https://github.com/dynatrace-oss/PTC-Windchill/releases/latest) the utility for your OS. Unzip the file. Move the file PTC-Windchill_x.x.x to root of your cloned directory and run it. You will need Dynatrace tenant URL and the token for this utility to work.

### 2. Manually upload the components
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

### Troubleshooting
Please open new issues if you encounter problems in uploading the components to Dynatrace or if there are missing functionalities for the PTC Windchill pack.

### Contact
Please reach out to PTC@dynatrace.com for assistance with deploying the starter set or for any questions/concerns
