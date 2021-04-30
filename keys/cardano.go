package keys

/*

S = sign    == private
V = verify  == public

stake pool cold key (node.cert)
stake pool hot key (kes.skey)        Key Evolving Signature
stake pool VRF key (vrf.skey)        Verifiable random function

Determine the number of slots per KES period
kesPeriod by dividing the slot tip number by the slotsPerKESPeriod
operational certificate

1. Copy kes.vkey to your cold environment.
2. Copy node.cert to your hot environment.


Payment keys are used to send and receive payments and stake keys are used to manage stake delegations.

payment.skey & payment.vkey
stake.skey & stake.vkey

stake.addr

Build a payment address for the payment key payment.vkey which will delegate to the stake address, stake.vkey

payment.addr

4. Copy payment.addr to your hot environment.
5. Create a certificate, stake.cert, using the stake.vkey

stake.cert

6. Copy stake.cert to your hot environment.

7. Copy tx.raw to your cold environment.
8. Copy tx.signed to your hot environment.

pool.cert

Copy pool.cert to your hot environment.

deleg.cert

Copy deleg.cert to your hot environment.

Copy tx.raw to your cold environment.
Copy tx.signed to your hot environment.

Copy stakepoolid.txt to your hot environment.


./topologyUpdater.sh

Complete this section after four hours when your relay node IP is properly registered.

./relay-topology_pull.sh


*/
