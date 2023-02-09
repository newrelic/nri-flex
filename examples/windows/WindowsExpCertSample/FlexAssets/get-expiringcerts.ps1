$Certs = get-childitem -path cert: -recurse -Expiringindays 300 | Select-Object Issuer, NotAfter, Subject, FriendlyName
$StartDate = (Get-Date)
$CertsToSend = @()
Foreach ($Cert in $Certs) {
    $CertsObject = @{
        ExpiringIn = New-TimeSpan -Start $StartDate -End $Cert.NotAfter
        ExpirationDate = $Cert.NotAfter | Get-Date -Uformat %s
        Issuer = $Cert.Issuer
        Subject = $Cert.Subject
        FriendlyName = $Cert.FriendlyName
    }
    $CertsToSend += $CertsObject

}
$CertsToSend | ConvertTo-Json