sudo su postgres -c 'psql -c "CREATE DATABASE abc_test;"'
sudo su postgres -c 'psql -d abc_test  -f ./db/schema/002-core.sql'
sudo su postgres -c 'psql -d abc_test  -f ./db/schema/999-perms.sql'
