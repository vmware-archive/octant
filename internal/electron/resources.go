/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

// AssetFunc loads an asset as bytes.
type AssetFunc func(name string) ([]byte, error)

// AssetDirFunc returns the files in an asset directory.
type AssetDirFunc func(name string) ([]string, error)

// RestoreAssetsFunc restores assets.
type RestoreAssetsFunc func(dir string, name string) error

func absoluteResourcesPath(a *astilectron.Astilectron, relativeResourcesPath string) string {
	return filepath.Join(a.Paths().DataDirectory(), relativeResourcesPath)
}

func restoreResources(l astikit.SeverityLogger, a *astilectron.Astilectron, asset AssetFunc, assetDir AssetDirFunc, assetRestorer RestoreAssetsFunc, relativeResourcesPath string) (err error) {
	// Check resources
	var restore bool
	var computedChecksums map[string]string
	var checksumsPath string
	if restore, computedChecksums, checksumsPath, err = checkResources(l, a, asset, assetDir, relativeResourcesPath); err != nil {
		err = fmt.Errorf("checking resources failed: %w", err)
		return
	}

	// Restore resources
	if restore {
		if err = restoreResourcesFunc(l, a, relativeResourcesPath, assetRestorer, computedChecksums, checksumsPath); err != nil {
			err = fmt.Errorf("restoring resources failed: %w", err)
			return
		}
	} else {
		l.Debug("Skipping restoring resources...")
	}
	return
}

func checkResources(l astikit.SeverityLogger, a *astilectron.Astilectron, asset AssetFunc, assetDir AssetDirFunc, relativeResourcesPath string) (restore bool, computedChecksums map[string]string, checksumsPath string, err error) {
	// Compute checksums
	arp := absoluteResourcesPath(a, relativeResourcesPath)
	checksumsPath = filepath.Join(arp, "checksums.json")
	if asset != nil && assetDir != nil {
		computedChecksums = make(map[string]string)
		if err = checksumAssets(asset, assetDir, relativeResourcesPath, computedChecksums); err != nil {
			err = fmt.Errorf("getting checksum of assets failed: %w", err)
			return
		}
	}

	// Stat resources
	if _, err = os.Stat(arp); err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("stating %s failed: %w", arp, err)
		return
	} else if os.IsNotExist(err) {
		l.Debug("Resources folder doesn't exist, restoring resources...")
		err = nil
		restore = true
		return
	}

	// No computed checksums
	if computedChecksums == nil {
		l.Debug("No computed checksums, restoring resources...")
		restore = true
		return
	}

	// Stat checksums file
	if _, err = os.Stat(checksumsPath); err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("stating %s failed: %w", checksumsPath, err)
		return
	} else if os.IsNotExist(err) {
		l.Debug("Checksums file doesn't exist, restoring resources...")
		err = nil
		restore = true
		return
	}

	// Open resources checksums
	var f *os.File
	if f, err = os.Open(checksumsPath); err != nil {
		err = fmt.Errorf("opening %s failed: %w", checksumsPath, err)
		return
	}
	defer f.Close()

	// Unmarshal checksums
	var unmarshaledChecksums map[string]string
	if err = json.NewDecoder(f).Decode(&unmarshaledChecksums); err != nil {
		err = fmt.Errorf("unmarshaling checksums failed: %w", err)
		return
	}

	// Check number of paths
	if len(unmarshaledChecksums) != len(computedChecksums) {
		l.Debugf("%d paths in unmarshaled checksums != %d paths in computed checksums, restoring resources...", len(unmarshaledChecksums), len(computedChecksums))
		restore = true
		return
	}

	// Loop through computed checksums
	for p, c := range computedChecksums {
		// Path doesn't exist in unmarshaled checksums
		v, ok := unmarshaledChecksums[p]
		if !ok {
			l.Debugf("Path %s doesn't exist in unmarshaled checksums, restoring resources...", p)
			restore = true
			return
		}

		// Checksums are different
		if c != v {
			l.Debugf("Unmarshaled checksum (%s) != computed checksum (%s) for path %s, restoring resources...", v, c, p)
			restore = true
			return
		}
	}
	return
}

func checksumAssets(asset AssetFunc, assetDir AssetDirFunc, name string, m map[string]string) (err error) {
	// Get children
	children, errDir := assetDir(name)

	// File
	if errDir != nil {
		// Get checksum
		var h string
		if h, err = checksumAsset(asset, name); err != nil {
			err = fmt.Errorf("getting checksum of %s failed: %w", name, err)
			return
		}
		m[name] = h
		return
	}

	// Dir
	for _, child := range children {
		if err = checksumAssets(asset, assetDir, filepath.Join(name, child), m); err != nil {
			err = fmt.Errorf("getting checksum of assets in %s failed: %w", name, err)
			return
		}
	}
	return
}

func checksumAsset(asset AssetFunc, name string) (o string, err error) {
	// Get data
	var b []byte
	if b, err = asset(name); err != nil {
		err = fmt.Errorf("getting data from asset %s failed: %w", name, err)
		return
	}

	// Hash
	h := md5.New()
	if _, err = h.Write(b); err != nil {
		err = fmt.Errorf("writing data of asset %s to hash failed: %w", name, err)
		return
	}
	o = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

func restoreResourcesFunc(l astikit.SeverityLogger, a *astilectron.Astilectron, relativeResourcesPath string, assetRestorer RestoreAssetsFunc, computedChecksums map[string]string, checksumsPath string) (err error) {
	// Remove resources
	arp := absoluteResourcesPath(a, relativeResourcesPath)
	l.Debugf("Removing %s", arp)
	if err = os.RemoveAll(arp); err != nil {
		err = fmt.Errorf("removing %s failed: %w", arp, err)
		return
	}

	// Restore resources
	l.Debugf("Restoring resources in %s", arp)
	if err = assetRestorer(a.Paths().DataDirectory(), relativeResourcesPath); err != nil {
		err = fmt.Errorf("restoring resources in %s failed: %w", arp, err)
		return
	}

	// Write checksums
	if computedChecksums != nil {
		// Create checksums file
		var f *os.File
		if f, err = os.Create(checksumsPath); err != nil {
			err = fmt.Errorf("creating %s failed: %w", checksumsPath, err)
			return
		}
		defer f.Close()

		// Marshal
		if err = json.NewEncoder(f).Encode(computedChecksums); err != nil {
			err = fmt.Errorf("marshaling checksums failed: %w", err)
			return
		}
	}
	return
}
