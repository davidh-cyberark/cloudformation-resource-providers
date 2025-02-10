{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "Template for Creating a CyberArk Account in Privilege Cloud Vault",
  "Resources":
    {
      "DemoSafe":
        {
          "Type": "CyberArk::PrivilegeCloud::Account",
          "Properties":
            {
              "SafeName": "new-demo-safe-1",
              "PlatformID": "AuroraPostgresSQLRDS-ASM",
              "Address": "127.0.0.1",
              "UserName": "Oscar",
              "Secret": "Hello!Oscar!",
              "SecretType": "Password",
              "PlatformAccountProperties":
                [
                  { "Key": "Database", "Value": "my-database-name" },
                  { "Key": "Port", "Value": "9876" }
                ]
            }
        }
    },
  "Outputs":
    {
      "NewAccountId":
        {
          "Value": { "Fn::GetAtt": ["DemoSafe", "AccountResourceId"] },
          "Description": "New Account ID for use with PCloud Account custom resource"
        }
    }
}
