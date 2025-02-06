# HaveIBeenPwned alternative download tool

This tool is alternative to HIBP's official [PwnedPasswordsDownloader](https://github.com/HaveIBeenPwned/PwnedPasswordsDownloader), 
but does not requires you to install SDK or tool itself and works better and faster.

## Command line options

```
Usage: HIBPDL [options] [file]

Options:
  -n,    --ntlm              Download NTLM hashes instead of SHA1
                               Default: download SHA1
  -o,   --overwrite          Overwrite output file if exists
  -q,   --no-progress        Do not output progress bar.
  -p=N, --parallelism=N      Use N parallel jobs.
                               Default: number of CPUs.
  file                       Output file name or full path.
                               Default: pwnedpasswords.sha1.txt for SHA1
                               Default: pwnedpasswords.ntlm.txt for NTLM
```


## Build

On Windows : 
```powershell
$env:GOARCH="amd64"
$env:GOOS="windows"
$ext=".exe"
go build --ldflags '-extldflags "-static"' -o HIBPDL$($ext) .\cmd\HIBPDL
```
