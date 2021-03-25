
$ErrorActionPreference = 'Stop';
$packageName = 'octant'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vmware-tanzu/octant/releases/download/v0.18.0/octant_0.18.0_Windows-64bit.zip'
$checksum64 = '490c8524fa7fcfa364b5200da255e0265df41ea61646f12a8a508b67d5127ea3'
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

