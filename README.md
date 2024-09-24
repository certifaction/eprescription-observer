# E-Prescription observer
This is a reference implementation of an E-Prescription observer. It is a simple application that fetches consistency
proofs and validates that only new events were added and no event was modified. This is a simplified version that does not
have any persistence layer and uses in-memory storage. Because of that, it assumes the first consistency proof as valid and 
only validates proofs that are fetch afterward.
# How it works
This tool fetches and verifies consistency proofs from the E-Prescription API. Consistency proofs ensure:
* Later versions contain all events from previous versions
* No events were modified
* All new events are appended after old events
# How to run
### build:
```shell
make build
```

# Run 
```shell
./eprescription-observer
```
