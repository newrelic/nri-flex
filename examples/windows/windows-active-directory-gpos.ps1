 #region Top of Script

#requires -version 4
#requires -Module ActiveDirectory

<#
.SYNOPSIS
	Queries Active Directory for all GPOs
.DESCRIPTION
    This script uses a static list of queries
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	9-May-2023
	Purpose/Change:	Initial script development
#>

#endregion Top of Script

#####-----------------------------------------------------------------------------------------#####

#region Execution

$gpos = Get-GPO -All

# Build an empty array to add our results to
$results = @()

# Iterate through our GPO checks and grab the results
foreach( $g in $gpos ){

  # Run a GPO Report to get the link status
  $gpoReport = [xml](Get-GPOReport -Guid $g.id -ReportType xml) | Select-Object -ExpandProperty GPO

  # Build a custom object to pass into the results
  $item = @(

    [ PSCustomObject ]@{ 

      gpoName = $g.DisplayName
      gpoDomain = $g.DomainName
      gpoStatus = $g.GpoStatus.ToString()
      gpoCreated = ( [DateTimeOffset ]$g.CreationTime ).ToUnixTimeSeconds()
      gpoLastModified = ( [DateTimeOffset ]$g.ModificationTime ).ToUnixTimeSeconds()
      gpoEnabled = if ( $gpoReport.LinksTo ) { $gpoReport | Select-Object -ExpandProperty LinksTo | Select-Object -ExpandProperty Enabled } else { "false" }
      gpoLinkedFrom = if ( $gpoReport.LinksTo ) { $gpoReport | Select-Object -ExpandProperty LinksTo | Select-Object -ExpandProperty SOMPath } else { "none" }

      }

  )

  $results += $item

}

# Print the results to STDOUT in JSON for Flex to pickup
$results | ConvertTo-Json

#endregion Execution
