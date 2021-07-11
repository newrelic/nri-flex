[timespan]$downTime = New-TimeSpan -start 0 -end 0
[timespan]$totalDownTime = New-TimeSpan -start 0 -end 0
Get-EventLog -LogName system | 
Where-Object `
{ $_.eventid -eq 6005 -OR $_.eventID -eq 6006 -AND $_.timegenerated -gt (get-date).adddays(-1) } | 
Sort-Object -Property timegenerated |
Foreach-Object `
{
  if ($_.EventID -eq 6006)
     { 
       $down = $_.TimeGenerated
     } #end if eventID
  Else
     { 
      $up = $_.TimeGenerated 
     } #end else
   if($down -AND $up)
     {
      if($down -ge $up) 
         { 
           Write-Host -foregroundColor Red "*** Invalid data. Ignoring $($up)"
           $up = $down 
          } #end if down is greater than up
       [timespan]$CurrentDownTime = new-TimeSpan -start $down -end $up
       [timeSpan]$TotalDownTime = $currentDownTime + $TotalDownTime
       $down = $up = $null
     } #end if down and up
} #end foreach
#"Total down time on $env:computername is:"
#$TotaldownTime
$minutesInMonth = (24*60)*30
$minutesInDay = 24*60
$downTimeMinutes = $TotaldownTime.TotalMinutes
$percentUpTime = "{0:n2}" -f (100 - ($downTimeMinutes/$minutesInDay)*100)
"{
    ""UptimePercent1Day"":  $percentUptime
}"