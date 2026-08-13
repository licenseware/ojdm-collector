# Changelog

## [1.1.1](https://github.com/licenseware/ojdm-collector/compare/v1.1.0...v1.1.1) (2026-08-13)


### Features

* **CI:** add macOS acceptance job and update readme ([#28](https://github.com/licenseware/ojdm-collector/issues/28)) ([c188880](https://github.com/licenseware/ojdm-collector/commit/c18888030f7bbb34797700b1e0f3640b12a4e5a3))
* timestamp report and logs ([#30](https://github.com/licenseware/ojdm-collector/issues/30)) ([9e70a90](https://github.com/licenseware/ojdm-collector/commit/9e70a90c6b267ccf9ed8c951af526396e69d5fe5))


### Bug Fixes

* **docs:** remove bash sign from snippet ([5109eb3](https://github.com/licenseware/ojdm-collector/commit/5109eb38f0806abe3b23b2cf7906bc29d59638ce))


### Build & Dependencies

* **deps:** update actions/attest-build-provenance digest to 4d10147 ([#27](https://github.com/licenseware/ojdm-collector/issues/27)) ([a89ad8b](https://github.com/licenseware/ojdm-collector/commit/a89ad8b07d4297cdd8fb3237f1d0c7f8f6b4b35c))


### Chores

* restructure packages into internal ([#32](https://github.com/licenseware/ojdm-collector/issues/32)) ([b518b40](https://github.com/licenseware/ojdm-collector/commit/b518b40242489977f0bed8d03e55a829e32c5404))

## [1.1.0](https://github.com/licenseware/ojdm-collector/compare/v1.0.1...v1.1.0) (2026-08-11)


### Features

* **CI:** add release-please with complete goreleaser setup ([#25](https://github.com/licenseware/ojdm-collector/issues/25)) ([dd16d35](https://github.com/licenseware/ojdm-collector/commit/dd16d35fe428a1689145986c3c9e62055f288b65))
* **logging:** write a structured debug log beside the report ([#24](https://github.com/licenseware/ojdm-collector/issues/24)) ([5fa9906](https://github.com/licenseware/ojdm-collector/commit/5fa99061906c171a402b422b8542697fc16eb08b))
* **search:** specify search paths through a file and update docs ([#19](https://github.com/licenseware/ojdm-collector/issues/19)) ([1312dab](https://github.com/licenseware/ojdm-collector/commit/1312dab354fad6c37715b3034abc4e5e87d50440))


### Bug Fixes

* search for java binaries as well as dlls ([7102519](https://github.com/licenseware/ojdm-collector/commit/71025193ea7ffb1de5fc1639cb730eeec0851e5f))
* search for java binaries as well as dlls ([95eb67c](https://github.com/licenseware/ojdm-collector/commit/95eb67cbd190340237aa7606529a021eacdca41f))
* **search:** skip unreadable paths, locate jvm.dll on windows ([#23](https://github.com/licenseware/ojdm-collector/issues/23)) ([8b7b774](https://github.com/licenseware/ojdm-collector/commit/8b7b77486862627353d98a2c06b5b0e9259dc0fe))


### Chores

* **deps:** update actions/setup-go action to v5 ([#11](https://github.com/licenseware/ojdm-collector/issues/11)) ([62486a2](https://github.com/licenseware/ojdm-collector/commit/62486a2d9310490bc11cdf0d567a1f43242d7f44))
* **deps:** update goreleaser/goreleaser-action action to v6 ([#12](https://github.com/licenseware/ojdm-collector/issues/12)) ([b7809bc](https://github.com/licenseware/ojdm-collector/commit/b7809bc42878877b32cbf6868d2be823b83c1e38))
