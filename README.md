# Account pool

[![Build Status](https://travis-ci.org/dolph/ibm-cloud-account-pool.svg?branch=master)](https://travis-ci.org/dolph/ibm-cloud-account-pool)

`acctpool` is a web service for temporarily borrowing IBM Cloud accounts. It
works using an authenticated reservation system, where consumers can
authenticate themselves with `acctpool`, make a reservation request for a
particular type of cloud account, and a unique set of credentials will be
returned when a matching account is available. All reservations are
time-limited (but should be explicitly canceled ASAP for efficiency), and any
resources in use after the reservation expires will be forcefully cleaned up.

## Background

IBM Cloud accounts carry strict quota limits, which makes automated testing
difficult (you need multiple accounts to test more than one deployment at a
time, for example). Generally, this means that test systems either need to be
limited to a single account (and therefore limited to running one test at a
time), or rotate through an array of predefined test accounts.

This project solves that limitation by allowing multiple test systems to share
a pool of single-use accounts. The solution is inspired by [OpenStack's
nodepool](https://zuul-ci.org/docs/nodepool/).

## Design

The system does not currently use an external datastore, and is designed to be
run as a single process. It does not currently support high availability.
Credentials are loaded into the service at startup. Upon startup, accounts are
verified and forcefully cleaned before making them available for reservation.

## Workflow

1. Obtain an API token for the service.

2. Use the API token to create a reservation request for the type of cloud
   account you need.

3. Poll the reservation endpoint until credentials are available. An expiration
   time for the reservation will be provided, as well.

4. Use the credentials to consume IBM Cloud resources.

5. Cancel the reservation when you're done with the cloud account. If the
   reservation is not explicitly cancelled, it will be forcefully terminated at
   the expiration time specified along with the credentials.

## API

### `GET /`

Returns statistics about and status of the system. This request does not
require authentication.

### `POST /reservations`

Successful requests returns a `HTTP 302 Found` redirect to a unique reservation
URL, which can be polled until an account is available.

Query string parameters:

- `token` (required): authentication token used to identify the consumer.

- `duration` (optional): duration of the reservation request, in minutes.

- `type` (optional): the type of account being requested (defaults to `free`),
  but `lite` is also supported.

### `GET /reservations/{reservation_id}`

Returns the details of a reservation, include account credentials for the
reservation (when available). While the reservation is pending, there will not
be an expiration for the reservation, and no credentials will be provided.

Query string parameters:

- `token` (required): authentication token used to identify the consumer. The
  token must be the same as the token used to create the reservation.

### `DELETE /reservations/{reservation_id}`

Query string parameters:

- `token` (required): authentication token used to identify the consumer. The
  token must be the same as the token used to create the reservation.

## Data model

Tokens:

- ID

Accounts:

- ID
- Credentials
- Account Type
- Expiration
- Dirty

Requests:

- Token ID
- Account Type
- Duration
- Expiration

Reservations:

- Token ID
- Account ID
- Duration
- Expiration

## Recurring jobs

- Clean up accounts after reservations expire.
- Verify credentials.
- Rotate credentials.
