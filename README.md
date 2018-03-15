# Exercise in fault tolerance #

Take a server with a single bank account.

Expects multiple clients asking for balance or movements over TCP.

Collect all balances of individual clients that make transactions on behalf of that account.

Sum and collect a total balance.

Check if summed balance matches account balance.
