# Dfinity Oracles

[Dfinity Oracles](https://github.com/hyplabs/dfinity-oracle-framework) is a framework for building [blockchain oracles](https://en.wikipedia.org/wiki/Blockchain_oracle) for the [Internet Computer](https://dfinity.org/).

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Background](#background)
- [Tutorial and Examples](#tutorial-and-examples)
- [Framework Reference](#framework-reference)
  - [`oracle.Bootstrap()`](#oraclebootstrap)
  - [`oracle.Run()`](#oraclerun)
- [The Oracle Update Lifecycle](#the-oracle-update-lifecycle)
  - [Acquiring data](#acquiring-data)
  - [Summarizing data](#summarizing-data)
- [Updating the canister](#updating-the-canister)
- [Testing the Framework](#testing-the-framework)
- [Oracle Revocation](#oracle-revocation)

## Background

Applications on the Internet Computer often require information that comes from outside sources. For example, an application might require the current temperature of various cities in the world, or it might require the current prices of various stocks and cryptocurrencies.

Oracles are useful here in order to move information from traditional APIs to software running on the Internet Computer.

## Tutorial and Examples

- [Tutorial](docs/tutorial.md) - step-by-step setup of a sample oracle for retrieving the temperature in different cities.
- [Dfinity Weather Oracle Example](https://github.com/hyplabs/dfinity-oracle-weather) - a complete working example oracle for retrieving various weather conditions in 20 different cities.
- [Dfinity Crypto Oracle Example](https://github.com/hyplabs/dfinity-oracle-crypto) - a complete working example oracle for retrieving current ETH and BTC prices.

## Framework Reference

Now that we have an example app working with the library, let's take a deeper look at what the DFINITY Oracle Framework is doing behind the scenes.

### `oracle.Bootstrap()`

The oracle framework bootstraps the oracle canister by calling into `dfx` for the following setup tasks:

- Creating the oracle canister DFX project if it doesn't exist.
- Starting the DFX network locally in the background.
- Creating the oracle writer identity if it doesn't exist.
- Creating the oracle canister if it doesn't exist.
- Building and installing the canister (if the canister already exists, this performs an upgrade instead).
- Starting the oracle canister if it isn't already running.
- Claiming the oracle owner role if not already claimed, allowing it to manage canister roles.
- Assigning the oracle writer identity the writer role, allowing it to update canister data.

### `oracle.Run()`

The oracle framework starts the oracle service and periodically updates the mappings in the canister. Once this service is running, it will update the canister at the configured time interval.

## The Oracle Update Lifecycle

The oracle framework performs updates at the interval configured via `UpdateInterval`. This update consists of several steps, which are described below.

### Acquiring data

Through the information provided in `config`, we have a list of API endpoints (`Endpoints`), and for each endpoint, its URL (`Endpoint`), a collection of field names and their corresponding JSONPaths within API responses (`JSONPaths`), and an optional normalization function (`NormalizeFunc`).

For more details about JSONPath syntax, see this [JSONPath reference](https://restfulapi.net/json-jsonpath).

When the oracle framework performs an update, it will first make `GET` requests to every endpoint in `config`, resulting in one JSON response per endpoint. The framework then extracts the desired fields from these responses by each field's JSONPath expression, resulting in a `map[string]interface{}` (a mapping from strings to anything).

In our example, we make requests to both WeatherAPI and WeatherBit, and extract the temperature, pressure, and humidity. We have two different configurations of these endpoints - one for Tokyo, and one for Delhi.

This value is then passed to `NormalizeFunc`, which is responsible for turning it into a `map[string]float64`. If no `NormalizeFunc` is specified, then a default one will be used - every field's value will simply be casted as a `float64`.

### Summarizing data

Oracles generally acquire redundant data from many independent sources, then combine them into one trustworthy value.

From the previous step, every endpoint resulted in a `map[string]float64` result, and each of these results is supposed to represent a redundant copy of the same piece of data (e.g., the current temperature in Tokyo). We now need to turn those into a single value.

This is the job of `SummarizeFunc`, which accepts a collection of `map[string]float64` values, and is responsible for converting that into a single `map[string]float64` value. In our example, we've created a custom summarization function called `summarizeWeatherData`. There are two particularly useful oracle framework utility functions that it calls:

- `summary.GroupByKey([]map[string]float64) map[string][]float64`: takes a list of mappings, and converts it into a mapping where each key contains a list of all values in those mappings under those keys.
- `summary.MeanWithoutOutliers([]float64) float64`: takes a list of numbers, removes outliers (outliers are values that are more than 2 standard deviations from the median, so a 95% confidence interval), and takes the mean of the remaining values. This is more stable than the median while still rejecting rogue values.

### Updating the canister

As part of the bootstrap step, the oracle framework created a `writer` identity for the oracle to use - an identity that is allowed to write new values to the mappings stored in the canister.

After the previous step, we now have a `map[string]float64` for Tokyo, and a `map[string]float64` for Delhi. We've written some simple Candid IDL serialization functions in Go that then turn this into a string suitable for passing into DFX.

The final step is then to write this serialized string to the canister using the `writer` identity.

## Testing the Framework

Although the `writer` identity is the only one capable of changing the values inside the oracle canister, reading those values can be done by anyone. To test it out, enter the following command:

```bash
cd weather_oracle
```

```bash
dfx canister call weather_oracle get_map_field_value ("London", "temperature_celsius")
```

This would return the currently stored temperature in Celsius of London within the canister.

Similarly,

```bash
dfx canister call weather_oracle get_map_field_value ("Tokyo", "humidity_pct")
```

would return the currently stored humidity percentage of Tokyo within the canister.

## Oracle Revocation

The framework assigns the `owner` identity to the user who deploys the canister, and this cannot be edited once deployed.

The `owner` identity is capable of assigning the `writer` role to other identities, and also of permanently disabling the oracle if they believe it has been compromised.

If the `writer` identity is compromised, but not the `owner` identity, the `owner` identity can simply assign a new `writer`, and users can continue to use the oracle as they did before.

If the `owner` identity is compromised (e.g., private key leaked), they cannot take control away from the original `owner`, because the `owner` field is not editable. However, the original owner can permanently disable the oracle, preventing attackers from using it for their own gain. This disincentivizes attackers from attempting to take control of the oracle in the first place.

As an oracle owner, we can permanently disable the oracle using the following command:

```bash
cd weather_oracle
dfx canister call weather_oracle self_destruct
```

After this, if you try to get the value through the command given in the Testing section, you will receive an error.
