#region Top of Script

#requires -version 3

<#
.SYNOPSIS
	Tests connection to target endpoint using Test-Connection cmdlet
.DESCRIPTION
	This script will emulate the output of the 'ping' utility in Linux systems
  The -Count argument allows you to specify the number of packets being sent
.EXAMPLE
	.\windows-ping-script.ps1 -Target "www.newrelic.com" -Count 5
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	03/24/2020
	Purpose/Change:	Initial script development
#>

#endregion

#####-----------------------------------------------------------------------------------------#####

#region Script Parameters

Param
	(
		
        [ Parameter( Mandatory = $true ) ] [ string ] $Target,
		[ Parameter( Mandatory = $true ) ] [ int ] $Count

	)

#endregion Script Parameters

#####-----------------------------------------------------------------------------------------#####

#region Execution

# Build an empty object to hold our results
$results = New-Object -TypeName psobject

# Test the connection
$ping = Test-Connection -ComputerName $target -Count $count -ErrorAction SilentlyContinue -WarningAction SilentlyContinue

# If the test fails, manually set the output
if( !( $ping ) ) { 

    $results | Add-Member -MemberType NoteProperty -Name "target" -Value $target 
    $results | Add-Member -MemberType NoteProperty -Name "packetLoss" -Value 100
    $results | Add-Member -MemberType NoteProperty -Name "packetsReceived" -Value 0
    $results | Add-Member -MemberType NoteProperty -Name "packetsTransmitted" -Value $count 
    $results | Add-Member -MemberType NoteProperty -Name "avgResponse" -Value 0
    $results | Add-Member -MemberType NoteProperty -Name "minResponse" -Value 0
    $results | Add-Member -MemberType NoteProperty -Name "maxResponse" -Value 0
    
} 

# Otherwise, add data to our results
else {

    # Calculate the response time summary
    $responseTime = $ping | Measure-Object -Property ResponseTime -Minimum -Maximum -Average
    
    # Calculate packet loss percentage
    $packetLoss = ( ( $count - $responseTime.Count ) / $count ) * 100

    # Fill our results array
    $results | Add-Member -MemberType NoteProperty -Name "target" -Value $target 
    $results | Add-Member -MemberType NoteProperty -Name "packetLoss" -Value $packetLoss 
    $results | Add-Member -MemberType NoteProperty -Name "packetsReceived" -Value $responseTime.Count 
    $results | Add-Member -MemberType NoteProperty -Name "packetsTransmitted" -Value $count 
    $results | Add-Member -MemberType NoteProperty -Name "avgResponse" -Value $responseTime.Average
    $results | Add-Member -MemberType NoteProperty -Name "minResponse" -Value $responseTime.Minimum
    $results | Add-Member -MemberType NoteProperty -Name "maxResponse" -Value $responseTime.Maximum

}

# Print the results to STDOUT in JSON format
$results | ConvertTo-Json

#endregion Execution
