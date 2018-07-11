# Account pool

`acctpool` is a web service for temporarily borrowing IBM Cloud accounts. It
works using an authenticated reservation system, where consumers can
authenticate themselves with `acctpool`, make a reservation request for a
particular type of cloud account, and a unique set of credentials will be
returned when a matching account is available. All reservations are
time-limited (but should be explicitly canceled ASAP for efficiency), and any
resources in use after the reservation expires will be forcefully cleaned up.
