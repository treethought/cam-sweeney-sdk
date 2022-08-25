# SDK Design

This document provides a general overview of the goals and design of the SDK

## Goals

The primary goals of the SDK are to provide:

- Simple to use
- Provide interfaces to all API resources
- Support all request options including pagination, sorting, and filtering
- Ergonomic API - use of the SDK should feel natural and intuitive
- Transparent and helpful error information
- Provide data returned by API as ready to use structs, without requiring the user to handle deserialization

## Design

The SDK is modular in nature. The base client provides only the core functionality of making requests to the API, as well as applying any type of query parameters and authorization headers. Additionally, the client is meant to be configurable, allowing users to provide their own http.Client configuration and modifying API keys and API URL.

To provide methods for interacting with the various API resources, the client provides convenient methods to resource specific "clients" which provide their own respective methods for making requests to the resource-specific endpoints.
You can think of these resource-specific "clients" as representing a namespace within the API.

While these namespaces are named "clients", they really utilize the base client for making http requests rather than doing so themselves. Additionally, they utilize the base client's ability to deserialize the API response payloads into the appropriate API structs. This allows the user to receive known types from client methods, so they can perform operations on the data without needing to examine the response and determine how to unmarshal into a usable type.

The SDK is also deigned to provide insight into any errors returned from the API. This includes both HTTP errors for failed requests, as well as API responses that indicate an error associated with a resource. Since the API returns 200 even when something goes wrong, the base client handles this case by inspecting the response to determine if it was unsuccessful. If so, it provides an SDKError, which provides information about the request to indicate if an APIError was returned, for which endpoint an error occurred, as well as exposing the underlying APIError message.
