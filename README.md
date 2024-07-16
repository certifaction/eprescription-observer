# ePrescription events auditor
This is reference implementation of ePrescription events auditor. It is a simple application that fetches consistency
proofs and validates that only new events was added and no event was modified. This is simplified version that does not
have any persistence layer and uses in-memory storage. Because of that, it assumes first consistency proof as valid and 
only validates proofs that are fetch afterward.