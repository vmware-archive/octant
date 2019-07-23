
$ErrorActionPreference = 'Stop';
$packageName = 'octant'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vmware/octant/releases/download/v0.5.1/octant_0.5.1_Windows-64bit.zip'
$checksum64 = '7d2029f2388ebf141e9252a62dae10de32a4cd4dfbdad2ee345c04d5a14f3c5c'
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

