#Update AZ Location to the location you'd like to search by based on this doc:
## https://docs.microsoft.com/en-us/powershell/module/az.compute/get-azvmusage?view=azps-7.5.0
$AZLocation = "West US 3"
$VMUsage = Get-AzVMUsage -Location $AZLocation
$VMUsage | ConvertTo-Json