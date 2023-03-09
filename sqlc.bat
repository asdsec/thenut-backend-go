:: This file has been written for windows user.
:: Because sqlc tool does not work properly with Git Bash.
:: In the project path, you can run `sqlc.bat` via cmd.
:: To initialize sqlc, you can run 'sqlc init' in cmd.
:: Or it is possible to use 'Makefile'.

docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc %*