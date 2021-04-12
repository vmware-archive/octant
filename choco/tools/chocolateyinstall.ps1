
$ErrorActionPreference = 'Stop';
$packageName = 'octant'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vmware-tanzu/octant/releases/download/v0.19.0/octant_0.19.0_Windows-64bit.zip'
$checksum64 = 'f99f6083edc63d3e73f27495b3b878914a39c68950ddc3d385224a594ca910d9'
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

