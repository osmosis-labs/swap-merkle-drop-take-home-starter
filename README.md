## Task Description

### Summary: Fee Rebate System

**Background**

Osmosis is implementing a fee rebate system to allow entities to offer rebates based on frequently changing conditions. For example:

- The Stride team may want to offer $10,000 USDC in rebates to users who have swapped in pools containing Stride LST tokens over an epoch.
- Astroport may want to offer 10 ASTRO to any user who has swapped over $1,000 worth of value in an Astroport PCL pool in a specific timeframe.

Given these dynamic requirements and conditions, we aim to create a Merkle-drop style system that indexes swap events in every block. The system will expose an API for external fee rebate providers (like Stride and Astroport) to construct a Merkle tree of users who have met specific volume criteria between two blockchain heights. The amount to distribute and the distribution strategy (flat fee rebate or proportional to volume) will be inputs to the system.

**Informational and Out-of-Scope**

Fee rebate providers will upload the Merkle root on-chain while keeping the Merkle tree structure in their centralized storage. Clients will query their web service APIs for Merkle tree proofs and submit them on-chain. Upon verification, the fee rebate will be distributed to them according to the allocation specified by the fee rebate provider address. Note that this is not in-scope of the service we are building. Rather, this is how the service we are building will be used.

**Goals for Osmosis**

Develop a service around a database that indexes swap events, allowing fee rebate providers to easily retrieve a ready-to-use Merkle tree.

### Requirements

Create an off-chain service that:

1. Collects transaction information for every block, indexing only the swap events (see `main.go`)
2. Exposes an endpoint `GET /swap-merkle-tree` with the following query parameters:
    1. `startHeight` (int): The starting blockchain height.
    2. `endHeight` (int): The ending blockchain height.
    3. `volumeThreshold` (int): The minimum swap volume required for eligibility.
    4. `poolID` (int): The ID of the pool to consider.
    5. `strategy` (int): The distribution strategy.
        - `1` = Flat fee rebate
        - `2` = Proportional to volume
    6. `totalDistrCoin` (e.g. `100000000uosmo`) - total amount to split between users based on `strategy`
3. The service is started as a single binary with a flag `startHeight` to specify the initial height to start indexing from.
4. The service should index the data from the `startHeight` and up until the tip of the chain. Once the tip is reached, the indexing worker should continue running in the background.

By following these requirements, the off-chain service will efficiently index and provide Merkle tree data for fee rebate providers to use in their rebate distribution schemes.

For simplicity, assume that clients would be querying at max 300 height intervals, making the payload size be appropriate for the data transfer.

### Steps to Approach the Task

---

1. Open the [main.go](https://github.com/osmosis-labs/swap-merkle-drop-take-home-starter/blob/main/main.go)
    - We created a starter project for you to build upon that already queries the Osmosis API for block data
    between two heights.
    - Get familiar with TODOs
    - Run the starter `go run main.go`
    - Implement the remaining requirements of the service in Go.
2. **Data Collection:**
    - Retrieve swap event details from each block per earlier requirements.
    - Convert the token in amounts into USDC value
        - Note an additional API endpoint for getting the pricing data for converting token amounts into USDC. This is [its swagger](https://sqs.osmosis.zone/swagger/index.html#/default/get_tokens_prices).
    - Index the data in the format that you deem is appropriate for the requirements.
3. **Store Data Locally:**
    - Store the collected data in a local database (e.g., SQLite, PostgreSQL). Your choice.
4. **Expose the API Endpoint:**
    - Implement an endpoint `GET /swap-merkle-tree`
    - Refer to query parameters earlier in the document
    - Create a mapping from user address to total USDC volume in a given `poolID` between the `startHeight` and `endHeight`.
    - Filter out the users based on `volumeThreshold`
    - Create a distribution allocation for each user based on `strategy` and `totalDistrCoin`
    - Create a Merkle tree of the unique user allocations and return the full structure as part of the response.
        - Note the earlier assumption that the payload size is appropriate.

### Deliverables

---

1. Source code for the off-chain service. Create fork of the starter project.
2. Instructions on how to run the service.
3. A short design document explaining the architecture, key components, and reasoning behind major decisions. Answer questions left at the bottom of the starter.
4. Bonus:
    1. Cloud Deployment and Infrastructure Automation. Working endpoint that we can query.

### Notes

---

- Ensure your code is well-documented and follows best practices.
- Consider edge cases and error handling.
- Write unit tests for key components.
- Include deployment plan, monitoring stack and details about performance optimizations.
- Include a `RESULTS.md` file with clear instructions on how to set up and run your service.
- Feel free to use any framework of your choice.
- Feel free to structure the service in the way you see fit.
- Please do not use external Merkle tree libraries and create your own implementation.
- Feel free to use ChatGPT for any step.
- Reach out to us if you have any questions or need clarifications.

### Questions

---

- What was the motivation for the backend of your choice? Why did you structure the data in the storage per the selected way?
- Assume the limitation of querying  the REST endpoint `/swap-merkle-tree` in 300 height intervals is removed. What protocol would you use to query for the merkle tree given the possibility of payload size being large?
