# Dynatrace PTC Windchill Starter Set
This repo contains files to create PTC Windchill assets in a Dynatrace tenant including application, metrics, dashboards and synthetic monitors. Check out our Dynatrace blog on PTC Windchill/ThingWorx [here](https://www.dynatrace.com/news/blog/ai-driven-observability-for-ptc-windchill-thingworx/).
In order to create these entities in Dynatrace, you will need two items:

Name | Description
------------ | -------------
Dynatrace tenant url | `Managed` https://{your-domain}/e/{your-environment-id}  <br/>`SaaS` https://{your-environment-id}.live.dynatrace.com
API Token | Required permissions outlined in the screenshot below

The API Token needs to have these minimum permissions:

![GitHub Logo](/images/TokenPermissions.png)

All the relevant files required to create the Dynatrace objects are in the **assets** directory. Under this directory, there are subfolders containing the assets by category. Each subfolder is prefixed by a numeric value to show the order in which these objects need to be created. Note that some dashboards may require minor configuration after deployment - such as filtering a tile to a particular process group.

Example synthetic monitors are available in the **synthetic** directory, please refer to [README-Synthetic.md](https://github.com/dynatrace-oss/PTC-Windchill/tree/main/synthetic/README-Synthetic.md) - examples include remote file server, queue and server status monitors.

There are two ways to deploy the assets to your Dynatrace tenant, with the utility tool or uploading the them manually. Each option is covered in more detail below.

### 1. Utility Tool

A utility tool is available to automatically upload all the assets to the Dynatrace tenant. In order to use the utility:
* Clone or download this repo. 
* [Download](https://github.com/dynatrace-oss/PTC-Windchill/releases/latest) the utility for your OS. Unzip and move the file PTC-Windchill_x.x.x to the root of your cloned directory. To run the executable - open command prompt, navigate to the executable and run it.

### 2. Manually upload the components
* Please use any tool of your convenience (Postman, curl etc.) to make POST REST calls to your Dynatrace tenant or utilize the API explorer in the Dynatrace UI. Apply all REST calls of the files in a given subfolder.
  
Method | Endpoint | Subfolder Name | Notes
------------| ----------------------------------- | --------------- | -----------------------
POST | /api/config/v1/managementZones | 01-ManagementZones |
POST | /api/config/v1/applications/web | 02-Application |
POST | /api/config/v1/applicationDetectionRules | 03-DetectionRules | Replace the application id in the json with the application id 
POST | /api/config/v1/autoTags | 04-AutoTags | 
POST | /api/config/v1/service/requestAttributes | 05-RequestAttributes |  
POST | /api/config/v1/calculatedMetrics/service | 06-MetricsService |  
POST | /api/config/v1/service/customServices/java?position=APPEND | 07-CustomServices |
POST | /api/config/v1/extensions | 08-Extensions | Extensions require uploading the zip files 
POST | /api/config/v1/dashboards | 11-Dashboards

### Contact
Please reach out to PTC@dynatrace.com for assistance with deploying the starter set or for any questions/concerns
