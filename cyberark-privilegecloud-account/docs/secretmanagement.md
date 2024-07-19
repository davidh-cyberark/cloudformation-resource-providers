# CyberArk::PrivilegeCloud::Account SecretManagement

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#automaticmanagementenabled" title="AutomaticManagementEnabled">AutomaticManagementEnabled</a>" : <i>Boolean</i>,
    "<a href="#manualmanagementreason" title="ManualManagementReason">ManualManagementReason</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#automaticmanagementenabled" title="AutomaticManagementEnabled">AutomaticManagementEnabled</a>: <i>Boolean</i>
<a href="#manualmanagementreason" title="ManualManagementReason">ManualManagementReason</a>: <i>String</i>
</pre>

## Properties

#### AutomaticManagementEnabled

Whether the account secret is automatically managed by the CPM.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ManualManagementReason

Reason for disabling automatic secret management.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

