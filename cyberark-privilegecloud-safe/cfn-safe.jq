{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "Template for Creating a CyberArk Safe in Privilege Cloud Vault",
  "Resources": {
    "DemoSafe": {
      "Type": "CyberArk::PrivilegeCloud::Safe",
      "Properties": {
        "NewSafeName": "new-demo-safe-1"
      }
    }
  },
  "Outputs": {
    "NewSafeUrl": {
      "Value": {
        "Fn::GetAtt": [
          "DemoSafe",
          "SafeResourceId"
        ]
      },
      "Description": "New Safe ID for use with PCloud Safe custom resource"
    }
  }
}