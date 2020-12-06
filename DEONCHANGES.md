## Migration to DEON Controller Package
These changes are up to date with commit [1453a06cb6bf888727fe0c1b7973d98476277cbd](https://github.com/off-grid-block/fabric/commit/1453a06cb6bf888727fe0c1b7973d98476277cbd).
- add controller package in ```common/controller```
    - provide a set of functions enabling communication with DEON ACA-Py agent instances.
    - see more details at [the DEON Github repository](https://github.com/off-grid-block/controller).

- changes in ```core/common/msgvalidation```
    - changes inside ```ValidateProposalMessage()```
    - replace ```indyverify.Indyverify()``` calls with calls to ```VerifySignature()``` from ```common/controller```.
    - create admin agent controller with ```NewAdminController()```
    - retrieve admin-client agent connection details with ```GetConnection()```
    - direct admin agent to request proof from client agent with ```RequireProof()```.

- changes in ```orderer/common/msgprocessor/sigfilter.go```
    - changes inside ```Apply()```
    - replace ```indyverify.Indyverify()``` calls with calls to ```VerifySignature()``` from ```common/controller```.

- changes in ```common/deliver/deliver.go```
    - remove calls to ```indyverify.Indyverify()```
