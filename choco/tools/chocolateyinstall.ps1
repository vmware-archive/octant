
$ErrorActionPreference = 'Stop';
$packageName = 'hcli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/heptio/developer-dash/releases/download/v0.1.1/hcli_0.1.1_Windows-64bit.zip'
$checksum64 = '9150f16fe49834ebd3ad065ac5a77acd506cd6ec5232fbe3b2fa191588670a47'
$checksumType64= 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64

  softwareName  = 'hcli*'

  checksum64    = $checksum64
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

