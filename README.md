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
Appends additional keys to your Info.plist. The resulting file will be encoded in binary *(even when using the `-unsigned` flag)* and the app signature will still be invalidated.

---

### Codesigning your app:

```shell
go-ios sign -profile=myprofile.mobileprovision
```

```shell
go-ios sign -profile=myprofile.mobileprovision -cert=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```
To codesign your app, you'll have to provide a mobileprovision file yourself *(which requires an [apple developer membership](http://developer.apple.com))*

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
*Note: If the app is already installed and does not match the current build's signature, you must uninstall the old version before deploying. Read the note #4 for more details*

---

View available commands:
```shell
go-ios export -help
go-ios build -help
go-ios sign -help
```
*The export command just builds, signs, and installs all at once. Unsigned builds are created automatically when the -profile flag is used. You can use the -no-install flag if you want don't want it to automatically install to your connected device.*

# Notes
1. There is a flag to build for an ios simulator but it's unstable and doesn't work with cgo.
2. You **MAY** be able to sign your app with a mobileprovision provided by xcode but the app will only stay on your device for up to 7 days. Afterwards, it'll become unavailable and require a rebuild. I'm not sure how often xcode rotates the provisioning profiles but it may require you to uninstall the app before it can be deployed again. I haven't tested this though so take that with a grain of salt.
3. This tool was intended to be used with [gioui](https://gioui.org) and uses gomobile under the hood so you **MUST** import gomobile into your main file or it will not compile.
    -   This means some libraries *(like fyne.io)* will not work because they already import gomobile somewhere internally and importing it twice results in linker errors.
4. When deploying, if your app is already installed *(if there's an with the same bundle identifier)*, and the signature does not match the build you're deploying, apple won't treat it as an upgrade and the app must be uninstalled first. This includes:
    1. Signing your app with a different provisioning profile or certificate.
       - For example, the gomobile build command embeds its own provisioning profile provided by xcode. If you decide to use go-ios and codesign with your own profile, apple won't replace/upgrade the build because the application identifiers don't match.
    2. Switching to an unsigned build or vice versa.
        - This is because gomobile is used internally which automatically signs the build. (read above)
    3. Building your app on a different device.
5. I have not looked much into signing certificates and do not know if changing certificates requires a rebuild.
### TLDR
---
To keep it simple, I recommend using the **SAME** mobile provisioning profile and signing certificate for each build. If you have installed your app **BEFORE** using go-ios, you must uninstall it or your app won't deploy.