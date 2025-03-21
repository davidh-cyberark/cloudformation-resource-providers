# CyberArk::PrivilegeCloud::Safe

Manage CyberArk Privilege Cloud safe.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "CyberArk::PrivilegeCloud::Safe",
    "Properties" : {
        "<a href="#newsafename" title="NewSafeName">NewSafeName</a>" : <i>String</i>,
        "<a href="#newsaferole" title="NewSafeRole">NewSafeRole</a>" : <i>String</i>,
    }
}
</pre>

### YAML

<pre>
Type: CyberArk::PrivilegeCloud::Safe
Properties:
    <a href="#newsafename" title="NewSafeName">NewSafeName</a>: <i>String</i>
    <a href="#newsaferole" title="NewSafeRole">NewSafeRole</a>: <i>String</i>
</pre>

## Properties

#### NewSafeName

Name of new safe to create

_Required_: Yes

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### NewSafeRole

Name of new role to create added to the safe

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### SafeResourceId

Example: PCLOUDURL=https://{YOUR-SUBDOMAIN}.privilegecloud.cyberark.cloud|SAFEURLID=New-Safe-Url-Id

#### RoleId

The unique ID of the Role. The unique ID is generated when the role is created. Ex: ROLENAME=myrolename|ROLEID=79d6xxxx_9bxx_46xx_9fxx_d5db6cxxxxxx

