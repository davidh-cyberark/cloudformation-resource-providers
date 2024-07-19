# CyberArk::PrivilegeCloud::Account

Manage CyberArk Privilege Cloud account.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "CyberArk::PrivilegeCloud::Account",
    "Properties" : {
        "<a href="#safename" title="SafeName">SafeName</a>" : <i>String</i>,
        "<a href="#platformid" title="PlatformID">PlatformID</a>" : <i>String</i>,
        "<a href="#name" title="Name">Name</a>" : <i>String</i>,
        "<a href="#address" title="Address">Address</a>" : <i>String</i>,
        "<a href="#username" title="UserName">UserName</a>" : <i>String</i>,
        "<a href="#secret" title="Secret">Secret</a>" : <i>String</i>,
        "<a href="#secrettype" title="SecretType">SecretType</a>" : <i>String</i>,
        "<a href="#secretmanagement" title="SecretManagement">SecretManagement</a>" : <i><a href="secretmanagement.md">SecretManagement</a></i>,
        "<a href="#platformaccountproperties" title="PlatformAccountProperties">PlatformAccountProperties</a>" : <i>[ <a href="kvpair.md">KVpair</a>, ... ]</i>,
        "<a href="#remotemachinesaccess" title="RemoteMachinesAccess">RemoteMachinesAccess</a>" : <i><a href="remotemachinesaccess.md">RemoteMachinesAccess</a></i>,
    }
}
</pre>

### YAML

<pre>
Type: CyberArk::PrivilegeCloud::Account
Properties:
    <a href="#safename" title="SafeName">SafeName</a>: <i>String</i>
    <a href="#platformid" title="PlatformID">PlatformID</a>: <i>String</i>
    <a href="#name" title="Name">Name</a>: <i>String</i>
    <a href="#address" title="Address">Address</a>: <i>String</i>
    <a href="#username" title="UserName">UserName</a>: <i>String</i>
    <a href="#secret" title="Secret">Secret</a>: <i>String</i>
    <a href="#secrettype" title="SecretType">SecretType</a>: <i>String</i>
    <a href="#secretmanagement" title="SecretManagement">SecretManagement</a>: <i><a href="secretmanagement.md">SecretManagement</a></i>
    <a href="#platformaccountproperties" title="PlatformAccountProperties">PlatformAccountProperties</a>: <i>
      - <a href="kvpair.md">KVpair</a></i>
    <a href="#remotemachinesaccess" title="RemoteMachinesAccess">RemoteMachinesAccess</a>: <i><a href="remotemachinesaccess.md">RemoteMachinesAccess</a></i>
</pre>

## Properties

#### SafeName

The Safe name where the account is created.

_Required_: Yes

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### PlatformID

The platform assigned to this account. Platform ID.

_Required_: Yes

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Name

The name of the account.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Address

The name or address of the machine where the account will be used.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### UserName

Account user's name.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Secret

The password value.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### SecretType

The type of password. Valid values: Password, key

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### SecretManagement

_Required_: No

_Type_: <a href="secretmanagement.md">SecretManagement</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### PlatformAccountProperties

_Required_: No

_Type_: List of <a href="kvpair.md">KVpair</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### RemoteMachinesAccess

_Required_: No

_Type_: <a href="remotemachinesaccess.md">RemoteMachinesAccess</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the AccountResourceId.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### AccountResourceId

Example: ACCTNAME=Account-Name|ACCTID=Account-ID based on values returned from PCloud

