#region Top of Script

#requires -version 3

<#
.SYNOPSIS
	Collects details on SSL Certificates found on target URL
.DESCRIPTION
	Requires URL parameters to have protocol identifier (http://, https://)
.EXAMPLE
	.\windows-ssl-certificate-expiration-script.ps1 -URL "https://www.newrelic.com"
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	04/30/2020
	Purpose/Change:	Initial script development
#>

#endregion

#####-----------------------------------------------------------------------------------------#####

#region Script Parameters

Param
	(
		
        [ Parameter( Mandatory = $true ) ] [ string ] $URL

	)

#endregion Script Parameters

#####-----------------------------------------------------------------------------------------#####

#region Execution

# Suppress errors from output
$ErrorActionPreference = "SilentlyContinue"

# Disable Certificate Validation so we don't fail on expired certs
[Net.ServicePointManager]::ServerCertificateValidationCallback = { $true }

# Build a .NET HTTP Web Request
$request = $null
$request = [ Net.HttpWebRequest ]::Create( $URL )
# Set the Timeout value (ms)
# Flex default timeout is 10000ms so you'll want to stay under that or adjust the Flex config
$request.Timeout = 5000
# Disable redirects
$request.AllowAutoRedirect = $false

# Execute the request
$request.GetResponse() | Out-Null

# If there's no certificate used, tell us
If( $request.ServicePoint.Certificate -eq $null ) {

    # Build a set of data with empty keys
    $results = New-Object -TypeName PSObject
    $results | Add-Member -MemberType NoteProperty -Name "targetURL" -Value $URL
    $results | Add-Member -MemberType NoteProperty -Name "expirationDate" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "daysLeft" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "isExpired" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certDistinguishedName" -Value "No Certificate Found"
    $results | Add-Member -MemberType NoteProperty -Name "certCommonName" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certOrganization" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certCountry" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certState" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certLocality" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "certificateEffectiveDate" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "issuerDistinguishedName" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "issuerCommonName" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "issuerOrganization" -Value $null
    $results | Add-Member -MemberType NoteProperty -Name "issuerCountry" -Value $null
    #$results | Add-Member -MemberType NoteProperty -Name "certificatePublicKeyString" -Value $null
    #$results | Add-Member -MemberType NoteProperty -Name "certificateSerialNumber" -Value $null
    #$results | Add-Member -MemberType NoteProperty -Name "certificateThumbprint" -Value $null

    $results | ConvertTo-Json

}

Else {

    # Breakdown the Certificate DN into human-readable chunks
    # This Regex Pattern uses a Negative Lookahead to ignore data between `=` and `"`
    # This is necessary to account for situations where the value has a comma
    $certDN = $request.ServicePoint.Certificate.GetName()
    $dnHash = $certDN -split ",\s(?![^=]+`")" | ConvertFrom-StringData

    # Break down the Certificate Issuer into human-readable chunks
    $issuerDN = $request.ServicePoint.Certificate.GetIssuerName()
    $issuerHash = $issuerDN -split ", " | ConvertFrom-StringData

    # Set the dates to epoch
    $expirationDate = [ int ]( New-TimeSpan -Start ( Get-Date '01/01/1970' ) -End $( [ System.DateTime ]::Parse( $request.ServicePoint.Certificate.GetExpirationDateString() ) ) ).TotalSeconds
    $effectiveDate = [ int ]( New-TimeSpan -Start ( Get-Date '01/01/1970' ) -End $( [ System.DateTime ]::Parse( $request.ServicePoint.Certificate.GetEffectiveDateString() ) ) ).TotalSeconds
        
    # Convert expired measurements (negatives) to zero
    $daysLeft = [ math ]::Max( 0, [ int ]( New-TimeSpan -Start ( Get-Date ) -End $( [ System.DateTime ]::Parse( $request.ServicePoint.Certificate.GetExpirationDateString() ) ) ).TotalDays )

    # Collect and export our results
    $results = New-Object -TypeName PSObject
    $results | Add-Member -MemberType NoteProperty -Name "targetURL" -Value $URL
    $results | Add-Member -MemberType NoteProperty -Name "expirationDate" -Value $expirationDate
    $results | Add-Member -MemberType NoteProperty -Name "daysLeft" -Value $daysLeft
    $results | Add-Member -MemberType NoteProperty -Name "isExpired" -Value $( if( $daysLeft -eq 0 ){ $true } else{ $false } )
    $results | Add-Member -MemberType NoteProperty -Name "certDistinguishedName" -Value $certDN 
    $results | Add-Member -MemberType NoteProperty -Name "certCommonName" -Value $( $dnHash."CN" )
    $results | Add-Member -MemberType NoteProperty -Name "certOrganization" -Value $( $dnHash."O" )
    $results | Add-Member -MemberType NoteProperty -Name "certCountry" -Value $( $dnHash."C" )
    $results | Add-Member -MemberType NoteProperty -Name "certState" -Value $( $dnHash."S" )
    $results | Add-Member -MemberType NoteProperty -Name "certLocality" -Value $( $dnHash."L" )
    $results | Add-Member -MemberType NoteProperty -Name "certificateEffectiveDate" -Value $effectiveDate
    $results | Add-Member -MemberType NoteProperty -Name "issuerDistinguishedName" -Value $issuerDN
    $results | Add-Member -MemberType NoteProperty -Name "issuerCommonName" -Value $( $issuerHash."CN" )
    $results | Add-Member -MemberType NoteProperty -Name "issuerOrganization" -Value $( $issuerHash."O" )
    $results | Add-Member -MemberType NoteProperty -Name "issuerCountry" -Value $( $issuerHash."C" )
    #$results | Add-Member -MemberType NoteProperty -Name "certificatePublicKeyString" -Value $( $request.ServicePoint.Certificate.GetPublicKeyString() )
    #$results | Add-Member -MemberType NoteProperty -Name "certificateSerialNumber" -Value $( $request.ServicePoint.Certificate.GetSerialNumberString() )
    #$results | Add-Member -MemberType NoteProperty -Name "certificateThumbprint" -Value $( $request.ServicePoint.Certificate.GetCertHashString() )

    $results | ConvertTo-Json

}

#endregion Execution
