# wolfservers
go cli to help create various crypto systems at places like aws,
digitalocean,vultr,linode,google,and more.

# example

./wolfservers ed255

./wolfservers sqlite

./wolfservers fresh2linode --sure

./wolfservers update-ips --producer=ip1 --relay=ip2 --name=wolf-ABCD

./wolfservers ls

./wolfservers relay --producer=ip1 --relay=ip2

./wolfservers ssh --ip=ip2 --root
setup.sh
. .bashrc
relay.sh

./wolfservers producer --producer=ip1 --relay=ip2

./wolfservers ssh --ip=ip1 --root
setup.sh
. .bashrc
producer.sh


./wolfservers add-a-record --name=wolf-ABCD-1 --ip=ip1

# read more

https://andrewarrow.substack.com/p/in-order-to-bank-in-the-modern-era
