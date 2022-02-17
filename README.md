# Sonarcon

Uses the SonarQube API to interact and extract sources from public instances.

## List project

```sh
$ ./sonarcon http://127.0.0.1:9000 lp
app_lhg_key
AppsApi
AppsClient
CommonApi
Core
GoldenTulipApi
GoldenTulipSite
```

## List project's files

```sh
$ ./sonarcon http://127.0.0.1:9000 ls app_lhg_key
app_lhg_key:src/views/AboutUsView/AboutUs.js
app_lhg_key:src/views/AboutUsView/__tests__/AboutUs.test.js
app_lhg_key:src/views/AboutUsView/styles/aboutUsStyle.js
app_lhg_key:src/controllers/tracking/objects/AddToCartObject.js
...
```

## Dump project's files

```sh
$ ./sonarcon http://127.0.0.1:9000 dump app_lhg_key ./app_directory
2022/02/17 02:44:36 Downloading to app_directory/src/views/ResultSearchView/components/SearchNoResult/styles/searchNoResultStyle.js
2022/02/17 02:44:36 Downloading to app_directory/src/controllers/search/searchReducer.js
2022/02/17 02:44:36 Downloading to app_directory/src/controllers/search/__tests__/searchsActions.test.js
2022/02/17 02:44:36 Downloading to app_directory/src/controllers/search/searchSagas.js
2022/02/17 02:44:36 Downloading to app_directory/src/controllers/search/searchSelectors.js
2022/02/17 02:44:36 Downloading to app_directory/src/services/SearchServices.js
2022/02/17 02:44:36 Downloading to app_directory/src/services/__tests__/SearchServices.test.js
2022/02/17 02:44:36 Downloading to app_directory/src/controllers/search/__tests__/searchsReducer.test.js
...
```

