# NOTES
 - cgo can use objective c as long as its wrapped in c headers
    ```cpp
        // ifdef __cplusplus makes sure the compiler (if its obj-c or cpp compiler) treats these functions as normal unmangled c functions so they work with the bindings

        #ifdef __cplusplus
        extern "C" {
        #endif

        void CameraStartSession(void);
        void CameraStopSession(void);
        int CameraIsRunning(void);
        unsigned char* CameraGetCurrentFrame(int* width, int* height);
        #ifdef __cplusplus
        }
        #endif
    ```

    ```m
        #import <AVFoundation/AVFoundation.h>

        #import "camera.h"

        static AVCaptureSession *session = nil;
        static AVCaptureVideoDataOutput *output = nil;
        static dispatch_queue_t cameraQueue;
        static CVPixelBufferRef latestFrame = nil;

        void CameraStopSession(void) {
            if (session == nil) return;
            
            [session stopRunning];
            session = nil;
            
            if (latestFrame != nil) {
                CVPixelBufferRelease(latestFrame);
                latestFrame = nil;
            }
        }
    ```

# BUILDING THE IOS APP

Replace ./cmd with your directory (use . if in the current directory)
```shell
gomobile build -target=ios -bundleid=com.heytyshawn ./cmd
```

# IF BUILDING WITH FYNE
Disclaimer: apple developer portal provisioning profile is required. still waiting on support to get that sorted
```shell
fyne package --os ios --app-id com.heytyshawn --app-build 1 --app-version 1 --cert "Apple Development: heytyshawn@gmail.com (M378RYRS4B)"   
```

# IF BUILDING WITH GIOUI
You'll have to edit the Info.plist with necessary permisisons and then resign

Convert plist to readible xml (use full path and put in quotations)
```shell
plutil -convert xml1 "/Users/heytyshawn/Library/Mobile Documents/com~apple~CloudDocs/Documents/Old Documents/Code/go/gioui/ioslibrary/ioslibrary.app/Info.plist"  
```
Optional: converting it back
```shell
plutil -convert binary1 "/Users/heytyshawn/Library/Mobile Documents/com~apple~CloudDocs/Documents/Old Documents/Code/go/gioui/ioslibrary/ioslibrary.app/Info.plist"  
```

Clearing any finder files because ios wont let you sign with them
```shell
xattr -cr ioslibrary.app
```

Find cert to sign with
```shell
security find-identity -v -p codesigning
```

# Resigning the app

```shell
# Step 1: Embed provisioning profile
cp PROVISIONINGPROFILE.mobileprovision APPNAME.app/embedded.mobileprovision

# Step 2: Extract entitlements from the provisioning profile
security cms -D -i PROVISIONINGPROFILE.mobileprovision > temp.plist
/usr/libexec/PlistBuddy -x -c 'Print:Entitlements' temp.plist > entitlements.plist

# Step 3: Sign the inner binary (main)
codesign -f -s "Apple Development: ..." --entitlements entitlements.plist APPNAME.app/main

# Step 4: Sign the full app bundle
codesign -f -s "Apple Development: ..." --entitlements entitlements.plist APPNAME.app

# Step 5: Verify the signature
codesign --verify --deep --strict --verbose=2 APPNAME.app
```

# INSTALLING THE IOS APP

```shell
# Step 6: Install to device
ios-deploy --bundle APPNAME.app
```