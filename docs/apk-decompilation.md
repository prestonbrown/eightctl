# Eight Sleep APK Decompilation Guide

This document describes how to extract API endpoint information from the Eight Sleep Android application.

## APK Sources

### APKPure
- URL: https://apkpure.com/eight-sleep/com.eightsleep.eight/versions
- Provides historical versions
- Latest available: 7.41.65+

### APKMirror
- URL: https://www.apkmirror.com/apk/eight-sleep-inc/eight-sleep/
- Trusted source with verification
- Multiple architecture variants (arm64-v8a, armeabi-v7a, x86_64)

### Direct Download
APK files can be extracted from a rooted device or Android emulator:
```bash
adb shell pm path com.eightsleep.eight
adb pull /data/app/.../base.apk
```

## Decompilation Tools

### JADX (Recommended)

JADX converts Android DEX bytecode to readable Java source code.

**Installation:**
```bash
# macOS
brew install jadx

# Linux (via GitHub releases)
wget https://github.com/skylot/jadx/releases/latest/download/jadx-x.x.x.zip
unzip jadx-*.zip -d jadx
export PATH=$PATH:$(pwd)/jadx/bin
```

**Usage:**
```bash
# GUI mode (recommended for exploration)
jadx-gui eightsleep.apk

# CLI mode (for scripted extraction)
jadx -d output_dir eightsleep.apk
```

### APKtool

APKtool decodes resources and can rebuild APKs. Useful for extracting strings and XML resources.

**Installation:**
```bash
# macOS
brew install apktool

# Linux
wget https://raw.githubusercontent.com/iBotPeaches/Apktool/master/scripts/linux/apktool
chmod +x apktool
wget https://bitbucket.org/iBotPeaches/apktool/downloads/apktool_x.x.x.jar
mv apktool_*.jar apktool.jar
```

**Usage:**
```bash
apktool d eightsleep.apk -o output_dir
```

## Search Targets

### API URLs and Endpoints

Search for these strings in decompiled code:

```
8slp.net
client-api.8slp.net
auth-api.8slp.net
app-api.8slp.net
/v1/
/v2/
/users/
/devices/
/alarms
/routines
/temperature
/autopilot
/audio
/base
/travel
/household
```

### Relevant Classes

Look for these patterns in class names:

```
*Api.java
*ApiService.java
*Endpoint*.java
*Network*.java
*Retrofit*.java
*Repository.java
*Client.java
```

### BuildConfig

Check `BuildConfig.java` for:
- API base URLs
- OAuth client IDs
- Feature flags

### Retrofit Annotations

If the app uses Retrofit (common for Android), look for:

```java
@GET("/users/{userId}/temperature")
@POST("/v1/users/{userId}/alarms")
@PUT("/v2/users/{userId}/routines/{routineId}")
@DELETE("/devices/{deviceId}/schedules/{scheduleId}")
```

## Known Findings

### From APK v7.39.17

OAuth credentials (currently in use):
- `client_id`: `0894c7f33bb94800a03f1f4df13a4f38`
- `client_secret`: `f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76`

### API Version Migration

Eight Sleep appears to be migrating from v1 to v2 APIs:
- v1: `/v1/users/{userId}/alarms`
- v2: `/v2/users/{userId}/routines`

The v2 routines endpoint combines alarms, bedtime, and other scheduling features.

## Extraction Script

Example script to search for API endpoints:

```bash
#!/bin/bash
# extract-endpoints.sh

APK_DIR="$1"
if [ -z "$APK_DIR" ]; then
    echo "Usage: $0 <jadx-output-dir>"
    exit 1
fi

echo "=== API Base URLs ==="
grep -r "8slp.net" "$APK_DIR" --include="*.java"

echo ""
echo "=== HTTP Endpoints ==="
grep -rE '@(GET|POST|PUT|DELETE|PATCH)\(' "$APK_DIR" --include="*.java" -A1

echo ""
echo "=== URL Patterns ==="
grep -rE '"/v[12]/|"/users/|"/devices/' "$APK_DIR" --include="*.java"
```

## Security Notes

- APK decompilation is for interoperability research
- Do not redistribute decompiled code
- OAuth credentials should be treated as semi-public (embedded in published apps)
- Respect Eight Sleep's terms of service

## Related Resources

- JADX GitHub: https://github.com/skylot/jadx
- APKtool GitHub: https://github.com/iBotPeaches/Apktool
- pyEight (reference implementation): https://github.com/lukas-clarke/pyEight
- Android RE guide: https://mobile-security.gitbook.io/mobile-security-testing-guide/android-testing-guide/0x05c-reverse-engineering-and-tampering
