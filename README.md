# wolfservers
go cli to help create various crypto systems at places like aws,
digitalocean,vultr,linode,google,and more.

# about

![image](https://cdn.substack.com/image/fetch/w_1456,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F0dd7a8b4-77b6-4859-88fd-510c105a16fc_1280x696.jpeg)

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
