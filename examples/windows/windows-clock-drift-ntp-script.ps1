#region Top of Script

#requires -version 3

<#
.SYNOPSIS
	Tests clock drift from target NTP server
.DESCRIPTION
    This script uses the "w32tm" tool from Windows
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
$ntpServer = 'time.windows.com'

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
