#region Top of Script

#requires -version 3

<#
.SYNOPSIS
	Scrapes a target URL for embedded links and tests each one
.DESCRIPTION
	This script can take >60s to run on pages with a large number of links
    It is advised to setup a distinct config with an extended interval and timeout
    Note that -UseBasicParsing is being used to avoid issues with systems that don't have Internet Explorer installed
.EXAMPLE
	.\windows-web-scrape-script.ps1 -Target "www.swapi.co"
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
		
        [ Parameter( Mandatory = $true ) ] [ string ] $Target

	)

#endregion Script Parameters

#####-----------------------------------------------------------------------------------------#####

#region Execution

# Setup acceptance of TLS 1.2 during this script to account for outdated sites
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12

# Build an empty object to hold our results
$results = @()

# Test the page and suppress error output to STDOUT
$page = Invoke-WebRequest -Method Get -UseBasicParsing -Uri $Target -ErrorAction SilentlyContinue -WarningAction SilentlyContinue

# Validate that we loaded our page correctly

If( $page.StatusCode -ne 200 ) {

    # Build a new psobject with results for the page
    $fail = New-Object -TypeName psobject
    $fail | Add-Member -MemberType NoteProperty -Name 'targetURL' -Value $Target
    $fail | Add-Member -MemberType NoteProperty -Name 'targetSuccess' -Value $false # Did my target page actually load?
    $fail | Add-Member -MemberType NoteProperty -Name 'linkURL' -Value $Target 
    $fail | Add-Member -MemberType NoteProperty -Name 'statusCode' -Value $page.StatusCode

    # Print our object to STDOUT in JSON format
    $fail | ConvertTo-Json
        
    # Exit with a failed status
    Exit 1  

}

else {

    # Filter out links that aren't http* and links without names
    $links = $page.Links | Where-Object { ( $_.href -like 'http*' ) }

    # Iterate through the links and test each one
    foreach( $l in $links ) {

        # Test the link following up to 99 redirects and suppress error output to STDOUT
        $test = Invoke-WebRequest -Method Get -UseBasicParsing -Uri $l.href -MaximumRedirection 99 -ErrorAction SilentlyContinue -WarningAction SilentlyContinue

        # Build a new psobject with results for this link
        $item = New-Object -TypeName psobject
    	$item | Add-Member -MemberType NoteProperty -Name 'targetURL' -Value $Target
        $item | Add-Member -MemberType NoteProperty -Name 'targetSuccess' -Value $true # Did my target page actually load?
        $item | Add-Member -MemberType NoteProperty -Name 'linkURL' -Value $l.href
        $item | Add-Member -MemberType NoteProperty -Name 'statusCode' -Value $test.StatusCode

        # Add the psobject to our Results array
        $results += $item

    }

    # Print the results to STDOUT in JSON format
    $results | ConvertTo-Json

}
#endregion Execution
