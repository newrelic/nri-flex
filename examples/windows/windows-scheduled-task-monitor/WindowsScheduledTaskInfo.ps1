$taskInfos = Get-ScheduledTask | Where-Object {
	($_.TaskPath -notlike '\Microsoft\*' -and $_.TaskPath -ne '\')
   } | ForEach-Object {
	   (Get-ScheduledTaskInfo -TaskName $_.TaskName -TaskPath $_.TaskPath | Select-Object -Property *, 
		   @{Name='LastTaskResultHex'; Expression={'0x'+([System.Convert]::ToString($_.LastTaskResult,16)).ToUpper() }}, 
		   @{Name='LastRunTimeInMillis'; Expression={([int] (Get-Date $_.LastRunTime.ToUniversalTime() -UFormat '%s'))*1000}},
		   @{Name='NextRunTimeInMillis'; Expression={([int] (Get-Date $_.NextRunTime.ToUniversalTime() -UFormat '%s'))*1000}},
		   @{Name='TaskFullPath'; Expression={$_.TaskPath + $_.TaskName}},
		   @{Name='TaskFullPathStringLiteral'; Expression={'\' + $_.TaskPath + '\' + $_.TaskName}
		   } -ExcludeProperty PSComputerName, CimClass, CimInstanceProperties, CimSystemProperties, CimInstance)
   } | ConvertTo-Json

Write-Output $taskInfos