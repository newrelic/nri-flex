#region Top of Script

#requires -version 3

<#
.SYNOPSIS
	Tests clock drift from target NTP server
.DESCRIPTION
  This script uses the "w32tm" tool from Windows.
  There is the option to use a hard-coded time server (see like 24-26), or as in most enterprise
  settings where a primary domain controller serves as the time server, the option to have the
  $ntpServer variable set dynamically based on the source of the NTP configuration.  It should be
  noted however, that if more than one time server is set, either manually or through group policy
  this dynamic option will not work.
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	09/24/2020
	Purpose/Change:	Initial script development
#>

#endregion Top of Script

#####-----------------------------------------------------------------------------------------#####

#region Execution

# Target NTP Server
# Hard coded option
# $ntpServer = 'time.windows.com'

# Target NTP Server
# Dynamic option to pickup configured NTP Server on local machine
$server = [Net.Dns]::GetHostName()
$ntpServer = w32tm /query /computer:$server /source | Out-String

# REGEX to find Skew later
$findSkew = [regex]"(?:NTP\: )(?<Value>/?[^s]+)"

# Grab the local server time in Epoch
$localEpoch = [int64](Get-Date -UFormat %s)

# Query the current skew using the w32tm tool
$ntpQuery = Invoke-Expression "w32tm /monitor /computers:$ntpServer" | Out-String

# Check to see if there is a match
If ( $ntpQuery -match $findSkew ) {

    # Extract the skew from the resulting string match
    $ntpSkew = [decimal]$Matches['Value']

    # Set the results object, making sure we pass the +/- sign to show the direction of skew from NTP
    If( $ntpSkew -lt 0 ) {
        $skewValue = [math]::Abs( $ntpSkew )
        $skewSign = "-"
        $results = New-Object -TypeName psobject
        $results | Add-Member -MemberType NoteProperty -Name 'ntpServer' -Value $ntpServer
        $results | Add-Member -MemberType NoteProperty -Name 'localTime' -Value $localEpoch
        $results | Add-Member -MemberType NoteProperty -Name 'skewString' -Value $ntpSkew.ToString()
        $results | Add-Member -MemberType NoteProperty -Name 'skewSign' -Value $skewSign
        $results | Add-Member -MemberType NoteProperty -Name 'skewValue' -Value $skewValue
    }
    Else {
        $skewValue = [math]::Abs( $ntpSkew )
        $skewSign = "+"
        $results = New-Object -TypeName psobject
        $results | Add-Member -MemberType NoteProperty -Name 'ntpServer' -Value $ntpServer
        $results | Add-Member -MemberType NoteProperty -Name 'localTime' -Value $localEpoch
        $results | Add-Member -MemberType NoteProperty -Name 'skewString' -Value $ntpSkew.ToString()
        $results | Add-Member -MemberType NoteProperty -Name 'skewSign' -Value $skewSign
        $results | Add-Member -MemberType NoteProperty -Name 'skewValue' -Value $skewValue
    }
}
# Otherwise exit with the $nulls
Else {
    $results = New-Object -TypeName psobject
    $results | Add-Member -MemberType NoteProperty -Name 'ntpServer' -Value $ntpServer
    $results | Add-Member -MemberType NoteProperty -Name 'localTime' -Value $localEpoch
    $results | Add-Member -MemberType NoteProperty -Name 'skewString' -Value "Error Collecting NTP Skew Value"
    $results | Add-Member -MemberType NoteProperty -Name 'skewSign' -Value $null
    $results | Add-Member -MemberType NoteProperty -Name 'skewValue' -Value $null
}

# Print the results to STDOUT in JSON for Flex to pickup
$results | ConvertTo-Json

#endregion Execution
