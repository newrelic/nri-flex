$colProcs = Get-wmiobject win32_process | select *,@{Name="Owner";Expression={($_.GetOwner()).User}}
$colPerfs = Get-wmiobject win32_perfformatteddata_perfproc_process  

# in case you run this script for a remote computer, pls. specify computername and credentials
#$computername = "<COMPUTERNAME>"
#$colProcs = Get-wmiobject win32_process -computername $computername  -Credential $credential | select *, @{Name = "Owner"; Expression = { ($_.GetOwner()).User } }
#$colPerfs = Get-wmiobject win32_perfformatteddata_perfproc_process -computername $computername  -Credential $credential

$colTasklist = @()
 
foreach ($perf in $colPerfs) {
	if ($perf.Name.StartsWith("w3wp"))
	{
		$proc = $colProcs | Where-Object { $_.ProcessId -eq $perf.IDProcess }
		$process = New-Object System.Object
		$process | Add-Member -MemberType NoteProperty -name "PerfName" -value $perf.Name 
		$perfCaption = $perf.Caption
		if ( $perfCaption -eq $null )
		{
			$perfCaption = ""
		}
		$process | Add-Member -MemberType NoteProperty -name "PerfCaption" -value $perfCaption
		$perfDescription = $perf.Description
		if ( $perfDescription -eq $null )
		{
			$perfDescription = ""
		}		
		$process | Add-Member -MemberType NoteProperty -name "PerfDescription" -value $perfDescription 
		$process | Add-Member -MemberType NoteProperty -name "PerfElapsedTime" -value $perf.ElapsedTime 
		$process | Add-Member -MemberType NoteProperty -name "PerfHandleCount" -value $perf.HandleCount
		$process | Add-Member -MemberType NoteProperty -name "PerfIDProcess" -value $perf.IDProcess	
		$process | Add-Member -MemberType NoteProperty -name "PerfIODataBytesPersec" -value $perf.IODataBytesPersec	
		$process | Add-Member -MemberType NoteProperty -name "PerfIODataOperationsPersec" -value $perf.IODataOperationsPersec
		$process | Add-Member -MemberType NoteProperty -name "PerfIOOtherBytesPersec" -value $perf.IOOtherBytesPersec		
		$process | Add-Member -MemberType NoteProperty -name "PerfIOOtherOperationsPersec" -value $perf.IOOtherOperationsPersec
		$process | Add-Member -MemberType NoteProperty -name "PerfIOReadBytesPersec" -value $perf.IOReadBytesPersec		
		$process | Add-Member -MemberType NoteProperty -name "PerfIOReadOperationsPersec" -value $perf.IOReadOperationsPersec
		$process | Add-Member -MemberType NoteProperty -name "PerfIOWriteBytesPersec" -value $perf.IOWriteBytesPersec		
		$process | Add-Member -MemberType NoteProperty -name "PerfIOWriteOperationsPersec" -value $perf.IOWriteOperationsPersec		
		$process | Add-Member -MemberType NoteProperty -name "PerfPageFaultsPersec" -value $perf.PageFaultsPersec	
		$process | Add-Member -MemberType NoteProperty -name "PerfPageFileBytes" -value $perf.PageFileBytes			
		$process | Add-Member -MemberType NoteProperty -name "PerfPageFileBytesPeak" -value $perf.PageFileBytesPeak		
		$process | Add-Member -MemberType NoteProperty -name "PerfPercentPrivilegedTime" -value $perf.PercentPrivilegedTime	
		$process | Add-Member -MemberType NoteProperty -name "PerfPercentProcessorTime" -value $perf.PercentProcessorTime			
		$process | Add-Member -MemberType NoteProperty -name "PerfPercentUserTime" -value $perf.PercentUserTime			
		$process | Add-Member -MemberType NoteProperty -name "PerfPoolNonpagedBytes" -value $perf.PoolNonpagedBytes			
		$process | Add-Member -MemberType NoteProperty -name "PerfPoolPagedBytes" -value $perf.PoolPagedBytes			
		$process | Add-Member -MemberType NoteProperty -name "PerfPrivateBytes" -value $perf.PrivateBytes		
		$process | Add-Member -MemberType NoteProperty -name "PerfThreadCount" -value $perf.ThreadCount				
		$process | Add-Member -MemberType NoteProperty -name "PerfVirtualBytes" -value $perf.VirtualBytes				
		$process | Add-Member -MemberType NoteProperty -name "PerfWorkingSet" -value $perf.WorkingSet
		$process | Add-Member -MemberType NoteProperty -name "PerfWorkingSetPrivate" -value $perf.WorkingSetPrivate
		$procOwner = $proc.Owner
		if ( $procOwner -eq $null )
		{
			$procOwner = ""
		}
		$process | Add-Member -MemberType NoteProperty -name "ProcOwner" -value $procOwner
		$procExecutablePath = $proc.ExecutablePath
		if ( $procExecutablePath -eq $null )
		{
			$procExecutablePath = ""
		}
		$process | Add-Member -MemberType NoteProperty -name "ProcExecutablePath" -value $procExecutablePath

		$colTasklist += $process
	}
}

$colTasklist | ConvertTo-Json
