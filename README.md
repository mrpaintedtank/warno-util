# warno-util
 
## Install 

Download the appropriate version for your system by going to [the releases page](https://github.com/mrpaintedtank/warno-util/releases). Unzip and drop it on the desktop. 

### Shortcut Config

You can create a shortcut to the exe and add the following to the target field to run it without needing to use the CLI. Create a shortcut to it, then open the shortcut and replace the target field with the following

```batch
C:\Windows\System32\cmd.exe /k "WHEREVER\YOU\PUT\IT\warno-util.exe switch"
````

Be sure to change WHEREVER\YOU\PUT\IT to the path to the warno-util.exe

### Non Default Install Location For Warno

If you've installed the game in a random location never fear! You can create a config.yml file in the same directory you extracted the binary to and it will pull from that whenever you run the switcher

Update this with the correct values for you and then you're good to go.
```yaml
switch:
  steamAppsPath: "C:\\Program Files (x86)\\Steam\\steamapps"
  steamUserDataPath: "C:\\Program Files (x86)\\Steam\\userdata"
  steamExecutablePath: "C:\\Program Files (x86)\\Steam\\steam.exe"
  savedGamesPath: "C:\\Users\\<username>\\Saved Games\EugenSystems\WARNO"
```

### Using it with CMD

You can quickly run it with

```batch
warno-util.exe switch
```

There is also help and update functionality. The updater will grab the latest version from the release page of the repo based on your os/arch.

## Development

There is a makefile that handles the usual build/lint/etc tasks. PRs are welcome and will be reviewed eventually. 