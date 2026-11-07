# Sample diff run

Looks like a lot, but you can quickly identify the pieces that changed -- some
growth in `/Users/oklischat/Library/Caches/` and a new directory
`/Users/oklischat/smith/glue-cli/target/` contributing to a +12.3 GB change in
my home directory; a removed macOS software update accounting for most of the
-3.2 GB in `/System/`; removed temporary copies of the Chrome app bundle
(`com.google.Chrome.code_sign_clone`) accounting for most of the -6.9 GB in
`/private/`; plus some smaller changes, all adding up to a +2.1 GB change in
scanned files. You can also reduce the output by adding an argument like
`--depth 2`, which will print the same changes, but summed up for the topmost
two directory levels.

Chrome most likely creates the app bundle copies as APFS clones, which share
their storage with the original app and take up almost no extra disk space.
GrandPerspective counts each copy at its full size, though, so deleting them
probably freed much less space than the -5.8 GB shown (and the real change
outside the scan is correspondingly bigger than the -29.9 GB below).

The scanned change is reversed by a -29.9 GB change in "Not in scan" data --
most likely due to a reboot between the two scans, which frees swap space,
files that were deleted but still held open, and leftover update data. On
macOS, swap and update data are stored in separate volumes of the same APFS
container, which the scan doesn't include; you can get an overview of them
with `diskutil apfs list`.

Each change also shows up in all of its parent folders, so read the list from
the top and look for where a chain of equal numbers ends -- that's where the
change actually happened.

## Default output

```
#> gpscanutil diff 2026-10-31-2.gpscan 2026-11-07.gpscan
read 2026-10-31-2.gpscan (9,442,123 items) in 5.9s
read 2026-11-07.gpscan (9,762,364 items) in 6.1s
Scan 1:  2026-10-31-2.gpscan  (2026-10-31T21:35:41Z, physical sizes)
Scan 2:  2026-11-07.gpscan  (2026-11-07T14:17:36Z, physical sizes)

                       scan 1           scan 2           change
Free space            34.2 GB          62.0 GB         +27.8 GB
Used                 460.1 GB         432.3 GB         -27.8 GB
Scanned              368.5 GB         370.7 GB          +2.1 GB
Not in scan           91.6 GB          61.7 GB         -29.9 GB
Files               7,824,782        8,132,958         +308,176
Folders             1,617,341        1,629,406          +12,065

    +12.3 GB     279.5 GB → 291.8 GB     /Users/
    +12.3 GB     278.2 GB → 290.4 GB     /Users/oklischat/
     +7.1 GB      12.8 GB → 19.9 GB      /Users/oklischat/Library/Caches/
     -6.9 GB      14.6 GB → 7.7 GB       /private/
     +6.5 GB      96.6 GB → 103.1 GB     /Users/oklischat/Library/
     -6.4 GB      14.0 GB → 7.7 GB       /private/var/
     -5.8 GB       8.1 GB → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/
     -5.8 GB       8.1 GB → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/
     -5.8 GB       8.6 GB → 2.8 GB       /private/var/folders/
     -5.8 GB       8.5 GB → 2.7 GB       /private/var/folders/9b/
     -5.8 GB       8.5 GB → 2.7 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/
     +4.6 GB       1.9 GB → 6.5 GB       /Users/oklischat/Library/Caches/go-build/
     +3.9 GB       2.1 GB → 6.0 GB       /Users/oklischat/Library/Caches/Google/
     +3.9 GB       2.1 GB → 6.0 GB       /Users/oklischat/Library/Caches/Google/Chrome/
     -3.2 GB      21.1 GB → 17.9 GB      /System/
     -3.2 GB      20.0 GB → 16.8 GB      /System/Library/
     -3.2 GB       9.5 GB → 6.4 GB       /System/Library/AssetsV2/
     -3.1 GB       3.1 GB → (gone)       /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/f57b47f170b54c43687df44e34759b47a6c38b4f.asset/
     -3.1 GB       3.1 GB → (gone)       /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/f57b47f170b54c43687df44e34759b47a6c38b4f.asset/AssetData/
     -3.1 GB       3.1 GB → 213.0 kB     /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/
     +3.0 GB      32.5 GB → 35.4 GB      /Users/oklischat/smith/
     +2.8 GB       9.5 MB → 2.8 GB       /Users/oklischat/smith/glue-cli/
     +2.8 GB        (new) → 2.8 GB       /Users/oklischat/smith/glue-cli/target/
     -2.6 GB       2.6 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/
     -2.6 GB       2.6 GB → 3.1 MB       /Users/oklischat/Library/Caches/com.openai.codex/
     -2.6 GB       2.6 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/
     -2.6 GB       2.6 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/
     -2.4 GB       2.4 GB → (gone)       /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/f57b47f170b54c43687df44e34759b47a6c38b4f.asset/AssetData/boot/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/Frameworks/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/
     +2.3 GB        (new) → 2.3 GB       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/Frameworks/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/Google Chrome.app.bundle/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/Google Chrome.app.bundle/Contents/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/Google Chrome.app.bundle/Contents/Frameworks/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.wshMMs/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/
     +2.2 GB        (new) → 2.2 GB       /Users/oklischat/smith/glue-cli/target/debug/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/Google Chrome.app.bundle/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/Google Chrome.app.bundle/Contents/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/Google Chrome.app.bundle/Contents/Frameworks/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/
     -2.2 GB       2.2 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.cug8mt/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/
     +2.1 GB     368.5 GB → 370.7 GB     /
     -1.8 GB       1.8 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/Ppqs9nvfW/
     -1.8 GB       1.8 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/Ppqs9nvfW/ChatGPT.app/
     -1.8 GB       1.8 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/Ppqs9nvfW/ChatGPT.app/Contents/
     -1.7 GB       1.7 GB → (gone)       /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/f57b47f170b54c43687df44e34759b47a6c38b4f.asset/AssetData/boot/Firmware/
     +1.5 GB      53.2 kB → 1.5 GB       /Users/oklischat/Library/Caches/Google/Chrome/Profile 2/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/Google Chrome.app.bundle/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/Google Chrome.app.bundle/Contents/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/Google Chrome.app.bundle/Contents/Frameworks/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/
     -1.5 GB       1.5 GB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.Em0Xtp/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/
     -1.4 GB       1.4 GB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/Ppqs9nvfW/ChatGPT.app/Contents/Resources/
     +1.4 GB       7.7 MB → 1.4 GB       /Users/oklischat/Library/Caches/Google/Chrome/Default/
     +1.3 GB      16.4 kB → 1.3 GB       /Users/oklischat/Library/Caches/Google/Chrome/Profile 2/Cache/
     +1.3 GB       8.2 kB → 1.3 GB       /Users/oklischat/Library/Caches/Google/Chrome/Profile 2/Cache/Cache_Data/
     +1.2 GB        (new) → 1.2 GB       /Users/oklischat/smith/glue-cli/target/debug/deps/
     +1.1 GB       6.3 MB → 1.1 GB       /Users/oklischat/Library/Caches/Google/Chrome/Default/Cache/
     +1.1 GB       6.3 MB → 1.1 GB       /Users/oklischat/Library/Caches/Google/Chrome/Default/Cache/Cache_Data/
     +1.1 GB      15.1 GB → 16.2 GB      /Users/oklischat/doc/
     +1.1 GB      11.6 GB → 12.6 GB      /Users/oklischat/doc/mydocs/
     +1.1 GB      78.2 MB → 1.2 GB       /Users/oklischat/doc/mydocs/ai/
     +1.1 GB      44.7 MB → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/micrograd/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/micrograd/.venv/
     +1.1 GB      20.5 kB → 1.1 GB       /Users/oklischat/Library/Caches/com.bitwarden.desktop.ShipIt/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/Library/Caches/com.bitwarden.desktop.ShipIt/update.gJbispL/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/Library/Caches/com.bitwarden.desktop.ShipIt/update.gJbispL/Bitwarden.app/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/Library/Caches/com.bitwarden.desktop.ShipIt/update.gJbispL/Bitwarden.app/Contents/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/micrograd/.venv/lib/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/micrograd/.venv/lib/python3.13/
     +1.1 GB        (new) → 1.1 GB       /Users/oklischat/doc/mydocs/ai/basics/karpathy/micrograd/.venv/lib/python3.13/site-packages/
   +960.3 MB       7.0 GB → 8.0 GB       /Users/oklischat/.cache/uv/
   +958.0 MB      11.2 GB → 12.2 GB      /Users/oklischat/.cache/
   +956.9 MB       6.7 GB → 7.7 GB       /Users/oklischat/.cache/uv/archive-v0/
   +811.8 MB     667.6 kB → 812.5 MB     /Users/oklischat/Library/Caches/Google/Chrome/Profile 11/
   +788.9 MB        (new) → 788.9 MB     /Users/oklischat/smith/glue-cli/target/debug/incremental/
   -766.0 MB     766.0 MB → (gone)       /Users/oklischat/Library/Application Support/Google/GoogleUpdater/crx_cache/e8412d0225e49f4b40fa318ee56d3e9b26cfec22d482349b560fcf37cc10e6aa
   +766.0 MB        (new) → 766.0 MB     /Users/oklischat/Library/Application Support/Google/GoogleUpdater/crx_cache/444a3d2281f6ac5fbcac55e06a6478e976f515c5470ee31d40027ee158bd9791
   +759.9 MB        (new) → 759.9 MB     /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/156.0.8078.81/
   +759.9 MB        (new) → 759.9 MB     /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.3NTeWr/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/156.0.8078.80/
   +759.9 MB        (new) → 759.9 MB     /Applications/Google Chrome.app/Contents/Frameworks/Google Chrome Framework.framework/Versions/156.0.8078.81/
   -759.5 MB     759.5 MB → (gone)       /Applications/Google Chrome.app/Contents/Frameworks/Google Chrome Framework.framework/Versions/156.0.8078.74/
   -748.0 MB     748.0 MB → (gone)       /Users/oklischat/Library/Caches/com.openai.codex/org.sparkle-project.Sparkle/Installation/EBUdwFhg7/EDB3A004-27AE-415D-B3DC-702A3CC7EF33.1023.50603.zip
   -746.6 MB     746.6 MB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/153.0.8010.53/
   -746.6 MB     746.6 MB → (gone)       /private/var/folders/9b/wykd1_9n5b53vzjpp76mhws00000gs/X/com.google.Chrome.code_sign_clone/code_sign_clone.DYCjCF/Google Chrome.app.bundle/Contents/Frameworks/Google Chrome Framework.framework/Versions/153.0.8010.52/
```

## Tree output

```
#> gpscanutil diff 2026-10-31-2.gpscan 2026-11-07.gpscan --tree
read 2026-10-31-2.gpscan (9,442,123 items) in 6.8s
read 2026-11-07.gpscan (9,762,364 items) in 6.9s
Scan 1:  2026-10-31-2.gpscan  (2026-10-31T21:35:41Z, physical sizes)
Scan 2:  2026-11-07.gpscan  (2026-11-07T14:17:36Z, physical sizes)

                       scan 1           scan 2           change
Free space            34.2 GB          62.0 GB         +27.8 GB
Used                 460.1 GB         432.3 GB         -27.8 GB
Scanned              368.5 GB         370.7 GB          +2.1 GB
Not in scan           91.6 GB          61.7 GB         -29.9 GB
Files               7,824,782        8,132,958         +308,176
Folders             1,617,341        1,629,406          +12,065

     +2.1 GB     368.5 GB → 370.7 GB     /
    +12.3 GB     279.5 GB → 291.8 GB     ├── Users/
    +12.3 GB     278.2 GB → 290.4 GB     │   └── oklischat/
     +6.5 GB      96.6 GB → 103.1 GB     │       ├── Library/
     +3.0 GB      32.5 GB → 35.4 GB      │       ├── smith/
     +1.1 GB      15.1 GB → 16.2 GB      │       ├── doc/
   +958.0 MB      11.2 GB → 12.2 GB      │       ├── .cache/
   +659.7 MB      45.8 GB → 46.5 GB      │       ├── go/
   +130.2 MB        (new) → 130.2 MB     │       ├── screenrec.mpg
   -103.2 MB       1.4 GB → 1.3 GB       │       ├── .vscode/
    +47.2 MB       1.5 GB → 1.6 GB       │       ├── Downloads/
    +31.5 MB        (new) → 31.5 MB      │       ├── .Trash/
    -19.0 MB     673.0 MB → 654.0 MB     │       ├── .codex/
     +1.1 MB                             │       └── … 12 more changed (net)
     -6.9 GB      14.6 GB → 7.7 GB       ├── private/
     -6.4 GB      14.0 GB → 7.7 GB       │   ├── var/
     -5.8 GB       8.6 GB → 2.8 GB       │   │   ├── folders/
   -452.4 MB     529.4 MB → 77.0 MB      │   │   ├── log/
    -86.2 MB       2.7 GB → 2.6 GB       │   │   ├── db/
     -5.7 MB      51.7 MB → 46.0 MB      │   │   ├── protected/
     -1.7 MB       1.7 MB → (gone)       │   │   ├── FileWave/
   +131.1 kB       7.9 MB → 8.1 MB       │   │   ├── logs/
    -24.6 kB       1.2 MB → 1.2 MB       │   │   └── run/
   -527.1 MB     527.1 MB → 8.2 kB       │   └── tmp/
   -525.0 MB     525.0 MB → (gone)       │       ├── glue-head/
   -528.4 kB     528.4 kB → (gone)       │       ├── api1
   -524.3 kB     524.3 kB → (gone)       │       ├── api2-stage
   -524.3 kB     524.3 kB → (gone)       │       ├── api2-stage-2
   -471.0 kB     471.0 kB → (gone)       │       ├── api2
     -4.1 kB       4.1 kB → (gone)       │       ├── new-failures.txt
     -4.1 kB       4.1 kB → (gone)       │       ├── head-failures.txt
     -4.1 kB       4.1 kB → (gone)       │       ├── old_apismod.rs
     -4.1 kB       4.1 kB → (gone)       │       └── nix-89069-2234953807/
     -3.2 GB      21.1 GB → 17.9 GB      ├── System/
     -3.2 GB      20.0 GB → 16.8 GB      │   └── Library/
     -3.2 GB       9.5 GB → 6.4 GB       │       └── AssetsV2/
   -165.7 MB       1.2 GB → 1.1 GB       ├── usr/
   -165.7 MB     422.1 MB → 256.3 MB     │   └── local/
   -165.7 MB     165.7 MB → (gone)       │       ├── sbin/
    -36.9 kB     215.6 MB → 215.6 MB     │       ├── bin/
    -12.3 kB      16.4 kB → 4.1 kB       │       ├── share/
     -4.1 kB      12.3 kB → 8.2 kB       │       └── etc/
   +130.7 MB      24.9 GB → 25.0 GB      ├── Applications/
    +99.1 MB       1.8 GB → 1.9 GB       │   ├── ChatGPT.app/
    +99.1 MB       1.8 GB → 1.9 GB       │   │   └── Contents/
    +24.1 MB     985.6 MB → 1.0 GB       │   ├── Visual Studio Code.app/
    +24.1 MB     985.6 MB → 1.0 GB       │   │   └── Contents/
     +7.5 MB     202.9 MB → 210.4 MB     │   ├── iTerm.app/
     +7.5 MB     202.9 MB → 210.4 MB     │   │   └── Contents/
     +8.2 kB       2.3 GB → 2.3 GB       │   └── Google Chrome.app/
     +8.2 kB       2.3 GB → 2.3 GB       │       └── Contents/
    -37.8 MB       7.7 GB → 7.6 GB       └── Library/
    -51.4 MB      53.3 MB → 1.9 MB           ├── Frameworks/
    -51.4 MB      51.4 MB → (gone)           │   └── FileWaveCore.framework/
    +18.2 MB      31.2 MB → 49.5 MB          ├── Logs/
    +20.3 MB      14.6 MB → 34.9 MB          │   ├── CrashReporter/
     -2.3 MB       9.4 MB → 7.1 MB           │   ├── DiagnosticReports/
   +249.9 kB       7.2 MB → 7.5 MB           │   └── MCXTools.log
     -4.6 MB       6.6 MB → 2.1 MB           ├── Bluetooth/
     -2.3 MB       3.1 MB → 856.1 kB         │   ├── com.apple.MobileBluetooth.ledevices.other.db-wal
     -2.3 MB       3.1 MB → 856.1 kB         │   └── com.apple.MobileBluetooth.ledevices.paired.db-wal
    +45.1 kB      68.0 MB → 68.1 MB          ├── Apple/
    +45.1 kB      65.0 MB → 65.1 MB          │   └── System/
    -16.4 kB      35.6 MB → 35.6 MB          ├── Trial/
    -16.4 kB      41.0 kB → 24.6 kB          │   └── NamespaceDescriptors/
    -12.3 kB     593.9 kB → 581.6 kB         ├── Preferences/
     +8.2 kB     102.4 kB → 110.6 kB         │   ├── com.apple.networkextension.plist
     -8.2 kB     135.2 kB → 127.0 kB         │   ├── SystemConfiguration/
     -4.1 kB       4.1 kB → (gone)           │   ├── com.filewave.plist
     -4.1 kB       4.1 kB → (gone)           │   ├── com.filewave.Client.plist
     -4.1 kB      16.4 kB → 12.3 kB          │   └── Audio/
     -8.2 kB       2.3 GB → 2.3 GB           ├── Application Support/
    -16.4 kB     134.1 MB → 134.0 MB         │   ├── Tunnelblick/
     +8.2 kB     106.8 MB → 106.8 MB         │   └── Fortinet/
     -4.1 kB      61.4 kB → 57.3 kB          ├── LaunchDaemons/
     -4.1 kB       4.1 kB → (gone)           │   └── com.filewave.fwcld.plist
     -4.1 kB      41.0 kB → 36.9 kB          ├── LaunchAgents/
     -4.1 kB       4.1 kB → (gone)           │   └── com.filewave.fwGUI.plist
     -4.1 kB       3.3 MB → 3.3 MB           └── Keychains/
     -4.1 kB      90.1 kB → 86.0 kB              └── System.keychain
```
