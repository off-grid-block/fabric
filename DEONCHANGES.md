## Migration to DEON Controller Package
These changes are up to date with commit [1453a06cb6bf888727fe0c1b7973d98476277cbd](https://github.com/off-grid-block/fabric/commit/1453a06cb6bf888727fe0c1b7973d98476277cbd).
- ```common/controller```
    - new package inside Fabric source code
    - provide a set of functions enabling communication with DEON ACA-Py agent instances.
    - see more details at [the DEON Github repository](https://github.com/off-grid-block/controller).

- ```core/common/msgvalidation```
    - changes inside ```ValidateProposalMessage()``` and ```ValidateTransaction()```
    - replace ```indyverify.Indyverify()``` calls with calls to ```VerifySignature()``` from ```common/controller```.
    - create admin agent controller with ```NewAdminController()```
    - retrieve admin-client agent connection details with ```GetConnection()```
    - direct admin agent to request proof from client agent with ```RequireProof()``` (inside ```ValidateProposalMessage()```)

- ```orderer/common/msgprocessor/sigfilter.go```
    - changes inside ```Apply()```
    - verify transaction signature using signing DID and communication with admin agent
    - replace ```indyverify.Indyverify()``` calls with calls to ```VerifySignature()``` from ```common/controller```.

- ```common/deliver/deliver.go```
    - events signed by client agent were verified using Indy verifier, now disabled: removed calls ```indyverify.Indyverify()```

- ```protos/common/common.proto```
    - add DID as additional parameter in transaction header

- ```core/policy/policy.go```
    - disable channel policy check for transactions signed by client agent