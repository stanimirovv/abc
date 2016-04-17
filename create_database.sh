sudo su postgres -c 'psql -c "CREATE DATABASE abc;"'
sudo su postgres -c 'psql -d abc  -f ./db/schema/002-core.sql'
sudo su postgres -c 'psql -d abc  -f ./db/schema/999-permissions.sql'
