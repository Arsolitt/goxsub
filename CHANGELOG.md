# Changelog

## [1.1.0](https://github.com/Arsolitt/goxsub/compare/v1.0.0...v1.1.0) (2026-03-31)


### Features

* add SetRemarks method to VLESSProxy ([44f535c](https://github.com/Arsolitt/goxsub/commit/44f535c3940c9d956284ca9150a3aeb9a5086688))
* add sing-box format and --keep-remark, --singbox-* flags to CLI ([3f8851c](https://github.com/Arsolitt/goxsub/commit/3f8851c4f62c0662da4176a469ec39cd1622f194))
* add sing-box outbound formatter ([6561cff](https://github.com/Arsolitt/goxsub/commit/6561cff5a04c6d7c3d2e72082d3369d85c590e03))
* change Podkop return type from string to []string ([0f2ea26](https://github.com/Arsolitt/goxsub/commit/0f2ea2692f521550481049e9a6a8a80f5fbf0a38))

## 1.0.0 (2026-03-29)


### Features

* add --exclude-by-remark flag for glob-based remark exclusion ([2ff66e0](https://github.com/Arsolitt/goxsub/commit/2ff66e08c9e52bae11346cf53c2748668c3da478))
* add CLI for fetching subscription and printing vless URIs ([7f18f2d](https://github.com/Arsolitt/goxsub/commit/7f18f2d72d6432afb7cec7bdc16e27fd3b92d22c))
* add ExtractVLESSOutbounds with tests ([d5ba031](https://github.com/Arsolitt/goxsub/commit/d5ba0314b11146d8882f3c8b56a96628d1948096))
* add FilterByRemark for glob-based remark exclusion ([dc40095](https://github.com/Arsolitt/goxsub/commit/dc4009505c870afa4e60543298413134ea6ed399))
* add ParseSubscription with tests ([8447f17](https://github.com/Arsolitt/goxsub/commit/8447f17ff87a7be43942a320d3035dfe549a1aed))
* add podkop output format with CLI flags ([c8af81e](https://github.com/Arsolitt/goxsub/commit/c8af81e1e950b646cede20a1d1b97d0bf0588b1f))
* add ToVLESSURI with tests for all transport types ([be32724](https://github.com/Arsolitt/goxsub/commit/be327249fa69abbf8083bb69f6df6f52a7559042))
* add xray subscription type definitions ([442d0b3](https://github.com/Arsolitt/goxsub/commit/442d0b34ec6a1facdffa87cc135d59f843bd53ed))


### Bug Fixes

* URL-encode all URI parameters and fragment per spec ([d254619](https://github.com/Arsolitt/goxsub/commit/d25461945d88c923be9f473823a83cedf5b9434c))
