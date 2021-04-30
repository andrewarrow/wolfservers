# wolfservers
go cli to help create cardano proof-of-stake nodes at places like aws,
digitalocean,vultr,linode,google,and more.

# about

![image](https://cdn.substack.com/image/fetch/w_1456,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F0dd7a8b4-77b6-4859-88fd-510c105a16fc_1280x696.jpeg)

You want to be running these proof-of-stake nodes on the new blockchain.

# FAQ

Q. Why not just use coinbase?
A. You can and many people should. But if you are a programmer and grok SSH keys, or think you can learn, you'll make more money this way.

Q. New blockchain?
A. Yeah it's [cardano, aka ADA](https://roadmap.cardano.org/) and it changes the whole bitcoin or etherum mining concepts.

Q. Can I mess this up?
A. Absolutely. [read this](https://andrewarrow.substack.com/p/in-order-to-bank-in-the-modern-era)

Q. How does wolfservers store my keys?
A. We use a sqlite database on your local hard drive but all the private keys we write to it are encrypted with a > 36 character phrase you have to memorize. Each morning I open my .bash_profile and write `export WOLF_PHRASE="something long and very secret and definitely something I will never forget"` and then every night I erase it from my .bash_profile. The actual wolf.db sqlite file is safe to email to yourself, store on dropbox, etc. You want many copies of this wolf.db file. You can never ever lose it.

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
