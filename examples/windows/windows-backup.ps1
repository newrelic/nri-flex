#region Top of Script

#requires -version 4

<#
.SYNOPSIS
	Get latest Windows backup data
.DESCRIPTION
    This script gets the latest Windows backup data and converts it for NR
.NOTES
	Version:		1.0
	Author:			Samuel Vandamme
	Creation Date:	2020-02-02
	Purpose/Change:	Initial script development
#>

#endregion Top of Script

#####-----------------------------------------------------------------------------------------#####

# Get backup summary data
$backupData = Get-WBSummary

# Convert and output as JSON
$props = @{
    NextBackupTime = $backupData.NextBackupTime | Get-Date -UFormat %s
    NumberOfVersions = $backupData.NumberOfVersions
    LastSuccessfulBackupTime = $backupData.LastSuccessfulBackupTime | Get-Date -UFormat %s
    LastSuccessfulBackupTargetPath = $backupData.LastSuccessfulBackupTargetPath
    LastSuccessfulBackupTargetLabel = $backupData.LastSuccessfulBackupTargetLabel
    LastBackupTime = $backupData.LastBackupTime | Get-Date -UFormat %s
    LastBackupTarget = $backupData.LastBackupTarget
    DetailedMessage = $backupData.DetailedMessage
    LastBackupResultHR = $backupData.LastBackupResultHR
    LastBackupResultDetailedHR = $backupData.LastBackupResultDetailedHR
    CurrentOperationsStatus = $backupData.CurrentOperationStatus
}

$props | ConvertTo-Json
