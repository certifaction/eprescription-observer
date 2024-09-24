# E-Prescription events auditor
This is a reference implementation of an E-Prescription events auditor. It is a simple application that fetches consistency
proofs and validates that only new events were added and no event was modified. This is a simplified version that does not
have any persistence layer and uses in-memory storage. Because of that, it assumes the first consistency proof as valid and 
only validates proofs that are fetch afterward.
# How to run
### build:
```shell
make build
```

# Run 
```shell
./eprescription-observer
```
