
$ErrorActionPreference = 'Stop';
$packageName = 'octant'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vmware-tanzu/octant/releases/download/v0.10.0/octant_0.10.0_Windows-64bit.zip'
$checksum64 = 'dcbf6993dbd4aaa0933e85cf1c20b59b9c05f271d59abfc9317677bb62be48f2'
$checksumType64= 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64

  softwareName  = 'octant*'

  checksum64    = $checksum64
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

