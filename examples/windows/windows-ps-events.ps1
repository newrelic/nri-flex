###
# Parameters (a.k.a. Command Line Arguments)
# Usage: -LogName "LogName" -EventIds 4608:4609:4946
# -LogName      The name of the event log to gather events from.  Ex: System, Appication, etc. (Required)
# -EventIds     An optional, colon delimited, list of event ids to gather.
#
# Output:  JSON Array of Requested Windows Event Log Entries with a field name prefix of "winev."
# Example:  winev.Id, winev.Message
###

param (
    [string]$LogName=$(throw "-LogName is mandatory"),
    [string]$EventIds
)

###
# Logic to handle getting new log entries by saving current date to file
# to use as -After argument of Get-Date in next pull. On first run we use current date.
# On subsequent runs it will use last date written to file.
#
# Uses LogName param to create timestamp for each LogName
###

$LAST_PULL_TIMESTAMP_FILE = "c:\Program Files\New Relic\newrelic-infra\integrations.d\windows-ps-events-$LogName.timestamp"


###
# If timestamp file exists, use it; otherwise,
# set timestamp to 15 minutes ago to pull some data on
# first run.
###

if (Test-Path $LAST_PULL_TIMESTAMP_FILE -PathType Leaf) {

    $timestamp = Get-Content -Path $LAST_PULL_TIMESTAMP_FILE -Encoding String | Out-String
    $timestamp = [DateTime] $timestamp.ToString()

} else{

    $timestamp = (Get-Date).AddMinutes(-240)

}

###
# Write timestamp to file to pull on next run.
###
Set-Content -Path $LAST_PULL_TIMESTAMP_FILE -Value (Get-Date -Format o)

###
# Pull events using -After param with timestamp
###
#Write-output $LogName,$EventIds,$timestamp
$events = $(Get-WinEvent -FilterHashtable @{LogName=$LogName;StartTime=$timestamp}) 2>out-null

###
# If event ids were given, filter to keep only the events having an id in our event id list.
###
if ($EventIds) {
    $eventIdStrings = $EventIds.Split(":")
    $eventIdNums = @()
    foreach ($eventId in $eventIdStrings) {
        $eventIdNums += [convert]::ToInt32($eventId)
    }

    # Iterate over the events and copy only those we want to keep into the filteredEvents.
    $filteredEvents = @()
    foreach ($event in $events) {
        if (-Not ($eventIdNums -Contains $event.Id)) {
            continue
        }
        $filteredEvents += $event
    }
    $events = $filteredEvents
}
$modevents = @()
foreach ($event in $events) {
  $newev = @{}
  $event.PSObject.Properties | foreach {
    $newev.add("winev." + $_.Name, $_.Value)
  }
  $o = New-Object psobject -Property $newev
  $modevents += $o
}
$events = $modevents

###
# Add required 'event_type' to objects from Get-EventLog.
# Add optional 'log_name' value to object.
###
$events.ForEach({
    Add-Member -NotePropertyName 'log_name' -NotePropertyValue $LogName -InputObject $_
});

###
# Create hash table in required format for Infrastructure, populated
# with event object log data and pipe to ConvertTo-Json with
# -Compress argument required in order for Infrastructure to consume.
###
$payload = @($events) | ConvertTo-Json -Compress


###
# Output json string created above with regex to normalize date strings
# post json string conversion. Alternatively, you could create a
# new -NotePropertyName with the proper date string and remove
# the original object property. 
###
Write-Output ($payload -replace '"\\\/Date\((\d+)\)\\\/\"' ,'$1')
