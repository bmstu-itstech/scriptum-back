docker cp /Users/ekaterinafedina/scriptum-back/migrations/01_init.down.sql my_postgres:/down.migrations.sql
docker cp /Users/ekaterinafedina/scriptum-back/migrations/01_init.up.sql my_postgres:/up.migrations.sql    

docker exec -it my_postgres psql -U app_user -d dev -f /down.migrations.sql
docker exec -it my_postgres psql -U app_user -d dev -f /up.migrations.sql