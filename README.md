# CockroachDB Financial Ledger
This is a concept validation for a financial ledger using cockroachdb.

- [X] Balance should be consistent for every transaction
- [X] Balance should not end-up negative
- [X] It should be able to process multiple accounts at the same time

### Notes

The code was executed on a Apple M1 Chip with 8GB of memory (Macbook Pro M1)

The dedicated cluster is hosted on CockRoachDB Cloud (3 node cluster) on AWS with 4 vCPU, 16 GiB RAM, 15 GiB disk, 225 IOPS
The serverless cluster is hosted on CockRoachDB Cloud on AWS

A lot can be improved before something like this goes to production, but this test is enough to validate my idea.
### Time
!! THIS IS NOT A REAL WORLD SCIENTIFIC TEST !!
Please note this was executed using limited resources, it is only a concept for validating the idea.

The dedicated cluster has a better performance than the serverless one (as expected)

Times to process a transfer (updated_at - created_at)

Tests executed with a single account

##### Dedicated
Transfers per second: up to 16 (per account)

Min: 0.418274 second (418ms)
Max: 1.480113 second (1480ms)
Average: 0.810465 second (810ms)

##### Serverless
Transfers per second: up to 5

Min: 1.666849 second (1666ms)
Max: 13.134003 seconds (131340ms)
Average: 7.47914 seconds (7479ms)