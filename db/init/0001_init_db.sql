CREATE USER keeper
    PASSWORD 'password';

CREATE DATABASE keeper
    OWNER keeper
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';