# Contribute
## Clone
```
git clone git@github.com:ybonjour/atr.git
```

## Build
atr can be built from source using the `./scripts/build.sh` script:
```
./scripts/build.sh
```
This will put a binary for your local architecture and operating system under `build/bin/localOS_localARCH/atr`.


### Build for a specific architecture or operating system
The architecture and the operating system can be customized using the `GOOS` and `GOARCH` environment variables:
```
GOARCH=amd64 GOOS=linux ./scripts/build.sh
```
This will put a binary for `amd64` on `linux` under `build/bin/linux_amd64/atr`.
