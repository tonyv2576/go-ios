# Getting started
### Installation
```shell
go install github.com/tonyv2576/go-ios
```

# Dependencies
 - Golang 1.18+ or later
 - MACOS Sequoia 10.11 or later
 - ios-deploy
    - If using [Homebrew](https://brew.sh) do: `brew install ios-deploy`
 - PlistBuddy
    - Should come pre-installed on your Mac.


# Usage

### Building your project:
```shell
go-ios build -bundle=com.example ./cmd
```
Creates a module.app folder with "module" being the name of your go module.

---

### Building with an editable Info.plist:
```shell
go-ios build -bundle=com.example -unsigned ./cmd
```
Re-encodes your Info.plist in xml1 format. This, however, invalidates your apps signature and you won't be able to install it to your device without codesigning it.

---

### Building with an external plist/xml file:
```shell
go-ios build -bundle=com.example -append=example.xml ./cmd
```
```xml
<!-- example.xml -->
<dict>
    <key>NSCameraUsageDescription</key>
    <string>Need camera access to record video</string>
</dict>
```
Appends additional keys to your Info.plist. The resulting file will be encoded in binary (*even with the `-unsigned` flag*) and the app signature will still be invalidated.

---

### Codesigning your app:

```shell
go-ios sign -profile=myprofile.mobileprovision
```

```shell
go-ios sign -profile=myprofile.mobileprovision -cert=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```
To codesign your app, you'll have to provide a mobileprovision file yourself (*which requires an [apple developer membership](http://developer.apple.com)*)

---

If your device has multiple certificates installed, you'll have to provide the hash well which you can find by doing:
```shell
security find-identity -v -p codesigning

# 1) XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX "Apple Development: Team Name (TEAM)"
# 2) XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX "Apple Development: Team Name (TEAM)"
#   2 valid identities found
```

---

Installing the app to your device:
```shell
go-ios install
```
Installs directly onto your device. Plugging your device into your mac will speed up this step greatly.

---

View available commands:
```shell
go-ios build -help
go-ios sign -help
```

# Notes

1. There is a flag to build for an ios simulator but it's unstable and doesn't work with cgo.
2. You **MAY** be able to sign your app with a mobileprovision provided by xcode but this has not been tested.
3. This tool was intended to be used with [gioui](https://gioui.org) and uses gomobile under the hood so you **MUST** import gomobile into your main file or it will not compile.
    -   This means some libraries (like fyne.io) will not work because they already import gomobile somewhere internally and importing it twice results in linker errors.