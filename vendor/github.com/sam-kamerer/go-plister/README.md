# go-plister [![Build Status](https://travis-ci.org/sam-kamerer/go-plister.svg?branch=master)](https://travis-ci.org/sam-kamerer/go-plister) [![Coverage Status](https://coveralls.io/repos/github/sam-kamerer/go-plister/badge.svg?branch=master)](https://coveralls.io/github/sam-kamerer/go-plister?branch=master)
A simple Apple Property List generator

### Usage
```go
package app

import (
	"log"
	"os"
	
	"github.com/sam-kamerer/go-plister"
)

var dict = map[string]interface{}{
    "CFBundlePackageType":     "APPL",
    "CFBundleInfoDictionaryVersion": "6.0",
    "CFBundleIconFile":        "icon.icns",
    "CFBundleDisplayName":     "Best App",
    "CFBundleExecutable":      "app_binary",
    "CFBundleName":            "BestApp",
    "CFBundleIdentifier":      "com.company.BestApp",
    "LSUIElement":             "NO",
    "LSMinimumSystemVersion":  "10.11",
    "NSHighResolutionCapable": true,
    "NSAppTransportSecurity": map[string]interface{}{
        "NSAllowsArbitraryLoads": true,
    },
    "CFBundleURLTypes": []map[string]interface{}{
        {
            "CFBundleTypeRole":   "Viewer",
            "CFBundleURLName":    "com.developer.testapp",
            "CFBundleURLSchemes": []interface{}{"testappscheme"},
        }, {
            "CFBundleTypeRole":   "Reader",
            "CFBundleURLName":    "com.developer.testapp",
            "CFBundleURLSchemes": []interface{}{"testappscheme-read"},
        },
    },
}

func main() {
    infoPlist := plister.MapToInfoPlist(dict)
    if err := plister.Generate("path/to/Info.plist", infoPlist); err != nil {
    	log.Fatal(err)
    }
    
    // or
    
    if err := plister.GenerateFromMap("path/to/Info.plist", dict); err != nil {
    	log.Fatal(err)
    }
    
    // or
    
    fp, err := os.Open("path/to/Info.plist")
    if err != nil {
    	log.Fatal(err)
    }
    if err := plister.Fprint(fp, infoPlist); err != nil {
    	log.Fatal(err)
    }
}
```
