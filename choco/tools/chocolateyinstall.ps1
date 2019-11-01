
$ErrorActionPreference = 'Stop';
$packageName = 'octant'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vmware-tanzu/octant/releases/download/v0.9.1/octant_0.9.1_Windows-64bit.zip'
$checksum64 = '3d6131a5d5e68f8e7b3977b8b614d97daa3251f20280c486f0e124eb1e5cf458'
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

