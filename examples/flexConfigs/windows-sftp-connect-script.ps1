#region Top of Script

#requires -version 3 -modules Posh-SSH

<#
.SYNOPSIS
	Tests connection to SFTP Server using User/Pass and outputs to JSON
.DESCRIPTION
	Requires the Posh-SSH module
    https://www.powershellgallery.com/packages/Posh-SSH/
.EXAMPLE
	.\flexSFTP.ps1 -TargetHost "sftpTestServer" -Username "sftpUsername" -Password "password123!"
    -- NOTE -- Dollar Signs ($) are reserved in PowerShell and must be escaped with a backtick { ` } if used in a password
    Ex: "P@$$w0rd!" == "P@`$`$w0rd!"
.NOTES
	Version:		1.0
	Author:			Zack Mutchler
	Creation Date:	02/28/2020
	Purpose/Change:	Initial script development

	Version:		2.0
	Author:			Zack Mutchler
	Creation Date:	02/28/2020
	Purpose/Change:	Added script parameter to support multiple hosts in nri-flex
	
	Version:		2.1
	Author:			Brad Schmitt
	Creation Date: 	03/05/2020
	Purpose/Change: Added user/password parameters

    Version:		2.2
	Author:			Brad Schmitt
	Creation Date: 	03/05/2020
	Purpose/Change: Added port test

    Version:		2.3
	Author:			Zack Mutchler
	Creation Date: 	03/06/2020
	Purpose/Change: Removed port test and added error capture on New-SFTPSession
#>

#endregion

#####-----------------------------------------------------------------------------------------#####

#region Script Parameters

Param
	(
		
        [ Parameter( Mandatory = $true ) ] [ string ] $TargetHost,
		[ Parameter( Mandatory = $true ) ] [ string ] $Username,
		[ Parameter( Mandatory = $true ) ] [ string ] $Password

	)

#endregion Script Parameters

#####-----------------------------------------------------------------------------------------#####
#region Variables

# Create our credential object
$securePass = ConvertTo-SecureString -String $Password -AsPlainText -Force
$creds = New-Object System.Management.Automation.PSCredential ( $Username, $securePass )

# Build our results object
$results = New-Object -TypeName psobject

#endregion Variables

#####-----------------------------------------------------------------------------------------#####

#region Execution

# Build our connection to the SFTP Server using User/Pass
$sftpTest = New-SFTPSession -ComputerName $TargetHost -Credential $creds -AcceptKey -ErrorAction SilentlyContinue -ErrorVariable errorCatch

# If the new-SFTPSession had a problem connecting establishing a TCP connection, tell us why...
if( $errorCatch ) { 

    $results | Add-Member -MemberType NoteProperty -Name "sftpServer" -Value $TargetHost -Force
    $results | Add-Member -MemberType NoteProperty -Name "sftpConnected" -Value $( $errorCatch.Exception.Message ) -Force
    
} 

# Otherwise, add data to our results
else {

    $results | Add-Member -MemberType NoteProperty -Name "sftpServer" -Value $TargetHost -Force
    $results | Add-Member -MemberType NoteProperty -Name "sftpConnected" -Value $sftpTest.Connected -Force

    # Clean up the open SFTP session
    Remove-SFTPSession -SessionId $($sftpTest.SessionId) | Out-Null

}

# Print the results to STDOUT in JSON format
$results | ConvertTo-Json

#endregion Execution
