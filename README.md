# EVAA GO SDK

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/evaafi/evaa-go-sdk)

## Overview

The EVAA GO SDK is designed to easily integrate with the EVAA lending protocol on TON blockchain.

## Table of contents

* [Installing](#installing)
* [Packages](#packages)
  * [Config](#config)
  * [Asset](#asset)
  * [Price](#price)
  * [Principal](#principal)
  * [Transaction](#transaction)

### Installing

```bash
go get github.com/evaafi/evaa-go-sdk@latest
```

### Packages

#### Config

The [config](/config) package is an instruction to interacting with different pools such as Main, LP and Testnet.

#### Asset

The [asset](/asset) package is a tool for obtaining information about the data and configuration of the assets used in the selected version of the protocol.

#### Price

The [price](/price) package is a tool to receive and package prices obtained from oracles used in a selected pool.

#### Principal

The [principal](/principal) package is a guide for working with the user's principals, such as calculation and prediction of the health factor, checking of liquidatablity.

#### Transaction

The [transaction](/transaction) package is an assistant to working with different types of protocol operations such as supply, withdrawal and liquidation.
