 #region Top of Script

#requires -version 4
#requires -Module ActiveDirectory

<#
.SYNOPSIS
	Queries for target Active Directory health statistics
.DESCRIPTION
  This script uses a static list of queries
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	13-Mar-2023
	Purpose/Change:	Initial script development

	Version:		1.1
	Author:			Zack Mutchler
	Creation Date:	9-May-2023
	Purpose/Change:	Remove GPO queries and setup for filtered AD OU search
#>

#endregion Top of Script

#####-----------------------------------------------------------------------------------------#####

#region Script Parameters

Param(

    [ Parameter( Mandatory = $true ) ] [ValidateNotNullOrEmpty() ] [ string ] $OU

)

#endregion Script Parameters

#region Execution

# Build an empty PSObject to add our results to
$results = New-Object -TypeName PSCustomObject

# Iterate through our health checks, adding the results to our custom object
  # Add the OU Distinguished Name
  $results | Add-Member -MemberType NoteProperty -Name "ouDistinguishedName" -Value $OU
  # Inactive Accounts
  $results | Add-Member -MemberType NoteProperty -Name "adInactiveAccounts" -Value (Search-ADAccount -UsersOnly -AccountInactive -TimeSpan 365 -SearchBase $OU | Where-Object {$_.enabled -eq $true -and $_.lockedOut -eq $false -and $_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Expired Accounts
  $results | Add-Member -MemberType NoteProperty -Name "adExpiredAccounts" -Value (Search-ADAccount -UsersOnly -AccountExpired -SearchBase $OU | Where-Object {$_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Locked Accounts
  $results | Add-Member -MemberType NoteProperty -Name "adLockedAccounts" -Value (Search-ADAccount -UsersOnly -LockedOut -SearchBase $OU | Where-Object {$_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Disabled Accounts
  $results | Add-Member -MemberType NoteProperty -Name "adDisabledAccounts" -Value (Search-ADAccount -UsersOnly -AccountDisabled -SearchBase $OU | Where-Object{$_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Passwords have expired
  $results | Add-Member -MemberType NoteProperty -Name "adExpiredPasswords" -Value (Search-ADAccount -UsersOnly -PasswordExpired -SearchBase $OU | Where-Object{$_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Passwords that never expire
  $results | Add-Member -MemberType NoteProperty -Name "adPasswordsThatNeverExpire" -Value (Search-ADAccount -UsersOnly -PasswordNeverExpires -SearchBase $OU | Where-Object{$_.distinguishedname -notlike "*Microsoft Exchange System Objects*"}).count
  # Disabled Computers
  $results | Add-Member -MemberType NoteProperty -Name "adDisabledComputers" -Value (Search-ADAccount -ComputersOnly -AccountDisabled -SearchBase $OU).count
  # Is this OU empty? true|false
  $results | Add-Member -MemberType NoteProperty -Name "adEmptyOrganizationalUnit" -Value ((-not (Get-ADOrganizationalUnit -SearchBase $OU -Filter * -SearchScope OneLevel)) -or ((Get-ADOrganizationalUnit -SearchBase $OU -Filter * -SearchScope OneLevel).Count -eq 0))
  # Empty Groups
  $results | Add-Member -MemberType NoteProperty -Name "adEmptyGroups" -Value (Get-ADGroup -Filter * -SearchBase $OU -Properties Members | Where-Object {-not $_.members -and $_.distinguishedname -notlike "*deploy*" -and $_.distinguishedname -notlike "*Builtin*" -and $_.distinguishedname -notlike "*CN=Users*"}).count

# Print the results to STDOUT in JSON for Flex to pickup
$results | ConvertTo-Json

#endregion Execution
