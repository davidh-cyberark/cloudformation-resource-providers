# CyberArk::PrivilegeCloud::Account RemoteMachinesAccess

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#remotemachines" title="RemoteMachines">RemoteMachines</a>" : <i>String</i>,
    "<a href="#accessrestrictedtoremotemachines" title="AccessRestrictedToRemoteMachines">AccessRestrictedToRemoteMachines</a>" : <i>Boolean</i>
}
</pre>

### YAML

<pre>
<a href="#remotemachines" title="RemoteMachines">RemoteMachines</a>: <i>String</i>
<a href="#accessrestrictedtoremotemachines" title="AccessRestrictedToRemoteMachines">AccessRestrictedToRemoteMachines</a>: <i>Boolean</i>
</pre>

## Properties

#### RemoteMachines

List of remote machines, separated by semicolons.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### AccessRestrictedToRemoteMachines

Whether or not to restrict access only to specified remote machines.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

