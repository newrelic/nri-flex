#region Top of Script

#requires -version 4

<#
.SYNOPSIS
	Queries for target Perfmon counter values
.DESCRIPTION
    This script uses a static list of values locally
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	12/21/2020
	Purpose/Change:	Initial script development
#>

#endregion Top of Script

#####-----------------------------------------------------------------------------------------#####

#region Execution 

# Build an empty array to add our results to
$results = @()

# Build an array of pscustomobjects to hold our target counters
$counters = @(

    [ PSCustomObject ]@{ CounterName = "percentProcessorTime"; CounterPath = "\Processor(*)\% Processor Time" }
    [ PSCustomObject ]@{ CounterName = "memoryPercentCommitedBytes"; CounterPath = "\Memory\% Committed Bytes In Use" }
    [ PSCustomObject ]@{ CounterName = "diskWritesPerSecond"; CounterPath = "\PhysicalDisk(*)\Disk Writes/sec" }
    [ PSCustomObject ]@{ CounterName = "currentDiskQueueLength"; CounterPath = "\LogicalDisk(*)\Current Disk Queue Length" }

)

# Iterate through our target counters and grab our results
foreach( $c in $counters ) {

    $query = (Get-Counter -MaxSamples 1 -Counter $c.CounterPath).CounterSamples
    foreach( $q in $query ) {
        
        $item = @(

            [ PSCustomObject ]@{ CounterName = $c.CounterName; CounterPath = $c.CounterPath; CounterInstance = $q.InstanceName; CounterValue = $q.CookedValue }

        )

        $results += $item 

    }

}

# Print the results to STDOUT in JSON for Flex to pickup
$results | ConvertTo-Json

#endregion Execution
