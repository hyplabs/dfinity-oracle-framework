# Dfinity Oracles Tutorial

For this tutorial, we will create a sample oracle for fetching the current temperature in different cities, using Dfinity Oracles. To read more about the framework itself, see the [README](../README.md).

This sample oracle is also available as a complete project, [Dfinity Weather Oracle Example](https://github.com/hyplabs/dfinity-weather-oracle)!

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Step 0: Before you begin](#step-0-before-you-begin)
- [Step 1: Create a new project for the sample oracle](#step-1-create-a-new-project-for-the-sample-oracle)
- [Step 2: Get API endpoints for external data sources](#step-2-get-api-endpoints-for-external-data-sources)
- [Step 3: Configuring the sample oracle](#step-3-configuring-the-sample-oracle)
- [Step 4: Running the sample oracle](#step-4-running-the-sample-oracle)

## Step 0: Before you begin

Before starting this tutorial, double-check the following:

- You have downloaded and installed the [DFINITY Canister SDK](https://sdk.dfinity.org/docs/quickstart/local-quickstart.html#download-and-install).
- You have downloaded and installed the [Go programming language](https://golang.org/).
- You have stopped any Internet Computer network processes running on the local computer.

## Step 1: Create a new project for the sample oracle

Create a new folder named `sample-oracle`, then create a new Go project in this folder:

```bash
go mod init github.com/YOUR_USERNAME/sample-oracle
```

(NOTE: replace `YOUR_USERNAME` with your GitHub username)

Install the framework, `dfinity-oracles`, as a Go Module dependancy:

```bash
go get github.com/hyplabs/dfinity-oracles
```

Note that if you are trying this tutorial before Dfinity Oracles has been publicly released, you'll want to use the following command instead: `GOPRIVATE=github.com/hyplabs/dfinity-oracles GIT_TERMINAL_PROMPT=1 go get github.com/hyplabs/dfinity-oracles`.

## Step 2: Get API endpoints for external data sources

To fetch the current temperature in a city, we need to call various weather monitoring APIs and parse the responses. We will be using the following APIs, but you can choose your own if you prefer:

- [WeatherAPI](https://weatherapi.com) - the `http://api.weatherapi.com/v1/current.json?key=WEATHERAPI_API_KEY&q=CITY` endpoint.
- [Weatherbit](https://weatherbit.io) - the `https://api.weatherbit.io/v2.0/current?city=CITY&country=COUNTRY_CODE&key=WEATHERBIT_API_KEY` endpoint.

Sign up for both of the above services, and obtain a developer API key for each. You will need it in the next step!

## Step 3: Configuring the sample oracle

We can now configure the sample oracle. Create a new file in the `sample-oracle` folder named `main.go`, with the following contents:

```go
package main

import (
	"time"

	framework "github.com/hyplabs/dfinity-oracles"
	"github.com/hyplabs/dfinity-oracles/models"
)

func main() {
	tokyoEndpoints := []models.Endpoint{
		{
			Endpoint: "http://api.weatherapi.com/v1/current.json?key=WEATHERAPI_API_KEY&q=Tokyo,JP",
			JSONPaths: map[string]string{
				"temperature_celsius": "$.current.temp_c",
			},
		},
		{
			Endpoint: "https://api.weatherbit.io/v2.0/current?key=WEATHERBIT_API_KEY&city=Tokyo&country=JP",
			JSONPaths: map[string]string{
				"temperature_celsius": "$.data[0].temp",
			},
		},
	}
	delhiEndpoints := []models.Endpoint{
		{
			Endpoint: "http://api.weatherapi.com/v1/current.json?key=WEATHERAPI_API_KEY&q=Delhi,IN",
			JSONPaths: map[string]string{
				"temperature_celsius": "$.current.temp_c",
			},
		},
		{
			Endpoint: "https://api.weatherbit.io/v2.0/current?key=WEATHERBIT_API_KEY&city=Delhi&country=IN",
			JSONPaths: map[string]string{
				"temperature_celsius": "$.data[0].temp",
			},
		},
	}
	config := models.Config{
		CanisterName:   "sample_oracle",
		UpdateInterval: 5 * time.Minute,
	}
	engine := models.Engine{
		Metadata: []models.MappingMetadata{
			{Key: "Tokyo", Endpoints: tokyoEndpoints},
			{Key: "Delhi", Endpoints: delhiEndpoints},
		},
	}
	oracle := framework.NewOracle(&config, &engine)
	oracle.Bootstrap()
	oracle.Run()
}
```

(NOTE: replace `WEATHERAPI_API_KEY` with your WeatherAPI API key, and `WEATHERBIT_API_KEY` with your Weatherbit API key)

Let's take a look at the key parts of this sample oracle:

- `tokyoEndpoints` and `delhiEndpoints` specify the URLs where the temperature data can be found, as well as [JSONPath](https://www.baeldung.com/guide-to-jayway-jsonpath) expressions that extract just the temperature (in Celsius) out of the JSON response.
    - `Endpoint` is the URL. In this case, we've entered the WeatherAPI and WeatherBit API endpoint URLs here.
    - `JSONPaths` is a map of JSONPath expressions and the relevant info that they retrieve. In this case, the only piece of relevant info is `temperature_celsius`.
- `metadata` specifies the two pieces of data that we care about - the temperature in Tokyo, and the temperature in Delhi.
    - In this example, since the `SummaryFunc` option of `metadata` isn't specified, a default summarization function will be applied, which simply takes the `temperature_celsius` key from the API call results, eliminates outliers (outside 2 standard deviations), and takes the average of the remaining values to obtain the final temperature.
- `config` specifies the `CanisterName` and `UpdateInterval` - what the canister should be called, and how often it should update.
- We create a new oracle struct by calling `NewOracle(&config)`.
- We bootstrap the new oracle by calling `oracle.Bootstrap()`.
- Finally, we start the oracle by calling `oracle.Run()`.

For more details about what `oracle.Bootstrap()` and `oracle.Run()` do, see "Framework Reference" in the README.

## Step 4: Running the sample oracle

Since we've finished configuring the sample oracle, we can now build and run it:

```bash
go build
```

This should create an executable, `sample-oracle` (or `sample-oracle.exe` on Windows). Now let's run the sample oracle:

```bash
./sample-oracle
```

This will generate a new DFX project in the folder, bootstrap the project, and start the oracle service.
